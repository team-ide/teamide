package toolbox

import (
	"gitee.com/chunanyong/zorm"
	"github.com/wxnacy/wgo/arrays"
	"strconv"
	"strings"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

func CreateMysqlService(config db.DatabaseConfig) (service *MysqlService, err error) {
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

type MysqlService struct {
	config         db.DatabaseConfig
	lastUseTime    int64
	DatabaseWorker *db.DatabaseWorker
}

func (this_ *MysqlService) init() (err error) {
	this_.DatabaseWorker, err = db.NewDatabaseWorker(this_.config)
	if err != nil {
		return
	}
	return
}

func (this_ *MysqlService) GetDatabaseWorker() *db.DatabaseWorker {
	return this_.DatabaseWorker
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
	_ = this_.DatabaseWorker.Close()
}

func (this_ *MysqlService) Databases() (databases []*db.DatabaseModel, err error) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.SCHEMATA", "SCHEMA_NAME name")

	finder.Append("ORDER BY SCHEMA_NAME")
	//执行查询
	listMap, err := this_.DatabaseWorker.FinderQueryMap(finder)
	if err != nil { //标记测试失败
		return
	}
	for _, one := range listMap {
		database := &db.DatabaseModel{
			Name: one["name"].(string),
		}
		databases = append(databases, database)
	}
	return
}

func (this_ *MysqlService) Tables(database string) (tables []*db.TableModel, err error) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.tables", "TABLE_NAME AS name,TABLE_COMMENT AS comment")

	finder.Append("WHERE TABLE_SCHEMA=?", database)

	finder.Append("ORDER BY TABLE_NAME")
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &tables)
	if err != nil { //标记测试失败
		return
	}
	return
}

func (this_ *MysqlService) TableDetails(database string, table string) (tableDetails []*db.TableModel, err error) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.tables", "TABLE_NAME AS name,TABLE_COMMENT AS comment")

	finder.Append("WHERE TABLE_SCHEMA=?", database)
	if table != "" {
		finder.Append(" AND TABLE_NAME=?", table)
	}
	finder.Append(" ORDER BY TABLE_NAME")
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &tableDetails)
	if err != nil { //标记测试失败
		return
	}

	for _, one := range tableDetails {

		one.ColumnList, err = this_.TableColumnList(database, one.Name)
		if err != nil {
			return
		}

		one.IndexList, err = this_.TableIndexList(database, one.Name)
		if err != nil {
			return
		}

	}
	return
}

func (this_ *MysqlService) TableColumnList(database string, table string) (columnList []*db.TableColumnModel, err error) {

	keys, err := this_.TablePrimaryKeys(database, table)
	if err != nil {
		return
	}

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.columns", "COLUMN_NAME AS name,IS_NULLABLE,COLUMN_TYPE AS type,COLUMN_COMMENT AS comment")

	finder.Append(" WHERE TABLE_SCHEMA=?", database)
	finder.Append(" AND TABLE_NAME=?", table)
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &columnList)
	if err != nil { //标记测试失败
		return
	}
	for _, one := range columnList {
		if one.ISNullable == "NO" {
			one.NotNull = true
		}
		if arrays.ContainsString(keys, one.Name) >= 0 {
			one.PrimaryKey = true
		}
		columnTypeStr := one.Type
		columnType := columnTypeStr
		if strings.Contains(columnTypeStr, "(") {
			columnType = columnTypeStr[0:strings.Index(columnTypeStr, "(")]
			lengthStr := columnTypeStr[strings.Index(columnTypeStr, "(")+1 : strings.Index(columnTypeStr, ")")]
			if strings.Contains(lengthStr, ",") {
				length, _ := strconv.Atoi(lengthStr[0:strings.Index(lengthStr, ",")])
				decimal, _ := strconv.Atoi(lengthStr[strings.Index(lengthStr, ",")+1:])
				one.Length = length
				one.Decimal = decimal
			} else {
				length, _ := strconv.Atoi(lengthStr)
				one.Length = length
			}
		}
		one.Type = columnType
	}
	return
}

func (this_ *MysqlService) TablePrimaryKeys(database string, table string) (keys []string, err error) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.table_constraints t", "k.COLUMN_NAME")

	finder.Append(" JOIN information_schema.key_column_usage k USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME) ")
	finder.Append(" WHERE t.TABLE_SCHEMA=? AND t.TABLE_NAME=? AND t.CONSTRAINT_TYPE=? ", database, table, "PRIMARY KEY")
	//执行查询
	listMap, err := this_.DatabaseWorker.FinderQueryMap(finder)
	if err != nil { //标记测试失败
		return
	}

	for _, one := range listMap {
		keys = append(keys, one["COLUMN_NAME"].(string))
	}
	return
}

func (this_ *MysqlService) TableIndexList(database string, table string) (indexList []*db.TableIndexModel, err error) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.statistics", "INDEX_NAME name,NON_UNIQUE,INDEX_COMMENT comment,COLUMN_NAME columns")

	finder.Append("WHERE TABLE_SCHEMA=?", database)
	finder.Append(" AND TABLE_NAME=?", table)
	finder.Append(" AND INDEX_NAME != ?", "PRIMARY")
	var indexList_ []*db.TableIndexModel
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &indexList_)
	if err != nil { //标记测试失败
		return
	}

	for _, one := range indexList_ {

		var info *db.TableIndexModel
		if one.NONUnique == "0" {
			one.Type = "UNIQUE"
		}

		for _, in := range indexList {
			if in.Name == one.Name {
				info = in
				break
			}
		}
		if info == nil {
			indexList = append(indexList, one)
		} else {
			info.Columns += "," + one.Columns
		}

	}
	return
}

func (this_ *MysqlService) DataList(param *db.GenerateParam, dataListParam DataListParam) (dataListResult DataListResult, err error) {

	sql, values, err := db.DataListSelectSql(param, dataListParam.Database, dataListParam.Table, dataListParam.ColumnList, dataListParam.Wheres, dataListParam.Orders)
	if err != nil {
		return
	}

	finder := zorm.NewFinder()
	finder.InjectionCheck = false

	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = dataListParam.PageSize
	page.PageNo = dataListParam.PageIndex
	listMap, err := this_.DatabaseWorker.FinderQueryMapPage(finder, page)
	if err != nil {
		return
	}
	for _, one := range listMap {
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
	dataListResult.Sql = sql
	dataListResult.Params = values
	dataListResult.Total = page.TotalCount
	dataListResult.DataList = listMap
	return
}

func (this_ *MysqlService) Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error) {
	res, err = this_.DatabaseWorker.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}
