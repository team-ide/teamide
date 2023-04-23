package module_thrift

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/thrift"
	"sync"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"time"
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
	Power                 = base.AppendPower(&base.PowerAction{Action: "thrift", Text: "Thrift", ShouldLogin: true, StandAlone: true})
	context               = base.AppendPower(&base.PowerAction{Action: "context", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	getMethodArgFields    = base.AppendPower(&base.PowerAction{Action: "getMethodArgFields", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeByServerAddress = base.AppendPower(&base.PowerAction{Action: "invokeByServerAddress", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower            = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: context, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Power: getMethodArgFields, Do: this_.getMethodArgFields})
	apis = append(apis, &base.ApiWorker{Power: invokeByServerAddress, Do: this_.invokeByServerAddress})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

type Config struct {
	ThriftDir string `json:"thriftDir"`
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *Config, err error) {
	config = &Config{}
	err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

type BaseRequest struct {
	RelativePath  string        `json:"relativePath"`
	ServiceName   string        `json:"serviceName"`
	MethodName    string        `json:"methodName"`
	Args          []interface{} `json:"args"`
	ServerAddress string        `json:"serverAddress"`
	Reload        bool          `json:"reload"`
}

func (this_ *api) context(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getOrCreateWorkspace(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Reload {
		service.Reload()
	}

	data := map[string]interface{}{}
	res = data

	data["serviceList"] = service.ServiceList
	return
}

func (this_ *api) getMethodArgFields(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getOrCreateWorkspace(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	data := map[string]interface{}{}
	res = data

	filename := service.GetFormatDir() + "/" + request.RelativePath
	data["argFields"], data["structCache"], err = service.GetMethodArgFields(filename, request.ServiceName, request.MethodName)
	return
}

func (this_ *api) invokeByServerAddress(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getOrCreateWorkspace(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	filename := service.GetFormatDir() + "/" + request.RelativePath

	data := map[string]interface{}{}
	res = data
	data["start"] = time.Now().UnixMilli()
	defer func() {
		data["end"] = time.Now().UnixMilli()
	}()
	param, err := service.InvokeByServerAddress(request.ServerAddress, filename, request.ServiceName, request.MethodName, request.Args...)
	if param != nil {
		data["writeStart"] = param.WriteStart.UnixMilli()
		data["writeEnd"] = param.WriteEnd.UnixMilli()
		data["readStart"] = param.ReadStart.UnixMilli()
		data["readEnd"] = param.ReadEnd.UnixMilli()
		data["result"] = param.Result
		if param.Result != nil {
			bs, e := json.MarshalIndent(param.Result, "", "  ")
			if e == nil {
				data["result"] = string(bs)
			}
		}
	}
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}

	removeWorkspace(config.ThriftDir)
	return
}

var (
	workspaceCache     = map[string]*thrift.Workspace{}
	workspaceCacheLock = &sync.Mutex{}
)

func getOrCreateWorkspace(config *Config) (res *thrift.Workspace, err error) {
	workspaceCacheLock.Lock()
	defer workspaceCacheLock.Unlock()

	res = workspaceCache[config.ThriftDir]
	if res != nil {
		return
	}
	res = thrift.NewWorkspace(config.ThriftDir)
	res.Load()

	workspaceCache[config.ThriftDir] = res

	return
}

func removeWorkspace(dir string) {
	workspaceCacheLock.Lock()
	defer workspaceCacheLock.Unlock()

	find := workspaceCache[dir]
	if find != nil {
		find.Clean()
	}
	delete(workspaceCache, dir)
	return
}
