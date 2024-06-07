package golang

import (
	"sort"
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentZkCode = `
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

	common.OnEvent(common.EventAppStop, func(args ...any) {
		service.Close()
	}, math.MaxInt)

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
	path := dir + "zk.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	imports = append(imports, this_.golang.GetCommonImport())
	pack := this_.golang.GetComponentPack("zk", name)

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"math"
	"github.com/team-ide/go-tool/zookeeper"
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
		as := strings.Split(im, " ")
		if len(as) == 2 {
			builder.AppendTabLine(as[1] + " \"" + as[0] + "\"")
		} else {
			builder.AppendTabLine("\"" + im + "\"")
		}
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(componentZkCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
