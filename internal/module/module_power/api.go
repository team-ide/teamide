package module_power

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
)

type Api struct {
	*context.ServerContext
	PowerRoleService *PowerRoleService
}

func NewApi(PowerRoleService *PowerRoleService) *Api {
	return &Api{
		ServerContext:    PowerRoleService.ServerContext,
		PowerRoleService: PowerRoleService,
	}
}

var (
	// 用户 权限

	// Power 用户基本 权限
	Power     = base.AppendPower(&base.PowerAction{Action: "power", Text: "权限", ShouldLogin: true, StandAlone: true})
	dataPower = base.AppendPower(&base.PowerAction{Action: "data", Text: "权限基本数据", Parent: Power, ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: dataPower, Do: this_.data})

	return
}

type DataResponse struct {
	Powers []*base.PowerAction `json:"powers"`
}

func (this_ *Api) data(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	response := &DataResponse{}

	response.Powers = base.GetPowers()

	res = response
	return
}
