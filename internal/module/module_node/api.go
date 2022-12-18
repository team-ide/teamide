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
	PowerNode         = base.AppendPower(&base.PowerAction{Action: "node", Text: "节点", ShouldLogin: true, StandAlone: true})
	contextPower      = base.AppendPower(&base.PowerAction{Action: "context", Text: "节点上下文", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	listPower         = base.AppendPower(&base.PowerAction{Action: "list", Text: "节点列表", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	startPower        = base.AppendPower(&base.PowerAction{Action: "start", Text: "节点启动", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	stopPower         = base.AppendPower(&base.PowerAction{Action: "stop", Text: "节点停止", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	insertPower       = base.AppendPower(&base.PowerAction{Action: "insert", Text: "节点新增", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	updatePower       = base.AppendPower(&base.PowerAction{Action: "update", Text: "节点修改", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	updateOptionPower = base.AppendPower(&base.PowerAction{Action: "updateOption", Text: "节点修改配置", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	enablePower       = base.AppendPower(&base.PowerAction{Action: "enable", Text: "节点启用", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	disablePower      = base.AppendPower(&base.PowerAction{Action: "disable", Text: "节点停用", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	deletePower       = base.AppendPower(&base.PowerAction{Action: "delete", Text: "节点删除", Parent: PowerNode, ShouldLogin: true, StandAlone: true})

	systemPower                 = base.AppendPower(&base.PowerAction{Action: "system", Text: "节点服务器信息", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	systemInfoPower             = base.AppendPower(&base.PowerAction{Action: "info", Text: "节点服务器信息", Parent: systemPower, ShouldLogin: true, StandAlone: true})
	systemMonitorDataPower      = base.AppendPower(&base.PowerAction{Action: "monitorData", Text: "节点服务器监控数据", Parent: systemPower, ShouldLogin: false, StandAlone: true})
	systemCleanMonitorDataPower = base.AppendPower(&base.PowerAction{Action: "cleanMonitorData", Text: "节点服务器清理监控数据", Parent: systemPower, ShouldLogin: true, StandAlone: true})

	PowerNetProxy             = base.AppendPower(&base.PowerAction{Action: "netProxy", Text: "节点代理", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	netProxyListPower         = base.AppendPower(&base.PowerAction{Action: "list", Text: "节点代理列表", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyInsertPower       = base.AppendPower(&base.PowerAction{Action: "insert", Text: "节点代理新增", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyUpdatePower       = base.AppendPower(&base.PowerAction{Action: "update", Text: "节点代理修改", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyUpdateOptionPower = base.AppendPower(&base.PowerAction{Action: "updateOption", Text: "节点代理修改配置", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyMonitorDataPower  = base.AppendPower(&base.PowerAction{Action: "monitorData", Text: "节点代理监控数据", Parent: PowerNetProxy, ShouldLogin: false, StandAlone: true})
	netProxyEnablePower       = base.AppendPower(&base.PowerAction{Action: "enable", Text: "节点代理启用", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyDisablePower      = base.AppendPower(&base.PowerAction{Action: "disable", Text: "节点代理停用", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
	netProxyDeletePower       = base.AppendPower(&base.PowerAction{Action: "delete", Text: "节点代理删除", Parent: PowerNetProxy, ShouldLogin: true, StandAlone: true})
)

func (this_ *NodeApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: contextPower, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Power: listPower, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Power: startPower, Do: this_.start})
	apis = append(apis, &base.ApiWorker{Power: stopPower, Do: this_.stop})
	apis = append(apis, &base.ApiWorker{Power: insertPower, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Power: updatePower, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Power: updateOptionPower, Do: this_.updateOption})
	apis = append(apis, &base.ApiWorker{Power: enablePower, Do: this_.enable})
	apis = append(apis, &base.ApiWorker{Power: disablePower, Do: this_.disable})
	apis = append(apis, &base.ApiWorker{Power: deletePower, Do: this_.delete})

	apis = append(apis, &base.ApiWorker{Power: systemInfoPower, Do: this_.nodeSystemInfo})
	apis = append(apis, &base.ApiWorker{Power: systemMonitorDataPower, Do: this_.nodeSystemQueryMonitorData, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: systemCleanMonitorDataPower, Do: this_.nodeSystemCleanMonitorData})

	apis = append(apis, &base.ApiWorker{Power: netProxyListPower, Do: this_.netProxyList})
	apis = append(apis, &base.ApiWorker{Power: netProxyInsertPower, Do: this_.netProxyInsert})
	apis = append(apis, &base.ApiWorker{Power: netProxyUpdatePower, Do: this_.netProxyUpdate})
	apis = append(apis, &base.ApiWorker{Power: netProxyUpdateOptionPower, Do: this_.netProxyUpdateOption})
	apis = append(apis, &base.ApiWorker{Power: netProxyMonitorDataPower, Do: this_.netProxyMonitorData, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: netProxyEnablePower, Do: this_.netProxyEnable})
	apis = append(apis, &base.ApiWorker{Power: netProxyDisablePower, Do: this_.netProxyDisable})
	apis = append(apis, &base.ApiWorker{Power: netProxyDeletePower, Do: this_.netProxyDelete})

	return
}

type ContextResponse struct {
	LocalIpList  []string         `json:"localIpList,omitempty"`
	NodeList     []*NodeModel     `json:"nodeList,omitempty"`
	NetProxyList []*NetProxyModel `json:"netProxyList,omitempty"`
}

func (this_ *NodeApi) context(requestBean *base.RequestBean, _ *gin.Context) (res interface{}, err error) {

	response := &ContextResponse{}

	ipList := util.GetLocalIPList()
	for _, ip := range ipList {
		response.LocalIpList = append(response.LocalIpList, ip.String())
	}

	var nodeModelList = this_.NodeService.nodeContext.getUserNodeModelList(requestBean.JWT.UserId)
	var netProxyModelList = this_.NodeService.nodeContext.getUserNetProxyModelList(requestBean.JWT.UserId)

	response.NodeList = nodeModelList
	response.NetProxyList = netProxyModelList
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
