package invokers

import (
	"errors"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

func NewInvoker(app *modelers.Application) (runner *Invoker) {
	runner = &Invoker{
		app: app,
	}
	return
}

type Invoker struct {
	app *modelers.Application
}

func (this_ *Invoker) InvokeServiceByName(name string, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("invoke service by name error", zap.Any("error", e))
		}
	}()

	service := this_.app.GetService(name)
	if service == nil {
		err = errors.New("service [" + name + "] is not exist")
		return
	}
	res, err = this_.InvokeService(service, invokeData)
	return
}

func (this_ *Invoker) InvokeService(service *modelers.ServiceModel, invokeData *InvokeData) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("invoke service error", zap.Any("error", e))
		}
	}()
	if service == nil {
		err = errors.New("invoke service error,service is null")
		return
	}
	defer func() {
		util.Logger.Info("invoke service end", zap.Any("name", service.Name), zap.Any("args", invokeData.args))
	}()
	if invokeData == nil {
		invokeData = NewInvokeData(this_.app)
	}
	if invokeData.app == nil {
		invokeData.app = this_.app
	}
	util.Logger.Info("invoke service start", zap.Any("name", service.Name), zap.Any("args", invokeData.args))

	return
}

func (this_ *Invoker) InvokeStep(step interface{}, invokeData *InvokeData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("invoke step error", zap.Any("error", e))
		}
	}()
	if step == nil {
		err = errors.New("invoke step error,step is null")
		return
	}

	var stepModel *modelers.StepModel
	var ok bool

	switch s := step.(type) {
	case *modelers.StepModel:
		stepModel = s
		ok, err = this_.invokeStep(stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		break
	default:
		err = errors.New("invoke step [" + util.GetRefType(step).Name() + "] can not be support")
		return
	}

	defer func() {
		util.Logger.Info("invoke step end", zap.Any("args", invokeData.args))
	}()
	util.Logger.Info("invoke service start", zap.Any("args", invokeData.args))

	return
}

func (this_ *Invoker) invokeStep(step *modelers.StepModel, invokeData *InvokeData) (ok bool, err error) {
	ok = true
	return
}
