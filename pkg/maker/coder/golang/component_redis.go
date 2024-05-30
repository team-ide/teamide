package golang

import (
	"strings"
	"teamide/pkg/maker/modelers"
)

var (
	componentRedisCode = `package {pack}

import (
	"app/config"
	"github.com/team-ide/go-tool/redis"
)

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
	config.AddOnStop(service.Close)
	return
}

func GetService() redis.IService {
	return service
}

func Get[T any](key string, obj T) (res T, err error) {
	return
}

func Set(key string, obj any, ex int64) (err error) {
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

	code := strings.ReplaceAll(componentRedisCode, "{pack}", this_.golang.GetComponentPack("redis", name))

	builder.AppendCode(code)
	return
}
