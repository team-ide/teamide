package module_database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"sync"
	"teamide/pkg/base"
	"time"
)

func (this_ *api) testStart(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	data := map[string]interface{}{}
	res = data
	data["start"] = time.Now().UnixMilli()
	defer func() {
		data["end"] = time.Now().UnixMilli()
	}()
	testOptions := &db.TestTaskOptions{
		Param:       this_.getParam(requestBean, c),
		OwnerName:   request.OwnerName,
		TestSql:     request.TestSql,
		Username:    request.Username,
		Password:    request.Password,
		IsBatch:     request.IsBatch,
		BatchSize:   request.BatchSize,
		ScriptVars:  request.ScriptVars,
		MaxOpenConn: request.MaxOpenConn,
		MaxIdleConn: request.MaxIdleConn,
	}
	executor, err := service.NewTestExecutor(testOptions)
	if err != nil {
		return
	}

	util.Logger.Info("test start", zap.Any("request", request), zap.Any("options", testOptions))

	if request.IsCallOnce {
		testOptions.IsBatch = false
		testOptions.BatchSize = 0
		data["request"] = request
		defer executor.Close()
		testOptions.OnExec = func(sqlList *[]string, sqlArgsList *[][]interface{}) {
			if sqlList != nil && len(*sqlList) > 0 {
				data["exeSql"] = (*sqlList)[0]
				data["exeSqlArgs"] = (*sqlArgsList)[0]
			}
		}
		var param = &task.ExecutorParam{}
		err = executor.Before(param)
		if err != nil {
			return
		}
		err = executor.Execute(param)
		if err != nil {
			return
		}
		return
	}

	var t *task.Task
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
	if request.CountSecond > 0 {
		t.Metric.SetCountSecond(request.CountSecond)
	}
	if request.CountTop {
		t.Metric.SetCountTop(request.CountTop)
	}

	taskDir := this_.getTestDirByTaskKey(request, t.Key)
	exists, err := util.PathExists(taskDir)
	if err != nil {
		return
	}
	if !exists {
		err = os.MkdirAll(taskDir, 0777)
		if err != nil {
			return
		}
	}
	taskJsonPath := taskDir + "task.json"
	saveTask := func() {
		var saveData = make(map[string]interface{})
		saveData["task"] = t
		saveData["request"] = request
		saveData["executor"] = executor
		saveData["count"] = t.Metric.GetCount()

		f, e := os.Create(taskJsonPath)
		if e != nil {
			util.Logger.Error("create task file error", zap.Any("taskJsonPath", taskJsonPath), zap.Error(e))
			return
		}
		defer func() { _ = f.Close() }()
		bs, e := json.Marshal(saveData)
		if e != nil {
			util.Logger.Error("task data to json error", zap.Any("taskJsonPath", taskJsonPath), zap.Error(e))
			return
		}
		_, _ = f.Write(bs)
	}
	t.OnExecute = func(param *task.ExecutorParam) {
		if param == nil {
			return
		}
		if param.Error != nil {
			var sqlList []string
			var sqlArgsList [][]interface{}
			if param.Extend != nil {
				p, ok := param.Extend.(*db.TestWorkerParam)
				if ok {
					sqlList = p.SqlList
					sqlArgsList = p.SqlArgsList
				}
			}

			util.Logger.Error("test task execute error", zap.Any("sqlList", sqlList), zap.Any("sqlArgsList", sqlArgsList), zap.Error(param.Error))
			return
		}
	}

	go func() {
		defer func() {
			saveTask()
			removeTestTask(t.Key)
			executor.Close()
		}()

		t.Run()
	}()
	saveTask()
	addTestTask(t)
	data["task"] = t
	data["executor"] = executor
	data["count"] = t.Metric.GetCount()
	return
}

func (this_ *api) testInfo(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.TaskKey == "" {
		err = errors.New("task key is empty")
		return
	}
	taskDir := this_.getTestDir(request)
	taskInfo, _ := this_.getTestTaskByDir(request.TaskKey, taskDir)
	if taskInfo == nil {
		return
	}
	res = taskInfo
	return
}

func (this_ *api) testList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	parentDir := this_.getTestParentDir(request)
	if exists, _ := util.PathExists(parentDir); !exists {
		return
	}
	var list []map[string]interface{}
	ds, err := os.ReadDir(parentDir)
	if err != nil {
		return
	}
	for _, d := range ds {
		if d.IsDir() {
			taskInfo, _ := this_.getTestTaskByDir(d.Name(), parentDir+d.Name()+"/")
			if taskInfo == nil {
				continue
			}
			list = append(list, taskInfo)
		}
	}
	res = list
	return
}

func (this_ *api) getTestTaskByDir(taskKey string, taskDir string) (res map[string]interface{}, err error) {

	if exists, _ := util.PathExists(taskDir + "/task.json"); !exists {
		return
	}
	bs, err := os.ReadFile(taskDir + "/task.json")
	if err != nil {
		return
	}
	res = map[string]interface{}{}
	err = util.JSONDecodeUseNumber(bs, &res)
	if err != nil {
		return
	}

	testTask := getTestTask(taskKey)
	if testTask != nil {
		res["task"] = testTask
		res["count"] = testTask.Metric.GetCount()
	}
	return
}

func (this_ *api) testStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.TaskKey == "" {
		err = errors.New("task key is empty")
		return
	}
	testTask := getTestTask(request.TaskKey)
	if testTask != nil {
		testTask.Stop()
	}

	return
}

func (this_ *api) testDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.TaskKey == "" {
		err = errors.New("task key is empty")
		return
	}
	testTask := getTestTask(request.TaskKey)
	if testTask != nil {
		testTask.Stop()
	}
	taskDir := this_.getTestDir(request)
	if exists, _ := util.PathExists(taskDir + "/task.json"); !exists {
		return
	}
	err = os.RemoveAll(taskDir)
	return
}

func (this_ *api) getTestDir(request *BaseRequest) (dir string) {
	dir = fmt.Sprintf("%s%s/", this_.getTestParentDir(request), request.TaskKey)
	return
}

func (this_ *api) getTestDirByTaskKey(request *BaseRequest, taskKey string) (dir string) {
	dir = fmt.Sprintf("%s%s/", this_.getTestParentDir(request), taskKey)
	return
}

func (this_ *api) getTestParentDir(request *BaseRequest) (dir string) {
	dir = fmt.Sprintf("%s%s/toolbox-%d/", this_.toolboxService.GetFilesDir(), "database-test", request.ToolboxId)
	return
}

var testTaskCache = map[string]*task.Task{}
var testTaskLocker = &sync.Mutex{}

func getTestTask(taskKey string) *task.Task {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	return testTaskCache[taskKey]
}

func addTestTask(task *task.Task) {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	testTaskCache[task.Key] = task
}

func removeTestTask(taskKey string) {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	delete(testTaskCache, taskKey)
}
