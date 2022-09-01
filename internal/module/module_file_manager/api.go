package module_file_manager

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/internal/module/module_toolbox"
)

type Api struct {
	*context.ServerContext
	ToolboxService *module_toolbox.ToolboxService
}

func NewApi(ToolboxService *module_toolbox.ToolboxService) *Api {
	return &Api{
		ServerContext:  ToolboxService.ServerContext,
		ToolboxService: ToolboxService,
	}
}

var (
	// 文件管理器 权限

	// Power 文件管理器 基本 权限
	Power = base.AppendPower(&base.PowerAction{Action: "file_manager", Text: "工具", ShouldLogin: false, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager"}, Power: Power, Do: this_.index})
	return
}

func (this_ *Api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}
