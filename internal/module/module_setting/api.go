package module_setting

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/context"
	"teamide/pkg/base"
)

type Api struct {
	*context.ServerContext
	SettingService *SettingService
}

func NewApi(SettingService *SettingService) *Api {
	return &Api{
		ServerContext:  SettingService.ServerContext,
		SettingService: SettingService,
	}
}

var (

	// Power 基本 权限
	Power     = base.AppendPower(&base.PowerAction{Action: "setting", Text: "设置", ShouldLogin: true, StandAlone: true, ShouldPower: true})
	getPower  = base.AppendPower(&base.PowerAction{Action: "get", Text: "保存", Parent: Power, ShouldLogin: true, StandAlone: true, ShouldPower: true})
	savePower = base.AppendPower(&base.PowerAction{Action: "save", Text: "保存", Parent: Power, ShouldLogin: true, StandAlone: true, ShouldPower: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: savePower, Do: this_.save})

	return
}

type GetResponse struct {
	*context.Setting
}

func (this_ *Api) get(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {

	response := &GetResponse{}

	response.Setting = this_.Setting
	res = response
	return
}

func (this_ *Api) save(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := map[string]interface{}{}
	if !base.RequestJSON(&request, c) {
		return
	}

	err = this_.SettingService.Save(request)
	if err != nil {
		return
	}

	return
}
