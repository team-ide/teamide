package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_database"
	"teamide/internal/module/module_datamove"
	"teamide/internal/module/module_elasticsearch"
	"teamide/internal/module/module_file_manager"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_javascript"
	"teamide/internal/module/module_kafka"
	"teamide/internal/module/module_log"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_power"
	"teamide/internal/module/module_redis"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_setting"
	"teamide/internal/module/module_terminal"
	"teamide/internal/module/module_thrift"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/module/module_tools"
	"teamide/internal/module/module_user"
	"teamide/internal/module/module_zookeeper"
	"teamide/pkg/base"
	"time"
)

func NewApi(ServerContext *context.ServerContext) (api *Api, err error) {

	api = &Api{
		ServerContext:          ServerContext,
		userService:            module_user.NewUserService(ServerContext),
		userSettingService:     module_user.NewUserSettingService(ServerContext),
		registerService:        module_register.NewRegisterService(ServerContext),
		loginService:           module_login.NewLoginService(ServerContext),
		installService:         NewInstallService(ServerContext),
		toolboxService:         module_toolbox.NewToolboxService(ServerContext),
		nodeService:            module_node.NewNodeService(ServerContext),
		powerRoleService:       module_power.NewPowerRoleService(ServerContext),
		powerRouteService:      module_power.NewPowerRouteService(ServerContext),
		powerUserService:       module_power.NewPowerUserService(ServerContext),
		settingService:         module_setting.NewSettingService(ServerContext),
		terminalCommandService: module_terminal.NewTerminalCommandService(ServerContext),
		idService:              module_id.NewIDService(ServerContext),
		logService:             module_log.NewLogService(ServerContext),
		apiCache:               make(map[string]*base.ApiWorker),
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

	err = api.installService.Install()
	if err != nil {
		return
	}

	err = api.InitSetting()
	if err != nil {
		return
	}

	if ServerContext.IsServer {
		err = api.initServer()
		if err != nil {
			return
		}
	} else {
		err = api.initStandAlone()
		if err != nil {
			return
		}
	}
	go api.nodeService.InitContext()

	err = api.logService.ServerReady()
	if err != nil {
		return
	}
	err = api.terminalCommandService.ServerReady()
	if err != nil {
		return
	}

	return
}

// InitSetting 查询
func (this_ *Api) InitSetting() (err error) {

	list, err := this_.settingService.Query()
	if err != nil {
		return
	}
	for _, one := range list {
		_, e := this_.Setting.Set(one.Name, one.Value)
		if e != nil {
			this_.Logger.Error("init setting error", zap.Any("setting", one), zap.Error(e))
		}
	}

	return
}

// Api ID服务
type Api struct {
	*context.ServerContext
	toolboxService         *module_toolbox.ToolboxService
	nodeService            *module_node.NodeService
	terminalCommandService *module_terminal.TerminalCommandService
	userService            *module_user.UserService
	userSettingService     *module_user.UserSettingService
	registerService        *module_register.RegisterService
	loginService           *module_login.LoginService
	powerRoleService       *module_power.PowerRoleService
	powerRouteService      *module_power.PowerRouteService
	powerUserService       *module_power.PowerUserService
	logService             *module_log.LogService
	settingService         *module_setting.SettingService
	idService              *module_id.IDService
	installService         *InstallService
	apiCache               map[string]*base.ApiWorker
}

var (

	//PowerRegister 基础权限
	PowerRegister    = base.AppendPower(&base.PowerAction{Action: "register", Text: "注册", StandAlone: false})
	PowerData        = base.AppendPower(&base.PowerAction{Action: "data", Text: "数据", StandAlone: true})
	showPlaintext    = base.AppendPower(&base.PowerAction{Action: "showPlaintext", Text: "数据", StandAlone: true})
	PowerSession     = base.AppendPower(&base.PowerAction{Action: "session", Text: "会话", StandAlone: true})
	PowerLogin       = base.AppendPower(&base.PowerAction{Action: "login", Text: "登录", StandAlone: false})
	PowerLogout      = base.AppendPower(&base.PowerAction{Action: "logout", Text: "登出", StandAlone: false})
	PowerAutoLogin   = base.AppendPower(&base.PowerAction{Action: "autoLogin", Text: "自动登录", StandAlone: false})
	PowerUpload      = base.AppendPower(&base.PowerAction{Action: "upload", Text: "上传", StandAlone: true})
	PowerUpdateCheck = base.AppendPower(&base.PowerAction{Action: "updateCheck", Text: "更新检测", ShouldPower: true, ShouldLogin: true, StandAlone: true})
	listenPower      = base.AppendPower(&base.PowerAction{Action: "listen", Text: "监听事件", StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker, err error) {
	apis = append(apis, &base.ApiWorker{Power: PowerData, Do: this_.apiData, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: showPlaintext, Do: this_.apiShowPlaintext})
	apis = append(apis, &base.ApiWorker{Power: PowerLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Power: PowerAutoLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Power: PowerLogout, Do: this_.apiLogout})
	apis = append(apis, &base.ApiWorker{Power: PowerRegister, Do: this_.apiRegister})
	apis = append(apis, &base.ApiWorker{Power: PowerSession, Do: this_.apiSession, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: PowerUpload, Do: this_.apiUpload, IsUpload: true})
	apis = append(apis, &base.ApiWorker{Power: PowerUpdateCheck, Do: this_.apiUpdateCheck, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: listenPower, Do: this_.listen, NotRecodeLog: true})

	apis = append(apis, module_toolbox.NewToolboxApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_node.NewNodeApi(this_.nodeService).GetApis()...)
	apis = append(apis, module_file_manager.NewApi(this_.toolboxService, this_.nodeService).GetApis()...)
	apis = append(apis, module_terminal.NewApi(this_.toolboxService, this_.nodeService, this_.terminalCommandService).GetApis()...)
	apis = append(apis, module_user.NewApi(this_.userService).GetApis()...)
	apis = append(apis, module_redis.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_database.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_datamove.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_zookeeper.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_kafka.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_elasticsearch.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_log.NewApi(this_.logService).GetApis()...)
	apis = append(apis, module_power.NewApi(this_.powerRoleService).GetApis()...)
	apis = append(apis, module_tools.NewApi(this_.ServerContext).GetApis()...)
	apis = append(apis, module_setting.NewApi(this_.settingService).GetApis()...)
	apis = append(apis, module_thrift.NewApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_javascript.NewApi(this_.toolboxService).GetApis()...)

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
		if api.Power.Action == "" {
			err = errors.New(fmt.Sprint("API未设置映射路径!", api))
			return
		}

		//fmt.Println("action:", api.Power.Action)
		_, find := this_.apiCache[api.Power.Action]
		if find {
			err = errors.New(fmt.Sprint("API映射路径[", api.Power.Action, "]已存在!", api))
			return
		}
		this_.apiCache[api.Power.Action] = api
	}
	return
}

func (this_ *Api) getRequestBean(c *gin.Context) (request *base.RequestBean) {
	request = &base.RequestBean{}
	request.JWT = this_.getJWT(c)
	request.ClientKey = c.GetHeader("key1")
	request.ClientTabKey = c.GetHeader("key2")
	if strings.EqualFold(c.Request.Method, "get") {
		request.ClientKey = c.Query("key1")
		request.ClientTabKey = c.Query("key2")
	}
	return
}

func (this_ *Api) DoApi(path string, c *gin.Context) bool {
	//fmt.Println("do api start path:", path)
	//defer func() {
	//if e := recover(); e != nil {
	//	fmt.Println(e)
	//}
	//fmt.Println("do api end path:", path)
	//}()

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	action := path[index+len("api/"):]

	api := this_.apiCache[action]
	//fmt.Println("do api start action:", action)
	if api == nil {
		return false
	}
	//fmt.Println("do api start api:", api)
	if api.IsGet && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	if api.IsWebSocket && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	requestBean := this_.getRequestBean(c)
	requestBean.Path = path
	requestBean.Power = api.Power
	if !this_.checkPower(api, requestBean.JWT, c) {
		this_.Logger.Warn("no power", zap.Any("action", action))
		return true
	}
	if api.Do != nil {
		var err error
		var startTime = util.GetNow()
		var logRecode *module_log.LogModel = nil
		if !api.NotRecodeLog {
			userAgentStr := c.Request.UserAgent()
			logRecode = &module_log.LogModel{
				Action:     action,
				Method:     c.Request.Method,
				StartTime:  startTime,
				CreateTime: startTime,
				Ip:         c.ClientIP(),
				UserAgent:  userAgentStr,
			}
			var param = make(map[string]interface{})
			_ = c.Request.ParseForm()
			f := c.Request.Form
			for k, v := range f {
				param[k] = v
			}
			f = c.Request.PostForm
			for k, v := range f {
				param[k] = v
			}
			if len(param) > 0 {
				bs, _ := json.Marshal(param)
				logRecode.Param = string(bs)
			}
			if !api.IsUpload {
				var data = make(map[string]interface{})
				_ = c.ShouldBindBodyWith(&data, binding.JSON)
				if len(data) > 0 {
					bs, _ := json.Marshal(data)
					logRecode.Data = string(bs)
				}
			}
			if requestBean.JWT != nil {
				logRecode.UserId = requestBean.JWT.UserId
				logRecode.UserName = requestBean.JWT.Name
				logRecode.UserAccount = requestBean.JWT.Account
				logRecode.LoginId = requestBean.JWT.LoginId
			}
		}

		defer func() {
			if logRecode != nil {
				logRecode.EndTime = util.GetNow()
				_ = this_.logService.Insert(logRecode, err)
			}
		}()

		//if !api.NotRecodeLog {
		//	this_.Logger.Info("处理操作", zap.String("action", action))
		//}
		res, err := api.Do(requestBean, c)
		useTime := time.Now().UnixMilli() - startTime.UnixMilli()
		if err != nil {
			this_.Logger.Error("处理操作异常", zap.Any("action", action), zap.Any("useTime", useTime), zap.Any("error", err.Error()))
		} else {
			if !api.NotRecodeLog {
				this_.Logger.Debug("处理操作", zap.Any("action", action), zap.Any("useTime", useTime))
			}
		}
		if res == base.HttpNotResponse {
			return true
		}
		base.ResponseJSON(res, err, c)
	}
	return true
}
