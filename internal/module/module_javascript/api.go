package module_javascript

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/javascript/context_map"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
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
	Power      = base.AppendPower(&base.PowerAction{Action: "javascript", Text: "Javascript", ShouldLogin: true, StandAlone: true})
	getModules = base.AppendPower(&base.PowerAction{Action: "getModules", Text: "getModules", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: getModules, Do: this_.getModules})

	return
}

type BaseRequest struct {
}

func (this_ *api) getModules(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res = context_map.ModuleList
	return
}
