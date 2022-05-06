package toolbox

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func CreateMysqlService(config DatabaseConfig) (service *MysqlService, err error) {
	service = &MysqlService{
		config: config,
	}
	service.lastUseTime = GetNowTime()
	err = service.init()
	return
}

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func ResultToMap(rows *sql.Rows) ([]map[string][]byte, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	list := []map[string][]byte{}
	for rows.Next() {

		values := []interface{}{}

		for range columnTypes {
			var value sql.RawBytes
			values = append(values, &value)
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		one := make(map[string][]byte)

		for index, column := range columns {
			v := values[index]
			value := v.(*sql.RawBytes)
			if value != nil {
				one[column] = (*value)
			} else {
				one[column] = nil
			}
		}

		list = append(list, one)
	}
	return list, err
}

type MysqlService struct {
	config      DatabaseConfig
	db          *sqlx.DB
	lastUseTime int64
}

func (this_ *MysqlService) init() (err error) {
	url := fmt.Sprint(this_.config.Username, ":", this_.config.Password, "@tcp(", this_.config.Host, ":", this_.config.Port, ")/?charset=utf8")
	var db *sqlx.DB
	db, err = sqlx.Open("mysql", url)
	this_.db = db
	return
}

func (this_ *MysqlService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *MysqlService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *MysqlService) SetLastUseTime() {
	this_.lastUseTime = GetNowTime()
}
func (this_ *MysqlService) Stop() {
	this_.db.Close()
}

func (this_ *MysqlService) Open() (err error) {
	err = this_.db.Ping()
	return
}

func (this_ *MysqlService) Close() (err error) {
	err = this_.db.Close()
	return
}

func (this_ *MysqlService) Databases() (databases []DatabaseInfo, err error) {
	sqlParam := SqlParam{
		Sql:    "show databases",
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		keys := make([]string, 0, len(one))
		for k := range one {
			keys = append(keys, k)
		}
		info := DatabaseInfo{}
		info.Name = string(one[keys[0]])
		databases = append(databases, info)
	}
	return
}

func (this_ *MysqlService) ShowCreateDatabase(database string) (create string, err error) {
	sqlParam := SqlParam{
		Sql:    "show create database `" + database + "`",
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		keys := make([]string, 0, len(one))
		for k := range one {
			keys = append(keys, k)
		}
		create = string(one[keys[1]])
	}
	return
}

func (this_ *MysqlService) ShowCreateTable(database string, table string) (create string, err error) {
	sqlParam := SqlParam{
		Sql:    "show create table `" + database + "`.`" + table + "`",
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		keys := make([]string, 0, len(one))
		for k := range one {
			keys = append(keys, k)
		}

		create = string(one[keys[1]])
	}
	return
}

func (this_ *MysqlService) Tables(database string) (tables []TableInfo, err error) {

	sql_ := "show table status from `" + database + "`"
	sqlParam := SqlParam{
		Sql:    sql_,
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		info := TableInfo{
			Name:    string(one["Name"]),
			Comment: string(one["Comment"]),
		}
		tables = append(tables, info)
	}
	return
}

func (this_ *MysqlService) TableDetail(database string, table string) (tableDetail TableDetailInfo, err error) {

	sql_ := "show table status from `" + database + "` where Name=?"
	sqlParam := SqlParam{
		Sql:    sql_,
		Params: []interface{}{table},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}
	tableDetail = TableDetailInfo{
		Name:    string(res[0]["Name"]),
		Comment: string(res[0]["Comment"]),
	}
	var columns []TableColumnInfo
	columns, err = this_.TableColumns(database, table)
	if err != nil {
		return
	}
	tableDetail.Columns = columns

	var indexs []TableIndexInfo
	indexs, err = this_.TableIndexs(database, table)
	if err != nil {
		return
	}
	tableDetail.Indexs = indexs
	return
}

func (this_ *MysqlService) TableColumns(database string, table string) (columns []TableColumnInfo, err error) {

	sql_ := "show full columns from  `" + database + "`.`" + table + "`"
	sqlParam := SqlParam{
		Sql:    sql_,
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		info := TableColumnInfo{
			Name:    string(one["Field"]),
			Comment: string(one["Comment"]),
		}
		if one["Key"] != nil {
			key := string(one["Key"])
			if key == "PRI" {
				info.PrimaryKey = true
			}
		}
		if one["Null"] != nil {
			null := string(one["Null"])
			if null == "NO" {
				info.NotNull = true
			}
		}
		columnTypeStr := string(one["Type"])
		columnType := columnTypeStr
		if strings.Contains(columnTypeStr, "(") {
			columnType = columnTypeStr[0:strings.Index(columnTypeStr, "(")]
			lengthStr := columnTypeStr[strings.Index(columnTypeStr, "(")+1 : strings.Index(columnTypeStr, ")")]
			if strings.Contains(lengthStr, ",") {
				length, _ := strconv.Atoi(lengthStr[0:strings.Index(lengthStr, ",")])
				decimal, _ := strconv.Atoi(lengthStr[strings.Index(lengthStr, ",")+1:])
				info.Length = length
				info.Decimal = decimal
			} else {
				length, _ := strconv.Atoi(lengthStr)
				info.Length = length
			}
		}
		info.Type = columnType
		columns = append(columns, info)
	}
	return
}

func (this_ *MysqlService) TableIndexs(database string, table string) (indexs []TableIndexInfo, err error) {

	sql_ := "show indexes from  `" + database + "`.`" + table + "`"
	sqlParam := SqlParam{
		Sql:    sql_,
		Params: []interface{}{},
	}
	res, err := this_.Query(sqlParam)
	if err != nil {
		return
	}
	for _, one := range res {
		Key_name := string(one["Key_name"])
		if Key_name == "PRIMARY" {
			continue
		}
		info := TableIndexInfo{
			Name:    Key_name,
			Comment: string(one["Comment"]),
		}
		indexs = append(indexs, info)
	}
	return
}

func (this_ *MysqlService) Datas(datasParam DatasParam) (datasResult DatasResult, err error) {

	sql_ := "SELECT "
	countSql := "SELECT "
	params_ := []interface{}{}

	columnMap := make(map[string]TableColumnInfo)

	for _, column := range datasParam.Columns {
		columnMap[column.Name] = column
		sql_ += "`" + column.Name + "`,"
	}
	sql_ = sql_[0 : len(sql_)-1]

	sql_ += " FROM `" + datasParam.Database + "`.`" + datasParam.Table + "` "
	countSql += " COUNT(1) AS total FROM `" + datasParam.Database + "`.`" + datasParam.Table + "` "

	if len(datasParam.Wheres) > 0 {
		whereSql := "WHERE "
		for index, where := range datasParam.Wheres {
			value := where.Value
			whereSql += "`" + where.Name + "` "
			switch where.SqlConditionalOperation {
			case "like":
				whereSql += "LIKE '%" + value + "%' "
			case "not like":
				whereSql += "NOT LIKE '%" + value + "%' "
			case "like start":
				whereSql += "LIKE '" + value + "%' "
			case "not like start":
				whereSql += "NOT LIKE '" + value + "%' "
			case "like end":
				whereSql += "LIKE '%" + value + "' "
			case "not like end":
				whereSql += "NOT LIKE '%" + value + "' "
			case "is null":
				whereSql += "IS NULL "
			case "is not null":
				whereSql += "IS NOT NULL "
			case "is empty":
				whereSql += "= '' "
			case "is not empty":
				whereSql += "<> '' "
			case "between":
				whereSql += "BETWEEN " + "'" + where.Before + "' AND '" + where.After + "' "
			case "not between":
				whereSql += "NOT BETWEEN " + "'" + where.Before + "' AND '" + where.After + "' "
			case "in":
				whereSql += "IN " + "(" + value + ") "
			case "not in":
				whereSql += "NOT IN " + "(" + value + ") "
			default:
				whereSql += where.SqlConditionalOperation + " '" + value + "' "
			}
			// params_ = append(params_, where.Value)
			if index < len(datasParam.Wheres)-1 {
				whereSql += " " + where.AndOr + " "
			}
		}
		sql_ += whereSql
		countSql += whereSql
	}
	sql_ = fmt.Sprint(sql_, " LIMIT ", (datasParam.PageIndex-1)*datasParam.PageSize, ",", datasParam.PageSize)

	totalRes, err := this_.Query(SqlParam{
		Sql:    countSql,
		Params: params_,
	})
	if err != nil {
		return
	}
	total := string(totalRes[0]["total"])
	res, err := this_.Query(SqlParam{
		Sql:    sql_,
		Params: params_,
	})
	if err != nil {
		return
	}
	datas := []map[string]interface{}{}
	for _, one := range res {
		data := make(map[string]interface{})
		for key, value := range one {
			column := columnMap[key]
			var value_ interface{}
			if value != nil {
				if column.Type != "" {
					value_ = string(value)
				} else {
					value_ = value
				}
			}
			data[key] = value_
		}
		datas = append(datas, data)
	}
	datasResult.Sql = sql_
	datasResult.Total = total
	datasResult.Params = params_
	datasResult.Datas = datas
	return
}

func (this_ *MysqlService) Query(sqlParam SqlParam) (res []map[string][]byte, err error) {
	rows, err := this_.db.Query(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		return
	}
	res, err = ResultToMap(rows)
	if err != nil {
		return
	}
	rows.Close()
	return
}

func (this_ *MysqlService) Insert(sqlParam SqlParam) (rowsAffected int64, err error) {

	result, err := this_.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (this_ *MysqlService) Update(sqlParam SqlParam) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (this_ *MysqlService) Delete(sqlParam SqlParam) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sqlParam.Sql, sqlParam.Params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}
