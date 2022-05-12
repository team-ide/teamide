package toolbox

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"go.uber.org/zap"
	"io"
	"teamide/pkg/application/base"
	"teamide/pkg/util"
	"teamide/pkg/vitess/sqlparser"
	"time"
)

func init() {
}

type executeSQLTask struct {
	Key         string    `json:"key,omitempty"`
	Database    string    `json:"database,omitempty"`
	ExecuteSQL  string    `json:"executeSQL,omitempty"`
	IsEnd       bool      `json:"isEnd,omitempty"`
	StartTime   time.Time `json:"startTime,omitempty"`
	EndTime     time.Time `json:"endTime,omitempty"`
	Error       string    `json:"error,omitempty"`
	UseTime     int64     `json:"useTime"`
	IsStop      bool      `json:"isStop"`
	service     DatabaseService
	ExecuteList []map[string]interface{} `json:"executeList,omitempty"`
}

func (this_ *executeSQLTask) Stop() {
	this_.IsStop = true
}

func (this_ *executeSQLTask) Start() (err error) {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("根据保存数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	tokens := sqlparser.NewStringTokenizer(this_.ExecuteSQL)

	ctx := this_.service.GetDatabaseWorker().GetContext()

	_, err = zorm.Transaction(ctx, func(ctx context.Context) (res interface{}, err error) {

		if this_.Database != "" {
			finder := zorm.NewFinder()
			finder.InjectionCheck = false
			finder.Append(`use ` + this_.Database)
			_, err = zorm.UpdateFinder(ctx, finder)
			if err != nil {
				return
			}
		}
		for {
			var stmt sqlparser.Statement
			stmt, err = sqlparser.ParseNext(tokens)
			if err != nil {
				if err == io.EOF {
					err = nil
					break
				}
				return
			}
			buf := sqlparser.NewTrackedBuffer(nil)
			stmt.Format(buf)
			sql := buf.String()

			var executeData map[string]interface{}

			switch stmt.(type) {
			case *sqlparser.Select:
				executeData, err = this_.doSelect(ctx, sql)
			case *sqlparser.Insert:
				executeData, err = this_.doInsert(ctx, sql)
			case *sqlparser.Update:
				executeData, err = this_.doUpdate(ctx, sql)
			case *sqlparser.Delete:
				executeData, err = this_.doDelete(ctx, sql)
			case *sqlparser.Use:
				executeData, err = this_.doUse(ctx, sql)
			default:
				err = errors.New("未解析SQL类型[" + base.GetRefType(stmt).Name() + "]，SQL[" + sql + "]")
			}
			if err != nil {
				return
			}
			this_.ExecuteList = append(this_.ExecuteList, executeData)

		}

		return
	})
	return
}

func (this_ *executeSQLTask) doSelect(ctx context.Context, sql string) (executeData map[string]interface{}, err error) {
	executeData = map[string]interface{}{}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["sql"] = sql
	executeData["isSelect"] = true
	dataList, err := zorm.QueryMap(ctx, finder, nil)

	if err != nil {
		Logger.Error("doSelect异常", zap.Error(err))
		return
	}
	for _, one := range dataList {
		for k, v := range one {
			t, tOk := v.(time.Time)
			if tOk {
				if t.IsZero() {
					one[k] = nil
				} else {
					one[k] = util.GetTimeTime(t)
				}
			}
		}
	}
	executeData["dataList"] = dataList
	return
}

func (this_ *executeSQLTask) doInsert(ctx context.Context, sql string) (executeData map[string]interface{}, err error) {
	executeData = map[string]interface{}{}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["sql"] = sql
	executeData["isInsert"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		Logger.Error("doInsert异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doUpdate(ctx context.Context, sql string) (executeData map[string]interface{}, err error) {
	executeData = map[string]interface{}{}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["sql"] = sql
	executeData["isUpdate"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		Logger.Error("doUpdate异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doDelete(ctx context.Context, sql string) (executeData map[string]interface{}, err error) {
	executeData = map[string]interface{}{}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["sql"] = sql
	executeData["isDelete"] = true
	rowsAffected, err := zorm.UpdateFinder(ctx, finder)

	if err != nil {
		Logger.Error("doDelete异常", zap.Error(err))
		return
	}
	executeData["rowsAffected"] = rowsAffected
	return
}

func (this_ *executeSQLTask) doUse(ctx context.Context, sql string) (executeData map[string]interface{}, err error) {
	executeData = map[string]interface{}{}
	finder := zorm.NewFinder()
	finder.InjectionCheck = false
	finder.Append(sql)
	executeData["sql"] = sql
	executeData["isUse"] = true
	_, err = zorm.UpdateFinder(ctx, finder)

	if err != nil {
		Logger.Error("doUse异常", zap.Error(err))
		return
	}
	return
}
