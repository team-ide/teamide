package module_toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/internal/base"
	"teamide/pkg/db/task"
	"teamide/pkg/ssh"
)

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
	err = ssh.WSSSHConnection(token, ws)
	if err != nil {
		_ = ws.Close()
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
	err = ssh.WSSFPTConnection(token, ws)
	if err != nil {
		_ = ws.Close()
		return
	}
	res = base.HttpNotResponse

	return
}

func (this_ *ToolboxApi) sshUpload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res, err = ssh.SFTPUpload(c)
	if err != nil {
		return
	}
	return
}

func (this_ *ToolboxApi) sshDownload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}
	err = ssh.SFTPDownload(data, c)
	if err != nil {
		return
	}
	res = base.HttpNotResponse
	return
}

func (this_ *ToolboxApi) sshFtpUpload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res, err = ssh.SFTPUpload(c)
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
	err = ssh.SFTPDownload(data, c)
	if err != nil {
		return
	}
	res = base.HttpNotResponse
	return
}

func (this_ *ToolboxApi) databaseExportDownload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}
	err = task.DatabaseExportDownload(data, c)
	if err != nil {
		return
	}
	res = base.HttpNotResponse
	return
}
