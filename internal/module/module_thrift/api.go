package module_thrift

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/metric"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/thrift"
	"github.com/team-ide/go-tool/util"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"sort"
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
	invokeReportDelete    = base.AppendPower(&base.PowerAction{Action: "invokeReportDelete", Text: "执行报告", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeStop            = base.AppendPower(&base.PowerAction{Action: "invokeStop", Text: "执行停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeInfo            = base.AppendPower(&base.PowerAction{Action: "invokeInfo", Text: "执行信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	downloadRecords       = base.AppendPower(&base.PowerAction{Action: "downloadRecords", Text: "执行信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	invokeMetric          = base.AppendPower(&base.PowerAction{Action: "invokeMetric", Text: "执行信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower            = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {

	apis = append(apis, &base.ApiWorker{Power: contextPower, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Power: getMethodArgFields, Do: this_.getMethodArgFields})
	apis = append(apis, &base.ApiWorker{Power: invokeByServerAddress, Do: this_.invokeByServerAddress})
	apis = append(apis, &base.ApiWorker{Power: invokeReports, Do: this_.invokeReports})
	apis = append(apis, &base.ApiWorker{Power: invokeReportDelete, Do: this_.invokeReportDelete})
	apis = append(apis, &base.ApiWorker{Power: downloadRecords, Do: this_.downloadRecords, IsGet: true})
	apis = append(apis, &base.ApiWorker{Power: invokeStop, Do: this_.invokeStop})
	apis = append(apis, &base.ApiWorker{Power: invokeInfo, Do: this_.invokeInfo})
	apis = append(apis, &base.ApiWorker{Power: invokeMetric, Do: this_.invokeMetric})
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
	Minute        bool     `json:"minute"`
	Second        bool     `json:"second"`
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

	var list = service.ServiceList
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name) //升序  即前面的值比后面的小  忽略大小写排序
	})
	for _, one := range list {
		sort.Slice(one.Methods, func(i, j int) bool {
			return strings.ToLower(one.Methods[i].Name) < strings.ToLower(one.Methods[j].Name) //升序  即前面的值比后面的小  忽略大小写排序
		})
	}

	data["serviceList"] = list
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
		str := string(bs)
		if str == `""` {
			argDemoDataList = append(argDemoDataList, "")
		} else {
			argDemoDataList = append(argDemoDataList, string(bs))
		}
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
			data["writeStart"] = param.WriteStart
			data["writeEnd"] = param.WriteEnd
			data["readStart"] = param.ReadStart
			data["readEnd"] = param.ReadEnd
			data["useTime"] = param.UseTime
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
		var parentDir string
		parentDir, err = this_.getTaskParentDir(request)
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
			Key:       fmt.Sprintf("%d", time.Now().UnixNano()),
			Worker:    request.Worker,
			Frequency: request.Frequency,
			Duration:  request.Duration,
			Executor:  executor,
		})
		if err != nil {
			return
		}
		executor.taskDir = parentDir + "" + t.Key
		_ = this_.saveTaskInfo(executor, request, t)
		go func() {
			defer func() {
				removeTask(t.Key)
				_ = this_.saveTaskInfo(executor, request, t)
				executor.stop()
			}()
			for !t.IsEnd {
				_ = this_.saveTaskInfo(executor, request, t)
				time.Sleep(time.Second * 1)
			}
		}()
		go t.Run()
		addTask(t)
	}
	return
}

func (this_ *api) getTaskParentDir(request *BaseRequest) (taskDir string, err error) {
	taskDir = this_.toolboxService.GetFilesDir()
	taskDir += this_.getTaskParentDirRelativePath(request)

	ex, err := util.PathExists(taskDir)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(taskDir, fs.ModePerm)
	}

	return
}

func (this_ *api) getTaskParentDirRelativePath(request *BaseRequest) (taskDir string) {
	taskDir = fmt.Sprintf("%s/toolbox-%d/%s", "thrift-tasks", request.ToolboxId, request.RelativePath) + "/" + request.ServiceName + "/" + request.MethodName + "/"

	return
}

func (this_ *api) saveTaskInfo(executor *invokeExecutor, request *BaseRequest, task *task.Task) (err error) {
	ex, err := util.PathExists(executor.taskDir)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(executor.taskDir, fs.ModePerm)
	}

	bs, _ := json.Marshal(request)
	requestMd5 := util.GetMD5(string(bs))
	data := map[string]interface{}{}
	data["requestMd5"] = requestMd5
	data["request"] = request
	data["task"] = task
	data["taskKey"] = task.Key
	c := task.Metric.Count()
	topItems := c.TopItems
	c.TopItems = []*metric.Item{}
	data["metric"] = c
	bs, _ = json.MarshalIndent(data, "", "  ")
	err = util.WriteFile(executor.taskDir+"/info.json", bs)
	if err != nil {
		return
	}
	c.TopItems = topItems
	bs, _ = json.MarshalIndent(c, "", "  ")
	_ = util.WriteFile(executor.taskDir+"/metric.json", bs)
	bs, _ = json.MarshalIndent(task.Metric.CountMinute(), "", "  ")
	_ = util.WriteFile(executor.taskDir+"/metric.minute.json", bs)
	bs, _ = json.MarshalIndent(task.Metric.CountSecond(), "", "  ")
	_ = util.WriteFile(executor.taskDir+"/metric.second.json", bs)

	paramList := executor.getAndCleanParamList()
	var recordsFile *os.File
	if ex, _ = util.PathExists(executor.taskDir + "/records.txt"); ex {
		recordsFile, _ = os.OpenFile(executor.taskDir+"/records.txt", os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		recordsFile, _ = os.Create(executor.taskDir + "/records.txt")
	}
	if recordsFile != nil {
		defer func() { _ = recordsFile.Close() }()

		for _, param := range paramList {
			param.ArgFields = nil
			param.ResultType = nil
			param.ExceptionFields = nil
			bs, _ = json.Marshal(param)
			_, _ = recordsFile.Write(bs)
			_, _ = recordsFile.WriteString("\n")
		}
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
	if ex, _ := util.PathExists(taskDir + "/info.json"); ex {
		if bs, err = os.ReadFile(taskDir + "/info.json"); err != nil {
			return
		}
		err = util.JSONDecodeUseNumber(bs, &data)
		if err != nil {
			return
		}
		if data["taskKey"] != nil {
			taskKey := util.GetStringValue(data["taskKey"])
			t := getTask(taskKey)
			if t != nil {
				data["isEnd"] = false
			} else {
				data["isEnd"] = true
			}
		}

	}

	return
}

func (this_ *api) invokeReports(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	parentDir, err := this_.getTaskParentDir(request)
	if err != nil {
		return
	}

	parentDirRelativePath := this_.getTaskParentDirRelativePath(request)

	fileList, err := os.ReadDir(parentDir)
	if err != nil {
		return
	}
	var taskInfo map[string]interface{}
	var taskList []map[string]interface{}
	var names []string
	for _, f := range fileList {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())

	}

	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j]) //升序  即前面的值比后面的小 忽略大小写排序
	})
	size := len(names)
	for i := size - 1; i >= 0; i-- {
		taskInfo, err = this_.loadTask(parentDir + names[i])
		if err != nil {
			return
		}
		if taskInfo != nil {
			taskInfo["taskRelativePath"] = parentDirRelativePath + names[i]
			taskList = append(taskList, taskInfo)
		}
	}
	res = taskList
	return
}

func (this_ *api) downloadRecords(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	this_.toolboxService.Logger.Info("下载执行记录 start")
	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	request := map[string]string{}

	err = c.Bind(&request)
	if err != nil {
		return
	}

	fileName := "" + request["serviceName"] + "." + request["methodName"] + "-执行记录.txt"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileName)))

	// 此处不设置 文件大小，如果设置文件大小，将无法终止下载
	//c.Header("Content-Length", fmt.Sprint(fileInfo.Size))
	c.Header("download-file-name", fileName)

	taskDir := this_.toolboxService.GetFilesDir() + request["taskRelativePath"]

	if ex, _ := util.PathExists(taskDir + "/records.txt"); ex {
		var f *os.File
		f, err = os.Open(taskDir + "/records.txt")
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(c.Writer, f)
	}
	c.Status(http.StatusOK)
	return
}
func (this_ *api) invokeReportDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	taskParentDir, err := this_.getTaskParentDir(request)
	if err != nil {
		return
	}

	t := getTask(request.TaskKey)

	if t != nil {
		t.Stop()
	}
	taskDir := taskParentDir + "" + request.TaskKey

	if ex, _ := util.PathExists(taskDir); ex {
		err = os.RemoveAll(taskDir)
	}
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

	taskParentDir, err := this_.getTaskParentDir(request)
	if err != nil {
		return
	}

	res, err = this_.loadTask(taskParentDir + request.TaskKey)

	return
}

func (this_ *api) invokeMetric(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	taskParentDir, err := this_.getTaskParentDir(request)
	if err != nil {
		return
	}

	taskDir := taskParentDir + request.TaskKey
	var data []*metric.Count
	var bs []byte
	if request.Minute {
		if ex, _ := util.PathExists(taskDir + "/metric.minute.json"); ex {
			if bs, err = os.ReadFile(taskDir + "/metric.minute.json"); err != nil {
				return
			}
			err = util.JSONDecodeUseNumber(bs, &data)
			if err != nil {
				return
			}
		}
	} else if request.Second {
		if ex, _ := util.PathExists(taskDir + "/metric.second.json"); ex {
			if bs, err = os.ReadFile(taskDir + "/metric.second.json"); err != nil {
				return
			}
			err = util.JSONDecodeUseNumber(bs, &data)
			if err != nil {
				return
			}
		}
	}
	res = data

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
