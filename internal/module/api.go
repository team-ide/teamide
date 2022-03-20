package module

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_user"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	"teamide/pkg/db"
)

func CacheApi(dbWorker db.DatabaseWorker) (api *Api) {

	api = NewApi(dbWorker)
	api.bindApi()

	var apiPowerMap = make(map[string]bool)
	for _, api := range api.apiCache {
		apiPowerMap[api.Power.Action] = true
	}
	ps := base2.GetPowers()
	for _, one := range ps {
		if base2.IsStandAlone {
			if !one.AllowNative {
				continue
			}
		}
		_, ok := apiPowerMap[one.Action]
		if !ok {
			component.Logger.Warn(component.LogStr("权限[", one.Action, "]未配置动作"))
		}
	}
	return
}

// NewApi 根据库配置创建IDService
func NewApi(dbWorker db.DatabaseWorker) (res *Api) {
	res = &Api{
		dbWorker:        dbWorker,
		userService:     module_user.NewUserService(dbWorker),
		registerService: module_register.NewRegisterService(dbWorker),
		loginService:    module_login.NewLoginService(dbWorker),
		apiCache:        make(map[string]*base2.ApiWorker),
	}
	return
}

// Api ID服务
type Api struct {
	dbWorker        db.DatabaseWorker
	userService     *module_user.UserService
	registerService *module_register.RegisterService
	loginService    *module_login.LoginService
	apiCache        map[string]*base2.ApiWorker
}

func (this_ *Api) bindApi() {
	this_.appendApi(&base2.ApiWorker{Apis: []string{"", "/", "data"}, Power: base2.PowerData, Do: apiData})
	this_.appendApi(&base2.ApiWorker{Apis: []string{"login"}, Power: base2.PowerLogin, Do: this_.apiLogin})
	this_.appendApi(&base2.ApiWorker{Apis: []string{"autoLogin"}, Power: base2.PowerAutoLogin, Do: this_.apiLogin})
	this_.appendApi(&base2.ApiWorker{Apis: []string{"logout"}, Power: base2.PowerLogout, Do: this_.apiLogout})
	this_.appendApi(&base2.ApiWorker{Apis: []string{"register"}, Power: base2.PowerRegister, Do: this_.apiRegister})
	this_.appendApi(&base2.ApiWorker{Apis: []string{"session"}, Power: base2.PowerSession, Do: this_.apiSession})

}

func (this_ *Api) appendApi(apis ...*base2.ApiWorker) {
	if len(apis) == 0 {
		return
	}
	for _, api := range apis {
		if api.Power == nil {
			panic(errors.New(fmt.Sprint("API未设置权限!", api)))
		}
		if len(api.Apis) == 0 {
			panic(errors.New(fmt.Sprint("API未设置映射路径!", api)))
		}

		if base2.IsStandAlone {
			if !api.Power.AllowNative {
				continue
			}
		}
		for _, apiName := range api.Apis {

			_, find := this_.apiCache[apiName]
			if find {
				panic(errors.New(fmt.Sprint("API映射路径[", apiName, "]已存在!", api)))
			}
			// println("add api path :" + apiName + ",action:" + api.Power.Action)
			this_.apiCache[apiName] = api
		}
	}
}

func (this_ *Api) getRequestBean(c *gin.Context) (request *base2.RequestBean) {
	request = &base2.RequestBean{}
	JWT := this_.getJWT(c)
	request.JWT = JWT
	return
}

func (this_ *Api) DoApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	requestBean := this_.getRequestBean(c)
	requestBean.Path = path
	name := path[index+len("api/"):]

	api := this_.apiCache[name]
	if api == nil {
		return false
	}
	if !this_.checkPower(api, requestBean.JWT, c) {
		return true
	}
	if api.Do != nil {
		res, err := api.Do(requestBean, c)
		base2.ResponseJSON(res, err, c)
	}
	if api.DoOther != nil {
		api.DoOther(requestBean, c)
	}
	return true
}
