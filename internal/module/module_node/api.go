package module_node

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
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
	PowerNode       = base.AppendPower(&base.PowerAction{Action: "node", Text: "节点", ShouldLogin: false, StandAlone: true})
	PowerNodePage   = base.AppendPower(&base.PowerAction{Action: "node_page", Text: "节点页面", Parent: PowerNode, ShouldLogin: true, StandAlone: true})
	PowerNodeList   = base.AppendPower(&base.PowerAction{Action: "node_list", Text: "节点列表", Parent: PowerNodePage, ShouldLogin: true, StandAlone: true})
	PowerNodeInsert = base.AppendPower(&base.PowerAction{Action: "node_insert", Text: "节点新增", Parent: PowerNodePage, ShouldLogin: true, StandAlone: true})
	PowerNodeUpdate = base.AppendPower(&base.PowerAction{Action: "node_update", Text: "节点修改", Parent: PowerNodePage, ShouldLogin: true, StandAlone: true})
	PowerNodeDelete = base.AppendPower(&base.PowerAction{Action: "node_delete", Text: "节点删除", Parent: PowerNodePage, ShouldLogin: true, StandAlone: true})
)

func (this_ *NodeApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"node"}, Power: PowerNode, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/page"}, Power: PowerNodePage, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/list"}, Power: PowerNodeList, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/insert"}, Power: PowerNodeInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/update"}, Power: PowerNodeUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Apis: []string{"node/delete"}, Power: PowerNodeDelete, Do: this_.delete})

	return
}

type IndexResponse struct {
}

func (this_ *NodeApi) index(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	response := &IndexResponse{}

	res = response
	return
}
