package module_database

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/goja"
	"sync"
	"teamide/pkg/base"
	"time"
)

func (this_ *api) testStart(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	data := map[string]interface{}{}
	res = data
	data["isCallOnce"] = request.IsCallOnce
	data["start"] = time.Now().UnixMilli()
	defer func() {
		data["end"] = time.Now().UnixMilli()
	}()

	if !request.IsCallOnce {

	} else {

	}
	return
}

func (this_ *api) testInfo(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	return
}

func (this_ *api) testList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	return
}

func (this_ *api) testStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	return
}

func (this_ *api) testDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	return
}

func (this_ *api) getTestDir(request *BaseRequest) (dir string) {
	dir = fmt.Sprintf("%s%s/", this_.getTestParentDir(request), request.TaskKey)
	return
}

func (this_ *api) getTestParentDir(request *BaseRequest) (dir string) {
	dir = fmt.Sprintf("%s/toolbox-%d/", "database-test", request.ToolboxId)
	return
}

type TestTask struct {
	*task.Task
	Request *BaseRequest `json:"request"`
}

var testTaskCache = map[string]*TestTask{}
var testTaskLocker = &sync.Mutex{}

func getTestTask(taskKey string) *TestTask {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	return testTaskCache[taskKey]
}

func addTestTask(task *TestTask) {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	testTaskCache[task.Key] = task
}

func removeTestTask(taskKey string) {
	testTaskLocker.Lock()
	defer testTaskLocker.Unlock()

	delete(testTaskCache, taskKey)
}

type testExecutor struct {
	task            *TestTask
	dbService       db.IService
	dataIndex       int64
	dataIndexLock   sync.Mutex
	workerParam     map[int]*TestWorkerParam
	workerParamLock sync.Mutex
}

type TestWorkerParam struct {
	sqlList       []string
	sqlParamsList [][]interface{}
	lock          sync.Mutex

	runtime       *goja.Runtime
	scriptContext map[string]interface{}
}

func (this_ *testExecutor) getWorkerParam(workerIndex int) (res *TestWorkerParam, err error) {
	this_.workerParamLock.Lock()
	defer this_.workerParamLock.Unlock()

	res = this_.workerParam[workerIndex]
	if res == nil {
		res = &TestWorkerParam{}
		res.runtime = goja.New()
		res.scriptContext = javascript.NewContext()
		if len(res.scriptContext) > 0 {
			for key, value := range res.scriptContext {
				err = res.runtime.Set(key, value)
				if err != nil {
					return
				}
			}
		}
		err = res.runtime.Set("workerIndex", workerIndex)
		if err != nil {
			return
		}
		this_.workerParam[workerIndex] = res
	}
	return
}

func (this_ *testExecutor) nextDataIndex(size int64) (dataIndex int64) {
	this_.dataIndexLock.Lock()
	defer this_.dataIndexLock.Unlock()

	dataIndex = this_.dataIndex
	this_.dataIndex += size
	return
}

func (this_ *TestWorkerParam) getScriptValue(param *task.ExecutorParam, dataIndex int64, workerParam *TestWorkerParam, script string) (res interface{}, err error) {

	err = workerParam.runtime.Set("index", dataIndex)
	if err != nil {
		return
	}

	v, err := workerParam.runtime.RunString(script)
	if err != nil {
		err = errors.New("get script [" + script + "] value error:" + err.Error())
		return
	}
	res = v.Export()
	return
}

func (this_ *TestWorkerParam) appendSql(param *task.ExecutorParam, dataIndex int64) (err error) {
	this_.lock.Lock()
	defer this_.lock.Unlock()

	return
}

func (this_ *testExecutor) initParam(param *task.ExecutorParam) (err error) {
	workerParam, err := this_.getWorkerParam(param.WorkerIndex)
	if err != nil {
		return
	}
	workerParam.sqlParamsList = [][]interface{}{}
	workerParam.sqlList = []string{}

	param.Extend = workerParam

	var genSize int64 = 1
	if this_.task.Request.IsBatch {
		genSize = this_.task.Request.BatchSize
	}
	if genSize <= 0 {
		return
	}
	dataIndex := this_.nextDataIndex(genSize)
	for i := int64(0); i < genSize; i++ {
		err = workerParam.appendSql(param, dataIndex)
		if err != nil {
			return
		}
		dataIndex++
	}
	return
}

func (this_ *testExecutor) Before(param *task.ExecutorParam) (err error) {

	err = this_.initParam(param)

	return
}

func (this_ *testExecutor) Execute(param *task.ExecutorParam) (err error) {

	workerParam := param.Extend.(*TestWorkerParam)

	_, err = this_.dbService.Execs(workerParam.sqlList, workerParam.sqlParamsList)

	return
}

func (this_ *testExecutor) After(param *task.ExecutorParam) (err error) {

	return
}
