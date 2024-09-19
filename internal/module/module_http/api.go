package module_http

import (
	"github.com/gin-gonic/gin"
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
	Power   = base.AppendPower(&base.PowerAction{Action: "http", Text: "HTTP", ShouldLogin: true, StandAlone: true})
	execute = base.AppendPower(&base.PowerAction{Action: "execute", Text: "执行", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: execute, Do: this_.execute})

	return
}

type BaseRequest struct {
	ToolboxId   int64             `json:"toolboxId,omitempty"`
	Header      map[string]string `json:"header"`
	Path        string            `json:"path"`
	UserSetting bool              `json:"userSetting"`
	Toolbox     bool              `json:"toolbox"`
	ExistsDo    int               `json:"existsDo"`
}

func (this_ *api) execute(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.Execute(request)
	return
}
