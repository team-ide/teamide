package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentDbCode = `package {pack}

func SelectOne[T any](columns string,table string, where string, obj T) (res T, err error){
	return
}

func Insert(table string, obj any) (res int64, err error){
	return
}

func Update(table string, update any, where string) (res int64, err error){
	return
}

func Delete(table string, where string) (res int64, err error){
	return
}
`
)

func (this_ *Generator) GenComponentDb(name string, model *modelers.ConfigDbModel) (err error) {
	dir := this_.golang.GetComponentDbDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "db.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentDbCode, "{pack}", this_.golang.GetComponentDbPack())

	builder.AppendCode(code)
	return
}
