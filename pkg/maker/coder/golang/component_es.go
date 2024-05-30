package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentEsCode = `package {pack}

import (
	"app/config"
	"github.com/team-ide/go-tool/elasticsearch"
)

var (
	service elasticsearch.IService
)

func Init(c *elasticsearch.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = elasticsearch.New(c)
	if err != nil {
		return
	}
	config.AddOnStop(service.Close)
	return
}

func GetService() elasticsearch.IService {
	return service
}

`
)

func (this_ *Generator) GenComponentEs(name string, model *modelers.ConfigEsModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "es", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "es.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentEsCode, "{pack}", this_.golang.GetComponentPack("es", name))

	builder.AppendCode(code)
	return
}
