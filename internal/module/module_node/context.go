package module_node

import (
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"teamide/node"
)

func (this_ *NodeService) InitContext() {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			nodeService: this_,
			wsCache:     make(map[string]*WSConn),

			nodeIdModelCache:   make(map[int64]*NodeModel),
			serverIdModelCache: make(map[string]*NodeModel),

			netProxyIdModelCache: make(map[int64]*NetProxyModel),
			codeModelCache:       make(map[string]*NetProxyModel),
		}
	}
	err := this_.nodeContext.initContext()
	if err != nil {
		this_.Logger.Error("节点上下文初始化异常", zap.Error(err))
		return
	}
	return
}

type NodeContext struct {
	server      *node.Server
	nodeService *NodeService
	root        *NodeModel

	wsCache     map[string]*WSConn
	wsCacheLock sync.Mutex

	nodeList               []*NodeInfo
	nodeListLock           sync.Mutex
	nodeIdModelCache       map[int64]*NodeModel
	nodeIdModelCacheLock   sync.Mutex
	serverIdModelCache     map[string]*NodeModel
	serverIdModelCacheLock sync.Mutex

	netProxyList             []*NetProxyInfo
	netProxyListLock         sync.Mutex
	netProxyIdModelCache     map[int64]*NetProxyModel
	netProxyIdModelCacheLock sync.Mutex
	codeModelCache           map[string]*NetProxyModel
	codeModelCacheLock       sync.Mutex
	runTimerLock             sync.Mutex
	runTimerRunning          bool
	onNetProxyListChangeIng  bool
	onNodeListChangeIng      bool
}

func (this_ *NodeContext) initContext() (err error) {
	go this_.runTimer()
	var list []*NodeModel
	list, _ = this_.nodeService.Query(&NodeModel{})
	for _, one := range list {
		this_.setNodeModel(one.NodeId, one)
		this_.setNodeModelByServerId(one.ServerId, one)
		if one.IsROOT() {
			this_.root = one
		}
	}

	if this_.root != nil {
		this_.onAddNodeModel(this_.root)
	}
	for _, one := range list {
		if !one.IsROOT() {
			this_.onAddNodeModel(one)
		}
	}

	var netProxyList []*NetProxyModel
	netProxyList, _ = this_.nodeService.QueryNetProxy(&NetProxyModel{})
	for _, one := range netProxyList {
		this_.onAddNetProxyModel(one)
	}

	return
}

func (this_ *NodeContext) initRoot(root *NodeModel) (err error) {
	if this_.server != nil {
		this_.server.Stop()
	}
	this_.root = root
	this_.server = &node.Server{
		Id:                   this_.root.ServerId,
		BindToken:            this_.root.BindToken,
		BindAddress:          this_.root.BindAddress,
		OnNodeListChange:     this_.onNodeListChange,
		OnNetProxyListChange: this_.onNetProxyListChange,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}

	return
}

func (this_ *NodeContext) runTimer() {
	if this_.runTimerRunning {
		return
	}
	this_.runTimerLock.Lock()
	defer this_.runTimerLock.Unlock()
	this_.runTimerRunning = true
	defer func() {
		//go func() {
		//	time.Sleep(time.Second * 10)
		//	this_.runTimerRunning = false
		//	go this_.runTimer()
		//}()
	}()
	var bs []byte
	var nodeList []*NodeInfo
	if !this_.onNodeListChangeIng && len(this_.nodeList) > 0 {
		nodeList = this_.nodeList
		bs, _ = json.Marshal(nodeList)
		oldStr := string(bs)
		for _, one := range nodeList {
			one.MonitorData = ToMonitorDataFormat(this_.server.GetNodeMonitorData(one.Info.Id))
		}
		bs, _ = json.Marshal(nodeList)
		newStr := string(bs)
		if newStr == oldStr {
			nodeList = []*NodeInfo{}
		}
	}
	var netProxyList []*NetProxyInfo
	if !this_.onNetProxyListChangeIng && len(this_.netProxyList) > 0 {
		netProxyList = this_.netProxyList
		bs, _ = json.Marshal(netProxyList)
		oldStr := string(bs)
		for _, one := range netProxyList {
			one.InnerMonitorData = ToMonitorDataFormat(this_.server.GetNetProxyInnerMonitorData(one.Info.Id))
			one.OuterMonitorData = ToMonitorDataFormat(this_.server.GetNetProxyOuterMonitorData(one.Info.Id))
		}
		bs, _ = json.Marshal(netProxyList)
		newStr := string(bs)
		if newStr == oldStr {
			netProxyList = []*NetProxyInfo{}
		}
	}

	if len(nodeList) > 0 || len(netProxyList) > 0 {
		this_.callMessage(&Message{
			Method:       "refresh_node_context",
			NodeList:     nodeList,
			NetProxyList: netProxyList,
		})
	}
	return
}

type MonitorDataFormat struct {
	ReadSize     string `json:"readSize,omitempty"`
	ReadSizeUnit string `json:"readSizeUnit,omitempty"`
	ReadTime     string `json:"readTime,omitempty"`
	ReadTimeUnit string `json:"readTimeUnit,omitempty"`
	ReadSleep    string `json:"readSleep,omitempty"`

	ReadLastSize      string `json:"readLastSize,omitempty"`
	ReadLastSizeUnit  string `json:"readLastSizeUnit,omitempty"`
	ReadLastTime      string `json:"readLastTime,omitempty"`
	ReadLastTimeUnit  string `json:"readLastTimeUnit,omitempty"`
	ReadLastTimestamp int64  `json:"readLastTimestamp,omitempty"`
	ReadLastSleep     string `json:"readLastSleep,omitempty"`
	ReadLastSleepUnit string `json:"readLastSleepUnit,omitempty"`

	WriteSize     string `json:"writeSize,omitempty"`
	WriteSizeUnit string `json:"writeSizeUnit,omitempty"`
	WriteTime     string `json:"writeTime,omitempty"`
	WriteTimeUnit string `json:"writeTimeUnit,omitempty"`
	WriteSleep    string `json:"writeSleep,omitempty"`

	WriteLastSize      string `json:"writeLastSize,omitempty"`
	WriteLastSizeUnit  string `json:"writeLastSizeUnit,omitempty"`
	WriteLastTime      string `json:"writeLastTime,omitempty"`
	WriteLastTimeUnit  string `json:"writeLastTimeUnit,omitempty"`
	WriteLastTimestamp int64  `json:"writeLastTimestamp,omitempty"`
	WriteLastSleep     string `json:"writeLastSleep,omitempty"`
	WriteLastSleepUnit string `json:"writeLastSleepUnit,omitempty"`
}

var (
	KBSize float64 = 1024
	MBSize float64 = 1024 * 1024
	GBSize float64 = 1024 * 1024 * 1024
)

func GetSizeAndUnit(size float64) (res float64, unit string) {
	res = size
	unit = "B"
	if res > GBSize {
		res = res / GBSize
		unit = "GB"
	} else if res > MBSize {
		res = res / MBSize
		unit = "MB"
	} else if res > KBSize {
		res = res / KBSize
		unit = "KB"
	}
	return
}

func ToMonitorDataFormat(monitorData *node.MonitorData) *MonitorDataFormat {
	if monitorData == nil {
		return nil
	}

	ReadSize, ReadSizeUnit := GetSizeAndUnit(float64(monitorData.ReadSize))
	ReadTime := float64(monitorData.ReadTime) / 1000000000
	ReadTimeUnit := "秒"
	ReadSleep := float64(0)
	if ReadSize > 0 && ReadTime > 0 {
		ReadSleep = ReadSize / ReadTime
	}

	ReadLastSize, ReadLastSizeUnit := GetSizeAndUnit(float64(monitorData.ReadLastSize))
	ReadLastTime := float64(monitorData.ReadLastTime) / 1000000000
	ReadLastTimeUnit := "秒"
	ReadLastSleep := float64(0)
	ReadLastSleepUnit := "B/秒"
	if monitorData.ReadLastSize > 0 && ReadLastTime > 0 {
		ReadLastSleep, ReadLastSleepUnit = GetSizeAndUnit(float64(monitorData.ReadLastSize) / ReadLastTime)
		ReadLastSleepUnit = ReadLastSleepUnit + "/秒"
	}

	WriteSize, WriteSizeUnit := GetSizeAndUnit(float64(monitorData.WriteSize))
	WriteTime := float64(monitorData.WriteTime) / 1000000000
	WriteTimeUnit := "秒"
	WriteSleep := float64(0)
	if WriteSize > 0 && WriteTime > 0 {
		WriteSleep = WriteSize / WriteTime
	}

	WriteLastSize, WriteLastSizeUnit := GetSizeAndUnit(float64(monitorData.WriteLastSize))
	WriteLastTime := float64(monitorData.WriteLastTime) / 1000000000
	WriteLastTimeUnit := "秒"
	WriteLastSleep := float64(0)
	WriteLastSleepUnit := "B/秒"
	if monitorData.WriteLastSize > 0 && WriteLastTime > 0 {
		WriteLastSleep, WriteLastSleepUnit = GetSizeAndUnit(float64(monitorData.WriteLastSize) / WriteLastTime)
		WriteLastSleepUnit = WriteLastSleepUnit + "/秒"
	}

	return &MonitorDataFormat{
		ReadSize:     strconv.FormatFloat(ReadSize, 'f', 2, 64),
		ReadSizeUnit: ReadSizeUnit,
		ReadTime:     strconv.FormatFloat(ReadTime, 'f', 2, 64),
		ReadTimeUnit: ReadTimeUnit,
		ReadSleep:    strconv.FormatFloat(ReadSleep, 'f', 2, 64),

		ReadLastSize:      strconv.FormatFloat(ReadLastSize, 'f', 2, 64),
		ReadLastSizeUnit:  ReadLastSizeUnit,
		ReadLastTime:      strconv.FormatFloat(ReadLastTime, 'f', 2, 64),
		ReadLastTimeUnit:  ReadLastTimeUnit,
		ReadLastTimestamp: monitorData.ReadLastTimestamp,
		ReadLastSleep:     strconv.FormatFloat(ReadLastSleep, 'f', 2, 64),
		ReadLastSleepUnit: ReadLastSleepUnit,

		WriteSize:     strconv.FormatFloat(WriteSize, 'f', 2, 64),
		WriteSizeUnit: WriteSizeUnit,
		WriteTime:     strconv.FormatFloat(WriteTime, 'f', 2, 64),
		WriteTimeUnit: WriteTimeUnit,
		WriteSleep:    strconv.FormatFloat(WriteSleep, 'f', 2, 64),

		WriteLastSize:      strconv.FormatFloat(WriteLastSize, 'f', 2, 64),
		WriteLastSizeUnit:  WriteLastSizeUnit,
		WriteLastTime:      strconv.FormatFloat(WriteLastTime, 'f', 2, 64),
		WriteLastTimeUnit:  WriteLastTimeUnit,
		WriteLastTimestamp: monitorData.WriteLastTimestamp,
		WriteLastSleep:     strconv.FormatFloat(WriteLastSleep, 'f', 2, 64),
		WriteLastSleepUnit: WriteLastSleepUnit,
	}
}
