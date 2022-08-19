package module_node

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/pkg/system"
	"teamide/pkg/util"
)

type NodeApi struct {
	*context.ServerContext
	NodeService *NodeService
}

func NewNodeApi(NodeService *NodeService) *NodeApi {
	return &NodeApi{
		ServerContext: NodeService.ServerContext,
		NodeService:   NodeService,
	}
}

var (
	// 节点 权限

	// PowerNode 节点基本 权限
	PowerNode             = base.AppendPower(&base.PowerAction{Action: "node", Text: "节点", ShouldLogin: false, StandAlone: true})
	PowerNodeContext      = base.AppendPower(&base.PowerAction{Action: "node_context", Text: "节点列表", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeList         = base.AppendPower(&base.PowerAction{Action: "node_list", Text: "节点列表", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeStart        = base.AppendPower(&base.PowerAction{Action: "node_start", Text: "节点启动", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeStop         = base.AppendPower(&base.PowerAction{Action: "node_stop", Text: "节点停止", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeInsert       = base.AppendPower(&base.PowerAction{Action: "node_insert", Text: "节点新增", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeUpdate       = base.AppendPower(&base.PowerAction{Action: "node_update", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeUpdateOption = base.AppendPower(&base.PowerAction{Action: "node_update_option", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeEnable       = base.AppendPower(&base.PowerAction{Action: "node_enable", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeDisable      = base.AppendPower(&base.PowerAction{Action: "node_disable", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeDelete       = base.AppendPower(&base.PowerAction{Action: "node_delete", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})

	PowerNetProxyList         = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_list", Text: "节点新增", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyInsert       = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_insert", Text: "节点新增", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyUpdate       = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_update", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyUpdateOption = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_update_option", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyMonitorData  = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_monitor_data", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyEnable       = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_enable", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyDisable      = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_disable", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNetProxyDelete       = base.AppendPower(&base.PowerAction{Action: "node_net_proxy_delete", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})

	PowerNodeSystemInfo             = base.AppendPower(&base.PowerAction{Action: "node_system_info", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeSystemQueryMonitorData = base.AppendPower(&base.PowerAction{Action: "node_system_queryMonitorData", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeSystemCleanMonitorData = base.AppendPower(&base.PowerAction{Action: "node_system_cleanMonitorData", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
)

func (this_ *NodeApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"node"}, Power: PowerNode, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/context"}, Power: PowerNodeContext, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/start"}, Power: PowerNodeStart, Do: this_.start})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/stop"}, Power: PowerNodeStop, Do: this_.stop})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/list"}, Power: PowerNodeList, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/insert"}, Power: PowerNodeInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/update"}, Power: PowerNodeUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/updateOption"}, Power: PowerNodeUpdateOption, Do: this_.updateOption})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/enable"}, Power: PowerNodeEnable, Do: this_.enable})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/disable"}, Power: PowerNodeDisable, Do: this_.disable})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/delete"}, Power: PowerNodeDelete, Do: this_.delete})

	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/list"}, Power: PowerNetProxyList, Do: this_.netProxyList})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/insert"}, Power: PowerNetProxyInsert, Do: this_.netProxyInsert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/update"}, Power: PowerNetProxyUpdate, Do: this_.netProxyUpdate})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/updateOption"}, Power: PowerNetProxyUpdateOption, Do: this_.netProxyUpdateOption})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/monitorData"}, Power: PowerNetProxyMonitorData, Do: this_.netProxyMonitorData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/enable"}, Power: PowerNetProxyEnable, Do: this_.netProxyEnable})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/disable"}, Power: PowerNetProxyDisable, Do: this_.netProxyDisable})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/netProxy/delete"}, Power: PowerNetProxyDelete, Do: this_.netProxyDelete})

	apis = append(apis, &base.ApiWorker{Apis: []string{"node/system/info"}, Power: PowerNodeSystemInfo, Do: this_.nodeSystemInfo})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/system/queryMonitorData"}, Power: PowerNodeSystemQueryMonitorData, Do: this_.nodeSystemQueryMonitorData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/system/cleanMonitorData"}, Power: PowerNodeSystemCleanMonitorData, Do: this_.nodeSystemCleanMonitorData})

	return
}

type IndexResponse struct {
}

func (this_ *NodeApi) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {

	response := &IndexResponse{}

	res = response
	return
}

type ContextResponse struct {
	LocalIpList  []string         `json:"localIpList,omitempty"`
	CountData    *NodeCountData   `json:"countData,omitempty"`
	NodeList     []*NodeModel     `json:"nodeList,omitempty"`
	NetProxyList []*NetProxyModel `json:"netProxyList,omitempty"`
}

func (this_ *NodeApi) context(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {

	response := &ContextResponse{}

	ipList := util.GetLocalIPList()
	for _, ip := range ipList {
		response.LocalIpList = append(response.LocalIpList, ip.String())
	}

	response.CountData = this_.NodeService.nodeContext.countData
	response.NodeList = this_.NodeService.nodeContext.nodeModelList
	response.NetProxyList = this_.NodeService.nodeContext.netProxyModelList
	res = response
	return
}

type StartRequest struct {
}

type StartResponse struct {
}

func (this_ *NodeApi) start(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &StartRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &StartResponse{}

	if err != nil {
		return
	}

	res = response
	return
}

type StopRequest struct {
}

type StopResponse struct {
}

func (this_ *NodeApi) stop(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &StopRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &StopResponse{}

	if err != nil {
		return
	}

	res = response
	return
}

type NodeSystemRequest struct {
	NodeId string `json:"nodeId,omitempty"`
	*system.QueryRequest
}

func (this_ *NodeApi) nodeSystemInfo(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NodeSystemRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if err != nil {
		return
	}

	res = this_.NodeService.nodeContext.SystemGetInfo(request.NodeId)
	return
}

func (this_ *NodeApi) nodeSystemQueryMonitorData(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NodeSystemRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if err != nil {
		return
	}

	res = this_.NodeService.nodeContext.SystemQueryMonitorData(request.NodeId, request.QueryRequest)
	return
}

func (this_ *NodeApi) nodeSystemCleanMonitorData(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &NodeSystemRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if err != nil {
		return
	}

	this_.NodeService.nodeContext.SystemCleanMonitorData(request.NodeId)
	return
}
