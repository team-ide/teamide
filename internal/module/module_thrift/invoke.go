package module_thrift

import (
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/thrift"
	"golang.org/x/net/context"
	"sync"
)

type invokeExecutor struct {
	*BaseRequest
	filename         string
	args             []interface{}
	workerClient     map[int]*thrift.ServiceClient
	workerClientLock sync.Mutex
	service          *thrift.Workspace
	taskDir          string
	paramList        []*thrift.MethodParam
	paramListLock    sync.Mutex
}

func (this_ *invokeExecutor) getAndCleanParamList() (paramList []*thrift.MethodParam) {
	this_.paramListLock.Lock()
	defer this_.paramListLock.Unlock()
	paramList = this_.paramList
	this_.paramList = []*thrift.MethodParam{}
	return
}

func (this_ *invokeExecutor) addParam(param *thrift.MethodParam) {
	this_.paramListLock.Lock()
	defer this_.paramListLock.Unlock()

	this_.paramList = append(this_.paramList, param)
	return
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
	if err != nil {
		methodParam.Error = err
	}
	this_.addParam(methodParam)
	return
}

func (this_ *invokeExecutor) After(param *task.ExecutorParam) (err error) {
	//util.Logger.Info("test After", zap.Any("param", param))
	return
}
