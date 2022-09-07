package module_node

import (
	"go.uber.org/zap"
	"strconv"
	"sync"
	"teamide/pkg/node"
	"teamide/pkg/util"
)

func (this_ *NodeService) InitContext() {
	if this_.nodeContext == nil {
		this_.nodeContext = &NodeContext{
			NodeService: this_,

			nodeIdModelCache:   make(map[int64]*NodeModel),
			serverIdModelCache: make(map[string]*NodeModel),

			netProxyIdModelCache: make(map[int64]*NetProxyModel),
			codeModelCache:       make(map[string]*NetProxyModel),

			lineNodeIdListCache: make(map[string][]string),

			countData: newNodeCountData(),
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
	server *node.Server
	*NodeService
	root *NodeModel

	oldNodeListStr         string
	nodeModelIdList        []int64
	nodeIdModelCache       map[int64]*NodeModel
	nodeIdModelCacheLock   sync.Mutex
	serverIdModelCache     map[string]*NodeModel
	serverIdModelCacheLock sync.Mutex

	lineNodeIdListCache     map[string][]string
	lineNodeIdListCacheLock sync.Mutex

	oldNetProxyListStr       string
	netProxyModelIdList      []int64
	netProxyIdModelCache     map[int64]*NetProxyModel
	netProxyIdModelCacheLock sync.Mutex
	codeModelCache           map[string]*NetProxyModel
	codeModelCacheLock       sync.Mutex
	runTimerLock             sync.Mutex
	runTimerRunning          bool
	onNetProxyListChangeIng  bool
	onNodeListChangeIng      bool

	doAliveIng  bool
	doAliveLock sync.Mutex

	oldCountDataStr string
	countData       *NodeCountData
}

func (this_ *NodeContext) GetServer() *node.Server {
	return this_.server

}

func (this_ *NodeContext) cleanNodeLine() {
	this_.lineNodeIdListCacheLock.Lock()
	defer this_.lineNodeIdListCacheLock.Unlock()

	this_.Logger.Info("cleanNodeLine")
	this_.lineNodeIdListCache = make(map[string][]string)

}

func (this_ *NodeContext) GetNodeLineTo(nodeId string) (lineIdList []string) {
	lineIdList = this_.GetNodeLineByFromTo(this_.root.ServerId, nodeId)
	return
}

func (this_ *NodeContext) GetNodeLineByFromTo(fromNodeId, toNodeId string) (lineIdList []string) {
	this_.lineNodeIdListCacheLock.Lock()
	defer this_.lineNodeIdListCacheLock.Unlock()

	key := "from:" + fromNodeId + " to:" + toNodeId
	lineIdList, ok := this_.lineNodeIdListCache[key]
	if ok {
		return
	}

	var nodeIdConnNodeIdListCache = make(map[string][]string)
	var list = this_.nodeModelIdList
	for _, one := range list {
		find := this_.getNodeModel(one)
		if find == nil {
			continue
		}
		nodeIdConnNodeIdListCache[find.ServerId] = find.ConnServerIdList
	}
	this_.Logger.Info("GetNodeLineByFromTo", zap.Any("nodeIdConnNodeIdListCache", nodeIdConnNodeIdListCache))
	lineIdList = getNodeLineByFromTo(fromNodeId, toNodeId, nodeIdConnNodeIdListCache)
	if len(lineIdList) > 0 &&
		util.ContainsString(lineIdList, fromNodeId) >= 0 &&
		util.ContainsString(lineIdList, toNodeId) >= 0 {
		this_.lineNodeIdListCache[key] = lineIdList
	}

	return
}

func (this_ *NodeContext) initContext() (err error) {
	var list []*NodeModel
	list, _ = this_.Query(&NodeModel{})
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
			this_.addNodeModel(one)
		}
	}

	var netProxyList []*NetProxyModel
	netProxyList, _ = this_.QueryNetProxy(&NetProxyModel{})
	for _, one := range netProxyList {
		this_.addNetProxyModel(one)
	}

	go this_.doAlive()
	return
}

func (this_ *NodeContext) initRoot(root *NodeModel) (err error) {
	if this_.server != nil {
		this_.server.Stop()
	}
	this_.root = root
	this_.server = &node.Server{
		Id:          this_.root.ServerId,
		BindToken:   this_.root.BindToken,
		BindAddress: this_.root.BindAddress,
	}
	err = this_.server.Start()
	if err != nil {
		return
	}

	return
}

func appendNodeLineList(loadedIdList *[]string, lineList *[][]string, parentLine []string, nodeIdList []string, nodeIdConnNodeIdListCache map[string][]string) {

	for _, nodeId := range nodeIdList {
		var line []string
		line = append(line, parentLine...)

		if util.ContainsString(line, nodeId) >= 0 {
			continue
		}
		line = append(line, nodeId)

		*lineList = append(*lineList, line)

		if util.ContainsString(*loadedIdList, nodeId) >= 0 {
			continue
		}
		*loadedIdList = append(*loadedIdList, nodeId)

		var connNodeIdList = nodeIdConnNodeIdListCache[nodeId]
		appendNodeLineList(loadedIdList, lineList, line, connNodeIdList, nodeIdConnNodeIdListCache)

		var parentIdList []string
		for cacheNodeId, cacheConnNodeIdList := range nodeIdConnNodeIdListCache {
			if util.ContainsString(cacheConnNodeIdList, nodeId) >= 0 {
				parentIdList = append(parentIdList, cacheNodeId)
			}
		}
		appendNodeLineList(loadedIdList, lineList, line, parentIdList, nodeIdConnNodeIdListCache)
	}
}

func findNodeLineList(nodeId string, nodeIdConnNodeIdListCache map[string][]string) (lineList [][]string) {

	var loadedIdList []string
	loadedIdList = append(loadedIdList, nodeId)
	var line []string
	line = append(line, nodeId)

	var connNodeIdList = nodeIdConnNodeIdListCache[nodeId]
	appendNodeLineList(&loadedIdList, &lineList, line, connNodeIdList, nodeIdConnNodeIdListCache)

	var parentIdList []string
	for cacheNodeId, cacheConnNodeIdList := range nodeIdConnNodeIdListCache {
		if util.ContainsString(cacheConnNodeIdList, nodeId) >= 0 {
			parentIdList = append(parentIdList, cacheNodeId)
		}
	}
	appendNodeLineList(&loadedIdList, &lineList, line, parentIdList, nodeIdConnNodeIdListCache)

	return
}

func getNodeLineByFromTo(fromNodeId, toNodeId string, nodeIdConnNodeIdListCache map[string][]string) (lineIdList []string) {

	if fromNodeId == toNodeId {
		lineIdList = append(lineIdList, fromNodeId)
		return
	}

	var lineList = findNodeLineList(fromNodeId, nodeIdConnNodeIdListCache)

	for _, line := range lineList {
		if util.ContainsString(line, toNodeId) >= 0 {
			if len(lineIdList) == 0 || len(line) < len(lineIdList) {
				lineIdList = line
			}
		}
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
		monitorData = &node.MonitorData{}
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
