package golang

import (
	"sort"
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentDbCode = `
var (
	service     db.IService
	options     *db.TemplateOptions
	mapTemplate *db.Template[map[string]interface{}]
)

func Init(c *db.Config) (err error) {
	if service != nil {
		service.Close()
	}
	service, err = db.New(c)
	if err != nil {
		return
	}
	options = &db.TemplateOptions{
		Service: service,
	}
	mapTemplate = db.WarpTemplate(map[string]interface{}{}, options)
	common.AddOnStop(service.Close)
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

func SelectOne[T any](ctx context.Context, sqlParamSql string, sqlParam any, obj T) (res T, err error) {
	template := db.WarpTemplate(obj, options)
	res, err = template.SelectOne(ctx, sqlParamSql, sqlParam)
	return
}

func SelectList[T any](ctx context.Context, sqlParamSql string, sqlParam any, obj T) (res []T, err error) {
	template := db.WarpTemplate(obj, options)
	res, err = template.SelectList(ctx, sqlParamSql, sqlParam)
	return
}

func Insert(ctx context.Context, table string, obj any) (res int64, err error) {
	res, err = mapTemplate.Insert(ctx, table, obj)
	return
}

func Update(ctx context.Context, table string, update any, whereSql string, whereParam any, ignoreColumns ...string) (res int64, err error) {
	res, err = mapTemplate.Update(ctx, table, update, whereSql, whereParam, ignoreColumns...)
	return
}

func Delete(ctx context.Context, table string, whereSql string, whereParam any) (res int64, err error) {
	res, err = mapTemplate.Delete(ctx, table, whereSql, whereParam)
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

	imports = append(imports, this_.golang.GetCommonImport())
	imports = append(imports, this_.golang.GetConfigImport())
	pack := this_.golang.GetComponentPack("db", name)

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"context"
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
