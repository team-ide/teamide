# Maker 相关说明

[DOT]

## 模型文件目录格式

* 根据 [`demo`](demo) 目录举例说明

```shell
# config
config # 配置目录
config/xxx # 某个组件配置
config/xxx/default # 某个组件默认配置
config/xxx/name1 # 某个组件 名称为 name1 的配置，使用时候 配置 config 即可指定

## config 举例
config/database # 数据库配置
config/database/default # 使用默认数据库的配置
config/database/TEST_DB # 指定 使用 TEST_DB

# constant
constant # 常量 用于存放全局可用的变量

## constant 举例
constant/user.yml  # 用户模块的常量 在使用时候 可以通过 USER_REDIS_EXPIRE 获取这个值

# error
error # 异常定义目录，可以使用 constant 中定义的变量

## error 举例
error/user.yml  # 用户模块的常量 在使用时候 可以通过 USER_IS_NULL 抛出 这个异常

# func
func # 自定义函数目录，可以使用 constant、error 中定义的变量

## func 举例
func/encryptPassword.yml  # 加密密码 在使用时候 可以通过 encryptPassword(arg1, arg2)

# struct
struct # 结构体目录，定义结构体、关联库表字段等，用于出入参，可以使用 constant 中定义的变量

## struct 举例
struct/user.yml  # 用户结构体 在使用时候 可以通过 user 定义类型

# storage
storage # 数据层，一般用于数据库读写

## storage 举例
storage/user/insert.yml  # 用户新增 在使用时候 可以通过 user/insert 指定调用

# service
service # 服务层，一般用于 业务逻辑处理 调用 storage、redis、zk、kafka、es等

## service 举例
service/user/insert.yml  # 用户新增 在使用时候 可以通过 user/insert 指定调用
```
