package module_toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/internal/base"
	"teamide/pkg/db/task"
	"teamide/pkg/ssh"
	"teamide/pkg/toolbox"
)

type WorkRequest struct {
	ToolboxId int64                  `json:"toolboxId,omitempty"`
	Work      string                 `json:"work,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func (this_ *ToolboxApi) work(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this_ *ToolboxApi) sshShell(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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

func (this_ *ToolboxApi) sshFtpUpload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	res, err = SFTPUpload(c)
	if err != nil {
		return
	}
	return
}

func (this_ *ToolboxApi) sshFtpDownload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}
	err = SFTPDownload(data, c)
	if err != nil {
		return
	}
	res = base.HttpNotResponse
	return
}

func SFTPUpload(c *gin.Context) (res interface{}, err error) {
	workerId := c.PostForm("workerId")
	if workerId == "" {
		err = errors.New("workerId获取失败")
		return
	}
	dir := c.PostForm("dir")
	if dir == "" {
		err = errors.New("dir获取失败")
		return
	}
	place := c.PostForm("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	workId := c.PostForm("workId")
	if workId == "" {
		err = errors.New("workId获取失败")
		return
	}
	client := toolbox.GetSftpClient(workerId)
	if client == nil {
		err = errors.New("FTP会话丢失")
		return
	}
	mF, err := c.MultipartForm()
	if err != nil {
		return
	}
	fileList := mF.File["file"]
	if err != nil {
		return
	}
	for _, file := range fileList {
		uploadFile := &ssh.UploadFile{
			Dir:      dir,
			Place:    place,
			WorkId:   workId,
			File:     file,
			FullPath: c.PostForm("fullPath"),
		}
		client.UploadFile <- uploadFile
	}

	return
}

func SFTPDownload(data map[string]string, c *gin.Context) (err error) {

	workerId := data["workerId"]
	if workerId == "" {
		err = errors.New("workerId获取失败")
		return
	}
	workId := data["workId"]
	if workId == "" {
		err = errors.New("workId获取失败")
		return
	}

	place := data["place"]
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	path := data["path"]
	if path == "" {
		err = errors.New("path获取失败")
		return
	}
	client := toolbox.GetSftpClient(workerId)
	if client == nil {
		err = errors.New("SSH会话丢失")
		return
	}
	if place == "local" {
		err = client.LocalDownload(c, workId, path)
	} else if place == "remote" {
		err = client.RemoteDownload(c, workId, path)
	}

	return
}

func (this_ *ToolboxApi) databaseExportDownload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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
