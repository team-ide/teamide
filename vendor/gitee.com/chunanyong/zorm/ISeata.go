package zorm

import "context"

// ISeataGlobalTransaction seata-golang的包装接口,隔离seata-golang的依赖
// 声明一个struct,实现这个接口,并配置实现 FuncSeataGlobalTransaction 函数
/**

//不使用proxy代理模式,全局托管,不修改业务代码,零侵入实现分布式事务
//tm.Implement(svc.ProxySvc)


// 分布式事务示例代码
_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

    // 获取当前分布式事务的XID.不用考虑怎么来的,如果是分布式事务环境,会自动设置值
    // xid := ctx.Value("XID").(string)

	// 把xid传递到第三方应用
	// req.Header.Set("XID", xid)

	// 如果返回的err不是nil,本地事务和分布式事务就会回滚
	return nil, err
})

///----------第三方应用-------///

// 第三方应用开启事务前,ctx需要绑定XID,例如使用了gin框架

// 接受传递过来的XID,绑定到本地ctx
// xid:=c.Request.Header.Get("XID")
// 获取到ctx
// ctx := c.Request.Context()
// ctx = context.WithValue(ctx,"XID",xid)

// ctx绑定XID之后,调用业务事务
_, err := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

    // 业务代码......

	// 如果返回的err不是nil,本地事务和分布式事务就会回滚
	return nil, err
})


// 建议以下代码放到单独的文件里
//................//

// ZormSeataGlobalTransaction 包装seata的*tm.DefaultGlobalTransaction,实现zorm.ISeataGlobalTransaction接口
type ZormSeataGlobalTransaction struct {
	*tm.DefaultGlobalTransaction
}

// MyFuncSeataGlobalTransaction zorm适配seata分布式事务的函数
// 重要!!!!需要配置zorm.DataSourceConfig.FuncSeataGlobalTransaction=MyFuncSeataGlobalTransaction 重要!!!
func MyFuncSeataGlobalTransaction(ctx context.Context) (zorm.ISeataGlobalTransaction, context.Context, error) {
	//获取seata的rootContext
	rootContext := seataContext.NewRootContext(ctx)
	//创建seata事务
	seataTx := tm.GetCurrentOrCreate(rootContext)
	//使用zorm.ISeataGlobalTransaction接口对象包装seata事务,隔离seata-golang依赖
	seataGlobalTransaction := ZormSeataGlobalTransaction{seataTx}

	return seataGlobalTransaction, rootContext, nil
}
//实现zorm.ISeataGlobalTransaction接口
func (gtx ZormSeataGlobalTransaction) SeataBegin(ctx context.Context) error {
	rootContext := ctx.(*seataContext.RootContext)
	return gtx.BeginWithTimeout(int32(6000), rootContext)
}

func (gtx ZormSeataGlobalTransaction) SeataCommit(ctx context.Context) error {
	rootContext := ctx.(*seataContext.RootContext)
	return gtx.Commit(rootContext)
}

func (gtx ZormSeataGlobalTransaction) SeataRollback(ctx context.Context) error {
	rootContext := ctx.(*seataContext.RootContext)
	//如果是Participant角色,修改为Launcher角色,允许分支事务提交全局事务.
	if gtx.Role != tm.Launcher {
		gtx.Role = tm.Launcher
	}
	return gtx.Rollback(rootContext)
}

func (gtx ZormSeataGlobalTransaction) GetSeataXID(ctx context.Context) string {
	rootContext := ctx.(*seataContext.RootContext)
	return rootContext.GetXID()
}
//................//
**/

type ISeataGlobalTransaction interface {
	//开启seata全局事务
	SeataBegin(ctx context.Context) error

	//提交seata全局事务
	SeataCommit(ctx context.Context) error

	//回滚seata全局事务
	SeataRollback(ctx context.Context) error

	//获取seata事务的XID
	GetSeataXID(ctx context.Context) string

	//重新包装为seata的context.RootContext
	//context.RootContext 如果后续使用了 context.WithValue,类型就是context.valueCtx 就会造成无法再类型断言为 context.RootContext
	//所以DBDao里使用了 seataRootContext变量,区分业务的ctx和seata的RootContext
	//SeataNewRootContext(ctx context.Context) context.Context
}
