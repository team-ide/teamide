package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) invokeDaoStep(step *modelers.StepDaoModel, invokeData *InvokeData) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke dao step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke dao step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke dao step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	var varValue interface{}

	varValue, err = this_.InvokeDaoByName(step.Dao, invokeData)
	if err != nil {
		util.Logger.Error("invoke dao step error", zap.Any("error", err))
		return
	}

	if step.SetVar != "" {
		util.Logger.Debug("invoke dao set var", zap.Any("setVar", step.SetVar), zap.Any("setVarType", step.SetVarType), zap.Any("varValue", varValue))
		err = invokeData.AddVar(step.SetVar, varValue, step.SetVarType)
		if err != nil {
			util.Logger.Error("invoke dao set var error", zap.Any("setVar", step.SetVar), zap.Any("varValue", varValue), zap.Any("error", err))
			return
		}
	}
	return
}
