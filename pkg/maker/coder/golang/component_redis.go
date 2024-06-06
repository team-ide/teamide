package golang

import (
	"sort"
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentRedisCode = `
var (
	service redis.IService
)

func Init(c *redis.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = redis.New(c)
	if err != nil {
		return
	}
	common.AddOnStop(service.Close)
	return
}

func GetService() redis.IService {
	return service
}

func Get[T any](key string, obj T) (res T, err error) {
	value, err := service.Get(key)
	if err != nil {
		return
	}
	if value != "" {
		res, err = util.StringTo(value, obj)
	}
	return
}

func Set(key string, obj any, ex int64) (err error) {
	value := util.GetStringValue(obj)
	err = service.Set(key, value)
	if err != nil {
		return
	}
	if ex > 0 {
		_, err = service.Expire(key, ex)
	}
	return
}

func Del(key string) (err error) {
	return
}

`
)

func (this_ *Generator) GenComponentRedis(name string, model *modelers.ConfigRedisModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "redis", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "redis.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	imports = append(imports, this_.golang.GetCommonImport())
	pack := this_.golang.GetComponentPack("redis", name)

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
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

	code := strings.ReplaceAll(componentRedisCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
