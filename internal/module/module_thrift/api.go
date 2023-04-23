package module_thrift

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/thrift"
	"github.com/team-ide/go-tool/util"
	"strings"
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
	RelativePath  string   `json:"relativePath"`
	ServiceName   string   `json:"serviceName"`
	MethodName    string   `json:"methodName"`
	Args          []string `json:"args"`
	ServerAddress string   `json:"serverAddress"`
	Reload        bool     `json:"reload"`
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

	filename = util.FormatPath(filename)

	methodNode := service.GetServiceMethod(filename, request.ServiceName, request.MethodName)
	if methodNode == nil {
		err = errors.New("service method node [" + filename + "][" + request.ServiceName + "][" + request.MethodName + "] not found")
		return
	}
	var structCache = map[string]*thrift.Struct{}

	argFields := service.GetFields(filename, methodNode.Params, structCache)
	var argDemoDataList []interface{}
	for _, argField := range argFields {
		bs, _ := json.MarshalIndent(service.GetFieldDemoData(filename, argField), "", "  ")
		argDemoDataList = append(argDemoDataList, string(bs))
	}
	data["argFields"] = argFields
	data["argDemoDataList"] = argDemoDataList
	data["structCache"] = structCache
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

	var argsJSON = "["
	for i, arg := range request.Args {
		if i > 0 {
			argsJSON += ","
		}
		trimS := strings.TrimSpace(arg)
		if strings.HasPrefix(trimS, "[") || strings.HasPrefix(trimS, "{") {
			argsJSON += arg
		} else {
			argsJSON += `"` + arg + `"`
		}
	}
	argsJSON += "]"

	var args []interface{}
	err = util.JSONDecodeUseNumber([]byte(argsJSON), &args)
	if err != nil {
		err = errors.New("args json " + argsJSON + " to args error:" + err.Error())
		return
	}

	param, err := service.InvokeByServerAddress(request.ServerAddress, filename, request.ServiceName, request.MethodName, args...)
	if param != nil {
		data["writeStart"] = param.WriteStart.UnixMilli()
		data["writeEnd"] = param.WriteEnd.UnixMilli()
		data["readStart"] = param.ReadStart.UnixMilli()
		data["readEnd"] = param.ReadEnd.UnixMilli()
		var result = map[string]interface{}{}
		result["args"] = param.Args
		result["result"] = param.Result
		result["exceptions"] = param.Exceptions
		bs, e := json.MarshalIndent(result, "", "  ")
		if e == nil {
			data["result"] = string(bs)
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
