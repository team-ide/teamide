package common

import (
	"fmt"
	"strings"
	model2 "teamide/pkg/application/model"

	databaseSql "database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type ISqlExecutor interface {
	Ping() (err error)
	Select(sql string, params []interface{}, columnFieldMap map[string]*model2.StructFieldModel) (res []map[string]interface{}, err error)
	SelectOne(sql string, params []interface{}, columnFieldMap map[string]*model2.StructFieldModel) (res interface{}, err error)
	SelectPage(sql string, params []interface{}, currentPage int64, pageSize int64, columnFieldMap map[string]*model2.StructFieldModel) (res []map[string]interface{}, err error)
	SelectCount(countSql string, countParams []interface{}) (count int64, err error)
	Insert(sql string, params []interface{}) (rowsAffected int64, err error)
	Update(sql string, params []interface{}) (rowsAffected int64, err error)
	Delete(sql string, params []interface{}) (rowsAffected int64, err error)
	ExecSqls(sqls []string) error
}

type SqlParam struct {
	Database    string        `json:"database,omitempty"`
	CurrentPage int64         `json:"currentPage,omitempty"` // 请求页码
	PageSize    int64         `json:"pageSize,omitempty"`    // 请求页大小
	Sql         string        `json:"sql,omitempty"`
	Params      []interface{} `json:"params,omitempty"`
	CountSql    string        `json:"countSql,omitempty"`
	CountParams []interface{} `json:"countParams,omitempty"`
}

type SqlSelectPageResult struct {
	CurrentPage int64                    `json:"currentPage,omitempty"` // 请求页码
	PageSize    int64                    `json:"pageSize,omitempty"`    // 请求页大小
	TotalSize   int64                    `json:"totalSize,omitempty"`   // 总记录数
	TotalPage   int64                    `json:"totalPage,omitempty"`   // 请求页大小
	Data        []map[string]interface{} `json:"data,omitempty"`        // 页数据
}

func DatabaseIsMySql(database *model2.DatasourceDatabase) bool {
	return strings.ToLower(database.Type) == "mysql"
}
func DatabaseIsOracle(database *model2.DatasourceDatabase) bool {
	return strings.ToLower(database.Type) == "oracle"
}

func GetColumnFieldMap(app IApplication, name string) (columnFieldMap map[string]*model2.StructFieldModel, err error) {
	columnFieldMap = make(map[string]*model2.StructFieldModel)
	structModel := app.GetContext().GetStruct(name)
	if structModel == nil {
		return
	}
	if len(structModel.Fields) == 0 {
		return
	}
	for _, one := range structModel.Fields {
		if one.Column == "" {
			continue
		}
		columnFieldMap[one.Column] = one
	}
	return
}

type SqlExecutorDefault struct {
	config *model2.DatasourceDatabase
	db     *databaseSql.DB
}

func CreateSqlExecutor(config *model2.DatasourceDatabase) (executor *SqlExecutorDefault, err error) {
	executor = &SqlExecutorDefault{
		config: config,
	}
	err = executor.init()
	return
}

func (this_ *SqlExecutorDefault) init() (err error) {
	var dbType string
	var dbUrl string
	config := this_.config
	if DatabaseIsMySql(config) {
		dbType = "mysql"
		dbUrl = fmt.Sprint(config.Username, ":", config.Password, "@tcp(", config.Host, ":", config.Port, ")/", config.Database, "?charset=utf8mb4&loc=Local&parseTime=true")
	}
	var db *databaseSql.DB
	db, err = databaseSql.Open(dbType, dbUrl)
	if err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	this_.db = db
	err = this_.Ping()
	if err != nil {
		return
	}
	return
}

func (this_ *SqlExecutorDefault) Ping() (err error) {
	err = this_.db.Ping()
	if err != nil {
		return
	}
	return
}

func (this_ *SqlExecutorDefault) SelectCount(sql string, params []interface{}) (count int64, err error) {
	var rows *databaseSql.Rows
	rows, err = this_.db.Query(sql, params...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *SqlExecutorDefault) SelectOne(sql string, params []interface{}, columnFieldMap map[string]*model2.StructFieldModel) (res interface{}, err error) {
	var list []map[string]interface{}
	list, err = this_.Select(sql, params, columnFieldMap)
	if err != nil {
		return
	}
	if len(list) > 0 {
		res = list[0]
	}
	return
}

func (this_ *SqlExecutorDefault) SelectPage(sql string, params []interface{}, pageNumber int64, pageSize int64, columnFieldMap map[string]*model2.StructFieldModel) (res []map[string]interface{}, err error) {
	sql += fmt.Sprint(" LIMIT ", (pageNumber-1)*pageSize, ",", pageSize)
	res, err = this_.Select(sql, params, columnFieldMap)
	if err != nil {
		return
	}
	return
}

func (this_ *SqlExecutorDefault) Select(sql string, params []interface{}, columnFieldMap map[string]*model2.StructFieldModel) (res []map[string]interface{}, err error) {

	var rows *databaseSql.Rows
	rows, err = this_.db.Query(sql, params...)
	if err != nil {
		return
	}
	defer rows.Close()
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	res = []map[string]interface{}{}
	if rows.Next() {
		var values = []interface{}{}
		for index := range columns {
			column := columns[index]
			columnField := columnFieldMap[column]
			if columnField != nil {
				switch columnField.DataType {
				case "long", "int64":
					var value databaseSql.NullInt64
					values = append(values, &value)
				case "int", "int32":
					var value databaseSql.NullInt32
					values = append(values, &value)
				case "short", "int16":
					var value databaseSql.NullInt16
					values = append(values, &value)
				case "byte", "int8":
					var value databaseSql.NullByte
					values = append(values, &value)
				case "date", "datetime", "time", "time.Time", "Time":
					var value databaseSql.NullTime
					values = append(values, &value)
				case "boolean", "bool":
					var value databaseSql.NullBool
					values = append(values, &value)
				case "float", "float64":
					var value databaseSql.NullFloat64
					values = append(values, &value)
				case "double", "float32":
					var value databaseSql.NullFloat64
					values = append(values, &value)
				default:
					var value databaseSql.NullString
					values = append(values, &value)
				}
			} else {
				var value interface{}
				values = append(values, &value)
			}
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}

		one := map[string]interface{}{}
		for index := range columns {
			column := columns[index]
			value := values[index]
			columnField := columnFieldMap[column]
			if value != nil && columnField != nil {
				switch columnField.DataType {
				case "long", "int64":
					value = value.(*databaseSql.NullInt64).Int64
				case "int", "int32":
					value = value.(*databaseSql.NullInt32).Int32
				case "short", "int16":
					value = value.(*databaseSql.NullInt16).Int16
				case "byte", "int8":
					value = value.(*databaseSql.NullByte).Byte
				case "date", "datetime", "time", "time.Time", "Time":
					value = value.(*databaseSql.NullTime).Time
				case "boolean", "bool":
					value = value.(*databaseSql.NullBool).Bool
				case "float", "float64":
					value = value.(*databaseSql.NullFloat64).Float64
				case "double", "float32":
					value = value.(*databaseSql.NullFloat64).Float64
				default:
					value = value.(*databaseSql.NullString).String
				}
			}
			if columnField != nil && columnField.Name != "" {
				one[columnField.Name] = value
			} else {
				one[column] = value
			}
		}
		res = append(res, one)
	}
	return
}

func (this_ *SqlExecutorDefault) Insert(sql string, params []interface{}) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sql, params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (this_ *SqlExecutorDefault) Update(sql string, params []interface{}) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sql, params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (this_ *SqlExecutorDefault) Delete(sql string, params []interface{}) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sql, params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}
func (this_ *SqlExecutorDefault) ExecSqls(sqls []string) error {

	tx, err := this_.db.Begin()
	var hasError = false
	if err != nil {
		return err
	}
	defer func() {
		if hasError {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	for _, sql_ := range sqls {
		if sql_ == "" {
			continue
		}
		// -- 定义执行sql语句
		_, err = tx.Exec(sql_)
		if err != nil {
			hasError = true
			return err
		}
	}

	return err
}
