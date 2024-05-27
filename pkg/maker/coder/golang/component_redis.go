package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentRedisCode = `package {pack}

func Get[T any](key string, obj T) (res T, err error){
	return
}

func Set(key string, obj any, ex int64) (err error){
	return
}

func Del(key string) (err error){
	return
}
`
)

func (this_ *Generator) GenComponentRedis(name string, model *modelers.ConfigRedisModel) (err error) {
	dir := this_.golang.GetComponentRedisDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "redis.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(componentRedisCode, "{pack}", this_.golang.GetComponentRedisPack())

	builder.AppendCode(code)
	return
}
