package module_http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
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
	Power   = base.AppendPower(&base.PowerAction{Action: "http", Text: "HTTP", ShouldLogin: true, StandAlone: true})
	execute = base.AppendPower(&base.PowerAction{Action: "execute", Text: "执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	close_  = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: execute, Do: this_.execute})
	apis = append(apis, &base.ApiWorker{Power: close_, Do: this_.close})

	return
}

type BaseRequest struct {
	ToolboxId int64 `json:"toolboxId,omitempty"`
	*Request
}

func (this_ *api) execute(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	extends, err := this_.toolboxService.QueryExtends(&module_toolbox.ToolboxExtendModel{
		ToolboxId:  request.ToolboxId,
		ExtendType: "http-config",
		UserId:     requestBean.JWT.UserId,
	})
	var extend = &Extend{}
	if len(extends) > 0 {
		_ = json.Unmarshal([]byte(extends[0].Value), extend)
	}
	dir := this_.getRequestDir(request.ToolboxId)
	//if e, _ := util.PathExists(dir); !e {
	//	_ = os.MkdirAll(dir, fs.ModePerm)
	//}
	request.Request.dir = dir
	request.Request.extend = extend
	request.toolboxService = this_.toolboxService

	res, err = this_.Execute(request.Request)
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

func (this_ *api) getRequestDir(toolboxId int64) (dir string) {
	dir = this_.getDir(toolboxId)
	dir += "request/"
	return
}
func (this_ *api) getDir(toolboxId int64) (dir string) {
	dir = this_.toolboxService.GetFilesDir()
	dir += fmt.Sprintf("%s/toolbox-%d/", "toolbox-http", toolboxId)
	return
}

type Extend struct {
	Secrets   []*Field `json:"secrets,omitempty"`
	Variables []*Field `json:"variables,omitempty"`
}
