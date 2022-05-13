// Package zorm 使用原生的sql语句,没有对sql语法做限制.语句使用Finder作为载体
// 占位符统一使用?,zorm会根据数据库类型,语句执行前会自动替换占位符,postgresql 把?替换成$1,$2...;mssql替换成@P1,@p2...;orace替换成:1,:2...
// zorm使用 ctx context.Context 参数实现事务传播,ctx从web层传递进来即可,例如gin的c.Request.Context()
// zorm的事务操作需要显示使用zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {})开启
// "package zorm" Use native SQL statements, no restrictions on SQL syntax. Statements use Finder as a carrier
// Use placeholders uniformly "?" "zorm" automatically replaces placeholders before statements are executed,depending on the database type. Replaced with $1, $2... ; Replace MSSQL with @p1,@p2... ; Orace is replaced by :1,:2...,
// "zorm" uses the "ctx context.Context" parameter to achieve transaction propagation,and ctx can be passed in from the web layer, such as "gin's c.Request.Context()",
// "zorm" Transaction operations need to be displayed using "zorm.transaction" (ctx, func(ctx context.context) (interface{}, error) {})
package zorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FuncReadWriteStrategy 单个数据库的读写分离的策略,用于外部复写实现自定义的逻辑,rwType=0 read,rwType=1 write
// 不能归属到DBDao里,BindContextDBConnection已经是指定数据库的连接了,和这个函数会冲突.就作为单数据库读写分离的处理方式
// 即便是放到DBDao里,因为是多库,BindContextDBConnection函数调用少不了,业务包装一个方法,指定一下读写获取一个DBDao效果是一样的,唯一就是需要根据业务指定一下读写,其实更灵活了
// FuncReadWriteStrategy Single database read and write separation strategy,used for external replication to implement custom logic, rwType=0 read, rwType=1 write.
// "BindContextDBConnection" is already a connection to the specified database and will conflict with this function. As a single database read and write separation of processing
var FuncReadWriteStrategy func(rwType int) *DBDao = getDefaultDao

// wrapContextStringKey 包装context的key,不直接使用string类型,避免外部直接注入使用
type wrapContextStringKey string

// context WithValue的key,不能是基础类型,例如字符串,包装一下
// The key of context WithValue cannot be a basic type, such as a string, wrap it
const contextDBConnectionValueKey = wrapContextStringKey("contextDBConnectionValueKey")

//事务选项设置TxOptions,主要是设置事务的隔离级别
const contextTxOptionsKey = wrapContextStringKey("contextTxOptionsKey")

// NewContextDBConnectionValueKey 创建context中存放DBConnection的key
// 故意使用一个公开方法,返回私有类型wrapContextStringKey,多库时禁止自定义contextKey,只能调用这个方法,不能接收也不能改变
// 例如:ctx = context.WithValue(ctx, zorm.NewContextDBConnectionValueKey(), dbConnection)
// func NewContextDBConnectionValueKey() wrapContextStringKey {
//   return contextDBConnectionValueKey
// }

//bug(springrain) 还缺少1对1的属性嵌套对象,sql别名查询,直接赋值的功能.

//不再处理日期零值,会干扰反射判断零值
//默认的零时时间1970-01-01 00:00:00 +0000 UTC,兼容数据库,避免0001-01-01 00:00:00 +0000 UTC的零值.数据库不让存值,加上1秒,跪了
//因为mysql 5.7后,The TIMESTAMP data type is used for values that contain both date and time parts. TIMESTAMP has a range of '1970-01-01 00:00:01' UTC to '2038-01-19 03:14:07' UTC.
//var defaultZeroTime = time.Date(1970, time.January, 1, 0, 0, 1, 0, time.UTC)

//var defaultZeroTime = time.Now()

//注释如果是 . 句号结尾,IDE的提示就截止了,注释结尾不要用 . 结束

//DBDao 数据库操作基类,隔离原生操作数据库API入口,所有数据库操作必须通过DBDao进行
//DBDao Database operation base class, isolate the native operation database API entry,all database operations must be performed through DB Dao
type DBDao struct {
	config     *DataSourceConfig
	dataSource *dataSource
}

var defaultDao *DBDao = nil

// NewDBDao 创建dbDao,一个数据库要只执行一次,业务自行控制
// 第一个执行的数据库为 defaultDao,后续zorm.xxx方法,默认使用的就是defaultDao
// NewDBDao Creates dbDao, a database must be executed only once, and the business is controlled by itself
// The first database to be executed is defaultDao, and the subsequent zorm.xxx method is defaultDao by default
func NewDBDao(config *DataSourceConfig) (*DBDao, error) {
	dataSource, err := newDataSource(config)

	if err != nil {
		err = fmt.Errorf("NewDBDao创建dataSource失败:%w", err)
		FuncLogError(err)
		return nil, err
	}

	if FuncReadWriteStrategy(1) == nil {
		defaultDao = &DBDao{config, dataSource}
		return defaultDao, nil
	}
	return &DBDao{config, dataSource}, nil
}

// getDefaultDao 获取默认的Dao,用于隔离读写的Dao
// getDefaultDao Get the default Dao, used to isolate Dao for reading and writing
func getDefaultDao(rwType int) *DBDao {
	return defaultDao
}

// newDBConnection 获取一个dbConnection
// 如果参数dbConnection为nil,使用默认的datasource进行获取dbConnection
// 如果是多库,Dao手动调用newDBConnection(),获得dbConnection,WithValue绑定到子context
// newDBConnection Get a db Connection
// If the parameter db Connection is nil, use the default datasource to get db Connection.
// If it is multi-database, Dao manually calls new DB Connection() to obtain db Connection, and With Value is bound to the sub-context
func (dbDao *DBDao) newDBConnection() (*dataBaseConnection, error) {
	if dbDao == nil || dbDao.dataSource == nil {
		return nil, errors.New("请不要自己创建dbDao,使用NewDBDao方法进行创建")
	}
	dbConnection := new(dataBaseConnection)
	dbConnection.db = dbDao.dataSource.DB
	dbConnection.config = dbDao.config
	return dbConnection, nil
}

// BindContextDBConnection 多库的时候,通过dbDao创建DBConnection绑定到子context,返回的context就有了DBConnection. parent 不能为空
// BindContextDBConnection In the case of multiple databases, create a DB Connection through db Dao and bind it to a sub-context,and the returned context will have a DB Connection. parent is not nil
func (dbDao *DBDao) BindContextDBConnection(parent context.Context) (context.Context, error) {
	if parent == nil {
		return nil, errors.New("BindContextDBConnection context的parent不能为nil")
	}
	dbConnection, errDBConnection := dbDao.newDBConnection()
	if errDBConnection != nil {
		return parent, errDBConnection
	}
	ctx := context.WithValue(parent, contextDBConnectionValueKey, dbConnection)
	return ctx, nil
}

// BindContextTxOptions 绑定事务的隔离级别,参考sql.IsolationLevel,如果txOptions为nil,使用默认的事务隔离级别.parent不能为空
//需要在事务开启前调用,也就是zorm.Transaction方法前,不然事务开启之后再调用就无效了
func (dbDao *DBDao) BindContextTxOptions(parent context.Context, txOptions *sql.TxOptions) (context.Context, error) {
	if parent == nil {
		return nil, errors.New("BindContextTxOptions context的parent不能为nil")
	}

	ctx := context.WithValue(parent, contextTxOptionsKey, txOptions)
	return ctx, nil
}

// CloseDB 关闭所有数据库连接
//请谨慎调用这个方法,会关闭所有数据库连接,用于处理特殊场景,正常使用无需手动关闭数据库连接
func (dbDao *DBDao) CloseDB() error {
	if dbDao == nil || dbDao.dataSource == nil {
		return errors.New("请不要自己创建dbDao,使用NewDBDao方法进行创建")
	}
	return dbDao.dataSource.Close()
}

/*
Transaction 的示例代码
  //匿名函数return的error如果不为nil,事务就会回滚
  zorm.Transaction(ctx context.Context,func(ctx context.Context) (interface{}, error) {

	  //业务代码


	  //return的error如果不为nil,事务就会回滚
      return nil, nil
  })
*/
// 事务方法,隔离dbConnection相关的API.必须通过这个方法进行事务处理,统一事务方式
// 如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
// 如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
// 如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
// 但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
// 在多库的场景,手动获取dbConnection,然后绑定到一个新的context,传入进来
// 不要去掉匿名函数的context参数,因为如果Transaction的context中没有dbConnection,会新建一个context并放入dbConnection,此时的context指针已经变化,不能直接使用Transaction的context参数
// bug(springrain)如果有大神修改了匿名函数内的参数名,例如改为ctx2,这样业务代码实际使用的是Transaction的context参数,如果为没有dbConnection,会抛异常,如果有dbConnection,实际就是一个对象.影响有限.也可以把匿名函数抽到外部
// 如果全局DefaultTxOptions配置不满足需求,可以在zorm.Transaction事务方法前设置事务的隔离级别,例如 ctx, _ := dbDao.BindContextTxOptions(ctx, &sql.TxOptions{Isolation: sql.LevelDefault}),如果txOptions为nil,使用全局DefaultTxOptions
// return的error如果不为nil,事务就会回滚
// 如果使用了分布式事务,需要设置分布式事务函数zorm.DataSourceConfig.FuncSeataGlobalTransaction,实现ISeataGlobalTransaction接口
// 如果是分布式事务开启方,需要在本地事务前开启分布事务,开启之后获取XID,设值到ctx的XID和TX_XID.XID是seata MySQL驱动需要,TX_XID是seataContext.NewRootContext需要
// 分布式事务需要传递XID,接收方context.WithValue(ctx, "XID", XID)绑定到ctx
// 如果分支事务出现异常或者回滚,会立即回滚分布式事务
// Transaction method, isolate db Connection related API. This method must be used for transaction processing and unified transaction mode
// If there is no db Connection in the input ctx, use default Dao to start the transaction and submit it finally
// If the input ctx has db Connection and no transaction, call db Connection.begin() to start the transaction and finally commit
// If the input ctx has a db Connection and a transaction, only use non-commit, and the open party submits the transaction
// If you encounter an error or exception, although it is not the initiator of the transaction, the transaction will be rolled back,
// so that the transaction can be rolled back as soon as possible
// In a multi-database scenario, manually obtain db Connection, then bind it to a new context and pass in
// Do not drop the anonymous function's context parameter, because if the Transaction context does not have a DBConnection,
// then a new context will be created and placed in the DBConnection
// The context pointer has changed and the Transaction context parameters cannot be used directly
// "bug (springrain)" If a great god changes the parameter name in the anonymous function, for example, change it to ctx 2,
// so that the business code actually uses the context parameter of Transaction. If there is no db Connection,
// an exception will be thrown. If there is a db Connection, the actual It is an object
// The impact is limited. Anonymous functions can also be extracted outside
// If the return error is not nil, the transaction will be rolled back
func Transaction(ctx context.Context, doTransaction func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	//是否是dbConnection的开启方,如果是开启方,才可以提交事务
	// Whether it is the opener of db Connection, if it is the opener, the transaction can be submitted
	localTxOpen := false
	//是否是seata全局事务的开启方.如果ctx中没有xid,认为是开启方
	seataTxOpen := false
	//如果dbConnection不存在,则会用默认的datasource开启事务
	// If db Connection does not exist, the default datasource will be used to start the transaction
	var checkerr error
	var dbConnection *dataBaseConnection
	ctx, dbConnection, checkerr = checkDBConnection(ctx, false, 1)
	if checkerr != nil {
		return nil, checkerr
	}

	//适配seata分布式事务的函数
	funcSeataTx := dbConnection.config.FuncSeataGlobalTransaction
	//实现ISeataGlobalTransaction接口的事务对象
	var seataGlobalTransaction ISeataGlobalTransaction
	//seata分布式事务的 rootContext,和业务的ctx区别开来,如果业务ctx使用WithValue,就会出现差异
	var seataRootContext context.Context
	//seata分布式事务的异常
	var seataErr error

	//如果没有事务,并且事务没有被禁用,开启事务
	//开启本地事务前,需要拿到分布式事务对象
	//if dbConnection.tx == nil && (!dbConnection.config.DisableTransaction) {
	if dbConnection.tx == nil {

		//需要开启分布式事务,初始化分布式事务对象,判断是否是分布式事务入口
		if funcSeataTx != nil {
			//获取Seata XID
			ctxXIDval := ctx.Value("XID")
			if ctxXIDval != nil { //如果本地ctx中有XID
				seataXID, _ := ctxXIDval.(string)
				//不知道为什么需要两个Key,还需要请教seata-golang团队
				//Seata mysql驱动需要 XID,seataContext.NewRootContext 需要 TX_XID
				ctx = context.WithValue(ctx, "TX_XID", seataXID)
			} else { //如果本地ctx中没有XID,也就是没有传递过来XID,认为是分布式事务的开启方.ctx中没有XID和TX_XID的值
				seataTxOpen = true
			}
			//获取seata事务实现对象,用于控制事务提交和回滚.分支事务需要ctx中TX_XID有值,将分支事务关联到主事务
			seataGlobalTransaction, seataRootContext, seataErr = funcSeataTx(ctx)
			if seataErr != nil {
				seataErr = fmt.Errorf("seata:Transaction FuncSeataGlobalTransaction获取ISeataGlobalTransaction接口实现失败:%w ", seataErr)
				FuncLogError(seataErr)
				return nil, seataErr
			}
			if seataGlobalTransaction == nil || seataRootContext == nil {
				seataErr = errors.New("seata:Transaction FuncSeataGlobalTransaction获取ISeataGlobalTransaction接口的实现为nil ")
				FuncLogError(seataErr)
				return nil, seataErr
			}

		}
		if seataTxOpen { //如果是分布事务开启方,启动分布式事务
			seataErr = seataGlobalTransaction.SeataBegin(seataRootContext)
			if seataErr != nil {
				seataErr = fmt.Errorf("seata:Transaction 分布式事务开启失败:%w ", seataErr)
				FuncLogError(seataErr)
				return nil, seataErr
			}

			//分布式事务开启成功,获取XID,设置到ctx的XID和TX_XID
			//Seata mysql驱动需要 XID,seataContext.NewRootContext 需要 TX_XID
			seataXID := seataGlobalTransaction.GetSeataXID(seataRootContext)
			if len(seataXID) < 1 {
				seataErr = errors.New("seata:seataGlobalTransaction.SeataBegin无异常开启后,获取的XID为空")
				FuncLogError(seataErr)
				return nil, seataErr
			}
			ctx = context.WithValue(ctx, "XID", seataXID)
			ctx = context.WithValue(ctx, "TX_XID", seataXID)
		}

		//开启本地事务/分支事务
		beginerr := dbConnection.beginTx(ctx)
		if beginerr != nil {
			beginerr = fmt.Errorf("Transaction 事务开启失败:%w ", beginerr)
			FuncLogError(beginerr)
			return nil, beginerr
		}
		//本方法开启的事务,由本方法提交
		//The transaction opened by this method is submitted by this method
		localTxOpen = true
	}

	defer func() {
		if r := recover(); r != nil {
			//err = fmt.Errorf("事务开启失败:%w ", err)
			//记录异常日志
			//if _, ok := r.(runtime.Error); ok {
			//	panic(r)
			//}
			err, errOk := r.(error)
			if errOk {
				err = fmt.Errorf("recover异常:%w", err)
				FuncLogPanic(err)
			}
			//if !txOpen { //如果不是开启方,也应该回滚事务,虽然可能造成日志不准确,但是回滚要尽早
			//	return
			//}
			//如果全局禁用了事务
			//if dbConnection.config.DisableTransaction {
			//	return
			//}
			rberr := dbConnection.rollback()
			if rberr != nil {
				rberr = fmt.Errorf("recover内事务回滚失败:%w", rberr)
				FuncLogError(rberr)
			}
			//任意一个分支事务回滚,分布式事务就整体回滚
			if seataGlobalTransaction != nil {
				seataErr = seataGlobalTransaction.SeataRollback(seataRootContext)
				if seataErr != nil {
					seataErr = fmt.Errorf("seata:recover内seataGlobalTransaction事务回滚失败:%w", seataErr)
					FuncLogError(seataErr)
				}
			}

		}
	}()

	//执行业务的事务函数
	info, err := doTransaction(ctx)

	if err != nil {
		err = fmt.Errorf("Transaction-->doTransaction业务执行异常:%w", err)
		FuncLogError(err)

		//如果全局禁用了事务
		//if dbConnection.config.DisableTransaction {
		//	return info, err
		//}

		//不是开启方回滚事务,有可能造成日志记录不准确,但是回滚最重要了,尽早回滚
		//It is not the start party to roll back the transaction, which may cause inaccurate log records,but rollback is the most important, roll back as soon as possible
		rberr := dbConnection.rollback()
		if rberr != nil {
			rberr = fmt.Errorf("Transaction-->rollback事务回滚失败:%w", rberr)
			FuncLogError(rberr)
		}
		//任意一个分支事务回滚,分布式事务就整体回滚
		if seataGlobalTransaction != nil {
			seataErr = seataGlobalTransaction.SeataRollback(seataRootContext)
			if seataErr != nil {
				seataErr = fmt.Errorf("seata:Transaction-->rollback seataGlobalTransaction事务回滚失败:%w", seataErr)
				FuncLogError(seataErr)
			}
		}
		return info, err
	}
	//如果是事务开启方,提交事务
	//If it is the transaction opener, commit the transaction
	if localTxOpen {
		commitError := dbConnection.commit()
		//本地事务提交成功,如果是全局事务的开启方,提交分布式事务
		if commitError == nil && seataTxOpen {
			seataErr = seataGlobalTransaction.SeataCommit(seataRootContext)
			if seataErr != nil {
				seataErr = fmt.Errorf("seata:Transaction-->commit seataGlobalTransaction 事务提交失败:%w", seataErr)
				FuncLogError(seataErr)
			}
		}
		if commitError != nil {
			commitError = fmt.Errorf("Transaction-->commit事务提交失败:%w", commitError)
			FuncLogError(commitError)

			//任意一个分支事务回滚,分布式事务就整体回滚
			if seataGlobalTransaction != nil {
				seataErr = seataGlobalTransaction.SeataRollback(seataRootContext)
				if seataErr != nil {
					seataErr = fmt.Errorf("seata:Transaction-->commit失败,然后回滚seataGlobalTransaction事务也失败:%w", seataErr)
					FuncLogError(seataErr)
				}
			}

			return info, commitError
		}

	}

	return nil, nil
}

// QueryRow 不要偷懒调用Query返回第一条,问题1.需要构建一个slice,问题2.调用方传递的对象其他值会被抛弃或者覆盖.
// 根据Finder和封装为指定的entity类型,entity必须是*struct类型或者基础类型的指针.把查询的数据赋值给entity,所以要求指针类型
// context必须传入,不能为空
// 如果数据库是null,基本类型不支持,会返回异常,不做默认值处理,Query因为是列表,会设置为默认值
// QueryRow Don't be lazy to call Query to return the first one
// Question 1. A selice needs to be constructed, and question 2. Other values ​​of the object passed by the caller will be discarded or overwritten
// context must be passed in and cannot be empty
func QueryRow(ctx context.Context, finder *Finder, entity interface{}) (bool, error) {

	has := false
	typeOf, checkerr := checkEntityKind(entity)
	if checkerr != nil {
		checkerr = fmt.Errorf("QueryRow-->checkEntityKind类型检查错误:%w", checkerr)
		FuncLogError(checkerr)
		return has, checkerr
	}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return has, errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return has, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//dbConnection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(0).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//获取到sql语句
	//Get the sql statement
	sqlstr, err := wrapQuerySQL(dbType, finder, nil)
	if err != nil {
		err = fmt.Errorf("QueryRow-->wrapQuerySQL获取查询SQL语句错误:%w", err)
		FuncLogError(err)
		return has, err
	}

	//检查dbConnection.有可能会创建dbConnection或者开启事务,所以要尽可能的接近执行时检查
	//Check db Connection. It is possible to create a db Connection or start a transaction, so check it as close as possible to the execution
	var dbConnectionerr error
	ctx, dbConnection, dbConnectionerr = checkDBConnection(ctx, false, 0)
	if dbConnectionerr != nil {
		return has, dbConnectionerr
	}

	//根据语句和参数查询
	//Query based on statements and parameters
	rows, e := dbConnection.queryContext(ctx, &sqlstr, finder.values)
	if e != nil {
		e = fmt.Errorf("QueryRow-->queryContext查询数据库错误:%w", e)
		FuncLogError(e)
		return has, e
	}
	//先判断error 再关闭
	defer rows.Close()

	//typeOf := reflect.TypeOf(entity).Elem()

	//数据库返回的列名
	//Column name returned by the database
	/*
		columns, cne := rows.Columns()
		if cne != nil {
			cne = fmt.Errorf("QueryRow-->rows.Columns数据库返回列名错误:%w", cne)
			FuncLogError(cne)
			return cne
		}
	*/
	//数据库字段类型
	columnTypes, cte := rows.ColumnTypes()
	if cte != nil {
		cte = fmt.Errorf("QueryRow-->rows.ColumnTypes数据库类型错误:%w", cte)
		FuncLogError(cte)
		return has, cte
	}

	//反射获取 []driver.Value的值,用于处理nil值和自定义类型
	var driverValue = reflect.Indirect(reflect.ValueOf(rows))
	driverValue = driverValue.FieldByName("lastcols")

	cdvMapHasBool := len(CustomDriverValueMap) > 0
	//就查询一个字段
	//If it is a basic type, query a field
	//if allowBaseTypeMap[typeOf.Kind()] && len(columns) == 1 {
	if len(columnTypes) == 1 {
		//类型转换的接口实现
		var converFunc CustomDriverValueConver
		//是否有需要的类型转换
		var converOK bool = false
		//类型转换的临时值
		var tempDriverValue driver.Value
		//循环遍历结果集
		for i := 0; rows.Next(); i++ {
			has = true
			if i > 0 {
				return has, errors.New("QueryRow查询出多条数据")
			}

			dv := driverValue.Index(0)
			if dv.IsValid() && dv.InterfaceData()[0] == 0 { // 该字段的数据库值是null,不再处理,使用默认值
				return has, nil
			}
			var scanerr error
			//判断是否有自定义扩展,避免无意义的反射
			if cdvMapHasBool {
				//根据接收的类型,获取到类型转换的接口实现
				converFunc, converOK = CustomDriverValueMap[dv.Elem().Type().String()]
			}

			var errGetDriverValue error
			if converOK { //如果有类型需要转换
				//获取需要转换的临时值
				tempDriverValue, errGetDriverValue = converFunc.GetDriverValue(columnTypes[0], &typeOf, finder)
				if errGetDriverValue != nil {
					errGetDriverValue = fmt.Errorf("QueryRow-->conver.GetDriverValue异常:%w", errGetDriverValue)
					FuncLogError(errGetDriverValue)
					return has, errGetDriverValue
				}

				//返回值为nil,不做任何处理
				if tempDriverValue == nil {
					scanerr = rows.Scan(entity)
				} else { //如果有值,需要类型转换
					scanerr = rows.Scan(tempDriverValue)
				}

			} else { //如果不需要类型转换
				scanerr = rows.Scan(entity)
			}

			if scanerr != nil {
				scanerr = fmt.Errorf("QueryRow-->rows.Scan异常:%w", scanerr)
				FuncLogError(scanerr)
				return has, scanerr
			}
		}

		//如果需要类型转换,需要把临时值转换成需要接收的类型值
		if converOK && tempDriverValue != nil {
			//根据接收的临时值,返回需要接收值的指针
			rightValue, errConverDriverValue := converFunc.ConverDriverValue(columnTypes[0], &typeOf, tempDriverValue, finder)
			if errConverDriverValue != nil {
				errConverDriverValue = fmt.Errorf("QueryRow-->converFunc.ConverDriverValue异常:%w", errConverDriverValue)
				FuncLogError(errConverDriverValue)
				return has, errConverDriverValue
			}

			//把返回的值复制给接收的对象,
			reflect.ValueOf(entity).Elem().Set(reflect.ValueOf(rightValue).Elem())

		}
		return has, nil
		//只查询一个字段的逻辑结束
	}

	//查询多个字段的逻辑开始

	//获取接收值的对象
	valueOf := reflect.ValueOf(entity).Elem()
	//获取到类型的字段缓存
	//Get the type field cache
	dbColumnFieldMap, exportFieldMap, dbe := getDBColumnExportFieldMap(&typeOf)
	if dbe != nil {
		dbe = fmt.Errorf("QueryRow-->getDBColumnFieldMap获取字段缓存错误:%w", dbe)
		FuncLogError(dbe)
		return has, dbe
	}

	//反射获取 []driver.Value的值
	//driverValue = reflect.Indirect(reflect.ValueOf(rows))
	//driverValue = driverValue.FieldByName("lastcols")

	//循环遍历结果集
	//Loop through the result set
	for i := 0; rows.Next(); i++ {
		has = true
		if i > 0 {
			return has, errors.New("QueryRow查询出多条数据")
		}

		//接收对象设置值
		scanerr := sqlRowsValues(rows, &driverValue, columnTypes, dbColumnFieldMap, exportFieldMap, &valueOf, finder, cdvMapHasBool)
		if scanerr != nil {
			scanerr = fmt.Errorf("QueryRow-->sqlRowsValues错误:%w", scanerr)
			FuncLogError(scanerr)
			return has, scanerr
		}

	}

	return has, nil
}

// Query 不要偷懒调用QueryMap,需要处理sql驱动支持的sql.Nullxxx的数据类型,也挺麻烦的
// 根据Finder和封装为指定的entity类型,entity必须是*[]struct类型,已经初始化好的数组,此方法只Append元素,这样调用方就不需要强制类型转换了
// context必须传入,不能为空.如果想不分页,查询所有数据,page传入nil
// Query:Don't be lazy to call QueryMap, you need to deal with the sql,Nullxxx data type supported by the sql driver, which is also very troublesome
// According to the Finder and encapsulation for the specified entity type, the entity must be of the *[]struct type, which has been initialized,This method only Append elements, so the caller does not need to force type conversion
// context must be passed in and cannot be empty
func Query(ctx context.Context, finder *Finder, rowsSlicePtr interface{}, page *Page) error {

	if rowsSlicePtr == nil { //如果为nil
		return errors.New("Query数组必须是*[]struct类型或者*[]*struct或者基础类型数组的指针")
	}

	pv1 := reflect.ValueOf(rowsSlicePtr)
	if pv1.Kind() != reflect.Ptr { //如果不是指针
		return errors.New("Query数组必须是*[]struct类型或者*[]*struct或者基础类型数组的指针")
	}

	//获取数组元素
	//Get array elements
	sliceValue := reflect.Indirect(pv1)

	//如果不是数组
	//If it is not an array.
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("Query数组必须是*[]struct类型或者*[]*struct或者基础类型数组的指针")
	}
	//获取数组内的元素类型
	//Get the element type in the array
	sliceElementType := sliceValue.Type().Elem()
	//slice数组里是否是指针,实际参数类似 *[]*struct,兼容这种类型
	sliceElementTypePtr := false
	//如果数组里还是指针类型
	if sliceElementType.Kind() == reflect.Ptr {
		sliceElementTypePtr = true
		sliceElementType = sliceElementType.Elem()
	}

	//如果不是struct
	//if !(sliceElementType.Kind() == reflect.Struct || allowBaseTypeMap[sliceElementType.Kind()]) {
	//	return errors.New("Query数组必须是*[]struct类型或者*[]*struct或者基础类型数组的指针")
	//}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return errDBConnection
	}

	var dbType string = ""
	if dbConnection == nil { //dbConnection为nil,使用defaultDao
		dbType = FuncReadWriteStrategy(0).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	sqlstr, err := wrapQuerySQL(dbType, finder, page)
	if err != nil {
		err = fmt.Errorf("Query-->wrapQuerySQL获取查询SQL语句错误:%w", err)
		FuncLogError(err)
		return err
	}

	//检查dbConnection.有可能会创建dbConnection或者开启事务,所以要尽可能的接近执行时检查
	//Check db Connection. It is possible to create a db Connection or start a transaction, so check it as close as possible to the execution
	var dbConnectionerr error
	ctx, dbConnection, dbConnectionerr = checkDBConnection(ctx, false, 0)
	if dbConnectionerr != nil {
		return dbConnectionerr
	}

	//根据语句和参数查询
	//Query based on statements and parameters
	rows, e := dbConnection.queryContext(ctx, &sqlstr, finder.values)
	if e != nil {
		e = fmt.Errorf("Query-->queryContext查询rows异常:%w", e)
		FuncLogError(e)
		return e
	}
	//先判断error 再关闭
	defer rows.Close()

	//数据库返回的字段类型
	columnTypes, cte := rows.ColumnTypes()
	if cte != nil {
		cte = fmt.Errorf("Query-->rows.ColumnTypes数据库类型错误:%w", cte)
		FuncLogError(cte)
		return cte
	}
	//反射获取 []driver.Value的值
	driverValue := reflect.Indirect(reflect.ValueOf(rows))
	driverValue = driverValue.FieldByName("lastcols")

	cdvMapHasBool := len(CustomDriverValueMap) > 0

	//如果是基础类型,就查询一个字段
	//If it is a basic type, query a field
	//if allowBaseTypeMap[sliceElementType.Kind()] {
	if len(columnTypes) == 1 {

		//循环遍历结果集
		//Loop through the result set
		for rows.Next() {

			//初始化一个基本类型,new出来的是指针
			//Initialize a basic type, and new is a pointer
			pv := reflect.New(sliceElementType)
			//列表查询单个字段要处理数据库为null的情况,如果是Query,会有错误异常,不需要处理null
			dv := driverValue.Index(0)
			if dv.IsValid() && dv.InterfaceData()[0] == 0 { // 该字段的数据库值是null,取默认值
				if sliceElementTypePtr { //如果数组里是指针地址,*[]*struct
					sliceValue.Set(reflect.Append(sliceValue, pv))
				} else {
					sliceValue.Set(reflect.Append(sliceValue, pv.Elem()))
				}
				continue
			}
			//类型转换的接口实现
			var converFunc CustomDriverValueConver
			//是否需要类型转换
			var converOK bool = false

			//根据接收的类型,获取到类型转换的接口实现
			if cdvMapHasBool {
				converFunc, converOK = CustomDriverValueMap[dv.Elem().Type().String()]
			}

			//类型转换的临时值
			var tempDriverValue driver.Value
			var errGetDriverValue error
			//如果需要转换
			if converOK {
				//获取需要转的临时值
				tempDriverValue, errGetDriverValue = converFunc.GetDriverValue(columnTypes[0], &sliceElementType, finder)
				if errGetDriverValue != nil {
					errGetDriverValue = fmt.Errorf("Query-->conver.GetDriverValue异常:%w", errGetDriverValue)
					FuncLogError(errGetDriverValue)
					return errGetDriverValue
				}
				if tempDriverValue != nil { //为nil,不做处理
					pv = reflect.ValueOf(tempDriverValue)
				}

			}

			//把数据库值赋给指针
			//Assign database value to pointer
			scanerr := rows.Scan(pv.Interface())

			if scanerr != nil {
				scanerr = fmt.Errorf("Query-->rows.Scan异常:%w", scanerr)
				FuncLogError(scanerr)
				return scanerr
			}
			//如果需要类型转换,需要把临时值转换成需要接收的类型值
			if converOK && tempDriverValue != nil {
				//根据接收的临时值,返回需要接收值的指针
				rightValue, errConverDriverValue := converFunc.ConverDriverValue(columnTypes[0], &sliceElementType, tempDriverValue, finder)
				if errConverDriverValue != nil {
					errConverDriverValue = fmt.Errorf("Query-->conver.ConverDriverValue异常:%w", errConverDriverValue)
					FuncLogError(errConverDriverValue)
					return errConverDriverValue
				}
				//把正确的值赋值给pv
				pv = reflect.ValueOf(rightValue)
			}

			//通过反射给slice添加元素.添加指针下的真实元素
			//Add elements to slice through reflection. Add real elements under the pointer
			if sliceElementTypePtr { //如果数组里是指针地址,*[]*struct
				sliceValue.Set(reflect.Append(sliceValue, pv))
			} else {
				sliceValue.Set(reflect.Append(sliceValue, pv.Elem()))
			}

		}

		//查询总条数
		//Query total number
		if page != nil && finder.SelectTotalCount {
			count, counterr := selectCount(ctx, finder)
			if counterr != nil {
				counterr = fmt.Errorf("Query-->selectCount查询总条数错误:%w", counterr)
				FuncLogError(counterr)
				return counterr
			}
			page.setTotalCount(count)
		}
		return nil
		//只查询一个字段的逻辑结束
	}

	//获取到类型的字段缓存
	//Get the type field cache
	dbColumnFieldMap, exportFieldMap, dbe := getDBColumnExportFieldMap(&sliceElementType)
	if dbe != nil {
		dbe = fmt.Errorf("Query-->getDBColumnFieldMap获取字段缓存错误:%w", dbe)
		FuncLogError(dbe)
		return dbe
	}

	//循环遍历结果集
	//Loop through the result set
	for rows.Next() {
		//deepCopy(a, entity)
		//反射初始化一个数组内的元素
		//new出来的是指针
		//Reflectively initialize the elements in an array
		pv := reflect.New(sliceElementType).Elem()
		//设置接收值
		scanerr := sqlRowsValues(rows, &driverValue, columnTypes, dbColumnFieldMap, exportFieldMap, &pv, finder, cdvMapHasBool)
		//scan赋值.是一个指针数组,已经根据struct的属性类型初始化了,sql驱动能感知到参数类型,所以可以直接赋值给struct的指针.这样struct的属性就有值了
		//scan assignment. It is an array of pointers that has been initialized according to the attribute type of the struct,The sql driver can perceive the parameter type,so it can be directly assigned to the pointer of the struct. In this way, the attributes of the struct have values
		//scanerr := rows.Scan(values...)
		if scanerr != nil {
			scanerr = fmt.Errorf("Query-->sqlRowsValues异常:%w", scanerr)
			FuncLogError(scanerr)
			return scanerr
		}

		//values[i] = f.Addr().Interface()
		//通过反射给slice添加元素
		//Add elements to slice through reflection
		if sliceElementTypePtr { //如果数组里是指针地址,*[]*struct
			sliceValue.Set(reflect.Append(sliceValue, pv.Addr()))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, pv))
		}

	}

	//查询总条数
	//Query total number
	if page != nil && finder.SelectTotalCount {
		count, counterr := selectCount(ctx, finder)
		if counterr != nil {
			counterr = fmt.Errorf("Query-->selectCount查询总条数错误:%w", counterr)
			FuncLogError(counterr)
			return counterr
		}
		page.setTotalCount(count)
	}

	return nil

}

// QueryRowMap 根据Finder查询,封装Map
// context必须传入,不能为空
// QueryRowMap encapsulates Map according to Finder query
// context must be passed in and cannot be empty
func QueryRowMap(ctx context.Context, finder *Finder) (map[string]interface{}, error) {

	if finder == nil {
		return nil, errors.New("QueryRowMap-->finder参数不能为nil")
	}
	resultMapList, listerr := QueryMap(ctx, finder, nil)
	if listerr != nil {
		listerr = fmt.Errorf("QueryRowMap-->QueryMap查询错误:%w", listerr)
		FuncLogError(listerr)
		return nil, listerr
	}
	if resultMapList == nil {
		return nil, nil
	}
	if len(resultMapList) > 1 {
		return resultMapList[0], errors.New("QueryRowMap查询出多条数据")
	} else if len(resultMapList) == 0 { //数据库不存在值
		return nil, nil
	}
	return resultMapList[0], nil
}

// QueryMap 根据Finder查询,封装Map数组
// 根据数据库字段的类型,完成从[]byte到golang类型的映射,理论上其他查询方法都可以调用此方法,但是需要处理sql.Nullxxx等驱动支持的类型
// context必须传入,不能为空
// QueryMap According to Finder query, encapsulate Map array
//According to the type of database field, the mapping from []byte to golang type is completed. In theory,other query methods can call this method, but need to deal with types supported by drivers such as sql.Nullxxx
//context must be passed in and cannot be empty
func QueryMap(ctx context.Context, finder *Finder, page *Page) ([]map[string]interface{}, error) {
	_, mapList, err := QueryMapAndColumnTypes(ctx, finder, page)
	return mapList, err
}

// QueryMap 根据Finder查询,封装Map数组
// 根据数据库字段的类型,完成从[]byte到golang类型的映射,理论上其他查询方法都可以调用此方法,但是需要处理sql.Nullxxx等驱动支持的类型
// context必须传入,不能为空
// QueryMap According to Finder query, encapsulate Map array
//According to the type of database field, the mapping from []byte to golang type is completed. In theory,other query methods can call this method, but need to deal with types supported by drivers such as sql.Nullxxx
//context must be passed in and cannot be empty
func QueryMapAndColumnTypes(ctx context.Context, finder *Finder, page *Page) ([]*sql.ColumnType, []map[string]interface{}, error) {

	if finder == nil {
		return nil, nil, errors.New("QueryMap-->finder参数不能为nil")
	}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return nil, nil, errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return nil, nil, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//db Connection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(0).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	sqlstr, err := wrapQuerySQL(dbType, finder, page)
	if err != nil {
		//err = fmt.Errorf("QueryMap -->wrapQuerySQL查询SQL语句错误:%w", err)
		FuncLogError(fmt.Errorf("QueryMap -->wrapQuerySQL查询SQL语句错误:%w", err))
		return nil, nil, err
	}

	//检查dbConnection.有可能会创建dbConnection或者开启事务,所以要尽可能的接近执行时检查
	//Check db Connection. It is possible to create a db Connection or start a transaction, so check it as close as possible to the execution
	var dbConnectionerr error
	ctx, dbConnection, dbConnectionerr = checkDBConnection(ctx, false, 0)
	if dbConnectionerr != nil {
		return nil, nil, dbConnectionerr
	}

	//根据语句和参数查询
	//Query based on statements and parameters
	rows, e := dbConnection.queryContext(ctx, &sqlstr, finder.values)
	if e != nil {
		//e = fmt.Errorf("QueryMap-->queryContext查询rows错误:%w", e)
		FuncLogError(fmt.Errorf("QueryMap-->queryContext查询rows错误:%w", e))
		return nil, nil, e
	}
	//先判断error 再关闭
	defer rows.Close()

	//数据库返回的列类型
	//The types returned by column Type.scan Type are all []byte, use column Type.database Type to judge one by one
	columnTypes, cne := rows.ColumnTypes()
	if cne != nil {
		//cne = fmt.Errorf("QueryMap-->rows.ColumnTypes数据库返回列名错误:%w", cne)
		FuncLogError(fmt.Errorf("QueryMap-->rows.ColumnTypes数据库返回列名错误:%w", cne))
		return nil, nil, cne
	}
	//反射获取 []driver.Value的值
	var driverValue reflect.Value
	cdvMapHasBool := len(CustomDriverValueMap) > 0
	//判断是否有自定义扩展,避免无意义的反射
	if cdvMapHasBool {
		driverValue = reflect.Indirect(reflect.ValueOf(rows))
		driverValue = driverValue.FieldByName("lastcols")
	}
	resultMapList := make([]map[string]interface{}, 0)
	//循环遍历结果集
	//Loop through the result set
	for rows.Next() {
		//接收数据库返回的数据,需要使用指针接收
		//To receive the data returned by the database, you need to use the pointer to receive
		values := make([]interface{}, len(columnTypes))
		//使用指针类型接收字段值,需要使用interface{}包装一下
		//To use the pointer type to receive the field value, you need to use interface() to wrap it
		result := make(map[string]interface{})

		//记录需要类型转换的字段信息
		var fieldTempDriverValueMap map[int]*driverValueInfo
		if cdvMapHasBool {
			fieldTempDriverValueMap = make(map[int]*driverValueInfo)
		}

		//给数据赋值初始化变量
		//Initialize variables by assigning values ​​to data
		for i, columnType := range columnTypes {

			//类型转换的接口实现
			var converFunc CustomDriverValueConver
			//是否需要类型转换
			var converOK bool = false
			//类型转换的临时值
			var tempDriverValue driver.Value
			//判断是否有自定义扩展,避免无意义的反射
			if cdvMapHasBool {
				dv := driverValue.Index(i)
				//根据接收的类型,获取到类型转换的接口实现
				converFunc, converOK = CustomDriverValueMap[dv.Elem().Type().String()]
			}
			var errGetDriverValue error
			//如果需要类型转换
			if converOK {
				//获取需要转的临时值
				tempDriverValue, errGetDriverValue = converFunc.GetDriverValue(columnType, nil, finder)
				if errGetDriverValue != nil {
					errGetDriverValue = fmt.Errorf("QueryMap-->conver.GetDriverValue异常:%w", errGetDriverValue)
					FuncLogError(errGetDriverValue)
					return nil, nil, errGetDriverValue
				}
				//返回值为nil,不做任何处理,使用原始逻辑
				if tempDriverValue == nil {
					values[i] = new(interface{})
				} else { //如果需要类型转换
					values[i] = tempDriverValue
					dvinfo := driverValueInfo{}
					dvinfo.converFunc = converFunc
					dvinfo.columnType = columnType
					dvinfo.tempDriverValue = tempDriverValue
					fieldTempDriverValueMap[i] = &dvinfo
				}

				continue
			}

			//不需要类型转换,正常赋值
			values[i] = new(interface{})
		}
		//scan赋值
		//scan assignment
		scanerr := rows.Scan(values...)
		if scanerr != nil {
			//scanerr = fmt.Errorf("QueryMap-->rows.Scan异常:%w", scanerr)
			FuncLogError( fmt.Errorf("QueryMap-->rows.Scan异常:%w", scanerr))
			return nil, nil, scanerr
		}

		//循环 需要类型转换的字段,把临时值赋值给实际的接收对象
		for i, driverValueInfo := range fieldTempDriverValueMap {
			//driverValueInfo := *driverValueInfoPtr
			//根据列名,字段类型,新值 返回符合接收类型值的指针,返回值是个指针,指针,指针!!!!
			rightValue, errConverDriverValue := driverValueInfo.converFunc.ConverDriverValue(driverValueInfo.columnType, nil, driverValueInfo.tempDriverValue, finder)
			if errConverDriverValue != nil {
				//errConverDriverValue = fmt.Errorf("QueryMap-->conver.ConverDriverValue异常:%w", errConverDriverValue)
				FuncLogError(fmt.Errorf("QueryMap-->conver.ConverDriverValue异常:%w", errConverDriverValue))
				return nil, nil, errConverDriverValue
			}
			//result[driverValueInfo.columnType.Name()] = reflect.ValueOf(rightValue).Elem().Interface()
			values[i] = rightValue
		}

		//获取每一列的值
		//Get the value of each column
		for i, columnType := range columnTypes {

			//取到指针下的值,[]byte格式
			//Get the value under the pointer, []byte format
			//v := *(values[i].(*interface{}))
			v := reflect.ValueOf(values[i]).Elem().Interface()
			//从[]byte转化成实际的类型值,例如string,int
			//Convert from []byte to actual type value, such as string, int
			v = converValueColumnType(v, columnType)
			//赋值到Map
			//Assign to Map
			result[columnType.Name()] = v

		}

		//添加Map到数组
		//Add Map to the array
		resultMapList = append(resultMapList, result)

	}

	//查询总条数
	//Query total number
	if page != nil && finder.SelectTotalCount {
		count, counterr := selectCount(ctx, finder)
		if counterr != nil {
			//counterr = fmt.Errorf("QueryMap-->selectCount查询总条数错误:%w", counterr)
			FuncLogError(fmt.Errorf("QueryMap-->selectCount查询总条数错误:%w", counterr))
			return columnTypes, resultMapList, counterr
		}
		page.setTotalCount(count)
	}

	return columnTypes, resultMapList, nil
}

// UpdateFinder 更新Finder语句
// ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
// affected影响的行数,如果异常或者驱动不支持,返回-1
// UpdateFinder Update Finder statement
// ctx cannot be nil, refer to zorm.Transaction method to pass in ctx. Don't build DB Connection yourself
// The number of rows affected by affected, if it is abnormal or the driver does not support it, return-1
func UpdateFinder(ctx context.Context, finder *Finder) (int, error) {
	affected := -1
	if finder == nil {
		return affected, errors.New("UpdateFinder-->finder不能为空")
	}
	sqlstr, err := finder.GetSQL()
	if err != nil {
		err = fmt.Errorf("UpdateFinder-->finder.GetSQL()错误:%w", err)
		FuncLogError(err)
		return affected, err
	}

	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}

	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//dbConnection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	sqlstr, err = reBindSQL(dbType, sqlstr)
	if err != nil {
		err = fmt.Errorf("UpdateFinder-->reBindSQL获取SQL语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}

	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	_, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, finder.values, nil)
	if errexec != nil {
		errexec = fmt.Errorf("UpdateFinder-->wrapExecUpdateValuesAffected执行更新错误:%w", errexec)
		FuncLogError(errexec)
	}

	return affected, errexec
}

// Insert 保存Struct对象,必须是IEntityStruct类型
// ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
// affected影响的行数,如果异常或者驱动不支持,返回-1
// Insert saves the Struct object, which must be of type IEntityStruct
// ctx cannot be nil, refer to zorm.Transaction method to pass in ctx. Don't build dbConnection yourself
// The number of rows affected by affected, if it is abnormal or the driver does not support it, return -1
func Insert(ctx context.Context, entity IEntityStruct) (int, error) {
	affected := -1
	if entity == nil {
		return affected, errors.New("对象不能为空")
	}
	typeOf, columns, values, columnAndValueErr := columnAndValue(entity)
	if columnAndValueErr != nil {
		columnAndValueErr = fmt.Errorf("Insert-->columnAndValue获取实体类的列和值异常:%w", columnAndValueErr)
		FuncLogError(columnAndValueErr)
		return affected, columnAndValueErr
	}
	if len(columns) < 1 {
		return affected, errors.New("Insert没有tag信息,请检查struct中 column 的tag")
	}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//dbConnection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//SQL语句
	//SQL statement
	sqlstr, autoIncrement, pktype, err := wrapInsertSQL(dbType, &typeOf, entity, &columns, &values)
	if err != nil {
		err = fmt.Errorf("Insert-->wrapInsertSQL获取保存语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}

	//oracle 12c+ 支持IDENTITY属性的自增列,因为分页也要求12c+的语法,所以数据库就IDENTITY创建自增吧
	//处理序列产生的自增主键,例如oracle,postgresql等
	var lastInsertID *int64
	var zormSQLOutReturningID *int64
	//如果是postgresql的SERIAL自增,需要使用 RETURNING 返回主键的值
	if autoIncrement > 0 {
		if dbType == "postgresql" || dbType == "kingbase" {
			var p int64 = 0
			lastInsertID = &p
			sqlstr = sqlstr + " RETURNING " + entity.GetPKColumnName()
		} else if dbType == "oracle" || dbType == "shentong" {
			var p int64 = 0
			zormSQLOutReturningID = &p
			sqlstr = sqlstr + " RETURNING " + entity.GetPKColumnName() + " INTO :zormSQLOutReturningID "
			v := sql.Named("zormSQLOutReturningID", sql.Out{Dest: zormSQLOutReturningID})
			values = append(values, v)
		}

	}

	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	res, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, lastInsertID)
	if errexec != nil {
		errexec = fmt.Errorf("Insert-->wrapExecUpdateValuesAffected执行保存错误:%w", errexec)
		FuncLogError(errexec)
		return affected, errexec
	}

	//如果是自增主键
	//If it is an auto-incrementing primary key
	if autoIncrement > 0 {
		//如果是oracle,shentong 的返回自增主键
		if lastInsertID == nil && zormSQLOutReturningID != nil {
			lastInsertID = zormSQLOutReturningID
		}

		var autoIncrementIDInt64 int64
		var e error
		if lastInsertID != nil {
			autoIncrementIDInt64 = *lastInsertID
		} else {
			//需要数据库支持,获取自增主键
			//Need database support, get auto-incrementing primary key
			autoIncrementIDInt64, e = (*res).LastInsertId()
		}

		//数据库不支持自增主键,不再赋值给struct属性
		//The database does not support self-incrementing primary keys, and no longer assigns values ​​to struct attributes
		if e != nil {
			e = fmt.Errorf("Insert-->LastInsertId数据库不支持自增主键,不再赋值给struct属性:%w", e)
			FuncLogError(e)
			return affected, nil
		}
		pkName := entity.GetPKColumnName()

		var seterr error

		if pktype == "int" {
			//int64 转 int
			//int64 to int
			autoIncrementIDInt, _ := typeConvertInt64toInt(autoIncrementIDInt64)
			//设置自增主键的值
			//Set the value of the auto-incrementing primary key
			seterr = setFieldValueByColumnName(entity, pkName, autoIncrementIDInt)
		} else if pktype == "int64" {
			//设置自增主键的值
			//Set the value of the auto-incrementing primary key
			seterr = setFieldValueByColumnName(entity, pkName, autoIncrementIDInt64)
		}

		if seterr != nil {
			seterr = fmt.Errorf("Insert-->setFieldValueByColumnName反射赋值数据库返回的自增主键错误:%w", seterr)
			FuncLogError(seterr)
			return affected, seterr
		}
	}

	return affected, nil

}

//InsertSlice 批量保存Struct Slice 数组对象,必须是[]IEntityStruct类型,golang目前没有泛型,使用IEntityStruct接口,兼容Struct实体类
//如果是自增主键,无法对Struct对象里的主键属性赋值
//ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
//affected影响的行数,如果异常或者驱动不支持,返回-1
func InsertSlice(ctx context.Context, entityStructSlice []IEntityStruct) (int, error) {
	affected := -1
	if entityStructSlice == nil || len(entityStructSlice) < 1 {
		return affected, errors.New("InsertSlice对象数组不能为空")
	}
	//第一个对象,获取第一个Struct对象,用于获取数据库字段,也获取了值
	entity := entityStructSlice[0]
	typeOf, columns, values, columnAndValueErr := columnAndValue(entity)
	if columnAndValueErr != nil {
		columnAndValueErr = fmt.Errorf("InsertSlice-->columnAndValue获取实体类的列和值异常:%w", columnAndValueErr)
		FuncLogError(columnAndValueErr)
		return affected, columnAndValueErr
	}
	if len(columns) < 1 {
		return affected, errors.New("InsertSlice没有tag信息,请检查struct中 column 的tag")
	}
	//从contxt中获取数据库连接,可能为nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}
	//自己构建的dbConnection
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	if dbConnection == nil { //dbConnection为nil,使用defaultDao
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//SQL语句
	sqlstr, _, err := wrapInsertSliceSQL(dbType, &typeOf, entityStructSlice, &columns, &values)
	if err != nil {
		err = fmt.Errorf("InsertSlice-->wrapInsertSliceSQL获取保存语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}
	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	_, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, nil)
	if errexec != nil {
		errexec = fmt.Errorf("InsertSlice-->wrapExecUpdateValuesAffected执行保存错误:%w", errexec)
		FuncLogError(errexec)
	}

	return affected, errexec

}

//Update 更新struct所有属性,必须是IEntityStruct类型
//ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
func Update(ctx context.Context, entity IEntityStruct) (int, error) {
	affected, err := updateStructFunc(ctx, entity, false)
	if err != nil {
		err = fmt.Errorf("Update-->updateStructFunc更新错误:%w", err)
		return affected, err
	}
	return affected, nil
}

//UpdateNotZeroValue 更新struct不为默认零值的属性,必须是IEntityStruct类型,主键必须有值
//ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
func UpdateNotZeroValue(ctx context.Context, entity IEntityStruct) (int, error) {
	affected, err := updateStructFunc(ctx, entity, true)
	if err != nil {
		err = fmt.Errorf("UpdateNotZeroValue-->updateStructFunc更新错误:%w", err)
		return affected, err
	}
	return affected, nil
}

//Delete 根据主键删除一个对象.必须是IEntityStruct类型
//ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
//affected影响的行数,如果异常或者驱动不支持,返回-1
func Delete(ctx context.Context, entity IEntityStruct) (int, error) {
	affected := -1
	typeOf, checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return affected, checkerr
	}

	pkName, pkNameErr := entityPKFieldName(entity, &typeOf)

	if pkNameErr != nil {
		pkNameErr = fmt.Errorf("Delete-->entityPKFieldName获取主键名称错误:%w", pkNameErr)
		FuncLogError(pkNameErr)
		return affected, pkNameErr
	}

	value, e := structFieldValue(entity, pkName)
	if e != nil {
		e = fmt.Errorf("Delete-->structFieldValue获取主键值错误:%w", e)
		FuncLogError(e)
		return affected, e
	}
	//从contxt中获取数据库连接,可能为nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}
	//自己构建的dbConnection
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	if dbConnection == nil { //dbConnection为nil,使用defaultDao
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//SQL语句
	sqlstr, err := wrapDeleteSQL(dbType, entity)
	if err != nil {
		err = fmt.Errorf("Delete-->wrapDeleteSQL获取SQL语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}
	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	values := make([]interface{}, 1)
	values[0] = value
	_, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, nil)
	if errexec != nil {
		errexec = fmt.Errorf("Delete-->wrapExecUpdateValuesAffected执行删除错误:%w", errexec)
		FuncLogError(errexec)
	}

	return affected, errexec

}

//InsertEntityMap 保存*IEntityMap对象.使用Map保存数据,用于不方便使用struct的场景,如果主键是自增或者序列,不要entityMap.Set主键的值
//ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
//affected影响的行数,如果异常或者驱动不支持,返回-1
func InsertEntityMap(ctx context.Context, entity IEntityMap) (int, error) {
	affected := -1
	//检查是否是指针对象
	_, checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return affected, checkerr
	}

	//从contxt中获取数据库连接,可能为nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}

	//自己构建的dbConnection
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	if dbConnection == nil { //dbConnection为nil,使用defaultDao
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//SQL语句
	sqlstr, values, autoIncrement, err := wrapInsertEntityMapSQL(dbType, entity)
	if err != nil {
		err = fmt.Errorf("InsertEntityMap-->wrapInsertEntityMapSQL获取SQL语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}

	//处理序列产生的自增主键,例如oracle,postgresql等
	var lastInsertID *int64
	var zormSQLOutReturningID *int64
	//如果是postgresql的SERIAL自增,需要使用 RETURNING 返回主键的值
	if autoIncrement && entity.GetPKColumnName() != "" {
		if dbType == "postgresql" || dbType == "kingbase" {
			var p int64 = 0
			lastInsertID = &p
			sqlstr = sqlstr + " RETURNING " + entity.GetPKColumnName()
		} else if dbType == "oracle" || dbType == "shentong" {
			var p int64 = 0
			zormSQLOutReturningID = &p
			sqlstr = sqlstr + " RETURNING " + entity.GetPKColumnName() + " INTO :zormSQLOutReturningID "
			v := sql.Named("zormSQLOutReturningID", sql.Out{Dest: zormSQLOutReturningID})
			values = append(values, v)
		}

	}

	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	res, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, lastInsertID)
	if errexec != nil {
		errexec = fmt.Errorf("InsertEntityMap-->wrapExecUpdateValuesAffected执行保存错误:%w", errexec)
		FuncLogError(errexec)
		return affected, errexec
	}

	//如果是自增主键
	if autoIncrement {
		//如果是oracle,shentong 的返回自增主键
		if lastInsertID == nil && zormSQLOutReturningID != nil {
			lastInsertID = zormSQLOutReturningID
		}

		var autoIncrementIDInt64 int64
		var e error
		if lastInsertID != nil {
			autoIncrementIDInt64 = *lastInsertID
		} else {
			//需要数据库支持,获取自增主键
			//Need database support, get auto-incrementing primary key
			autoIncrementIDInt64, e = (*res).LastInsertId()
		}
		if e != nil { //数据库不支持自增主键,不再赋值给struct属性
			e = fmt.Errorf("InsertEntityMap数据库不支持自增主键,不再赋值给IEntityMap:%w", e)
			FuncLogError(e)
			return affected, nil
		}
		//int64 转 int
		strInt64 := strconv.FormatInt(autoIncrementIDInt64, 10)
		autoIncrementIDInt, _ := strconv.Atoi(strInt64)
		//设置自增主键的值
		entity.Set(entity.GetPKColumnName(), autoIncrementIDInt)
	}

	return affected, nil

}

// UpdateEntityMap 更新IEntityMap对象.用于不方便使用struct的场景,主键必须有值
// ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
// affected影响的行数,如果异常或者驱动不支持,返回-1
// UpdateEntityMap Update IEntityMap object. Used in scenarios where struct is not convenient, the primary key must have a value
// ctx cannot be nil, refer to zorm.Transaction method to pass in ctx. Don't build DB Connection yourself
// The number of rows affected by "affected", if it is abnormal or the driver does not support it, return -1
func UpdateEntityMap(ctx context.Context, entity IEntityMap) (int, error) {
	affected := -1
	//检查是否是指针对象
	//Check if it is a pointer
	_, checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return affected, checkerr
	}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, it may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//dbConnection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	//SQL语句
	//SQL statement
	sqlstr, values, err := wrapUpdateEntityMapSQL(dbType, entity)
	if err != nil {
		err = fmt.Errorf("UpdateEntityMap-->wrapUpdateEntityMapSQL获取SQL语句错误:%w", err)
		FuncLogError(err)
		return affected, err
	}
	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	_, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, nil)
	if errexec != nil {
		errexec = fmt.Errorf("UpdateEntityMap-->wrapExecUpdateValuesAffected执行更新错误:%w", errexec)
		FuncLogError(errexec)
	}

	return affected, errexec

}

// updateStructFunc 更新对象
// ctx不能为nil,参照使用zorm.Transaction方法传入ctx.也不要自己构建DBConnection
// affected影响的行数,如果异常或者驱动不支持,返回-1
// updateStructFunc Update object
// ctx cannot be nil, refer to zorm.Transaction method to pass in ctx. Don't build DB Connection yourself
// The number of rows affected by "affected", if it is abnormal or the driver does not support it, return -1
func updateStructFunc(ctx context.Context, entity IEntityStruct, onlyUpdateNotZero bool) (int, error) {
	affected := -1
	if entity == nil {
		return affected, errors.New("updateStructFunc对象不能为空")
	}
	//从contxt中获取数据库连接,可能为nil
	//Get database connection from contxt, may be nil
	dbConnection, errFromContxt := getDBConnectionFromContext(ctx)
	if errFromContxt != nil {
		return affected, errFromContxt
	}
	//自己构建的dbConnection
	//dbConnection built by yourself
	if dbConnection != nil && dbConnection.db == nil {
		return affected, errDBConnection
	}

	var dbType string = ""
	//dbConnection为nil,使用defaultDao
	//dbConnection is nil, use default Dao
	if dbConnection == nil {
		dbType = FuncReadWriteStrategy(1).config.DBType
	} else {
		dbType = dbConnection.config.DBType
	}

	typeOf, columns, values, columnAndValueErr := columnAndValue(entity)
	if columnAndValueErr != nil {
		return affected, columnAndValueErr
	}

	//SQL语句
	//SQL statement
	sqlstr, err := wrapUpdateSQL(dbType, &typeOf, entity, &columns, &values, onlyUpdateNotZero)
	if err != nil {
		return affected, err
	}

	//包装update执行,赋值给影响的函数指针变量,返回*sql.Result
	_, errexec := wrapExecUpdateValuesAffected(ctx, &affected, &sqlstr, values, nil)
	if errexec != nil {
		errexec = fmt.Errorf("updateStruct-->wrapExecUpdateValuesAffected执行更新错误:%w", errexec)
		FuncLogError(errexec)
	}

	return affected, errexec

}

// selectCount 根据finder查询总条数
// context必须传入,不能为空
// selectCount Query the total number of items according to finder
// context must be passed in and cannot be empty
func selectCount(ctx context.Context, finder *Finder) (int, error) {

	if finder == nil {
		return -1, errors.New("selectCount参数为nil")
	}
	//自定义的查询总条数Finder,主要是为了在group by等复杂情况下,为了性能,手动编写总条数语句
	//Customized query total number Finder,mainly for the sake of performance in complex situations such as group by, manually write the total number of statements
	if finder.CountFinder != nil {
		count := -1
		_, err := QueryRow(ctx, finder.CountFinder, &count)
		if err != nil {
			return -1, err
		}
		return count, nil
	}

	countsql, counterr := finder.GetSQL()
	if counterr != nil {
		return -1, counterr
	}

	//查询order by 的位置
	//Query the position of order by
	locOrderBy := findOrderByIndex(countsql)
	//如果存在order by
	//If there is order by
	if len(locOrderBy) > 0 {
		countsql = countsql[:locOrderBy[0]]
	}
	s := strings.ToLower(countsql)
	gbi := -1
	locGroupBy := findGroupByIndex(countsql)
	if len(locGroupBy) > 0 {
		gbi = locGroupBy[0]
	}
	//特殊关键字,包装SQL
	//Special keywords, wrap SQL
	if strings.Contains(s, " distinct ") || strings.Contains(s, " union ") || gbi > -1 {
		countsql = "SELECT COUNT(*)  frame_row_count FROM (" + countsql + ") temp_frame_noob_table_name WHERE 1=1 "
	} else {
		locFrom := findSelectFromIndex(countsql)
		//没有找到FROM关键字,认为是异常语句
		//The FROM keyword was not found, which is considered an abnormal statement
		if len(locFrom) == 0 {
			return -1, errors.New("selectCount-->findFromIndex没有FROM关键字,语句错误")
		}
		countsql = "SELECT COUNT(*) " + countsql[locFrom[0]:]
	}

	countFinder := NewFinder()
	countFinder.Append(countsql)
	countFinder.values = finder.values

	count := -1
	_, cerr := QueryRow(ctx, countFinder, &count)
	if cerr != nil {
		return -1, cerr
	}
	return count, nil

}

// getDBConnectionFromContext 从Conext中获取数据库连接
// getDBConnectionFromContext Get database connection from Conext
func getDBConnectionFromContext(ctx context.Context) (*dataBaseConnection, error) {
	if ctx == nil {
		return nil, errors.New("getDBConnectionFromContext context不能为空")
	}
	//获取数据库连接
	//Get database connection
	value := ctx.Value(contextDBConnectionValueKey)
	if value == nil {
		return nil, nil
	}
	dbConnection, isdb := value.(*dataBaseConnection)
	if !isdb { //不是数据库连接
		return nil, errors.New("getDBConnectionFromContext context传递了错误的*DBConnection类型值")
	}
	return dbConnection, nil

}

//变量名建议errFoo这样的驼峰
//The variable name suggests a hump like "errFoo"
var errDBConnection = errors.New("更新操作需要使用zorm.Transaction开启事务.  读取操作如果ctx没有dbConnection,使用FuncReadWriteStrategy(rwType).newDBConnection(),如果dbConnection有事务,就使用事务查询")

// checkDBConnection 检查dbConnection.有可能会创建dbConnection或者开启事务,所以要尽可能的接近执行时检查
// context必须传入,不能为空.rwType=0 read,rwType=1 write
// checkDBConnection It is possible to create a db Connection or open a transaction, so check it as close as possible to execution
// The context must be passed in and cannot be empty. rwType=0 read, rwType=1 write
func checkDBConnection(ctx context.Context, hastx bool, rwType int) (context.Context, *dataBaseConnection, error) {

	dbConnection, errFromContext := getDBConnectionFromContext(ctx)
	if errFromContext != nil {
		return ctx, nil, errFromContext
	}
	//dbConnection为空
	//dbConnection is nil
	if dbConnection == nil {
		//是否禁用了事务
		//disabletx := FuncReadWriteStrategy(rwType).config.DisableTransaction
		//如果要求有事务,事务需要手动zorm.Transaction显示开启.如果自动开启,就会为了偷懒,每个操作都自动开启,事务就失去意义了
		//if hastx && (!disabletx) {
		if hastx {
			return ctx, nil, errDBConnection
		}

		//如果要求没有事务,实例化一个默认的dbConnection
		//If no transaction is required, instantiate a default db Connection
		var errGetDBConnection error
		dbConnection, errGetDBConnection = FuncReadWriteStrategy(rwType).newDBConnection()
		if errGetDBConnection != nil {
			return ctx, nil, errGetDBConnection
		}
		//把dbConnection放入context
		//Put db Connection into context
		ctx = context.WithValue(ctx, contextDBConnectionValueKey, dbConnection)

	} else { //如果dbConnection存在
		//If db Connection exists
		if dbConnection.db == nil { //禁止外部构建
			return ctx, dbConnection, errDBConnection
		}
		//if dbConnection.tx == nil && hastx && (!dbConnection.config.DisableTransaction) {
		if dbConnection.tx == nil && hastx { //如果要求有事务,事务需要手动zorm.Transaction显示开启.如果自动开启,就会为了偷懒,每个操作都自动开启,事务就失去意义了
			return ctx, dbConnection, errDBConnection
		}
	}
	return ctx, dbConnection, nil
}

// wrapExecUpdateValuesAffected 包装update执行,赋值给影响的函数指针变量,返回*sql.Result
func wrapExecUpdateValuesAffected(ctx context.Context, affected *int, sqlstrptr *string, values []interface{}, lastInsertID *int64) (*sql.Result, error) {
	//必须要有dbConnection和事务.有可能会创建dbConnection放入ctx或者开启事务,所以要尽可能的接近执行时检查
	//There must be a db Connection and transaction.It is possible to create a db Connection into ctx or open a transaction, so check as close as possible to the execution
	var dbConnectionerr error
	var dbConnection *dataBaseConnection
	ctx, dbConnection, dbConnectionerr = checkDBConnection(ctx, true, 1)
	if dbConnectionerr != nil {
		return nil, dbConnectionerr
	}

	// 数据库语法兼容处理
	sqlstr, reUpdateFinderSQLErr := reUpdateFinderSQL(dbConnection.config.DBType, sqlstrptr)
	if reUpdateFinderSQLErr != nil {
		reUpdateFinderSQLErr = fmt.Errorf("wrapExecUpdateValuesAffected-->reUpdateFinderSQL获取SQL语句错误:%w", reUpdateFinderSQLErr)
		FuncLogError(reUpdateFinderSQLErr)
		return nil, reUpdateFinderSQLErr
	}

	var res *sql.Result
	var errexec error
	if lastInsertID != nil {
		errexec = dbConnection.queryRowContext(ctx, sqlstr, values).Scan(lastInsertID)
		if errexec == nil { //如果插入成功,返回
			*affected = 1
			return res, errexec
		}
	} else {
		res, errexec = dbConnection.execContext(ctx, sqlstr, values)
	}

	if errexec != nil {
		return res, errexec
	}
	//影响的行数
	//Number of rows affected

	rowsAffected, errAffected := (*res).RowsAffected()
	if errAffected == nil {
		*affected, errAffected = typeConvertInt64toInt(rowsAffected)
	} else { //如果不支持返回条数,设置位nil,影响的条数设置成-1
		*affected = -1
		errAffected = nil
	}
	return res, errAffected
}
