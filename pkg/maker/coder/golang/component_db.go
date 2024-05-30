package golang

import (
	"sort"
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentDbCode = `
var (
	service db.IService
)

func Init(c *db.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = db.New(c)
	if err != nil {
		return
	}
	config.AddOnStop(service.Close)
	return
}

func GetService() db.IService {
	return service
}

func GetSqlDb() *sql.DB {
	if service == nil {
		return nil
	}
	return service.GetDb()
}

func SelectOne[T any](columns string, table string, where string, obj T) (res T, err error) {

	return
}

func Insert(table string, obj any) (res int64, err error) {
	return
}

func Update(table string, update any, where string) (res int64, err error) {
	return
}

func Delete(table string, where string) (res int64, err error) {
	return
}

`
)

func (this_ *Generator) GenComponentDb(name string, model *modelers.ConfigDbModel) (err error) {
	dir := this_.golang.GetComponentDir(this_.Dir, "db", name)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "db.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	imports = append(imports, this_.golang.GetConfigImport())
	pack := this_.golang.GetComponentPack("db", name)

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"database/sql"
	"github.com/team-ide/go-tool/db"
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
	switch model.Type {
	case "mysql":
		imports = append(imports, "github.com/team-ide/go-tool/db/db_type_mysql _")
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

	code := strings.ReplaceAll(componentDbCode, "{pack}", pack)

	builder.AppendCode(code)
	return
}
