package module_http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"io"
	"net/http"
	"net/url"
	"os"
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
	Power          = base.AppendPower(&base.PowerAction{Action: "http", Text: "HTTP", ShouldLogin: true, StandAlone: true})
	execute        = base.AppendPower(&base.PowerAction{Action: "execute", Text: "执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	history        = base.AppendPower(&base.PowerAction{Action: "history", Text: "历史执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	getExecute     = base.AppendPower(&base.PowerAction{Action: "getExecute", Text: "获取执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteExecute  = base.AppendPower(&base.PowerAction{Action: "deleteExecute", Text: "获取执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	getExecuteFile = base.AppendPower(&base.PowerAction{Action: "getExecuteFile", Text: "获取执行文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	close_         = base.AppendPower(&base.PowerAction{Action: "close", Text: "关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: execute, Do: this_.execute})
	apis = append(apis, &base.ApiWorker{Power: history, Do: this_.history})
	apis = append(apis, &base.ApiWorker{Power: getExecute, Do: this_.getExecute})
	apis = append(apis, &base.ApiWorker{Power: getExecuteFile, Do: this_.getExecuteFile})
	apis = append(apis, &base.ApiWorker{Power: deleteExecute, Do: this_.deleteExecute})
	apis = append(apis, &base.ApiWorker{Power: close_, Do: this_.close})

	return
}

func (this_ *api) execute(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &Request{}
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
	request.ExecuteId = util.GetUUID()
	request.dir = dir + "" + request.ExecuteId + "/"
	request.extend = extend
	request.toolboxService = this_.toolboxService

	res, err = this_.Execute(request)
	return
}

func (this_ *api) history(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	dir := this_.getRequestDir(request.ToolboxId)

	if e, _ := util.PathExists(dir); !e {
		return
	}

	fs, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var list []map[string]any
	for _, f := range fs {
		if f.IsDir() {
			bs, e := util.ReadFile(dir + "" + f.Name() + "/execute.json")
			if e != nil {
				continue
			}
			data := map[string]any{}
			_ = json.Unmarshal(bs, &data)
			list = append(list, data)
		}
	}
	res = list
	return
}

func (this_ *api) getExecute(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	dir := this_.getRequestDir(request.ToolboxId)

	bs, err := util.ReadFile(dir + "" + request.ExecuteId + "/execute.json")
	if err != nil {
		return
	}
	data := map[string]any{}
	_ = json.Unmarshal(bs, &data)
	res = data
	return
}

func (this_ *api) getExecuteFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := map[string]string{}

	err = c.Bind(&request)
	if err != nil {
		return
	}
	if request["isDownload"] == "1" || request["isDownload"] == "true" {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape("下载")))
	}
	dir := this_.getRequestDir(util.StringToInt64(request["toolboxId"]))
	dir = dir + "" + request["executeId"] + "/"

	bs, err := util.ReadFile(dir + "execute.json")
	if err != nil {
		return
	}
	data := &Execute{}
	_ = json.Unmarshal(bs, &data)
	if data.Response != nil && data.Response.FileName != "" {
		res = base.HttpNotResponse
		if request["isDownload"] == "1" || request["isDownload"] == "true" {
			c.Header("Content-Type", "application/octet-stream")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Content-Disposition", "")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(data.Response.FileName)))
		} else {
			c.Header("Content-Type", data.Response.ContentType)
			c.Header("Content-Disposition", data.Response.Header.Get("Content-Disposition"))
		}
		f, _ := os.Open(dir + data.Response.FileName)
		if f != nil {
			defer func() { _ = f.Close() }()
			_, _ = io.Copy(c.Writer, f)
		}
		c.Status(http.StatusOK)
	}
	return
}

func (this_ *api) deleteExecute(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	dir := this_.getRequestDir(request.ToolboxId)

	if e, _ := util.PathExists(dir + request.ExecuteId); !e {
		return
	}

	if e, _ := util.IsSubPath(dir, dir+request.ExecuteId); !e {
		return
	}
	err = os.RemoveAll(dir + request.ExecuteId)

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
