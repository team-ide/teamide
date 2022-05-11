package db

import (
	"context"
	"errors"
	"gitee.com/chunanyong/zorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type DatabaseType struct {
	DSNFormat         string
	DriverName        string
	DBType            string
	ColumnTypeInfoMap map[string]*ColumnTypeInfo
}

func (this_ *DatabaseType) GetColumnTypeInfo(name string) (c *ColumnTypeInfo) {
	c = this_.ColumnTypeInfoMap[name]
	if c == nil {
		for key, one := range this_.ColumnTypeInfoMap {
			if strings.ToLower(key) == strings.ToLower(name) {
				c = one
				break
			}
		}
	}
	return
}

var (
	DatabaseTypeMySql    = addDatabaseType(&DatabaseType{DSNFormat: "$username:$password@tcp($host:$port)/$database?charset=utf8mb4&parseTime=true", DriverName: "mysql", DBType: "mysql"})
	DatabaseTypeSqlite   = addDatabaseType(&DatabaseType{DSNFormat: "$database", DriverName: "sqlite3", DBType: "sqlite"})
	DatabaseTypeOracle   = addDatabaseType(&DatabaseType{})
	DatabaseTypeShenTong = addDatabaseType(&DatabaseType{})
	DatabaseTypeDM       = addDatabaseType(&DatabaseType{})
	DatabaseTypeKingBase = addDatabaseType(&DatabaseType{})
	DatabaseTypeKunLun   = addDatabaseType(&DatabaseType{})
	DatabaseTypeGBase    = addDatabaseType(&DatabaseType{})

	DatabaseTypes []*DatabaseType
)

func addDatabaseType(databaseType *DatabaseType) *DatabaseType {
	if databaseType.ColumnTypeInfoMap == nil {
		databaseType.ColumnTypeInfoMap = make(map[string]*ColumnTypeInfo)
	}
	DatabaseTypes = append(DatabaseTypes, databaseType)
	return databaseType
}

func GetDatabaseType(databaseType string) *DatabaseType {
	switch strings.ToLower(databaseType) {
	case "mysql":
		return DatabaseTypeMySql
	case "sqlite":
		return DatabaseTypeSqlite
	case "oracle":
		return DatabaseTypeOracle
	case "shentong":
		return DatabaseTypeShenTong
	case "dm", "dameng":
		return DatabaseTypeDM
	case "kingbase":
		return DatabaseTypeKingBase
	case "kunlun":
		return DatabaseTypeKunLun
	case "gbase":
		return DatabaseTypeGBase
	}
	return nil
}

var (
	Logger *zap.Logger
)

func init() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	Logger, _ = loggerConfig.Build()
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `json:"type,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// NewDatabaseWorker 根据数据库配置创建DatabaseWorker
func NewDatabaseWorker(config DatabaseConfig) (databaseWorker *DatabaseWorker, err error) {
	databaseWorker = &DatabaseWorker{config: config}
	err = databaseWorker.init()
	if err != nil {
		return nil, err
	}
	return databaseWorker, nil
}

// DatabaseWorker 基础操作
type DatabaseWorker struct {
	config       DatabaseConfig
	databaseType *DatabaseType
	dbDao        *zorm.DBDao
	baseContext  context.Context
}

func (this_ *DatabaseWorker) init() (err error) {
	this_.databaseType = GetDatabaseType(this_.config.Type)
	if this_.databaseType == nil {
		err = errors.New("数据库类型[" + this_.config.Type + "]暂不支持")
		return
	}

	dns := this_.databaseType.DSNFormat
	dns = strings.ReplaceAll(dns, `$username`, this_.config.Username)
	dns = strings.ReplaceAll(dns, `$password`, this_.config.Password)
	dns = strings.ReplaceAll(dns, `$host`, this_.config.Host)
	dns = strings.ReplaceAll(dns, `$port`, strconv.Itoa(this_.config.Port))
	dns = strings.ReplaceAll(dns, `$database`, this_.config.Database)
	// 自定义zorm日志输出
	// zorm.LogCallDepth = 4 //日志调用的层级

	// 记录异常日志的函数
	zorm.FuncLogError = func(err error) {
		//Logger.Error("zorm error", zap.Error(err))
	}

	// 记录panic日志,默认使用defaultLogError实现
	zorm.FuncLogPanic = func(err error) {
		Logger.Error("zorm panic error", zap.Error(err))
	}
	// 打印sql的函数
	zorm.FuncPrintSQL = func(sqlstr string, args []interface{}) {
		//Logger.Info("Exec Sql", zap.Any("sql", sqlstr), zap.Any("args", args))
	}
	//fmt.Println("dns:" + dns)
	// dbDaoConfig 数据库的配置.这里只是模拟,生产应该是读取配置配置文件,构造DataSourceConfig
	dbDaoConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: dns,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: this_.databaseType.DriverName,
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		DBType: this_.databaseType.DBType,
		//MaxOpenConns 数据库最大连接数 默认50
		MaxOpenConns: 100,
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
	this_.baseContext = context.Background()
	return
}

func (this_ *DatabaseWorker) GetContext() (ctx context.Context) {
	ctx, _ = this_.dbDao.BindContextDBConnection(this_.baseContext)
	return ctx
}

func (this_ *DatabaseWorker) GetConfig() (config DatabaseConfig) {
	config = this_.config
	return
}

func (this_ *DatabaseWorker) Open() (err error) {
	finder := zorm.NewFinder()
	finder.Append("SELECT 1")

	var count int
	_, err = zorm.QueryRow(this_.GetContext(), finder, &count)
	return
}

func (this_ *DatabaseWorker) Close() (err error) {
	err = this_.dbDao.CloseDB()
	return
}

func (this_ *DatabaseWorker) Exec(sql string, params []interface{}) (rowsAffected int64, err error) {

	rowsAffected, err = this_.Execs([]string{sql}, [][]interface{}{params})

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseWorker) Execs(sqlList []string, paramsList [][]interface{}) (rowsAffected int64, err error) {

	_, err = zorm.Transaction(this_.GetContext(), func(ctx context.Context) (interface{}, error) {

		var err error
		for index, _ := range sqlList {
			sql := sqlList[index]
			params := paramsList[index]
			finder := zorm.NewFinder()
			finder.Append(sql, params...)

			var updated int
			updated, err = zorm.UpdateFinder(ctx, finder)
			if err != nil {
				Logger.Error("Exec Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
				return nil, err
			}
			rowsAffected += int64(updated)
		}

		//如果返回的err不是nil,事务就会回滚
		return nil, err
	})

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseWorker) Count(sql string, params []interface{}) (count int64, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	count, err = this_.FinderCount(finder)

	if err != nil {
		Logger.Error("Count Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}
	return
}

func (this_ *DatabaseWorker) FinderCount(finder *zorm.Finder) (count int64, err error) {

	_, err = zorm.QueryRow(this_.GetContext(), finder, &count)

	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseWorker) QueryOne(sql string, params []interface{}, one interface{}) (find bool, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	find, err = this_.FinderQueryOne(finder, one)

	if err != nil {
		Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) FinderQueryOne(finder *zorm.Finder, one interface{}) (find bool, err error) {

	find, err = zorm.QueryRow(this_.GetContext(), finder, one)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseWorker) Query(sql string, params []interface{}, list interface{}) (err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	err = this_.FinderQuery(finder, list)

	if err != nil {
		Logger.Error("Query Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) FinderQuery(finder *zorm.Finder, list interface{}) (err error) {

	err = zorm.Query(this_.GetContext(), finder, list, nil)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryMap(sql string, params []interface{}) (list []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	list, err = this_.FinderQueryMap(finder)

	if err != nil {
		Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) FinderQueryMap(finder *zorm.Finder) (list []map[string]interface{}, err error) {

	list, err = zorm.QueryMap(this_.GetContext(), finder, nil)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryPage(sql string, params []interface{}, list interface{}, page *zorm.Page) (err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	err = this_.FinderQueryPage(finder, list, page)

	if err != nil {
		Logger.Error("QueryPage Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) FinderQueryPage(finder *zorm.Finder, list interface{}, page *zorm.Page) (err error) {

	err = zorm.Query(this_.GetContext(), finder, list, page)

	if err != nil {
		return
	}

	return
}

func (this_ *DatabaseWorker) QueryMapPage(sql string, params []interface{}, page *zorm.Page) (list []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.Append(sql, params...)

	list, err = this_.FinderQueryMapPage(finder, page)

	if err != nil {
		Logger.Error("QueryMap Error", zap.Any("sql", sql), zap.Any("params", params), zap.Error(err))
		return
	}

	return
}

func (this_ *DatabaseWorker) FinderQueryMapPage(finder *zorm.Finder, page *zorm.Page) (list []map[string]interface{}, err error) {

	list, err = zorm.QueryMap(this_.GetContext(), finder, page)

	if err != nil {
		return
	}

	return
}
