package module_zookeeper

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/util"
	"teamide/pkg/zookeeper"
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
	Power            = base.AppendPower(&base.PowerAction{Action: "zookeeper", Text: "Zookeeper", ShouldLogin: true, StandAlone: true})
	infoPower        = base.AppendPower(&base.PowerAction{Action: "info", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getPower         = base.AppendPower(&base.PowerAction{Action: "get", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	savePower        = base.AppendPower(&base.PowerAction{Action: "save", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getChildrenPower = base.AppendPower(&base.PowerAction{Action: "getChildren", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deletePower      = base.AppendPower(&base.PowerAction{Action: "delete", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower       = base.AppendPower(&base.PowerAction{Action: "close", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/info"}, Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/get"}, Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/save"}, Power: savePower, Do: this_.save})
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/getChildren"}, Power: getChildrenPower, Do: this_.getChildren})
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/delete"}, Power: deletePower, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Apis: []string{"zookeeper/close"}, Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *zookeeper.Config, err error) {
	config = &zookeeper.Config{}
	err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(zkConfig zookeeper.Config) (res *zookeeper.ZKService, err error) {
	key := "zookeeper-" + zkConfig.Address
	if zkConfig.Username != "" {
		key += "-" + util.GetMd5String(key+zkConfig.Username)
	}
	if zkConfig.Password != "" {
		key += "-" + util.GetMd5String(key+zkConfig.Password)
	}
	var service util.Service
	service, err = util.GetService(key, func() (res util.Service, err error) {
		var s *zookeeper.ZKService
		s, err = zookeeper.CreateZKService(zkConfig)
		if err != nil {
			util.Logger.Error("getZKService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		_, err = s.Exists("/")
		if err != nil {
			util.Logger.Error("getZKService error", zap.Any("key", key), zap.Error(err))
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
	res = service.(*zookeeper.ZKService)
	return
}

type BaseRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
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

	res, err = service.Info()
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
	res, err = service.Get(request.Path)
	if err != nil {
		return
	}
	return
}

func (this_ *api) save(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
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
	var isEx bool
	isEx, err = service.Exists(request.Path)
	if err != nil {
		return
	}
	if isEx {
		err = service.SetData(request.Path, []byte(request.Data))
	} else {
		err = service.CreateIfNotExists(request.Path, []byte(request.Data))
	}
	if err != nil {
		return
	}
	return
}

func (this_ *api) getChildren(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
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
	var isEx bool
	isEx, err = service.Exists(request.Path)
	if err != nil {
		return
	}
	if isEx {

		var children []map[string]interface{}

		var names []string
		names, err = service.GetChildren(request.Path)
		if err != nil {
			return
		}
		for _, name := range names {
			var one = map[string]interface{}{}
			one["name"] = name

			childrenPath := "/" + name
			if request.Path != "/" {
				childrenPath = request.Path + childrenPath
			}
			var statInfo *zookeeper.StatInfo
			statInfo, err = service.Stat(childrenPath)
			if err != nil {
				return
			}
			if statInfo != nil {
				one["hasChildren"] = statInfo.NumChildren > 0
			}

			children = append(children, one)
		}
		res = children
	}
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
	var isEx bool
	isEx, err = service.Exists(request.Path)
	if err != nil {
		return
	}
	if isEx {
		err = service.Delete(request.Path)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	return
}
