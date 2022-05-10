package toolbox

import (
	"gitee.com/chunanyong/zorm"
	"github.com/wxnacy/wgo/arrays"
	"strconv"
	"strings"
	"teamide/pkg/db"
	"teamide/pkg/sql_ddl"
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

func (this_ *MysqlService) Databases() (databases []*DatabaseInfo, err error) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.SCHEMATA", "SCHEMA_NAME name")

	finder.Append("ORDER BY SCHEMA_NAME")
	//执行查询
	listMap, err := this_.DatabaseWorker.FinderQueryMap(finder)
	if err != nil { //标记测试失败
		return
	}
	for _, one := range listMap {
		database := &DatabaseInfo{
			Name: one["name"].(string),
		}
		databases = append(databases, database)
	}
	return
}

func (this_ *MysqlService) Tables(database string) (tables []TableInfo, err error) {
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

func (this_ *MysqlService) TableDetails(database string, table string) (tableDetails []*sql_ddl.TableDetailInfo, err error) {

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

		one.Columns, err = this_.TableColumns(database, one.Name)
		if err != nil {
			return
		}

		one.Indexs, err = this_.TableIndexs(database, one.Name)
		if err != nil {
			return
		}

	}
	return
}

func (this_ *MysqlService) TableColumns(database string, table string) (columns []*sql_ddl.TableColumnInfo, err error) {

	keys, err := this_.TablePrimaryKeys(database, table)
	if err != nil {
		return
	}

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.columns", "COLUMN_NAME AS name,IS_NULLABLE,COLUMN_TYPE AS type,COLUMN_COMMENT AS comment")

	finder.Append(" WHERE TABLE_SCHEMA=?", database)
	finder.Append(" AND TABLE_NAME=?", table)
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &columns)
	if err != nil { //标记测试失败
		return
	}
	for _, one := range columns {
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

func (this_ *MysqlService) TableIndexs(database string, table string) (indexs []*sql_ddl.TableIndexInfo, err error) {

	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.statistics", "INDEX_NAME name,NON_UNIQUE,INDEX_COMMENT comment,COLUMN_NAME columns")

	finder.Append("WHERE TABLE_SCHEMA=?", database)
	finder.Append(" AND TABLE_NAME=?", table)
	finder.Append(" AND INDEX_NAME != ?", "PRIMARY")
	var indexs_ []*sql_ddl.TableIndexInfo
	//执行查询
	err = this_.DatabaseWorker.FinderQuery(finder, &indexs_)
	if err != nil { //标记测试失败
		return
	}

	for _, one := range indexs_ {

		var info *sql_ddl.TableIndexInfo
		if one.NONUnique == "0" {
			one.Type = "UNIQUE"
		}

		for _, in := range indexs {
			if in.Name == one.Name {
				info = in
				break
			}
		}
		if info == nil {
			indexs = append(indexs, one)
		} else {
			info.Columns += "," + one.Columns
		}

	}
	return
}

func (this_ *MysqlService) Datas(datasParam DatasParam) (datasResult DatasResult, err error) {

	var params []interface{}
	selectColumns := ""
	for _, column := range datasParam.Columns {
		selectColumns += "`" + column.Name + "`,"
	}
	selectColumns = selectColumns[0 : len(selectColumns)-1]
	//构造查询用的finder
	finder := zorm.NewSelectFinder(datasParam.Database+"."+datasParam.Table, selectColumns)
	if len(datasParam.Wheres) > 0 {
		finder.Append(" WHERE")
		for index, where := range datasParam.Wheres {
			value := where.Value
			switch where.SqlConditionalOperation {
			case "like":
				finder.Append(" "+where.Name+" LIKE ?", "%"+value+"%")
				params = append(params, "%"+value+"%")
			case "not like":
				finder.Append(" "+where.Name+" NOT LIKE ?", "%"+value+"%")
				params = append(params, "%"+value+"%")
			case "like start":
				finder.Append(" "+where.Name+" LIKE ?", ""+value+"%")
				params = append(params, ""+value+"%")
			case "not like start":
				finder.Append(" "+where.Name+" NOT LIKE ?", ""+value+"%")
				params = append(params, ""+value+"%")
			case "like end":
				finder.Append(" "+where.Name+" LIKE ?", "%"+value+"")
				params = append(params, "%"+value+"")
			case "not like end":
				finder.Append(" "+where.Name+" NOT LIKE ?", "%"+value+"")
				params = append(params, "%"+value+"")
			case "is null":
				finder.Append(" " + where.Name + " IS NULL")
			case "is not null":
				finder.Append(" " + where.Name + " IS NOT NULL")
			case "is empty":
				finder.Append(" "+where.Name+" = ?", "")
				params = append(params, "")
			case "is not empty":
				finder.Append(" "+where.Name+" <> ?", "")
				params = append(params, "")
			case "between":
				finder.Append(" "+where.Name+" BETWEEN ? AND ?", where.Before, where.After)
				params = append(params, where.Before, where.After)
			case "not between":
				finder.Append(" "+where.Name+" NOT BETWEEN ? AND ?", where.Before, where.After)
				params = append(params, where.Before, where.After)
			case "in":
				finder.Append(" "+where.Name+" IN (?)", value)
				params = append(params, value)
			case "not in":
				finder.Append(" "+where.Name+" NOT IN (?)", value)
				params = append(params, value)
			default:
				finder.Append(" "+where.Name+" "+where.SqlConditionalOperation+" ?", value)
				params = append(params, value)
			}
			// params_ = append(params_, where.Value)
			if index < len(datasParam.Wheres)-1 {
				finder.Append(" " + where.AndOr + " ")
			}
		}
	}
	if len(datasParam.Orders) > 0 {
		finder.Append(" ORDER BY")
		for index, order := range datasParam.Orders {
			finder.Append(" " + order.Name)
			if order.DescAsc != "" {
				finder.Append(" " + order.DescAsc)
			}
			// params_ = append(params_, where.Value)
			if index < len(datasParam.Orders)-1 {
				finder.Append(",")
			}
		}

	}
	page := zorm.NewPage()
	page.PageSize = datasParam.PageSize
	page.PageNo = datasParam.PageIndex
	listMap, err := this_.DatabaseWorker.FinderQueryMapPage(finder, page)
	if err != nil {
		return
	}
	for _, one := range listMap {
		for k, v := range one {
			t, tOk := v.(time.Time)
			if tOk {
				if t.IsZero() {
					one[k] = 0
				} else {
					one[k] = util.GetTimeTime(t)
				}
			}
		}
	}
	datasResult.Sql, err = finder.GetSQL()
	datasResult.Params = params
	datasResult.Total = page.TotalCount
	datasResult.Datas = listMap
	return
}

func (this_ *MysqlService) Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error) {
	res, err = this_.DatabaseWorker.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}
