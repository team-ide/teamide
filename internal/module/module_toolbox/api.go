package module_toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/pkg/toolbox"
)

type ToolboxApi struct {
	*context.ServerContext
	ToolboxService *ToolboxService
}

func NewToolboxApi(ToolboxService *ToolboxService) *ToolboxApi {

	return &ToolboxApi{
		ServerContext:  ToolboxService.ServerContext,
		ToolboxService: ToolboxService,
	}
}

var (
	// 工具 权限

	// PowerToolbox 工具基本 权限
	PowerToolbox               = base.AppendPower(&base.PowerAction{Action: "toolbox", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerToolboxPage           = base.AppendPower(&base.PowerAction{Action: "toolbox_page", Text: "工具页面", Parent: PowerToolbox, ShouldLogin: true, StandAlone: true})
	PowerToolboxContext        = base.AppendPower(&base.PowerAction{Action: "toolbox_context", Text: "工具上下文", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxInsert         = base.AppendPower(&base.PowerAction{Action: "toolbox_insert", Text: "工具新增", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxUpdate         = base.AppendPower(&base.PowerAction{Action: "toolbox_update", Text: "工具修改", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxRename         = base.AppendPower(&base.PowerAction{Action: "toolbox_rename", Text: "工具重命名", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxDelete         = base.AppendPower(&base.PowerAction{Action: "toolbox_delete", Text: "工具删除", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxWork           = base.AppendPower(&base.PowerAction{Action: "toolbox_work", Text: "工具工作", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHShell       = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_shell", Text: "工具SSH Shell连接", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtp         = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp", Text: "工具SSH FTP连接", Parent: PowerToolboxPage, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtpUpload   = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp_upload", Text: "工具SSH FTP上传", Parent: PowerToolboxSSHFtp, ShouldLogin: true, StandAlone: true})
	PowerToolboxSSHFtpDownload = base.AppendPower(&base.PowerAction{Action: "toolbox_ssh_ftp_download", Text: "工具SSH FTP下载", Parent: PowerToolboxSSHFtp, ShouldLogin: true, StandAlone: true})
)

func (this_ *ToolboxApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox"}, Power: PowerToolbox, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/page"}, Power: PowerToolboxPage, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/context"}, Power: PowerToolboxContext, Do: this_.context})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/insert"}, Power: PowerToolboxInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/update"}, Power: PowerToolboxUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/rename"}, Power: PowerToolboxRename, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/delete"}, Power: PowerToolboxDelete, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/work"}, Power: PowerToolboxWork, Do: this_.work})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/shell"}, Power: PowerToolboxSSHShell, Do: this_.sshShell, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp"}, Power: PowerToolboxSSHFtp, Do: this_.sshFtp, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp/upload"}, Power: PowerToolboxSSHFtpUpload, Do: this_.sshFtpUpload})
	apis = append(apis, &base.ApiWorker{Apis: []string{"toolbox/ssh/ftp/download"}, Power: PowerToolboxSSHFtpDownload, Do: this_.sshFtpDownload, IsGet: true})

	return
}

func (this_ *ToolboxApi) index(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *ToolboxApi) page(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	return
}

type ContextRequest struct {
}

type ContextResponse struct {
	Context map[string][]*ToolboxModel `json:"context,omitempty"`
}

func (this_ *ToolboxApi) context(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ContextResponse{}

	list, err := this_.ToolboxService.Query(&ToolboxModel{
		UserId: requestBean.JWT.UserId,
	})
	if err != nil {
		return
	}
	context := make(map[string][]*ToolboxModel)
	for _, one := range list {
		context[one.ToolboxType] = append(context[one.ToolboxType], one)
	}
	response.Context = context
	res = response
	return
}

type InsertRequest struct {
	*ToolboxModel
}

type InsertResponse struct {
}

func (this_ *ToolboxApi) insert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &InsertResponse{}

	toolbox := request.ToolboxModel
	toolbox.UserId = requestBean.JWT.UserId

	_, err = this_.ToolboxService.Insert(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateRequest struct {
	*ToolboxModel
}

type UpdateResponse struct {
}

func (this_ *ToolboxApi) update(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Update(toolbox)
	if err != nil {
		return
	}

	res = response
	return
}

type RenameRequest struct {
	*ToolboxModel
}

type RenameResponse struct {
}

func (this_ *ToolboxApi) rename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &RenameResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Rename(toolbox.ToolboxId, toolbox.Name)
	if err != nil {
		return
	}

	res = response
	return
}

type DeleteRequest struct {
	*ToolboxModel
}

type DeleteResponse struct {
}

func (this_ *ToolboxApi) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &RenameResponse{}

	toolbox := request.ToolboxModel

	_, err = this_.ToolboxService.Delete(toolbox.ToolboxId)
	if err != nil {
		return
	}

	res = response
	return
}

type WorkRequest struct {
	ToolboxId int64                  `json:"toolboxId,omitempty"`
	Work      string                 `json:"work,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func (this_ *ToolboxApi) work(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &WorkRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.Work(request.ToolboxId, request.Work, request.Data)
	if err != nil {
		return
	}

	return
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this_ *ToolboxApi) sshShell(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	token := c.Query("token")
	//fmt.Println("token=" + token)
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	err = toolbox.WSSSHConnection(token, ws, this_.Logger)
	if err != nil {
		ws.Close()
		return
	}
	res = base.HttpNotResponse
	return
}

func (this_ *ToolboxApi) sshFtp(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	token := c.Query("token")
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	err = toolbox.WSSFPTConnection(token, ws, this_.Logger)
	if err != nil {
		ws.Close()
		return
	}
	res = base.HttpNotResponse

	return
}

func (this_ *ToolboxApi) sshFtpUpload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res, err = toolbox.SFTPUpload(c)
	if err != nil {
		return
	}
	return
}

func (this_ *ToolboxApi) sshFtpDownload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}
	err = toolbox.SFTPDownload(data, c)
	if err != nil {
		return
	}
	res = base.HttpNotResponse
	return
}
