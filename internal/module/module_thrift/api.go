package module_thrift

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/thrift"
	"github.com/team-ide/go-tool/util"
	"io/fs"
	"os"
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
	contextPower          = base.AppendPower(&base.PowerAction{Action: "context", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	getMethodArgFields    = base.AppendPower(&base.PowerAction{Action: "getMethodArgFields", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeByServerAddress = base.AppendPower(&base.PowerAction{Action: "invokeByServerAddress", Text: "上下文", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeReports         = base.AppendPower(&base.PowerAction{Action: "invokeReports", Text: "执行报告", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeStop            = base.AppendPower(&base.PowerAction{Action: "invokeStop", Text: "执行停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeInfo            = base.AppendPower(&base.PowerAction{Action: "invokeInfo", Text: "执行信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower            = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: contextPower, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Power: getMethodArgFields, Do: this_.getMethodArgFields})
	apis = append(apis, &base.ApiWorker{Power: invokeByServerAddress, Do: this_.invokeByServerAddress})
	apis = append(apis, &base.ApiWorker{Power: invokeReports, Do: this_.invokeReports})
	apis = append(apis, &base.ApiWorker{Power: invokeStop, Do: this_.invokeStop})
	apis = append(apis, &base.ApiWorker{Power: invokeInfo, Do: this_.invokeInfo})
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
	TaskKey       string   `json:"taskKey"`
	ToolboxId     int64    `json:"toolboxId,omitempty"`
	IsTest        bool     `json:"isTest"`
	Worker        int      `json:"worker"`
	Duration      int      `json:"duration"`
	Frequency     int      `json:"frequency"`
	Timeout       int      `json:"timeout"`
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
	data["isTest"] = request.IsTest
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
	if !request.IsTest {
		var param *thrift.MethodParam
		param, err = service.InvokeByServerAddress(request.ServerAddress, filename, request.ServiceName, request.MethodName, args...)
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
	} else {
		var dir string
		dir, err = this_.invokeTaskDir(request)
		if err != nil {
			return
		}
		var t *task.Task

		executor := &invokeExecutor{
			BaseRequest:  request,
			filename:     filename,
			args:         args,
			workerClient: make(map[int]*thrift.ServiceClient),
			service:      service,
		}
		t, err = task.New(&task.Options{
			Key:       util.GetUUID(),
			Worker:    request.Worker,
			Frequency: request.Frequency,
			Duration:  request.Duration,
			Executor:  executor,
		})
		if err != nil {
			return
		}
		taskDir := dir + "" + t.Key
		_ = this_.saveTaskInfo(taskDir, t)
		go func() {
			defer func() {
				_ = this_.saveTaskInfo(taskDir, t)
				executor.stop()
			}()
			for !t.IsEnd {
				_ = this_.saveTaskInfo(taskDir, t)
				time.Sleep(time.Second * 1)
			}
		}()
		go t.Run()
		addTask(t)
	}
	return
}

type invokeExecutor struct {
	*BaseRequest
	filename         string
	args             []interface{}
	workerClient     map[int]*thrift.ServiceClient
	workerClientLock sync.Mutex
	service          *thrift.Workspace
}

func (this_ *invokeExecutor) stop() {
	this_.workerClientLock.Lock()
	defer this_.workerClientLock.Unlock()

	for _, client := range this_.workerClient {
		client.Stop()
	}

}
func (this_ *invokeExecutor) getClient(param *task.ExecutorParam) (client *thrift.ServiceClient, err error) {
	this_.workerClientLock.Lock()
	defer this_.workerClientLock.Unlock()

	client = this_.workerClient[param.WorkerIndex]
	if client != nil {
		return
	}

	client, err = thrift.NewServiceClientByAddress(this_.ServerAddress)
	if err != nil {
		return
	}
	this_.workerClient[param.WorkerIndex] = client
	return
}

func (this_ *invokeExecutor) Before(param *task.ExecutorParam) (err error) {
	_, err = this_.getClient(param)
	if err != nil {
		return
	}

	methodParam, err := this_.service.GetMethodParam(this_.filename, this_.ServiceName, this_.MethodName, this_.args...)
	if err != nil {
		return
	}
	param.Extend = methodParam

	//util.Logger.Info("test Before", zap.Any("param", param))
	return
}

func (this_ *invokeExecutor) Execute(param *task.ExecutorParam) (err error) {
	client, err := this_.getClient(param)
	if err != nil {
		return
	}
	//util.Logger.Info("test Execute", zap.Any("param", param))
	methodParam := param.Extend.(*thrift.MethodParam)

	_, err = client.Send(context.Background(), methodParam)
	return
}

func (this_ *invokeExecutor) After(param *task.ExecutorParam) (err error) {
	//util.Logger.Info("test After", zap.Any("param", param))
	return
}

func (this_ *api) invokeTaskDir(request *BaseRequest) (taskDir string, err error) {
	taskDir = this_.toolboxService.GetFilesDir()
	taskDir += fmt.Sprintf("%s/toolbox-%d/%s", "thrift-tasks", request.ToolboxId, request.RelativePath) + "/" + request.ServiceName + "/" + request.MethodName + "/"

	ex, err := util.PathExists(taskDir)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(taskDir, fs.ModePerm)
	}

	return
}

func (this_ *api) saveTaskInfo(taskDir string, task *task.Task) (err error) {
	ex, err := util.PathExists(taskDir)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(taskDir, fs.ModePerm)
	}
	bs, _ := json.MarshalIndent(task, "", "  ")
	err = util.WriteFile(taskDir+"/task.json", bs)
	if err != nil {
		return
	}
	bs, _ = json.MarshalIndent(task.Metric.Count(), "", "  ")
	err = util.WriteFile(taskDir+"/metric.json", bs)
	if err != nil {
		return
	}
	return
}
func (this_ *api) loadTask(taskDir string) (data map[string]interface{}, err error) {
	defer func() {
		if len(data) == 0 {
			data = nil
		}
	}()
	data = map[string]interface{}{}

	var bs []byte
	if ex, _ := util.PathExists(taskDir + "/task.json"); ex {
		if bs, err = os.ReadFile(taskDir + "/task.json"); err != nil {
			return
		}
		t := map[string]interface{}{}
		err = util.JSONDecodeUseNumber(bs, &t)
		if err != nil {
			return
		}
		data["task"] = t
	}
	if ex, _ := util.PathExists(taskDir + "/metric.json"); ex {
		if bs, err = os.ReadFile(taskDir + "/metric.json"); err != nil {
			return
		}
		metric := map[string]interface{}{}
		err = util.JSONDecodeUseNumber(bs, &metric)
		if err != nil {
			return
		}
		data["metric"] = metric
	}

	return
}

func (this_ *api) invokeReports(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	dir, err := this_.invokeTaskDir(request)
	if err != nil {
		return
	}

	fileList, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var taskInfo map[string]interface{}
	var taskList []map[string]interface{}
	for _, f := range fileList {
		if !f.IsDir() {
			continue
		}
		taskKey := f.Name()
		taskInfo, err = this_.loadTask(dir + f.Name())
		if err != nil {
			return
		}
		if taskInfo != nil {
			taskInfo["taskKey"] = taskKey
			taskList = append(taskList, taskInfo)
		}

	}
	res = taskList
	return
}

var taskCache = map[string]*task.Task{}
var taskLocker = &sync.Mutex{}

func getTask(taskKey string) *task.Task {
	taskLocker.Lock()
	defer taskLocker.Unlock()

	return taskCache[taskKey]
}

func addTask(task *task.Task) {
	taskLocker.Lock()
	defer taskLocker.Unlock()

	taskCache[task.Key] = task
}

func removeTask(taskKey string) {
	taskLocker.Lock()
	defer taskLocker.Unlock()

	delete(taskCache, taskKey)
}

func (this_ *api) invokeStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	t := getTask(request.TaskKey)

	if t != nil {
		t.Stop()
	}
	return
}

func (this_ *api) invokeInfo(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	dir, err := this_.invokeTaskDir(request)
	if err != nil {
		return
	}

	res, err = this_.loadTask(dir + request.TaskKey)

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
