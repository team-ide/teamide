package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) getDbSql(step *modelers.StepDbModel, invokeData *InvokeData) (sqlInfo string, sqlParams []interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("get db sql error", zap.Any("error", err))
		}
		util.Logger.Debug("get db sql end", zap.Any("sqlInfo", sqlInfo), zap.Any("sqlParams", sqlParams))
	}()
	util.Logger.Debug("get db sql start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	dbService, err := this_.GetDbServiceByName(step.Config)
	if err != nil {
		util.Logger.Error("get db sql get db service error", zap.Any("config", step.Config), zap.Any("error", err))
		return
	}
	dia := dbService.GetDialect()
	util.Logger.Debug("get db sql dialect info", zap.Any("dialectType", dia.DialectType()))

	switch step.GetType() {
	case modelers.DbSelectOneOne:
		sqlInfo = `SELECT * FROM TB_USER WHERE userId = ?`
		sqlParams = append(sqlParams, 1)

		break
	default:
		err = errors.New("get db sql [" + step.Db + "] can not be support")
		util.Logger.Error("get db sql error", zap.Any("error", err))
		return
	}
	return
}
