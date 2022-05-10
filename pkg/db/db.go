package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"go.uber.org/zap"
	"strings"
)

var (
	Logger *zap.Logger
)

func init() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	Logger, _ = loggerConfig.Build()
}

// DatabaseWorker 基础操作
type DatabaseWorker interface {
	GetConfig() (config DatabaseConfig)
	Open() (err error)
	Close() (err error)
	Exec(sql string, params []interface{}) (rowsAffected int64, err error)
	Execs(sqlList []string, paramsList [][]interface{}) (rowsAffected int64, err error)
	Count(sql string, params []interface{}) (count int64, err error)
	FinderCount(finder *zorm.Finder) (count int64, err error)
	Query(sql string, params []interface{}, list interface{}) (err error)
	FinderQuery(finder *zorm.Finder, list interface{}) (err error)
	QueryOne(sql string, params []interface{}, one interface{}) (find bool, err error)
	FinderQueryOne(finder *zorm.Finder, one interface{}) (find bool, err error)
	QueryMap(sql string, params []interface{}) (list []map[string]interface{}, err error)
	FinderQueryMap(finder *zorm.Finder) (list []map[string]interface{}, err error)
	QueryPage(sql string, params []interface{}, list interface{}, page *zorm.Page) (err error)
	FinderQueryPage(finder *zorm.Finder, list interface{}, page *zorm.Page) (err error)
	QueryMapPage(sql string, params []interface{}, page *zorm.Page) (list []map[string]interface{}, err error)
	FinderQueryMapPage(finder *zorm.Finder, page *zorm.Page) (list []map[string]interface{}, err error)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `json:"type,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int32  `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// NewDatabaseWorker 根据数据库配置创建DatabaseWorker
func NewDatabaseWorker(config DatabaseConfig) (databaseWorker DatabaseWorker, err error) {
	switch strings.ToLower(config.Type) {
	case "mysql":
		databaseWorker, err = NewDatabaseMysql(config)
		break
	case "sqlite":
		databaseWorker, err = NewDatabaseSqlite(config)
		break
	default:
		err = errors.New(fmt.Sprintf("数据库类型[%s]未适配", config.Type))
	}
	if err != nil {
		return nil, err
	}
	return databaseWorker, nil
}

func formatRowsValue(rows *sql.Rows, columnTypes map[string]string) (list []map[string]interface{}, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	list = []map[string]interface{}{}
	for rows.Next() {
		one := map[string]interface{}{}
		var values []interface{}

		for _, column := range columns {
			columnType := columnTypes[column]

			switch strings.ToLower(columnType) {
			case "int64", "long":
				var value sql.NullInt64
				values = append(values, &value)
			case "int32", "int":
				var value sql.NullInt32
				values = append(values, &value)
			case "int16":
				var value sql.NullInt16
				values = append(values, &value)
			case "int8":
				var value sql.NullByte
				values = append(values, &value)
			case "float64", "float32", "float":
				var value sql.NullFloat64
				values = append(values, &value)
			case "bool", "boolean":
				var value sql.NullBool
				values = append(values, &value)
			case "time", "time.Time":
				var value sql.NullTime
				values = append(values, &value)
			default:
				var value sql.NullString
				values = append(values, &value)
			}
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for index, column := range columns {
			columnValue := values[index]
			columnType := columnTypes[column]

			switch strings.ToLower(columnType) {
			case "int64", "long":
				sqlValue := columnValue.(*sql.NullInt64)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "int32", "int":
				sqlValue := columnValue.(*sql.NullInt32)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "int16":
				sqlValue := columnValue.(*sql.NullInt16)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "int8":
				sqlValue := columnValue.(*sql.NullByte)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "float64", "float32", "float":
				sqlValue := columnValue.(*sql.NullFloat64)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "bool", "boolean":
				sqlValue := columnValue.(*sql.NullBool)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			case "time", "time.Time":
				sqlValue := columnValue.(*sql.NullTime)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			default:
				sqlValue := columnValue.(*sql.NullString)
				if sqlValue.Valid {
					one[column], err = sqlValue.Value()
				}
			}
			if err != nil {
				return nil, err
			}
		}
		list = append(list, one)

	}
	return list, err
}
