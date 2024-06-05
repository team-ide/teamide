package golang

import (
	"sort"
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentKafkaCode = `
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
	common.AddOnStop(service.Close)
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

	var imports []string

	imports = append(imports, this_.golang.GetCommonImport())
	pack := this_.golang.GetComponentPack("kafka", name)

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"github.com/team-ide/go-tool/kafka"
`, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		imports = append(imports, s)
	}

	sort.Strings(imports)
	for _, im := range imports {
		if strings.HasSuffix(im, " _") {
			builder.AppendTabLine("_ \"" + strings.TrimSuffix(im, " _") + "\"")
		} else {
			builder.AppendTabLine("\"" + im + "\"")
		}
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(componentKafkaCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
