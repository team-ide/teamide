package toolbox

import (
	"context"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"go.uber.org/zap"
	"io"
	"strconv"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"teamide/pkg/vitess/sqlparser"
	"time"
)

func init() {
}

type executeSQLTask struct {
	Key           string    `json:"key,omitempty"`
	Database      string    `json:"database,omitempty"`
	ExecuteSQL    string    `json:"executeSQL,omitempty"`
	IsEnd         bool      `json:"isEnd,omitempty"`
	StartTime     time.Time `json:"startTime,omitempty"`
	EndTime       time.Time `json:"endTime,omitempty"`
	Error         string    `json:"error,omitempty"`
	UseTime       int64     `json:"useTime"`
	IsStop        bool      `json:"isStop"`
	service       DatabaseService
	ExecuteList   []map[string]interface{} `json:"executeList,omitempty"`
	generateParam *db.GenerateParam
}

func (this_ *executeSQLTask) Stop() {
	this_.IsStop = true
}

func (this_ *executeSQLTask) Start() {
	this_.StartTime = time.Now()
	var err error
	defer func() {
		if err != nil {
			util.Logger.Error("SQL执行异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		if err := recover(); err != nil {
			util.Logger.Error("SQL执行异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	tokens := sqlparser.NewStringTokenizer(this_.ExecuteSQL)

	ctx := this_.service.GetDatabaseWorker().GetContext()

	if this_.Database != "" {
		finder := zorm.NewFinder()
		finder.InjectionCheck = false
		finder.Append("use `" + this_.Database + "`")
		_, err = zorm.QueryMap(ctx, finder, nil)
		if err != nil {
			return
		}
	}

	if this_.generateParam.OpenTransaction {
		_, err = zorm.Transaction(ctx, func(ctx context.Context) (res interface{}, err error) {
			err = this_.do(ctx, tokens)
			return
		})
	} else {
		err = this_.do(ctx, tokens)
	}
	return
}

func (this_ *executeSQLTask) do(ctx context.Context, tokens *sqlparser.Tokenizer) (err error) {
	for {
		var stmt sqlparser.Statement
		stmt, err = sqlparser.ParseNext(tokens)

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return
		}
		// 如果已经开启过事务，则不用再次开启
		if this_.generateParam.OpenTransaction {
			err = this_.doExecute(ctx, stmt)
		} else {
			_, err = zorm.Transaction(ctx, func(ctx context.Context) (res interface{}, err error) {
				err = this_.doExecute(ctx, stmt)
				return
			})
		}

		if err != nil {
			err = nil
			if this_.generateParam.ErrorContinue {
				continue
			}
			return
		}
	}

	return
}

func (this_ *executeSQLTask) doExecute(ctx context.Context, stmt sqlparser.Statement) (err error) {
	buf := sqlparser.NewTrackedBuffer(nil)
	stmt.Format(buf)
	sql := buf.String()

	var executeData = map[string]interface{}{}
	this_.ExecuteList = append(this_.ExecuteList, executeData)

	var startTime = util.Now()
	executeData["sql"] = sql
	executeData["startTime"] = util.Format(startTime)

	switch stmt.(type) {
	case *sqlparser.Select:
		err = this_.doSelect(ctx, sql, executeData)
	case *sqlparser.Insert:
		err = this_.doInsert(ctx, sql, executeData)
	case *sqlparser.Update:
		err = this_.doUpdate(ctx, sql, executeData)
	case *sqlparser.Delete:
		err = this_.doDelete(ctx, sql, executeData)
	case *sqlparser.Use:
		err = this_.doUse(ctx, sql, executeData)
	case *sqlparser.Show:
		err = this_.doSelect(ctx, sql, executeData)
	default:
		err = this_.doExec(ctx, sql, executeData)
	}

	var endTime = time.Now()
	executeData["endTime"] = util.Format(endTime)
	executeData["isEnd"] = true
	executeData["useTime"] = util.GetTimeTime(endTime) - util.GetTimeTime(startTime)
	if err != nil {
		executeData["error"] = err.Error()
		return
	}

	return
}
func (this_ *executeSQLTask) doSelect(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isSelect"] = true
	columnTypes, dataList, err := zorm.QueryMapAndColumnTypes(ctx, finder, nil)

	if err != nil {
		util.Logger.Error("doSelect异常", zap.Error(err))
		return
	}
	var columnList []map[string]interface{}
	if len(columnTypes) > 0 {
		for _, columnType := range columnTypes {
			column := map[string]interface{}{}
			column["name"] = columnType.Name()
			column["type"] = columnType.DatabaseTypeName()
			columnList = append(columnList, column)
		}
	}
	executeData["columnList"] = columnList
	for _, one := range dataList {
		for k, v := range one {
			switch tV := v.(type) {
			case int64:
				one[k] = strconv.FormatInt(tV, 10)
			case uint64:
				one[k] = strconv.FormatInt(int64(tV), 10)
			case float64:
				one[k] = strconv.FormatFloat(tV, 'f', -1, 64)
			case time.Time:
				if tV.IsZero() {
					one[k] = nil
				} else {
					one[k] = util.GetTimeTime(tV)
				}
			default:
				one[k] = tV
			}
		}
	}
	executeData["dataList"] = dataList
	return
}

func (this_ *executeSQLTask) doInsert(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isInsert"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		util.Logger.Error("doInsert异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doUpdate(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isUpdate"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		util.Logger.Error("doUpdate异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doDelete(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isDelete"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		util.Logger.Error("doDelete异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doUse(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isUse"] = true
	_, err = zorm.QueryMap(ctx, finder, nil)

	if err != nil {
		util.Logger.Error("doUse异常", zap.Error(err))
		return
	}
	return
}

func (this_ *executeSQLTask) doExec(ctx context.Context, sql string, executeData map[string]interface{}) (err error) {
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["isExec"] = true
	_, err = zorm.UpdateFinder(ctx, finder)

	if err != nil {
		util.Logger.Error("doExec异常", zap.Error(err))
		return
	}
	return
}
