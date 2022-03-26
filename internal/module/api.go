package module

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/base"
	"teamide/internal/config"
	"teamide/internal/context"
	"teamide/internal/module/module_application"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/module/module_user"
)

func NewApi(ServerContext *context.ServerContext) (api *Api, err error) {

	api = &Api{
		ServerContext:      ServerContext,
		userService:        module_user.NewUserService(ServerContext),
		registerService:    module_register.NewRegisterService(ServerContext),
		loginService:       module_login.NewLoginService(ServerContext),
		installService:     NewInstallService(ServerContext),
		applicationService: module_application.NewApplicationService(ServerContext),
		toolboxService:     module_toolbox.NewToolboxService(ServerContext),

		apiCache: make(map[string]*base.ApiWorker),
	}
	var apis []*base.ApiWorker
	apis, err = api.GetApis()
	if err != nil {
		return
	}
	for _, one := range apis {
		err = api.appendApi(one)
		if err != nil {
			return
		}
	}

	var apiPowerMap = make(map[string]bool)
	for _, api := range api.apiCache {
		apiPowerMap[api.Power.Action] = true
	}
	ps := base.GetPowers()
	for _, one := range ps {
		if base.IsStandAlone {
			if !one.StandAlone {
				continue
			}
		}
		_, ok := apiPowerMap[one.Action]
		if !ok {
			ServerContext.Logger.Warn("权限[" + one.Action + "]未配置动作")
		}
	}

	err = api.installService.Install()
	if err != nil {
		return
	}
	// 如果是单机 初始化一些用户
	if base.IsStandAlone {
		err = api.InitStandAlone()
		if err != nil {
			return
		}
	}
	return
}

// Api ID服务
type Api struct {
	config.ServerConfig
	zap.Logger
	*context.ServerContext
	applicationService *module_application.ApplicationService
	toolboxService     *module_toolbox.ToolboxService
	userService        *module_user.UserService
	registerService    *module_register.RegisterService
	loginService       *module_login.LoginService
	installService     *InstallService
	apiCache           map[string]*base.ApiWorker
}

func (this_ *Api) GetApis() (apis []*base.ApiWorker, err error) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"", "/", "data"}, Power: base.PowerData, Do: apiData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"login"}, Power: base.PowerLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"autoLogin"}, Power: base.PowerAutoLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"logout"}, Power: base.PowerLogout, Do: this_.apiLogout})
	apis = append(apis, &base.ApiWorker{Apis: []string{"register"}, Power: base.PowerRegister, Do: this_.apiRegister})
	apis = append(apis, &base.ApiWorker{Apis: []string{"session"}, Power: base.PowerSession, Do: this_.apiSession})

	apis = append(apis, module_application.NewApplicationApi(this_.applicationService).GetApis()...)
	apis = append(apis, module_toolbox.NewToolboxApi(this_.toolboxService).GetApis()...)

	return
}

func (this_ *Api) appendApi(apis ...*base.ApiWorker) (err error) {
	if len(apis) == 0 {
		return
	}
	for _, api := range apis {
		if api.Power == nil {
			err = errors.New(fmt.Sprint("API未设置权限!", api))
			return
		}
		if len(api.Apis) == 0 {
			err = errors.New(fmt.Sprint("API未设置映射路径!", api))
			return
		}

		if base.IsStandAlone {
			if !api.Power.StandAlone {
				continue
			}
		}
		for _, apiName := range api.Apis {

			_, find := this_.apiCache[apiName]
			if find {
				err = errors.New(fmt.Sprint("API映射路径[", apiName, "]已存在!", api))
				return
			}
			// println("add api path :" + apiName + ",action:" + api.Power.Action)
			this_.apiCache[apiName] = api
		}
	}
	return
}

func (this_ *Api) getRequestBean(c *gin.Context) (request *base.RequestBean) {
	request = &base.RequestBean{}
	request.JWT = this_.getJWT(c)
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
		base.ResponseJSON(res, err, c)
	}
	if api.DoOther != nil {
		api.DoOther(requestBean, c)
	}
	return true
}
