package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) invokeDbStep(step *modelers.StepDbModel, invokeData *InvokeData) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke db step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke db step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke db step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	dbService, err := this_.GetDbServiceByName(step.Config)
	if err != nil {
		util.Logger.Error("invoke db step get db service error", zap.Any("config", step.Config), zap.Any("error", err))
		return
	}

	sqlInfo, sqlParams, err := this_.getDbSql(step, invokeData)
	if err != nil {
		util.Logger.Error("invoke db step get db sql error", zap.Any("error", err))
		return
	}

	var varValue interface{}
	switch step.GetType() {
	case modelers.DbSelectOneOne:

		util.Logger.Debug("invoke db select one start", zap.Any("sqlInfo", sqlInfo), zap.Any("sqlParams", sqlParams))
		var listMap []map[string]interface{}
		listMap, err = dbService.QueryMap(sqlInfo, sqlParams)
		if err != nil {
			util.Logger.Error("invoke db select one error", zap.Any("sqlInfo", sqlInfo), zap.Any("sqlParams", sqlParams), zap.Any("error", err))
			return
		}
		listSize := len(listMap)
		if listSize > 0 {
			if listSize > 1 {
				err = errors.New("invoke db select one has more result")
				util.Logger.Error("invoke db select error", zap.Any("listMap", listMap), zap.Any("error", err))
				return
			}
			varValue = listMap[0]
		}

		util.Logger.Debug("invoke db select one end", zap.Any("selectValue", varValue))
		break
	default:
		err = errors.New("invoke db [" + step.Db + "] can not be support")
		util.Logger.Error("invoke db error", zap.Any("error", err))
		return
	}

	if step.SetVar != "" {
		util.Logger.Debug("invoke db set var", zap.Any("setVar", step.SetVar), zap.Any("setVarType", step.SetVarType), zap.Any("varValue", varValue))
		err = invokeData.AddVar(step.SetVar, varValue, step.SetVarType)
		if err != nil {
			util.Logger.Error("invoke db set var error", zap.Any("setVar", step.SetVar), zap.Any("varValue", varValue), zap.Any("error", err))
			return
		}
	}
	return
}
