package module_maker

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/maker/modelers"
)

type api struct {
	*module_toolbox.ToolboxService
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService) *api {
	return &api{
		ToolboxService: toolboxService_,
	}
}

var (
	// Terminal 权限

	// Power 文件管理器 基本 权限
	Power        = base.AppendPower(&base.PowerAction{Action: "maker", Text: "Maker", ShouldLogin: true, StandAlone: true})
	contextPower = base.AppendPower(&base.PowerAction{Action: "context", Text: "context", ShouldLogin: true, StandAlone: true, Parent: Power})
	get          = base.AppendPower(&base.PowerAction{Action: "get", Text: "get", ShouldLogin: true, StandAlone: true, Parent: Power})
	getList      = base.AppendPower(&base.PowerAction{Action: "getList", Text: "getList", ShouldLogin: true, StandAlone: true, Parent: Power})
	insert       = base.AppendPower(&base.PowerAction{Action: "insert", Text: "insert", ShouldLogin: true, StandAlone: true, Parent: Power})
	save         = base.AppendPower(&base.PowerAction{Action: "save", Text: "save", ShouldLogin: true, StandAlone: true, Parent: Power})
	del          = base.AppendPower(&base.PowerAction{Action: "delete", Text: "delete", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower   = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: contextPower, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Power: get, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: getList, Do: this_.getList})
	apis = append(apis, &base.ApiWorker{Power: insert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Power: save, Do: this_.save})
	apis = append(apis, &base.ApiWorker{Power: del, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *Config, err error) {
	config = &Config{}
	_, err = this_.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func (this_ *api) getService(requestBean *base.RequestBean, c *gin.Context) (res *Service, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	key := "dir:" + config.Dir
	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s *Service
		s, err = createService(config)
		if err != nil {
			util.Logger.Error("createService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}

		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowMilli(),
			Service:     s,
			Stop:        s.Close,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(*Service)
	serviceInfo.SetLastUseTime()
	return
}

type Request struct {
	Key       string      `json:"key"`
	ModelType string      `json:"modelType"`
	ModelName string      `json:"modelName"`
	Model     interface{} `json:"model"`
}

func (this_ *api) context(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}
	context := make(map[string]interface{})
	context["app"] = service.app
	context["types"] = modelers.GetTypes()
	context["docTemplateCache"] = modelers.GetDocTemplateCache()

	res = context
	return
}

func (this_ *api) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = service.app.GetModelTypeModel(request.ModelType, request.ModelName)
	return
}

func (this_ *api) getList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = service.app.GetModelTypeModels(request.ModelType)
	return
}

func (this_ *api) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.ModelType == "" || request.ModelName == "" {
		err = errors.New("参数丢失")
		return
	}

	err = service.app.Save(request.ModelType, request.ModelName, request.Model, true)
	return
}

func (this_ *api) save(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.ModelType == "" || request.ModelName == "" {
		err = errors.New("参数丢失")
		return
	}

	err = service.app.Save(request.ModelType, request.ModelName, request.Model, false)
	return
}

func (this_ *api) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	service, err := this_.getService(requestBean, c)
	if err != nil {
		return
	}

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.ModelType == "" || request.ModelName == "" {
		err = errors.New("参数丢失")
		return
	}

	err = service.app.Remove(request.ModelType, request.ModelName)
	return
}

func (this_ *api) close(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
