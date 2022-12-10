package module_redis

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/redis"
	"teamide/pkg/util"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power              = base.AppendPower(&base.PowerAction{Action: "redis", Text: "Redis", ShouldLogin: true, StandAlone: true})
	infoPower          = base.AppendPower(&base.PowerAction{Action: "info", Text: "Redis信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getPower           = base.AppendPower(&base.PowerAction{Action: "get", Text: "Redis获取Key值", ShouldLogin: true, StandAlone: true, Parent: Power})
	keysPower          = base.AppendPower(&base.PowerAction{Action: "keys", Text: "Redis查询Keys", ShouldLogin: true, StandAlone: true, Parent: Power})
	setPower           = base.AppendPower(&base.PowerAction{Action: "set", Text: "Redis设置值", ShouldLogin: true, StandAlone: true, Parent: Power})
	saddPower          = base.AppendPower(&base.PowerAction{Action: "sadd", Text: "Redis SAdd", ShouldLogin: true, StandAlone: true, Parent: Power})
	sremPower          = base.AppendPower(&base.PowerAction{Action: "srem", Text: "Redis SRem", ShouldLogin: true, StandAlone: true, Parent: Power})
	lpushPower         = base.AppendPower(&base.PowerAction{Action: "lpush", Text: "Redis LPush", ShouldLogin: true, StandAlone: true, Parent: Power})
	rpushPower         = base.AppendPower(&base.PowerAction{Action: "rpush", Text: "Redis RPush", ShouldLogin: true, StandAlone: true, Parent: Power})
	lsetPower          = base.AppendPower(&base.PowerAction{Action: "lset", Text: "Redis LSet", ShouldLogin: true, StandAlone: true, Parent: Power})
	lremPower          = base.AppendPower(&base.PowerAction{Action: "lrem", Text: "Redis LRem", ShouldLogin: true, StandAlone: true, Parent: Power})
	hsetPower          = base.AppendPower(&base.PowerAction{Action: "hset", Text: "Redis HSet", ShouldLogin: true, StandAlone: true, Parent: Power})
	hdelPower          = base.AppendPower(&base.PowerAction{Action: "hdel", Text: "Redis HDel", ShouldLogin: true, StandAlone: true, Parent: Power})
	deletePower        = base.AppendPower(&base.PowerAction{Action: "delete", Text: "Redis删除Key", ShouldLogin: true, StandAlone: true, Parent: Power})
	deletePatternPower = base.AppendPower(&base.PowerAction{Action: "deletePattern", Text: "Redis删除匹配Key", ShouldLogin: true, StandAlone: true, Parent: Power})
	expirePower        = base.AppendPower(&base.PowerAction{Action: "expire", Text: "Redis设置过期", ShouldLogin: true, StandAlone: true, Parent: Power})
	ttlPower           = base.AppendPower(&base.PowerAction{Action: "ttl", Text: "Redis过期时间查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	persistPower       = base.AppendPower(&base.PowerAction{Action: "persist", Text: "Redis移除过期时间", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower         = base.AppendPower(&base.PowerAction{Action: "close", Text: "Redis关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: keysPower, Do: this_.keys})
	apis = append(apis, &base.ApiWorker{Power: setPower, Do: this_.set})
	apis = append(apis, &base.ApiWorker{Power: saddPower, Do: this_.sadd})
	apis = append(apis, &base.ApiWorker{Power: sremPower, Do: this_.srem})
	apis = append(apis, &base.ApiWorker{Power: lpushPower, Do: this_.lpush})
	apis = append(apis, &base.ApiWorker{Power: rpushPower, Do: this_.rpush})
	apis = append(apis, &base.ApiWorker{Power: lsetPower, Do: this_.lset})
	apis = append(apis, &base.ApiWorker{Power: lremPower, Do: this_.lrem})
	apis = append(apis, &base.ApiWorker{Power: hsetPower, Do: this_.hset})
	apis = append(apis, &base.ApiWorker{Power: hdelPower, Do: this_.hdel})
	apis = append(apis, &base.ApiWorker{Power: deletePower, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: deletePatternPower, Do: this_.deletePattern})
	apis = append(apis, &base.ApiWorker{Power: expirePower, Do: this_.expire})
	apis = append(apis, &base.ApiWorker{Power: ttlPower, Do: this_.ttl})
	apis = append(apis, &base.ApiWorker{Power: persistPower, Do: this_.persist})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *redis.Config, err error) {
	config = &redis.Config{}
	err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(redisConfig redis.Config) (res redis.Service, err error) {
	key := "redis-" + redisConfig.Address
	if redisConfig.Username != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Username)
	}
	if redisConfig.Auth != "" {
		key += "-" + util.GetMd5String(key+redisConfig.Auth)
	}
	if redisConfig.CertPath != "" {
		key += "-" + util.GetMd5String(key+redisConfig.CertPath)
	}
	var service util.Service
	service, err = util.GetService(key, func() (res util.Service, err error) {
		var s redis.Service
		s, err = redis.CreateService(redisConfig)
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}

		ctx := context.TODO()
		_, err = s.Exists(ctx, 0, "_")
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(redis.Service)
	return
}

type BaseRequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	ValueSize  int64  `json:"valueSize"`
	ValueStart int64  `json:"valueStart"`
	Pattern    string `json:"pattern"`
	Database   int    `json:"database"`
	Size       int64  `json:"size"`
	DoType     string `json:"doType"`
	Index      int64  `json:"index"`
	Count      int64  `json:"count"`
	Field      string `json:"field"`
	TaskKey    string `json:"taskKey,omitempty"`
	Expire     int64  `json:"expire"`
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	ctx := context.TODO()
	res, err = service.Info(ctx)
	if err != nil {
		return
	}
	return
}

func (this_ *api) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.Get(ctx, request.Database, request.Key, request.ValueStart, request.ValueSize)
	if err != nil {
		return
	}
	return
}

func (this_ *api) keys(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.Keys(ctx, request.Database, request.Pattern, request.Size)
	if err != nil {
		return
	}
	return
}

func (this_ *api) set(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.Set(ctx, request.Database, request.Key, request.Value)
	if err != nil {
		return
	}
	if request.Expire > 0 {
		_, err = service.Expire(ctx, request.Database, request.Key, request.Expire)
	}
	return
}

func (this_ *api) sadd(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.SAdd(ctx, request.Database, request.Key, request.Value)
	return
}

func (this_ *api) srem(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.SRem(ctx, request.Database, request.Key, request.Value)
	return
}
func (this_ *api) lpush(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.LPush(ctx, request.Database, request.Key, request.Value)
	return
}

func (this_ *api) rpush(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.RPush(ctx, request.Database, request.Key, request.Value)
	return
}

func (this_ *api) lset(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.LSet(ctx, request.Database, request.Key, request.Index, request.Value)
	return
}

func (this_ *api) lrem(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.LRem(ctx, request.Database, request.Key, request.Count, request.Value)
	return
}

func (this_ *api) hset(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.HSet(ctx, request.Database, request.Key, request.Field, request.Value)
	return
}

func (this_ *api) hdel(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	err = service.HDel(ctx, request.Database, request.Key, request.Field)
	return
}

func (this_ *api) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.Del(ctx, request.Database, request.Key)
	if err != nil {
		return
	}
	return
}

func (this_ *api) deletePattern(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.DelPattern(ctx, request.Database, request.Pattern)
	if err != nil {
		return
	}
	return
}

func (this_ *api) expire(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.Expire(ctx, request.Database, request.Key, request.Expire)
	if err != nil {
		return
	}
	return
}

func (this_ *api) ttl(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.TTL(ctx, request.Database, request.Key)
	if err != nil {
		return
	}
	return
}

func (this_ *api) persist(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	ctx := context.TODO()
	res, err = service.Persist(ctx, request.Database, request.Key)
	if err != nil {
		return
	}
	return
}

func (this_ *api) _import(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
	if err != nil {
		return
	}

	request := &redis.ImportTask{}
	if !base.RequestJSON(request, c) {
		return
	}
	taskKey := util.UUID()

	request.Key = taskKey
	request.Service = service

	redis.StartImportTask(request)

	res = taskKey
	return
}

func (this_ *api) importStatus(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res = redis.GetImportTask(request.TaskKey)
	return
}

func (this_ *api) importStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	redis.StopImportTask(request.TaskKey)
	return
}

func (this_ *api) importClean(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	redis.CleanImportTask(request.TaskKey)
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
