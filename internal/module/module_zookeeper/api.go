package module_zookeeper

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power     = base.AppendPower(&base.PowerAction{Action: "zookeeper", Text: "Zookeeper", ShouldLogin: true, StandAlone: true})
	PowerInfo = base.AppendPower(&base.PowerAction{Action: "zookeeper_info", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/info"}, Power: PowerInfo, Do: this_.info})

	return
}

func (this_ *api) info(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}
