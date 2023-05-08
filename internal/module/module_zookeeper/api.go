package module_zookeeper

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	goSSH "golang.org/x/crypto/ssh"
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
	Power            = base.AppendPower(&base.PowerAction{Action: "zookeeper", Text: "Zookeeper", ShouldLogin: true, StandAlone: true})
	infoPower        = base.AppendPower(&base.PowerAction{Action: "info", Text: "Zookeeper信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getPower         = base.AppendPower(&base.PowerAction{Action: "get", Text: "Zookeeper获取节点数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	savePower        = base.AppendPower(&base.PowerAction{Action: "save", Text: "Zookeeper保存节点数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	getChildrenPower = base.AppendPower(&base.PowerAction{Action: "getChildren", Text: "Zookeeper查询子节点", ShouldLogin: true, StandAlone: true, Parent: Power})
	deletePower      = base.AppendPower(&base.PowerAction{Action: "delete", Text: "Zookeeper删除节点", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower       = base.AppendPower(&base.PowerAction{Action: "close", Text: "Zookeeper关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: savePower, Do: this_.save})
	apis = append(apis, &base.ApiWorker{Power: getChildrenPower, Do: this_.getChildren})
	apis = append(apis, &base.ApiWorker{Power: deletePower, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *zookeeper.Config, sshConfig *ssh.Config, err error) {
	config = &zookeeper.Config{}
	sshConfig, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(zkConfig zookeeper.Config, sshConfig *ssh.Config) (res zookeeper.IService, err error) {
	key := "zookeeper-" + zkConfig.Address
	if zkConfig.Username != "" {
		key += "-" + base.GetMd5String(key+zkConfig.Username)
	}
	if zkConfig.Password != "" {
		key += "-" + base.GetMd5String(key+zkConfig.Password)
	}
	if sshConfig != nil {
		key += "-ssh-" + sshConfig.Address
		key += "-ssh-" + sshConfig.Username
	}
	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s zookeeper.IService
		if sshConfig != nil {
			var sshClient *goSSH.Client
			sshClient, err = ssh.NewClient(*sshConfig)
			if err != nil {
				util.Logger.Error("getZKService ssh NewClient error", zap.Any("key", key), zap.Error(err))
				return
			}
			s, err = zookeeper.NewForSSH(zkConfig, sshClient)
		} else {
			s, err = zookeeper.New(zkConfig)
		}
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
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowTime(),
			Service:     s,
			Stop:        s.Stop,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(zookeeper.IService)
	serviceInfo.SetLastUseTime()
	return
}

type BaseRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config, sshConfig)
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
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config, sshConfig)
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
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config, sshConfig)
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
		err = service.Set(request.Path, []byte(request.Data))
	} else {
		err = service.CreateIfNotExists(request.Path, []byte(request.Data))
	}
	if err != nil {
		return
	}
	return
}

func (this_ *api) getChildren(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config, sshConfig)
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
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config, sshConfig)
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
