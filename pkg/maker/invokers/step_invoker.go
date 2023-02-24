package invokers

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

func (this_ *Invoker) InvokeStep(step interface{}, invokeData *InvokeData) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke step error", zap.Any("error", err))
		}
	}()
	if step == nil {
		err = errors.New("invoke step error,step is null")
		util.Logger.Error("invoke step error", zap.Any("error", err))
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
	case *modelers.StepRedisModel:
		stepModel = s.StepModel
		ok, err = this_.invokeStep(stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		ok, err = this_.invokeRedisStep(s, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		break
	case *modelers.StepDaoModel:
		stepModel = s.StepModel
		ok, err = this_.invokeStep(stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		ok, err = this_.invokeDaoStep(s, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		break
	case *modelers.StepDbModel:
		stepModel = s.StepModel
		ok, err = this_.invokeStep(stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		ok, err = this_.invokeDbStep(s, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		break
	default:
		err = errors.New("invoke step [" + util.GetRefType(step).Name() + "] can not be support")
		util.Logger.Error("invoke step error", zap.Any("error", err))
		return
	}

	return
}

func (this_ *Invoker) invokeStep(step *modelers.StepModel, invokeData *InvokeData) (ok bool, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke base step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke base step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke base step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	ok = true
	return
}
