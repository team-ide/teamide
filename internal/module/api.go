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
	for _, one := range api.apiCache {
		apiPowerMap[one.Power.Action] = true
	}
	ps := base.GetPowers()
	for _, one := range ps {
		if !ServerContext.IsServer {
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
	if !ServerContext.IsServer {
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
	*context.ServerContext
	applicationService *module_application.ApplicationService
	toolboxService     *module_toolbox.ToolboxService
	userService        *module_user.UserService
	registerService    *module_register.RegisterService
	loginService       *module_login.LoginService
	installService     *InstallService
	apiCache           map[string]*base.ApiWorker
}

var (

	//PowerRegister 基础权限
	PowerRegister    = base.AppendPower(&base.PowerAction{Action: "register", Text: "注册", StandAlone: false})
	PowerData        = base.AppendPower(&base.PowerAction{Action: "data", Text: "数据", StandAlone: true})
	PowerSession     = base.AppendPower(&base.PowerAction{Action: "session", Text: "会话", StandAlone: true})
	PowerLogin       = base.AppendPower(&base.PowerAction{Action: "login", Text: "登录", StandAlone: false})
	PowerLogout      = base.AppendPower(&base.PowerAction{Action: "logout", Text: "登出", StandAlone: false})
	PowerAutoLogin   = base.AppendPower(&base.PowerAction{Action: "auto_login", Text: "自动登录", StandAlone: false})
	PowerUpload      = base.AppendPower(&base.PowerAction{Action: "upload", Text: "上传", StandAlone: true})
	PowerUpdateCheck = base.AppendPower(&base.PowerAction{Action: "updateCheck", Text: "更新检测", StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker, err error) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"data"}, Power: PowerData, Do: this_.apiData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"login"}, Power: PowerLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"autoLogin"}, Power: PowerAutoLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"logout"}, Power: PowerLogout, Do: this_.apiLogout})
	apis = append(apis, &base.ApiWorker{Apis: []string{"register"}, Power: PowerRegister, Do: this_.apiRegister})
	apis = append(apis, &base.ApiWorker{Apis: []string{"session"}, Power: PowerSession, Do: this_.apiSession})
	apis = append(apis, &base.ApiWorker{Apis: []string{"upload"}, Power: PowerUpload, Do: this_.apiUpload})
	apis = append(apis, &base.ApiWorker{Apis: []string{"updateCheck"}, Power: PowerUpdateCheck, Do: this_.apiUpdateCheck})

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

		if !this_.IsServer {
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
	name := path[index+len("api/"):]

	api := this_.apiCache[name]
	if api == nil {
		return false
	}
	if api.IsGet && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	if api.IsWebSocket && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	requestBean := this_.getRequestBean(c)
	requestBean.Path = path
	if !this_.checkPower(api, requestBean.JWT, c) {
		return true
	}
	if api.Do != nil {
		this_.Logger.Info("处理操作", zap.String("path", path))
		res, err := api.Do(requestBean, c)
		if err != nil {
			this_.Logger.Error("操作异常", zap.String("path", path), zap.Error(err))
		}
		if res == base.HttpNotResponse {
			return true
		}
		base.ResponseJSON(res, err, c)
	}
	return true
}
