package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentMongodbCode = `package {pack}

import (
	"app/config"
	"github.com/team-ide/go-tool/mongodb"
)

var (
	service mongodb.IService
)

func Init(c *mongodb.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = mongodb.New(c)
	if err != nil {
		return
	}
	config.AddOnStop(service.Close)
	return
}

func GetService() mongodb.IService {
	return service
}

`
)

func (this_ *Generator) GenComponentMongodb(name string, model *modelers.ConfigMongodbModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "mongodb", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "mongodb.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentMongodbCode, "{pack}", this_.golang.GetComponentPack("mongodb", name))

	builder.AppendCode(code)
	return
}
