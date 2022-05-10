package db

import (
	"gitee.com/chunanyong/zorm"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DatabaseSqlite Sqlite操作
type DatabaseSqlite struct {
	*DatabaseBase
}

// NewDatabaseSqlite 根据Sqlite配置创建DatabaseSqlite
func NewDatabaseSqlite(config DatabaseConfig) (res *DatabaseSqlite, err error) {
	DatabaseBase := &DatabaseBase{
		config: config,
	}
	res = &DatabaseSqlite{
		DatabaseBase: DatabaseBase,
	}
	err = res.init()
	return
}

func (this_ *DatabaseSqlite) init() (err error) {

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
	zorm.FuncPrintSQL = func(sql string, args []interface{}) {
		//Logger.Info("Exec Sql", zap.Any("sql", sqlstr), zap.Any("args", args))
	}

	DSN := this_.config.Database
	// dbDaoConfig 数据库的配置.这里只是模拟,生产应该是读取配置配置文件,构造DataSourceConfig
	dbDaoConfig := zorm.DataSourceConfig{
		//DSN 数据库的连接字符串
		DSN: DSN,
		//数据库驱动名称:mysql,postgres,oci8,sqlserver,sqlite3,clickhouse,dm,kingbase,aci 和DBType对应,处理数据库有多个驱动
		DriverName: "sqlite3",
		//数据库类型(方言判断依据):mysql,postgresql,oracle,mssql,sqlite,clickhouse,dm,kingbase,shentong 和 DriverName 对应,处理数据库有多个驱动
		DBType: "sqlite",
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

	return

}
