package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
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

	dbService, err := this_.GetDbServiceByName(step.Datasource)
	if err != nil {
		util.Logger.Error("invoke db step get db service error", zap.Any("source", step.Datasource), zap.Any("error", err))
		return
	}

	var res interface{}

	switch step.GetType() {
	case modelers.DbSelectOneOne:

		util.Logger.Debug("invoke db select one start", zap.Any("table", step.Table))
		res, err = dbService.DatabaseWorker.QueryMap(`select 1`, []interface{}{})
		if err != nil {
			util.Logger.Error("db select one error", zap.Any("datasource", step.Datasource), zap.Any("error", err))
			return
		}

		util.Logger.Debug("invoke db select one end", zap.Any("table", step.Table), zap.Any("res", res))
		break
	default:
		err = errors.New("invoke db [" + step.Db + "] can not be support")
		util.Logger.Error("invoke db error", zap.Any("error", err))
		return
	}

	ok = true
	return
}
