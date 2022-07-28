package module_node

import (
	"go.uber.org/zap"
	"strconv"
	"sync"
	"teamide/node"
	"time"
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
	this_.runTimer()
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
		go func() {
			time.Sleep(time.Second * 5)
			this_.runTimerRunning = false
			go this_.runTimer()
		}()
	}()

	if !this_.onNodeListChangeIng {
		this_.refreshNodeList(this_.nodeList)
	}
	if !this_.onNetProxyListChangeIng {
		this_.refreshNetProxyList(this_.netProxyList)
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

	ReadLastSize, ReadLastSizeUnit := GetSizeAndUnit(float64(monitorData.ReadLastSize))

	ReadLastTime := float64(monitorData.ReadLastTime) / 1000000000
	ReadLastTimeUnit := "秒"

	WriteSize, WriteSizeUnit := GetSizeAndUnit(float64(monitorData.WriteSize))

	WriteTime := float64(monitorData.WriteTime) / 1000000000
	WriteTimeUnit := "秒"

	WriteLastSize, WriteLastSizeUnit := GetSizeAndUnit(float64(monitorData.WriteLastSize))

	WriteLastTime := float64(monitorData.WriteLastTime) / 1000000000
	WriteLastTimeUnit := "秒"

	return &MonitorDataFormat{
		ReadSize:     strconv.FormatFloat(ReadSize, 'f', 2, 64),
		ReadSizeUnit: ReadSizeUnit,
		ReadTime:     strconv.FormatFloat(ReadTime, 'f', 2, 64),
		ReadTimeUnit: ReadTimeUnit,
		ReadSleep:    strconv.FormatFloat(ReadSize/ReadTime, 'f', 2, 64),

		ReadLastSize:      strconv.FormatFloat(ReadLastSize, 'f', 2, 64),
		ReadLastSizeUnit:  ReadLastSizeUnit,
		ReadLastTime:      strconv.FormatFloat(ReadLastTime, 'f', 2, 64),
		ReadLastTimeUnit:  ReadLastTimeUnit,
		ReadLastTimestamp: monitorData.ReadLastTimestamp,
		ReadLastSleep:     strconv.FormatFloat(ReadLastSize/ReadLastTime, 'f', 2, 64),

		WriteSize:     strconv.FormatFloat(WriteSize, 'f', 2, 64),
		WriteSizeUnit: WriteSizeUnit,
		WriteTime:     strconv.FormatFloat(WriteTime, 'f', 2, 64),
		WriteTimeUnit: WriteTimeUnit,
		WriteSleep:    strconv.FormatFloat(WriteSize/WriteTime, 'f', 2, 64),

		WriteLastSize:      strconv.FormatFloat(WriteLastSize, 'f', 2, 64),
		WriteLastSizeUnit:  WriteLastSizeUnit,
		WriteLastTime:      strconv.FormatFloat(WriteLastTime, 'f', 2, 64),
		WriteLastTimeUnit:  WriteLastTimeUnit,
		WriteLastTimestamp: monitorData.WriteLastTimestamp,
		WriteLastSleep:     strconv.FormatFloat(WriteLastSize/WriteLastTime, 'f', 2, 64),
	}
}
