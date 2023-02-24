package invokers

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

func (this_ *Invoker) invokeDbStep(step *modelers.StepDbModel, invokeData *InvokeData) (ok bool, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke db step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke db step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke db step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	switch step.GetType() {
	case modelers.DbGet:

		util.Logger.Debug("invoke db get", zap.Any("key", step.Table))

		break
	}

	ok = true
	return
}
