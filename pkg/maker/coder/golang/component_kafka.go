package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentKafkaCode = `package {pack}

import (
	"app/config"
	"github.com/team-ide/go-tool/kafka"
)

var (
	service kafka.IService
)

func Init(c *kafka.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = kafka.New(c)
	if err != nil {
		return
	}
	config.AddOnStop(service.Close)
	return
}

func GetService() kafka.IService {
	return service
}

`
)

func (this_ *Generator) GenComponentKafka(name string, model *modelers.ConfigKafkaModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "kafka", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "kafka.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentKafkaCode, "{pack}", this_.golang.GetComponentPack("kafka", name))

	builder.AppendCode(code)
	return
}
