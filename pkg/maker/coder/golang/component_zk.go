package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentZkCode = `package {pack}

import (
	"app/config"
	"github.com/team-ide/go-tool/zookeeper"
)

var (
	service zookeeper.IService
)

func Init(c *zookeeper.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = zookeeper.New(c)
	if err != nil {
		return
	}
	config.AddOnStop(service.Close)
	return
}

func GetService() zookeeper.IService {
	return service
}

`
)

func (this_ *Generator) GenComponentZk(name string, model *modelers.ConfigZkModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "zk", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "es.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentZkCode, "{pack}", this_.golang.GetComponentPack("zk", name))

	builder.AppendCode(code)
	return
}
