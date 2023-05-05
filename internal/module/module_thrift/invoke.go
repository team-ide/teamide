package module_thrift

import (
	"crypto/tls"
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

	var protocolFactory go_thrift.TProtocolFactory

	switch this_.ProtocolFactory {
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
	cfg := &go_thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if this_.Buffered {
		transportFactory = go_thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = go_thrift.NewTTransportFactory()
	}

	if this_.Framed {
		transportFactory = go_thrift.NewTFramedTransportFactoryConf(transportFactory, cfg)
	}

	transport := go_thrift.NewTSocketConf(this_.ServerAddress, &go_thrift.TConfiguration{
		ConnectTimeout: time.Millisecond * time.Duration(this_.Timeout),
		SocketTimeout:  time.Millisecond * time.Duration(this_.Timeout),
	})

	var useTransport go_thrift.TTransport
	useTransport, err = transportFactory.GetTransport(transport)
	if err != nil {
		err = errors.New("transportFactory.GetTransport error:" + err.Error())
		return
	}
	client = thrift.NewServiceClientFactory(useTransport, protocolFactory)
	if err = transport.Open(); err != nil {
		err = errors.New("opening socket to " + this_.ServerAddress + " error:" + err.Error())
		return
	}

	this_.workerClient[param.WorkerIndex] = client
	return
}

func (this_ *invokeExecutor) Before(param *task.ExecutorParam) (err error) {
	_, err = this_.getClient(param)
	if err != nil {
		param.Error = err
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
