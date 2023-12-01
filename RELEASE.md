# 2.5.6

## Version 2.5.6 （2023/12/1）

发布功能

* 数据库添加SQL文件上传、打开等功能
  * SQL 执行将不再自动保存SQL，需要手动点击保存
* 数据库 库列表默认将配置的库或匹配的模式置顶展示
* Docker启动服务添加 Oracle 驱动 `instantclient-basic-linuxx64`
* 服务端模式添加 TLS
```conf
server:
  tls:
    open: false # 是否开启 https 默认关闭 建议开启后使用 https 访问
    cert: ./conf/server.crt  # 证书
    key: ./conf/server.key   # 证书 密钥* 
```