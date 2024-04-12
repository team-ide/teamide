package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) InvokeStep(from string, step interface{}, invokeData *InvokeData) (res interface{}, isReturn bool, err error) {
	funcInvoke := invokeStart(from+" step", invokeData)
	if step == nil {
		err = errors.New(funcInvoke.name + " error,step is null")
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	util.Logger.Debug(funcInvoke.name + " start")

	var stepModel *modelers.StepModel
	var ok bool

	switch s := step.(type) {
	case *modelers.StepModel:
		stepModel = s
		ok, err = this_.invokeBaseStep(funcInvoke.name, stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		break
	case *modelers.StepRedisModel:
		stepModel = s.StepModel
		ok, err = this_.invokeBaseStep(funcInvoke.name, stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		err = this_.invokeRedisStep(funcInvoke.name, s, invokeData)
		if err != nil {
			return
		}
		break
	case *modelers.StepDaoModel:
		stepModel = s.StepModel
		ok, err = this_.invokeBaseStep(funcInvoke.name, stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		err = this_.invokeDaoStep(s, invokeData)
		if err != nil {
			return
		}
		break
	case *modelers.StepDbModel:
		stepModel = s.StepModel
		ok, err = this_.invokeBaseStep(funcInvoke.name, stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		err = this_.invokeDbStep(s, invokeData)
		if err != nil {
			return
		}
		break
	case *modelers.StepZkModel:
		stepModel = s.StepModel
		ok, err = this_.invokeBaseStep(funcInvoke.name, stepModel, invokeData)
		if err != nil {
			return
		}
		if !ok {
			return
		}
		err = this_.invokeZkStep(s, invokeData)
		if err != nil {
			return
		}
		break
	default:
		err = errors.New(funcInvoke.name + " [" + reflect.TypeOf(step).Name() + "] can not be support")
		util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		return
	}

	if len(stepModel.Steps) > 0 {
		err = this_.InvokeSteps(funcInvoke.name, stepModel.Steps, invokeData)
		if err != nil {
			return
		}
	}

	if stepModel.Return != "" {
		isReturn = true
		res, err = invokeData.InvokeScript(stepModel.Return)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" get return value error", zap.Any("return", stepModel.Return), zap.Any("error", err))
			return
		}
		return
	}

	return
}

func (this_ *Invoker) invokeBaseStep(from string, step *modelers.StepModel, invokeData *InvokeData) (ok bool, err error) {
	funcInvoke := invokeStart(from+" base", invokeData)

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	util.Logger.Debug(funcInvoke.name+" start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	// 组装 变量

	// 验证 变量

	// 条件
	if step.If != "" {
		util.Logger.Debug(funcInvoke.name+" check if", zap.Any("if", step.If))
		var ifValue interface{}
		ifValue, err = invokeData.InvokeScript(step.If)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" check if error", zap.Any("if", step.If), zap.Any("error", err))
			return
		}
		if util.IsFalse(ifValue) {
			ok = false
			return
		}
	}

	ok = true
	return
}
