
# 打包静态资源到html.go

```shell
## 需要go-bindata支持，安装一次即可
go get -u github.com/jteeuwen/go-bindata/...
```

```shell
## 删除../html.go文件
rm -rf web/html.go
# 执行以下命令打包，需要执行
cd web/html && go-bindata -pkg web -o ../html.go ./... && cd ../../
```
