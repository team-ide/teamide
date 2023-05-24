package module_thrift

import (
	"errors"
	go_thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/thrift"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type invokeExecutor struct {
	*BaseRequest
	*argFormat
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

func NewClient(request *BaseRequest) (client *thrift.ServiceClient, err error) {
	var protocolFactory go_thrift.TProtocolFactory

	switch request.ProtocolFactory {
	case "compact":
		protocolFactory = go_thrift.NewTCompactProtocolFactoryConf(nil)
	case "simpleJSON":
		protocolFactory = go_thrift.NewTSimpleJSONProtocolFactoryConf(nil)
	case "json":
		protocolFactory = go_thrift.NewTJSONProtocolFactory()
	case "binary":
		protocolFactory = go_thrift.NewTBinaryProtocolFactoryConf(nil)
	default:
		protocolFactory = go_thrift.NewTBinaryProtocolFactoryConf(nil)
	}

	var transportFactory go_thrift.TTransportFactory
	if request.Buffered {
		transportFactory = go_thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = go_thrift.NewTTransportFactory()
	}

	if request.Framed {
		transportFactory = go_thrift.NewTFramedTransportFactoryConf(transportFactory, nil)
	}

	transport := go_thrift.NewTSocketConf(request.ServerAddress, nil)
	_ = transport.SetConnTimeout(time.Millisecond * time.Duration(request.Timeout))
	_ = transport.SetSocketTimeout(time.Millisecond * time.Duration(request.Timeout))

	if err = transport.Open(); err != nil {
		err = errors.New("opening socket to " + request.ServerAddress + " error:" + err.Error())
		return
	}
	var useTransport go_thrift.TTransport
	useTransport, err = transportFactory.GetTransport(transport)
	if err != nil {
		err = errors.New("transportFactory.GetTransport error:" + err.Error())
		return
	}
	client = thrift.NewServiceClientFactory(useTransport, protocolFactory)
	return
}

func (this_ *invokeExecutor) getClient(param *task.ExecutorParam) (client *thrift.ServiceClient, err error) {
	this_.workerClientLock.Lock()
	defer this_.workerClientLock.Unlock()

	client = this_.workerClient[param.WorkerIndex]
	if client != nil {
		return
	}

	client, err = NewClient(this_.BaseRequest)
	if err != nil {
		return
	}

	this_.workerClient[param.WorkerIndex] = client
	return
}
func (this_ *invokeExecutor) Before(param *task.ExecutorParam) (err error) {
	args, err := this_.formatArgs(this_.args, param)
	if err != nil {
		return
	}

	methodParam, err := this_.service.GetMethodParam(this_.filename, this_.ServiceName, this_.MethodName, args...)
	if err != nil {
		return
	}
	param.Extend = methodParam

	_, err = this_.getClient(param)
	if err != nil {
		methodParam.Error = err.Error()
		return
	}

	//util.Logger.Info("test Before", zap.Any("param", param))
	return
}

func (this_ *invokeExecutor) Execute(param *task.ExecutorParam) (err error) {
	methodParam := param.Extend.(*thrift.MethodParam)
	client, err := this_.getClient(param)
	if err != nil {
		methodParam.Error = err.Error()
		return
	}
	//util.Logger.Info("test Execute", zap.Any("param", param))

	_, err = client.Send(context.Background(), methodParam)
	if err != nil {
		methodParam.Error = err.Error()
	}
	this_.addParam(methodParam)
	return
}

func (this_ *invokeExecutor) After(param *task.ExecutorParam) (err error) {
	//util.Logger.Info("test After", zap.Any("param", param))
	return
}
