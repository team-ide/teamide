v1.5.5
更新内容：
 - 增加CloseDB函数,关闭数据库连接池
 - 完善文档,注释

v1.5.4
更新内容：
 - QueryRow如果查询一个字段,而且这个字段数据库为null,会有异常,没有赋为默认值 
 - reflect.Type 类型的参数,修改为 *reflect.Type 指针,包括CustomDriverValueConver接口的参数
 - 完善文档,注释

v1.5.3
更新内容：
 - 感谢@Howard.TSE的建议,判断配置是否为空 
 - 感谢@haming123反馈性能问题.zorm 1.2.x 版本实现了基础功能,读性能比gorm和xorm快一倍.随着功能持续增加,造成性能下降,目前读性能只快了50%.  
 - 性能优化,去掉不必要的反射 
 - 完善文档,注释

v1.5.2
更新内容：
 - 感谢奔跑(@zeqjone)提供的正则,排除不在括号内的from,已经满足绝大部分场景
 - 感谢奔跑(@zeqjone) pr,修复 金仓数据库模型定义中tag数据库列标签与数据库内置关键词冲突时,加双引号处理
 - 升级 decimal 到1.3.1
 - 完善文档,注释

v1.5.1
更新内容：
 - 完善文档,注释
 - 注释未使用的代码
 - 先判断error,再执行defer rows.Close()
 - 增加微信社区支持(负责人是八块腹肌的单身小伙 @zhou-a-xing)


v1.5.0
更新内容：
 - 完善文档,注释
 - 支持clickhouse,更新,删除语句使用SQL92标准语法
 - ID默认使用时间戳+随机数,代替UUID实现
 - 优化SQL提取的正则表达式
 - 集成seata-golang,支持全局托管,不修改业务代码,零侵入分布式事务

v1.4.9
更新内容：
 - 完善文档,注释
 - 摊牌了,不装了,就是修改注释,刷刷版本活跃度

v1.4.8
更新内容：
 - 完善文档,注释
 - 数据库字段和实体类额外映射时,支持 _ 下划线转驼峰

v1.4.7
更新内容：
 - 情人节版本,返回map时,如果无法正常转换值类型,就返回原值,而不是nil

v1.4.6
更新内容：
 - 完善文档,注释
 - 千行代码,胜他十万,牛气冲天,zorm零依赖.(uuid和decimal这两个工具包竟然有1700行代码)
 - 在涉密内网开发环境中,零依赖能减少很多麻烦,做不到请不要说没必要......

v1.4.5
更新内容：
 - 增强自定义类型转换的功能
 - 完善文档,注释
 - 非常感谢 @anxuanzi 完善代码生成器
 - 非常感谢 @chien_tung 增加changelog,以后版本发布都会记录changelog

v1.4.4
更新内容：
 - 如果查询的字段在column tag中没有找到,就会根据名称(不区分大小写)映射到struct的属性上
 - 给QueryRow方法增加 has 的返回值,标识数据库是有一条记录的,各位已经使用的大佬,升级时注意修改代码,非常抱歉*3！

v1.4.3
更新内容：
 - 正式支持南大通用(gbase)数据库,完成国产四库的适配
 - 增加设置全局事务隔离级别和单个事务的隔离级别
 - 修复触发器自增主键的逻辑bug
 - 文档完善和细节调整

v1.4.2
更新内容:
 - 正式支持神州通用(shentong)数据库
 - 完善pgsql和kingbase的自增主键返回值支持
 - 七家公司的同学建议查询和golang sql方法命名保持统一.做了一个艰难的决定,修改zorm的部分方法名.全局依次替换字符串即可.
zorm.Query(                 替换为    zorm.QueryRow(
zorm.QuerySlice(          替换为    zorm.Query(
zorm.QueryMap(          替换为    zorm.QueryRowMap(
zorm.QueryMapSlice(   替换为    zorm.QueryMap(

v1.4.1
更新内容：
 - 支持自定义扩展字段映射逻辑
 
v1.4.0
更新内容：
 - 修改多条数据的判断逻辑

v1.3.9
更新内容：
 - 支持自定义数据类型,包括json/jsonb
 - 非常感谢 @chien_tung  同学反馈的问题, QuerySlice方法支持*[]*struct类型,简化从xorm迁移
 - 其他代码细节优化.

v1.3.7
更新内容：
 - 非常感谢 @zhou- a- xing 同学(八块腹肌的单身少年)的英文翻译,zorm的核心代码注释已经是中英双语了.
 - 非常感谢 @chien_tung  同学反馈的问题,修复主键自增int和init64类型的兼容性.
 - 其他代码细节优化.

v1.3.6
更新内容：
 - 完善注释文档
 - 修复Delete方法的参数类型错误
 - 其他代码细节优化.

v1.3.5
更新内容：
 - 完善注释文档
 - 兼容处理数据库为null时,基本类型取默认值,感谢@fastabler的pr
 - 修复批量保存方法的一个bug:如果slice的长度为1,在pgsql和oracle会出现异常
 - 其他代码细节优化.

v1.3.4
更新内容：
 - 完善注释文档
 - 取消分页语句必须有order by的限制
 - 支持人大金仓数据库
 - 人大金仓驱动说明: https://help.kingbase.com.cn/doc- view- 8108.html
 - 人大金仓kingbase 8核心是基于postgresql 9.6,可以使用 https://github.com/lib/pq 进行测试,生产环境建议使用官方驱动
 
v1.3.3
更新内容：
 - 完善注释文档
 - 增加批量保存Struct对象方法
 - 正式支持达梦数据库
 - 基于达梦官方驱动,发布go mod项目 https://gitee.com/chunanyong/dm
 
v1.3.2
更新内容：
 - 增加达梦数据的分页适配
 - 完善调整代码注释
 - 增加存储过程和函数的调用示例

v1.3.1
更新内容：
 - 修改方法名称,和gorm和xorm保持相似,降低迁移和学习成本
 - 更新测试用例文档

v1.3.0
更新内容：
 - 去掉zap日志依赖,通过复写  FuncLogError FuncLogPanic FuncPrintSQL 实现自定义日志
 - golang版本依赖调整为v1.13
 - 迁移测试用到readygo,zorm项目不依赖任何数据库驱动包

v1.2.9
更新内容：
 - IEntityMap支持主键自增或主键序列
 - 更新方法返回影响的行数affected
 - 修复 查询IEntityMap时数据库无记录出现异常的bug
 - 测试用例即文档 https://gitee.com/chunanyong/readygo/blob/master/test/testzorm/BaseDao_test.go

v1.2.8
更新内容：
 - 暴露FuncGenerateStringID函数,方便自定义扩展字符串主键ID
 - Finder.Append 默认加一个空格,避免手误出现语法错误
 - 缓存字段信息时,使用map代替sync.Map,提高性能
 - 第三方性能压测结果

v1.2.6
更新内容：
 - DataSourceConfig 配置区分 DriverName 和 DBType，兼容一种数据库的多个驱动包
 - 不再显示依赖数据库驱动，由使用者确定依赖的数据库驱动包
 
v1.2.5
更新内容：
 - 分页语句必须有明确的order by,避免数据库迁移时出现分页语法不兼容.
 - 修复列表查询时,page对象为nil的bug
 
v1.2.3
更新内容：
 - 完善数据库支持,目前支持MySQL,SQLServer,Oracle,PostgreSQL,SQLite3
 - 简化数据库读写分离实现,暴露zorm.FuncReadWriteBaseDao函数属性,用于自定义读写分离策略
 - 精简zorm.DataSourceConfig属性,增加PrintSQL属性

v1.2.2
更新内容：
 - 修改NewPage()返回Page对象指针,传递时少写一个 & 符号
 - 取消GetDBConnection()方法,使用BindContextConnection()方法进行多个数据库库绑定
 - 隐藏DBConnection对象,不再对外暴露数据库对象,避免手动初始化造成的异常
 
v1.1.8
更新内容：
 - 修复UUID支持
 - 数据库连接和事务隐藏到context.Context为统一参数,符合golang规范,更好的性能
 - 封装logger实现,方便更换log包
 - 增加zorm.UpdateStructNotZeroValue 方法,只更新不为零值的字段
 - 完善测试用例 
