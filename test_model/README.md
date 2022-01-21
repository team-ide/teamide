# application

## app 目录结构

app目录为应用程序模型目录

目录结构如`struct/userInfo.yaml`其中`userInfo`为模型名称
目录结构如`service/user/insert.yaml`其中`user/insert`为模型名称

``` shell
constant/**.yaml                 # 常量模型 定义一些常量键值对，全局可用，使用$xxx即可获取该值

error/**.yaml                    # 异常模型 定义错误码和错误信息，用于抛出自定义异常

struct/**.yaml                   # 结构体模型 所有服务参数使用结构体，可以配置结构体与表关联

datasource/database/**.yaml      # 数据库模型 配置数据库连接，用于服务中Sql操作 
                    default.yaml # 使用default.yaml作为文件名则属于默认连接

datasource/redis/**.yaml         # Redis模型 配置Redis连接，用于服务中使用Redis操作
                 default.yaml    # 使用default.yaml作为文件名则属于默认连接

service/**.yaml                  # 服务模型 可以在该模型配置业务操作，如参数验证、Sql操作、Redis操作、文件操作等

server/web/**.yaml               # 服务器模型 可以配置Web Http服务，为service提供Web Api接口能力

test/**.yaml                     # 测试模型 测试模型中定义变量，可以测试service

```

## 命令

``` shell
# 启动服务，将读取server下所有服务配置并启动
main server

# 初始化数据库（如果不存在则建库建表），将读取struct下所有结构体，如果结构体绑定表名称，则判断表是否存在，不存在则创建表
main init database

# 测试服务，将读取test下xxx xxxx xx测试配置并执行
main test xxx xxxx xx

# 输出Web Api文档，读取server/web配置，输出Web Api文档到doc目录
main doc web

# 参数 模型目录使用 dir=./app 读取main所在位置的相对路径
main xxx dir=./app

```