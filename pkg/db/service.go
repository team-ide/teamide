package db

import (
	"fmt"
	"gitee.com/teamide/zorm"
	"strconv"
	"strings"
	"teamide/pkg/util"
	"time"
)

func CreateService(config DatabaseConfig) (service *Service, err error) {
	service = &Service{
		config: config,
	}
	service.lastUseTime = util.GetNowTime()
	err = service.init()
	return
}

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

type Service struct {
	config         DatabaseConfig
	lastUseTime    int64
	DatabaseWorker *DatabaseWorker
}

func (this_ *Service) init() (err error) {
	this_.DatabaseWorker, err = NewDatabaseWorker(this_.config)
	if err != nil {
		return
	}
	return
}

func (this_ *Service) GetDatabaseWorker() *DatabaseWorker {
	return this_.DatabaseWorker
}

func (this_ *Service) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *Service) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *Service) SetLastUseTime() {
	this_.lastUseTime = util.GetNowTime()
}

func (this_ *Service) Stop() {
	_ = this_.DatabaseWorker.Close()
}

func (this_ *Service) Databases() (databases []*DatabaseModel, err error) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.SCHEMATA", "SCHEMA_NAME name")

	finder.Append("ORDER BY SCHEMA_NAME")
	//执行查询
	listMap, err := this_.DatabaseWorker.FinderQueryMap(finder)
	if err != nil { //标记测试失败
		return
	}
	for _, one := range listMap {
		database := &DatabaseModel{
			Name: one["name"].(string),
		}
		databases = append(databases, database)
	}
	return
}

func (this_ *Service) Tables(database string) (tables []*TableModel, err error) {
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

func (this_ *Service) TableDetails(database string, table string) (tableDetails []*TableModel, err error) {

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

func (this_ *Service) TableColumnList(database string, table string) (columnList []*TableColumnModel, err error) {

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
		if util.ContainsString(keys, one.Name) >= 0 {
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

func (this_ *Service) TablePrimaryKeys(database string, table string) (keys []string, err error) {

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

func (this_ *Service) TableIndexList(database string, table string) (indexList []*TableIndexModel, err error) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.statistics", "INDEX_NAME name,NON_UNIQUE,INDEX_COMMENT comment,COLUMN_NAME")

	finder.Append("WHERE TABLE_SCHEMA=?", database)
	finder.Append(" AND TABLE_NAME=?", table)
	finder.Append(" AND INDEX_NAME != ?", "PRIMARY")
	var indexList_ []*TableIndexModel
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &indexList_)
	if err != nil { //标记测试失败
		return
	}

	for _, one := range indexList_ {

		if one.NONUnique == "0" {
			one.Type = "unique"
		}
		one.Columns = append(one.Columns, one.COLUMNName)

		var find *TableIndexModel
		for _, in := range indexList {
			if in.Name == one.Name {
				find = in
				break
			}
		}
		if find == nil {
			indexList = append(indexList, one)
		} else {
			find.Columns = append(find.Columns, one.Columns...)
		}

	}
	return
}

type DataListResult struct {
	Sql      string                   `json:"sql"`
	Total    int                      `json:"total"`
	Params   []interface{}            `json:"params"`
	DataList []map[string]interface{} `json:"dataList"`
}

func (this_ *Service) DataList(param *GenerateParam, database string, table string, columnList []*TableColumnModel, whereList []*Where, orderList []*Order, pageSize int, pageIndex int) (dataListResult DataListResult, err error) {

	sql, values, err := DataListSelectSql(param, database, table, columnList, whereList, orderList)
	if err != nil {
		return
	}

	finder := zorm.NewFinder()
	finder.InjectionCheck = false

	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = pageSize
	page.PageNo = pageIndex
	listMap, err := this_.DatabaseWorker.FinderQueryMapPage(finder, page)
	if err != nil {
		return
	}
	for _, one := range listMap {
		for k, v := range one {
			if v == nil {
				continue
			}
			switch tV := v.(type) {
			case time.Time:
				if tV.IsZero() {
					one[k] = nil
				} else {
					one[k] = util.GetTimeTime(tV)
				}
			default:
				one[k] = fmt.Sprint(tV)
			}
		}
	}
	dataListResult.Sql, err = zorm.WrapPageSQL(this_.DatabaseWorker.GetDBType(), sql, page)
	dataListResult.Params = values
	dataListResult.Total = page.TotalCount
	dataListResult.DataList = listMap
	return
}

func (this_ *Service) Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error) {
	res, err = this_.DatabaseWorker.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}
