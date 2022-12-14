# Team · IDE

Team IDE 团队在线开发工具

[![Code](https://img.shields.io/badge/Code-TeamIDE-red)](https://github.com/team-ide/teamide)
[![License](https://img.shields.io/badge/License-Apache--2.0%20License-blue)](https://github.com/team-ide/teamide/blob/main/LICENSE)
[![GitHub Release Latest](https://img.shields.io/badge/Release-V1.8.7-brightgreen)](https://github.com/team-ide/teamide/releases)

## Team · IDE 功能模块

<table>
    <tr>
        <th>模块</th>
        <th>功能说明</th>
        <th>状态</th>
    </tr>
    <tr>
        <td rowspan="4">SSH</td>
        <td>配置SSH连接，连接远程服务器，执行命令，支持自定义快速指令</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>SSH支持rz、sz命令，rz支持批量上传，支持打开FTP</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>点击FTP连接方式查看本地目录、服务端目录</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>FTP在线编辑、上传、下载、移动、本地远程相互移动、重命名、删除、批量上传和下载等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td >Zookeeper</td>
        <td>支持单机、集群，增删改查等操作，批量删除等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="2">Kafka</td>
        <td>对Kafka主题增删改查等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>选择主题，推送、消费、删除数据等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="6">Redis</td>
        <td>Redis Key搜索、模糊查询、删除、新增等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>字符串值编辑</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>哈希值编辑</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>列表值编辑</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>集合值编辑</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>根据策略导入功能，配置Key、Value自动导入相应格式string、list、hash、set等数据</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="3">Elasticsearch</td>
        <td>索引增删改查等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>选择索引，增删改查数据等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>添加索引，设置字段，索引迁移等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="7">Database</td>
        <td>数据库库|用户|模式列表、表数据加载</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>数据库库表数据增删改查、批量新增、修改、删除等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>表格选择数据导出SQL（新增、修改、删除数据SQL）等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>根据策略批量导入数据，自定义导入数量，值格式，批量导入</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>自定义SQL执行面板，结果查看器</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>新建库，在线设计表，查看建表、更新表SQL语句</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>支持数据库，MySql、Oracle、达梦、金仓、神通、Sqlite等数据库</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="2">HTTP</td>
        <td>配置HTTP的GET，POST，DELETE，PUT等REST接口</td>
        <td>进行中</td>
    </tr>
    <tr>
        <td>配置策略，根据策略并发请求HTTP接口，汇总结果报文</td>
        <td>进行中</td>
    </tr>
    <tr>
        <td >格式转换</td>
        <td>XML、JSON、URL、YAML、TOML等格式相互转换</td>
        <td>进行中</td>
    </tr>
    <tr>
        <td >节点功能</td>
        <td>可以配置多服务器之间网络透传，内外网相互透传等</td>
        <td>完成</td>
    </tr>
</table>

## 目录结构

服务端：go开发

前端：vue开发


### 注意

> #### Team IDE 单机运行方式： 无需配置文件，数据和日志存储在`用户目录/temeide`下

> #### Team IDE 服务器运行方式： 需要配置文件，数据和日志存储在`程序同级目录`下

#### Docker 运行

```shell
docker run -itd --name teamide-21080 -p 21080:21080 -v /data/teamide/data:/opt/teamide/data teamide/teamide-server:1.9.0
```

```shell
conf/           # 配置文件
html/           # 前端，vue工程
internal/       # 服务源码
pkg/            # 工具等
```

### 源码调试运行

**前端调试运行**

```shell
# 前端打包

# 进入html目录
cd html

# 安装依赖
npm install

# 运行
npm run serve
```

**服务端调试运行**

```shell
# 安装依赖
go mod tidy

# 运行
# --isDev dev模式，自动打开到 前端调试页面，日志输出控制台

# 单机版调试运行，需要谷歌浏览器
go run . --isDev
```

### 打包

**前端打包**

```shell
# 前端打包

# 进入html目录
cd html

# 安装依赖
npm install

# 打包
npm run build
```

**静态资源打包为Go文件**

```shell
# 安装依赖
go mod tidy

# 前端文件发布到服务中
# 将自动将前端文件打包成到internal/static/html.go文件中
go test -v -timeout 3600s -run ^TestStatic$ teamide/internal/static
```

**单机版可执行文件打包，单机版运行需要谷歌浏览器**

```shell
# 安装依赖
go mod tidy

# 打包单机运行，需要本地安装谷歌浏览器，用于单个人员使用
# 不需要conf目录
go build .
```

**作为服务部署打包**

```shell
# 安装依赖
go mod tidy

# 作为服务端部署，通过浏览器打开，可供团队使用
# 需要conf目录
go build -ldflags "-X main.buildFlags=--isServer" .
```

## Toolbox 模块

工具箱，用于连接Redis、Zookeeper、Database、SSH、SFTP、Kafka、Elasticsearch等

### Toolbox 功能

![avatar](doc/toolbox-type.png)

#### Toolbox Redis（完成）

连接Redis，支持单机、集群，增删改查等操作，批量删除等

![avatar](doc/toolbox-redis.png)

![avatar](doc/toolbox-redis-set.png)

![avatar](doc/toolbox-redis-list.png)

![avatar](doc/toolbox-redis-hash.png)

![avatar](doc/toolbox-redis-import.png)

#### Toolbox Zookeeper（完成）

连接Zookeeper，支持单机、集群，增删改查等操作，批量删除等

![avatar](doc/toolbox-zookeeper.png)

#### Toolbox Kafka（完成）

连接Kafka，增删改查主题，推送主题消息，自定义消费主题消息等

![avatar](doc/toolbox-kafka.png)

![avatar](doc/toolbox-kafka-data.png)

#### Toolbox SSH、SFTP（完成）

配置Linux服务器SSH连接，在线连接服务执行命令

![avatar](doc/toolbox-ssh.png)

![avatar](doc/toolbox-ssh-upload.png)

![avatar](doc/toolbox-ssh-download.png)

SSH模块可以点击FTP，进行本地和远程文件管理 FTP：上传、下载、移动、本地远程相互移动、重命名、删除、批量上传和下载等功能

![avatar](doc/toolbox-ftp.png)

![avatar](doc/toolbox-ftp-edit-file.png)

#### Toolbox Database（完成）

连接Database，在线编辑库表，编辑库表记录，查看表结构等

![avatar](doc/toolbox-database.png)

![avatar](doc/toolbox-database-data.png)

![avatar](doc/toolbox-database-table.png)

![avatar](doc/toolbox-database-sql.png)

![avatar](doc/toolbox-database-ddl.png)

![avatar](doc/toolbox-database-export-data-sql.png)

![avatar](doc/toolbox-database-export-1.png)

![avatar](doc/toolbox-database-export-2.png)

![avatar](doc/toolbox-database-import-1.png)

#### Toolbox Elasticsearch（完成）

连接Elasticsearch，编辑索引，增删改查索引数据等

![avatar](doc/toolbox-elasticsearch.png)

![avatar](doc/toolbox-elasticsearch-data.png)

#### Toolbox 其它

![avatar](doc/toolbox-other-format.png)

## Node 模块

节点服务，用于不同网段通信，借助节点模块的网络代理实现内外网透传等

![avatar](doc/toolbox-node.png)

![avatar](doc/toolbox-node-net-proxy.png)