package db

import (
	"context"
	"database/sql"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/util"
	"time"
)

type executeTask struct {
	config       DatabaseConfig
	databaseType *DatabaseType
	dia          dialect.Dialect
	*Param
	ownerName string
}

func (this_ *executeTask) run(sqlContent string) (executeList []map[string]interface{}, errStr string, err error) {
	var executeData map[string]interface{}
	var query func(query string, args ...any) (*sql.Rows, error)
	var exec func(query string, args ...any) (sql.Result, error)

	config := this_.config
	config.MaxIdleConns = 2
	config.MaxIdleConns = 2
	if this_.ExecUsername != "" {
		config.Username = this_.ExecUsername
	}
	if this_.ExecPassword != "" {
		config.Password = this_.ExecPassword
	}
	switch this_.databaseType.DialectName {
	case "mysql":
		config.Database = this_.ownerName
		break
	default:
		config.Schema = this_.ownerName
		break
	}
	workDb, err := this_.databaseType.newDb(&config)
	if err != nil {
		util.Logger.Error("ExecuteSQL new db pool error", zap.Error(err))
		return
	}
	defer func() {
		_ = workDb.Close()
	}()
	cxt := context.Background()
	conn, err := workDb.Conn(cxt)
	if err != nil {
		util.Logger.Error("ExecuteSQL Conn error", zap.Error(err))
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	if this_.ownerName != "" {
		switch this_.databaseType.DialectName {
		case "mysql":
			_, err = conn.ExecContext(cxt, " USE "+this_.ownerName)
			if err != nil {
				util.Logger.Error("ExecuteSQL mysql change schema error", zap.Error(err))
				return
			}
			break
		case "oracle":
			_, err = conn.ExecContext(cxt, "ALTER SESSION SET CURRENT_SCHEMA="+this_.ownerName)
			if err != nil {
				util.Logger.Error("ExecuteSQL oracle change schema error", zap.Error(err))
				return
			}
			break
		}
	}
	var hasError bool
	if this_.OpenTransaction {
		var tx *sql.Tx
		tx, err = conn.BeginTx(cxt, nil)
		if err != nil {
			util.Logger.Error("ExecuteSQL BeginTx error", zap.Error(err))
			return
		}
		defer func() {
			if hasError {
				err = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()
		query = tx.Query
		exec = tx.Exec
	} else {
		query = func(query string, args ...any) (*sql.Rows, error) {
			return workDb.QueryContext(cxt, query, args...)
		}
		exec = func(query string, args ...any) (sql.Result, error) {
			return conn.ExecContext(cxt, query, args...)
		}
	}
	sqlList := this_.dia.SqlSplit(sqlContent)
	for _, executeSql := range sqlList {
		executeData, err = this_.execExecuteSQL(executeSql, query, exec)
		executeList = append(executeList, executeData)
		if err != nil {
			util.Logger.Error("ExecuteSQL execExecuteSQL error", zap.Any("executeSql", executeSql), zap.Error(err))
			errStr = err.Error()
			hasError = true
			if !this_.ErrorContinue {
				return
			}
			err = nil
		}
	}
	return
}

func (this_ *executeTask) execExecuteSQL(executeSql string,
	query func(query string, args ...any) (*sql.Rows, error),
	exec func(query string, args ...any) (sql.Result, error),
) (executeData map[string]interface{}, err error) {

	executeData = map[string]interface{}{}
	var startTime = util.Now()
	executeData["sql"] = executeSql
	executeData["startTime"] = util.Format(startTime)

	defer func() {
		var endTime = time.Now()
		executeData["endTime"] = util.Format(endTime)
		executeData["isEnd"] = true
		executeData["useTime"] = util.GetTimeTime(endTime) - util.GetTimeTime(startTime)
		if err != nil {
			executeData["error"] = err.Error()
			return
		}
	}()

	str := strings.ToLower(executeSql)
	if strings.HasPrefix(str, "select") ||
		strings.HasPrefix(str, "show") {
		executeData["isSelect"] = true
		// 查询
		var rows *sql.Rows
		rows, err = query(executeSql)
		if err != nil {
			return
		}
		defer func() {
			_ = rows.Close()
		}()
		var columnTypes []*sql.ColumnType
		columnTypes, err = rows.ColumnTypes()
		if err != nil {
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
		var dataList []map[string]interface{}
		for rows.Next() {
			cache := worker.GetSqlValueCache(columnTypes) //临时存储每行数据
			err = rows.Scan(cache...)
			if err != nil {
				return
			}
			item := make(map[string]interface{})
			for index, data := range cache {
				name := columnTypes[index].Name()
				switch tV := data.(type) {
				case time.Time:
					if tV.IsZero() {
						item[name] = nil
					} else {
						item[name] = util.GetTimeTime(tV)
					}
				default:
					item[name] = worker.GetSqlValue(columnTypes[index], data)
				}
			}
			dataList = append(dataList, item)
		}
		executeData["dataList"] = dataList
	} else if strings.HasPrefix(str, "insert") {
		executeData["isInsert"] = true
		var result sql.Result
		result, err = exec(executeSql)
		if err != nil {
			return
		}
		executeData["rowsAffected"], _ = result.RowsAffected()
	} else if strings.HasPrefix(str, "update") {
		executeData["isUpdate"] = true
		var result sql.Result
		result, err = exec(executeSql)
		if err != nil {
			return
		}
		executeData["rowsAffected"], _ = result.RowsAffected()
	} else if strings.HasPrefix(str, "delete") {
		executeData["isDelete"] = true
		var result sql.Result
		result, err = exec(executeSql)
		if err != nil {
			return
		}
		executeData["rowsAffected"], _ = result.RowsAffected()
	} else {
		executeData["isExec"] = true
		var result sql.Result
		result, err = exec(executeSql)
		if err != nil {
			return
		}
		executeData["rowsAffected"], _ = result.RowsAffected()
	}

	return
}
