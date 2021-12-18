package component

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"server/base"
	"server/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB MysqlService
)

func init() {
	var service interface{}
	var err error
	databaseConfig := DatabaseConfig{
		Host:     config.Config.Mysql.Host,
		Port:     config.Config.Mysql.Port,
		Database: config.Config.Mysql.Database,
		Username: config.Config.Mysql.Username,
		Password: config.Config.Mysql.Password,
	}
	Logger.Info(LogStr("数据库初始化:host:", databaseConfig.Host, ",port:", databaseConfig.Port, ",database:", databaseConfig.Database))
	service, err = CreateMysqlService(databaseConfig)
	if err != nil {
		panic(err)
	}
	DB = *service.(*MysqlService)

	_, err = DB.Exec(base.SqlParam{
		Sql:    "SELECT 1 FROM " + base.TABLE_INSTALL,
		Params: []interface{}{},
	})
	if err != nil {
		panic(err)
	}
	Logger.Info(LogStr("数据库连接成功!"))
}

type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MysqlService struct {
	config DatabaseConfig
	db     *sqlx.DB
}

func CreateMysqlService(config DatabaseConfig) (service *MysqlService, err error) {
	service = &MysqlService{
		config: config,
	}
	err = service.init()
	return
}

func (service *MysqlService) init() (err error) {
	url := fmt.Sprint(service.config.Username, ":", service.config.Password, "@tcp(", service.config.Host, ":", service.config.Port, ")/", service.config.Database, "?charset=utf8mb4&loc=Local&parseTime=true")
	var db *sqlx.DB
	db, err = sqlx.Open("mysql", url)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	service.db = db
	return
}

func (service *MysqlService) Open() (err error) {
	err = service.db.Ping()
	return
}

func (service *MysqlService) Close() (err error) {
	err = service.db.Close()
	return
}

func (service *MysqlService) InsertSqlByBean(table string, beans ...interface{}) (sqlParam base.SqlParam) {
	params := []interface{}{}
	if len(beans) == 0 {
		return
	}

	refType := base.GetRefType(beans[0])

	fieldCount := refType.NumField() // field count
	insertColumns := ""

	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		column := base.GetColumnNameByType(fieldType)
		if column != "" {
			insertColumns += column + ","
		}
	}
	insertValuesList := []string{}
	for _, bean := range beans {
		refValue := base.GetRefValue(bean) // value
		insertValues := ""
		for i := 0; i < fieldCount; i++ {
			fieldType := refType.Field(i)   // field type
			fieldValue := refValue.Field(i) // field vlaue
			column := base.GetColumnNameByType(fieldType)
			if column == "" {
				continue
			}
			value := base.GetFieldTypeValue(fieldType.Type, fieldValue)
			switch fieldType.Type.Name() {
			case "string":
				if value == "" {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			default:
				if value == nil {
					insertValues += "NULL,"
					continue
				} else if base.IsZero(value) {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			}
		}
		if len(insertValues) > 0 {
			insertValuesList = append(insertValuesList, insertValues)
		}
	}
	if len(insertColumns) > 0 {
		insertColumns = insertColumns[0 : len(insertColumns)-1]
	}
	sql := " INSERT INTO " + table + "(" + insertColumns + ") VALUES "
	for i, values := range insertValuesList {
		if i > 0 {
			sql += ","
		}
		values = values[0 : len(values)-1]
		sql += "("
		sql += values
		sql += ")"

	}
	return base.SqlParam{Sql: sql, Params: params}
}

func (service *MysqlService) UpdateSqlByBean(table string, keys []string, beans ...interface{}) (sqlParam base.SqlParam) {
	params := []interface{}{}
	if len(beans) == 0 {
		return
	}

	refType := base.GetRefType(beans[0])

	fieldCount := refType.NumField() // field count
	insertColumns := ""

	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		column := base.GetColumnNameByType(fieldType)
		if column != "" {
			insertColumns += column + ","
		}
	}
	insertValuesList := []string{}
	for _, bean := range beans {
		refValue := base.GetRefValue(bean) // value
		insertValues := ""
		for i := 0; i < fieldCount; i++ {
			fieldType := refType.Field(i)   // field type
			fieldValue := refValue.Field(i) // field vlaue
			column := base.GetColumnNameByType(fieldType)
			if column == "" {
				continue
			}
			value := base.GetFieldTypeValue(fieldType.Type, fieldValue)
			switch fieldType.Type.Name() {
			case "string":
				if value == "" {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			default:
				if value == nil {
					insertValues += "NULL,"
					continue
				} else if base.IsZero(value) {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			}
		}
		if len(insertValues) > 0 {
			insertValuesList = append(insertValuesList, insertValues)
		}
	}
	if len(insertColumns) > 0 {
		insertColumns = insertColumns[0 : len(insertColumns)-1]
	}
	sql := " INSERT INTO " + table + "(" + insertColumns + ") VALUES "
	for i, values := range insertValuesList {
		if i > 0 {
			sql += ","
		}
		values = values[0 : len(values)-1]
		sql += "("
		sql += values
		sql += ")"

	}
	return base.SqlParam{Sql: sql, Params: params}
}

func (service *MysqlService) GetColumnSqlByBean(refType reflect.Type, alias string) (columnSql string) {
	columnSql = ""
	fieldCount := refType.NumField() // field count
	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		column := base.GetColumnNameByType(fieldType)
		if column == "" {
			continue
		}
		if alias != "" {
			columnSql += alias + "."
		}
		columnSql += column + ","
	}
	columnSql = columnSql[0 : len(columnSql)-1]
	return
}

func (service *MysqlService) ResultToBeans(rows *sql.Rows, newBean func() interface{}) (list []interface{}, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	bean := newBean()
	columnTypes := base.GetColumnTypes(bean)

	list = []interface{}{}
	for rows.Next() {
		values := []interface{}{}
		beanColumnTypes := []base.ColumnType{}

		for _, column := range columns {
			columnType := columnTypes[column]

			if columnType.FieldType == nil {
				beanColumnTypes = append(beanColumnTypes, base.ColumnType{})
				continue
			}

			beanColumnTypes = append(beanColumnTypes, columnType)
			typeName := columnType.FieldType.Type.Name()
			if typeName == "int64" {
				var value sql.NullInt64
				values = append(values, &value)
			} else if typeName == "int32" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int16" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int8" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "float64" {
				var value sql.NullFloat64
				values = append(values, &value)
			} else if typeName == "float32" {
				var value sql.NullFloat64
				values = append(values, &value)
			} else if typeName == "bool" {
				var value sql.NullBool
				values = append(values, &value)
			} else if typeName == "time.Time" || typeName == "Time" {
				var value sql.NullTime
				values = append(values, &value)
			} else {
				var value sql.NullString
				values = append(values, &value)
			}
		}
		err = rows.Scan(values...)
		if err != nil {
			Logger.Error(LogStr("ResultToBeans error:", err))
			return nil, err
		}
		refValue := base.GetRefValue(bean)
		for index := range columns {
			beanColumnType := beanColumnTypes[index]
			if beanColumnType.Column == "" {
				continue
			}

			value := values[index]
			typeName := beanColumnType.FieldType.Type.Name()
			if value == nil {
				continue
			}
			if typeName == "int64" {
				v := value.(*sql.NullInt64)
				value = v.Int64
			} else if typeName == "int32" {
				v := value.(*sql.NullInt32)
				value = v.Int32
			} else if typeName == "int16" {
				v := value.(*sql.NullInt32)
				value = int16(v.Int32)
			} else if typeName == "int8" {
				v := value.(*sql.NullInt32)
				value = int8(v.Int32)
			} else if typeName == "int" {
				v := value.(*sql.NullInt32)
				value = int(v.Int32)
			} else if typeName == "float64" {
				v := value.(*sql.NullFloat64)
				value = v.Float64
			} else if typeName == "float32" {
				v := value.(*sql.NullFloat64)
				value = float32(v.Float64)
			} else if typeName == "bool" {
				v := value.(*sql.NullBool)
				value = v.Bool
			} else if typeName == "time.Time" || typeName == "Time" {
				v := value.(*sql.NullTime)
				value = v.Time
			} else {
				v := value.(*sql.NullString)
				value = v.String
			}
			val := reflect.ValueOf(value)
			refValue.FieldByName(beanColumnType.FieldType.Name).Set(val)
		}

		list = append(list, bean)

		bean = newBean()
	}
	return list, err
}

func (service *MysqlService) Query(sqlParam base.SqlParam, newBean func() interface{}) (list []interface{}, err error) {
	rows, err := service.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Query sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Query sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	list, err = service.ResultToBeans(rows, newBean)
	if err != nil {
		return
	}
	rows.Close()
	return
}

func (service *MysqlService) QueryOne(sqlParam base.SqlParam, newBean func() interface{}) (one interface{}, err error) {
	rows, err := service.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Query sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Query sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	var list []interface{}
	list, err = service.ResultToBeans(rows, newBean)
	if err != nil {
		return
	}

	if len(list) > 1 {
		err = errors.New("the result contains multiple pieces of data")
		return
	}
	if len(list) == 1 {
		one = list[0]
		return
	}
	rows.Close()
	return
}

func (service *MysqlService) Count(sqlParam base.SqlParam) (count int64, err error) {
	rows, err := service.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Count sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Count sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	rows.Close()
	return
}

func (service *MysqlService) QueryPage(sqlParam base.SqlParam, newBean func() interface{}) (page *base.PageBean, err error) {
	var total int64
	total, err = service.Count(base.SqlParam{CountSql: sqlParam.CountSql, CountParams: sqlParam.CountParams})
	if err != nil {
		return
	}
	page = &base.PageBean{
		Total:     total,
		PageSize:  sqlParam.PageSize,
		PageIndex: sqlParam.PageIndex,
	}
	if total > 0 {
		page.Init()
		if sqlParam.PageIndex > page.TotalPage {
			var res interface{}
			sqlParam.Sql += fmt.Sprint(" LIMIT ", (page.PageIndex-1)*page.PageSize, " , ", page.PageSize, " ")
			res, err = service.Query(sqlParam, newBean)
			if err != nil {
				return
			}
			page.Value = res
		}
	}
	return
}

func (service *MysqlService) Insert(sqlParam base.SqlParam) (rowsAffected int64, err error) {

	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Insert sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Insert sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Update(sqlParam base.SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Update sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Update sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Exec(sqlParam base.SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Exec sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Exec sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) Delete(sqlParam base.SqlParam) (rowsAffected int64, err error) {
	result, err := service.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		Logger.Error(LogStr("Delete sql error , sql:", sqlParam.Sql))
		Logger.Error(LogStr("Delete sql error , params:", base.ToJSON(sqlParam.Params)))
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (service *MysqlService) InsertBean(table string, one interface{}) (err error) {

	sqlParam := service.InsertSqlByBean(table, one)

	_, err = service.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func (service *MysqlService) BatchInsertBean(table string, list []interface{}) (err error) {

	sqlParam := service.InsertSqlByBean(table, list...)

	_, err = service.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func (service *MysqlService) UpdateBean(table string, keys []string, one interface{}) (err error) {
	if len(keys) == 0 {
		err = errors.New("update bean keys cannot be empty!")
	}
	sqlParam := service.UpdateSqlByBean(table, keys, one)

	_, err = service.Update(sqlParam)

	if err != nil {
		return
	}
	return
}

func (service *MysqlService) QueryBean(table string, one interface{}, newBean func() interface{}) (res []interface{}, err error) {
	sql := "SELECT * FROM " + table + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := base.NewSqlParam(sql, params)

	service.AppendWhere(one, &sqlParam)

	res, err = service.Query(sqlParam, newBean)

	if err != nil {
		return
	}
	return
}

func (service *MysqlService) CountBean(table string, one interface{}, newBean func() interface{}) (count int64, err error) {
	sql := "SELECT COUNT(*) FROM " + table + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := base.NewSqlParam(sql, params)

	service.AppendWhere(one, &sqlParam)

	count, err = service.Count(sqlParam)
	if err != nil {
		return
	}
	return
}
func (service *MysqlService) AppendWhere(one interface{}, sqlParam *base.SqlParam) {

}
