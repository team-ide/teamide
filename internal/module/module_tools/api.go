package module_tools

import (
	"teamide/internal/context"
	"teamide/pkg/base"
)

type Api struct {
	*context.ServerContext
}

func NewApi(ServerContext *context.ServerContext) *Api {
	return &Api{
		ServerContext: ServerContext,
	}
}

var (

	// Power 基本 权限
	Power = base.AppendPower(&base.PowerAction{Action: "tools", Text: "小工具", ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {

	return
}
