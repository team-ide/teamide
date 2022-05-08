package toolbox

import (
	"database/sql"
	"fmt"
	"gitee.com/chunanyong/zorm"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"teamide/pkg/sql_ddl"
	"teamide/pkg/util"
	"time"

	"context"
	_ "github.com/go-sql-driver/mysql"
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
	dbDao       *zorm.DBDao
	lastUseTime int64
	ctx         context.Context
}

func (this_ *MysqlService) init() (err error) {
	this_.ctx = context.Background()
	//自定义zorm日志输出
	//zorm.LogCallDepth = 4 //日志调用的层级
	zorm.FuncLogError = func(err error) {
		Logger.Error("Zorm Error", zap.Error(err))
	} //记录异常日志的函数
	zorm.FuncLogPanic = func(err error) {
		Logger.Error("Zorm Error", zap.Error(err))
	} //记录panic日志,默认使用defaultLogError实现
	zorm.FuncPrintSQL = func(sqlstr string, args []interface{}) {
		//Logger.Info("Zorm Error", zap.Error(err))
	} //打印sql的函数

	//自定义日志输出格式,把FuncPrintSQL函数重新赋值
	//log.SetFlags(log.LstdFlags)
	//zorm.FuncPrintSQL = zorm.FuncPrintSQL

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true", this_.config.Username, this_.config.Password, this_.config.Host, this_.config.Port)
	//dbDaoConfig 数据库的配置.这里只是模拟,生产应该是读取配置配置文件,构造DataSourceConfig
	dbDaoConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: DSN,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: "mysql",
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		DBType: "mysql",
		//MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: 10,
		//MaxIdleConns 数据库最大空闲连接数 默认50
		MaxIdleConns: 10,
		//ConnMaxLifetimeSecond 连接存活秒时间. 默认600(10分钟)后连接被销毁重建.避免数据库主动断开连接,造成死连接.MySQL默认wait_timeout 28800秒(8小时)
		ConnMaxLifetimeSecond: 600,
		//PrintSQL 打印SQL.会使用FuncPrintSQL记录SQL
		PrintSQL: true,
		//DefaultTxOptions 事务隔离级别的默认配置,默认为nil
		//DefaultTxOptions: nil,
		//如果是使用seata-golang分布式事务,建议使用默认配置
		//DefaultTxOptions: &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false},

		//FuncSeataGlobalTransaction seata-golang分布式的适配函数,返回ISeataGlobalTransaction接口的实现
		//FuncSeataGlobalTransaction : MyFuncSeataGlobalTransaction,
	}

	// 根据dbDaoConfig创建dbDao, 一个数据库只执行一次,第一个执行的数据库为 defaultDao,后续zorm.xxx方法,默认使用的就是defaultDao
	this_.dbDao, err = zorm.NewDBDao(&dbDaoConfig)
	if err != nil {
		return
	}
	this_.ctx, err = this_.dbDao.BindContextDBConnection(this_.ctx)
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
	_ = this_.dbDao.CloseDB()
}

func (this_ *MysqlService) Databases() (databases []*DatabaseInfo, err error) {
	//构造查询用的finder
	finder := zorm.NewSelectFinder("information_schema.SCHEMATA", "SCHEMA_NAME name")

	finder.Append("ORDER BY SCHEMA_NAME")
	//执行查询
	listMap, err := zorm.QueryMap(this_.ctx, finder, nil)
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
	err = zorm.Query(this_.ctx, finder, &tables, nil)
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
	err = zorm.Query(this_.ctx, finder, &tableDetails, nil)
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
	err = zorm.Query(this_.ctx, finder, &columns, nil)
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
	listMap, err := zorm.QueryMap(this_.ctx, finder, nil)
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
	err = zorm.Query(this_.ctx, finder, &indexs_, nil)
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
	listMap, err := zorm.QueryMap(this_.ctx, finder, page)
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

func (this_ *MysqlService) Execs(sqlParams []*SqlParam) (res int, err error) {
	_, err = zorm.Transaction(this_.ctx, func(ctx context.Context) (interface{}, error) {

		for _, sqlParam := range sqlParams {
			finder := zorm.NewFinder()

			finder.Append(sqlParam.Sql, sqlParam.Params...)

			num, err := zorm.UpdateFinder(ctx, finder)
			if err != nil {
				return nil, err
			}
			res += num
		}

		return nil, err
	})
	if err != nil {
		return
	}
	return
}
