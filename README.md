# Team · IDE

Team IDE 团队在线开发工具

[![Code](https://img.shields.io/badge/Code-TeamIDE-red)](https://github.com/team-ide/teamide)
[![License](https://img.shields.io/badge/License-Apache--2.0%20License-blue)](https://github.com/team-ide/teamide/blob/main/LICENSE)
[![Release Latest](https://img.shields.io/badge/Release-V1.0.0-brightgreen)](https://github.com/team-ide/teamide/releases)
[![Download Window](https://img.shields.io/badge/Download-Window-orange)](https://github.com/team-ide/teamide/releases/latest/download/teamide-windows-x64.zip)
[![Download Linux](https://img.shields.io/badge/Download-Linux-orange)](https://github.com/team-ide/teamide/releases/latest/download/teamide-linux-x64.zip)
[![Download MAC](https://img.shields.io/badge/Download-MAC-orange)](https://github.com/team-ide/teamide/releases/latest/download/teamide-darwin-x64.zip)

## 注意

> #### Team IDE 服务器运行方式： 需要配置文件，数据和日志存储在`程序同级目录`下

> #### Team IDE 单机运行方式： 无需配置文件，数据和日志存储在`用户目录/TeamIDE`下，本地`需要安装谷歌浏览器`


## Team · IDE 功能模块

<table>
    <tr>
        <th width="180px">模块</th>
        <th>功能说明</th>
        <th width="80px">状态</th>
    </tr>
    <tr>
        <td rowspan="2">Toolbox SSH</td>
        <td>配置SSH连接，连接远程服务器，执行命令</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>点击FTP连接方式，上传、下载、移动、本地远程相互移动、重命名、删除、批量上传和下载等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td >Toolbox Zookeeper</td>
        <td>支持单机、集群，增删改查等操作，批量删除等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="2">Toolbox Kafka</td>
        <td>对Kafka主题增删改查等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>选择主题，推送、消费、删除数据等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="5">Toolbox Redis</td>
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
        <td rowspan="2">Toolbox Elasticsearch</td>
        <td>索引增删改查等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>选择索引，增删改查数据等</td>
        <td>完成</td>
    </tr>
    <tr>
        <td rowspan="4">Toolbox Database</td>
        <td>MySql库列表、库表数据加载</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>MySql表数据增删改查、批量新增、修改、删除等操作</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>自定义SQL执行面板，结果查看器</td>
        <td>完成</td>
    </tr>
    <tr>
        <td>适配Oracle等主流数据库</td>
        <td>进行中</td>
    </tr>
</table>

## 目录结构

服务端：go开发

前端：vue开发

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
# --isServer 作为服务运行，读取配置文件

# 单机版调试运行，需要谷歌浏览器
go run . --isDev

# 服务端调试运行，需要配置conf
go run . --isDev --isServer
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

#### Toolbox Redis（完成）

连接Redis，支持单机、集群，增删改查等操作，批量删除等

![avatar](doc/toolbox-redis.png)

![avatar](doc/toolbox-redis-set.png)

![avatar](doc/toolbox-redis-list.png)

![avatar](doc/toolbox-redis-hash.png)

#### Toolbox Zookeeper（完成）

连接Zookeeper，支持单机、集群，增删改查等操作，批量删除等

![avatar](doc/toolbox-zookeeper.png)

#### Toolbox Kafka（完成）

连接Kafka，增删改查主题，推送主题消息，自定义消费主题消息等

![avatar](doc/toolbox-kafka.png)

#### Toolbox SSH、SFTP（完成）

配置Linux服务器SSH连接，在线连接服务执行命令

![avatar](doc/toolbox-ssh.png)

SSH模块可以点击FTP，进行本地和远程文件管理 FTP：上传、下载、移动、本地远程相互移动、重命名、删除、批量上传和下载等功能

![avatar](doc/toolbox-ftp.png)

#### Toolbox Database（完成）

连接Database，在线编辑库表，编辑库表记录，查看表结构等

![avatar](doc/toolbox-database.png)

![avatar](doc/toolbox-database-data.png)

![avatar](doc/toolbox-database-table.png)

![avatar](doc/toolbox-database-sql.png)

![avatar](doc/toolbox-database-ddl.png)

![avatar](doc/toolbox-database-export-1.png)

![avatar](doc/toolbox-database-import-1.png)

#### Toolbox Elasticsearch（完成）

连接Elasticsearch，编辑索引，增删改查索引数据等

![avatar](doc/toolbox-elasticsearch.png)