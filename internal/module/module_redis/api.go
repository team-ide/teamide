package module_redis

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	goSSH "golang.org/x/crypto/ssh"
	"sort"
	"strings"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
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
	check              = base.AppendPower(&base.PowerAction{Action: "check", Text: "Redis测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	infoPower          = base.AppendPower(&base.PowerAction{Action: "info", Text: "Redis信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getPower           = base.AppendPower(&base.PowerAction{Action: "get", Text: "Redis获取Key值", ShouldLogin: true, StandAlone: true, Parent: Power})
	keysPower          = base.AppendPower(&base.PowerAction{Action: "keys", Text: "Redis查询Keys", ShouldLogin: true, StandAlone: true, Parent: Power})
	scanPower          = base.AppendPower(&base.PowerAction{Action: "scan", Text: "Redis Scan", ShouldLogin: true, StandAlone: true, Parent: Power})
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
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: keysPower, Do: this_.keys})
	apis = append(apis, &base.ApiWorker{Power: scanPower, Do: this_.scan})
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

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *redis.Config, sshConfig *ssh.Config, err error) {
	config = &redis.Config{}
	sshConfig, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getServiceKey(redisConfig *redis.Config, sshConfig *ssh.Config) (key string) {
	key = "redis-" + redisConfig.Address
	if redisConfig.Username != "" {
		key += "-" + base.GetMd5String(key+redisConfig.Username)
	}
	if redisConfig.Auth != "" {
		key += "-" + base.GetMd5String(key+redisConfig.Auth)
	}
	if redisConfig.CertPath != "" {
		key += "-" + base.GetMd5String(key+redisConfig.CertPath)
	}
	if sshConfig != nil {
		key += "-ssh-" + sshConfig.Address
		key += "-ssh-" + sshConfig.Username
	}
	return
}

func getService(redisConfig *redis.Config, sshConfig *ssh.Config) (res redis.IService, err error) {
	key := getServiceKey(redisConfig, sshConfig)
	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s redis.IService
		if sshConfig != nil {
			var sshClient *goSSH.Client
			sshClient, err = ssh.NewClient(*sshConfig)
			if err != nil {
				util.Logger.Error("getRedisService ssh NewClient error", zap.Any("key", key), zap.Error(err))
				return
			}
			redisConfig.SSHClient = sshClient
		}
		s, err = redis.New(redisConfig)
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}

		_, err = s.Exists("_", redis.NewParam())
		if err != nil {
			util.Logger.Error("getRedisService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowMilli(),
			Service:     s,
			Stop:        s.Close,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(redis.IService)
	serviceInfo.SetLastUseTime()

	return
}

type BaseRequest struct {
	Key        string `json:"key"`
	KeyBase64  string `json:"keyBase64"`
	Value      string `json:"value"`
	ValueSize  int64  `json:"valueSize"`
	ValueStart int    `json:"valueStart"`
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

func (this_ *BaseRequest) getKey() string {
	if this_.KeyBase64 != "" {
		bs, e := base64.StdEncoding.DecodeString(this_.KeyBase64)
		if e == nil && len(bs) > 0 {
			return string(bs)
		}
	}
	return this_.Key
}

func (this_ *api) check(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	_, err = getService(config, sshConfig)
	if err != nil {
		return
	}

	return
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	res, err = service.Info(&redis.Param{})
	if err != nil {
		return
	}
	return
}

func (this_ *api) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = service.GetValueInfo(request.getKey(), redis.NewStartArg(request.ValueStart), redis.NewSizeArg(request.ValueSize), &redis.Param{Database: request.Database})
	if err != nil {
		return
	}

	return
}

type KeysResult struct {
	Database int        `json:"database"`
	Count    int64      `json:"count"`
	KeyList  []*KeyInfo `json:"keyList"`
}
type KeyInfo struct {
	Key       string `json:"key"`
	KeyBase64 string `json:"keyBase64"`
}

func (this_ *api) keys(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var size *redis.SizeArg = nil
	if request.Size > 0 {
		size = redis.NewSizeArg(request.Size)
	}

	kR, err := service.Keys(request.Pattern, size, &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	keysResult := &KeysResult{
		Count:    kR.Count,
		Database: kR.Database,
	}
	res = keysResult
	for _, one := range kR.KeyList {
		keyInfo := &KeyInfo{
			Key: one,
		}
		if keyInfo.Key != "" {
			keyInfo.KeyBase64 = base64.StdEncoding.EncodeToString([]byte(keyInfo.Key))
		}
		keysResult.KeyList = append(keysResult.KeyList, keyInfo)
	}
	return
}

func (this_ *api) scan(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	var size *redis.SizeArg = nil
	if request.Size > 0 {
		size = redis.NewSizeArg(request.Size)
	}
	var count *redis.CountArg = nil
	if request.Count > 0 {
		count = redis.NewCountArg(request.Count)
	}
	kR, err := service.Scan(request.Pattern, size, count, &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	keysResult := &KeysResult{
		Count:    kR.Count,
		Database: kR.Database,
	}
	res = keysResult
	sort.Slice(kR.KeyList, func(i, j int) bool {
		return strings.ToLower(kR.KeyList[i]) < strings.ToLower(kR.KeyList[j]) //升序  即前面的值比后面的小 忽略大小写排序
	})
	for _, one := range kR.KeyList {
		keyInfo := &KeyInfo{
			Key: one,
		}
		if keyInfo.Key != "" {
			keyInfo.KeyBase64 = base64.StdEncoding.EncodeToString([]byte(keyInfo.Key))
		}
		keysResult.KeyList = append(keysResult.KeyList, keyInfo)
	}
	return
}

func (this_ *api) set(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.Set(request.getKey(), request.Value, &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	if request.Expire > 0 {
		_, err = service.Expire(request.getKey(), request.Expire, &redis.Param{Database: request.Database})
	}
	return
}

func (this_ *api) sadd(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.SetAdd(request.getKey(), request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) srem(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.SetRem(request.getKey(), request.Value, &redis.Param{Database: request.Database})
	return
}
func (this_ *api) lpush(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.ListPush(request.getKey(), request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) rpush(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.ListRPush(request.getKey(), request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) lset(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.ListSet(request.getKey(), request.Index, request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) lrem(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.ListRem(request.getKey(), request.Count, request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) hset(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.HashSet(request.getKey(), request.Field, request.Value, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) hdel(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.HashDel(request.getKey(), request.Field, &redis.Param{Database: request.Database})
	return
}

func (this_ *api) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Del(request.getKey(), &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	return
}

func (this_ *api) deletePattern(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.DelPattern(request.Pattern, &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	return
}

func (this_ *api) expire(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Expire(request.getKey(), request.Expire, &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	return
}

func (this_ *api) ttl(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Ttl(request.getKey(), &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	return
}

func (this_ *api) persist(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Persist(request.getKey(), &redis.Param{Database: request.Database})
	if err != nil {
		return
	}
	return
}

func (this_ *api) _import(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	//config,sshConfig, err := this_.getConfig(requestBean, c)
	//if err != nil {
	//	return
	//}
	//service, err := getService(config,sshConfig)
	//if err != nil {
	//	return
	//}
	//
	//request := &redis.ImportTask{}
	//if !base.RequestJSON(request, c) {
	//	return
	//}
	//taskKey := util.GetUUID()
	//
	//request.Key = taskKey
	//request.Service = service
	//
	//redis.StartImportTask(request)

	//res = taskKey
	return
}

func (this_ *api) importStatus(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	//request := &BaseRequest{}
	//if !base.RequestJSON(request, c) {
	//	return
	//}
	//res = redis.GetImportTask(request.TaskKey)
	return
}

func (this_ *api) importStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	//request := &BaseRequest{}
	//if !base.RequestJSON(request, c) {
	//	return
	//}
	//redis.StopImportTask(request.TaskKey)
	return
}

func (this_ *api) importClean(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	//request := &BaseRequest{}
	//if !base.RequestJSON(request, c) {
	//	return
	//}
	//redis.CleanImportTask(request.TaskKey)
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
