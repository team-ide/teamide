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
	ToolboxId   int64                  `json:"toolboxId,omitempty"`
	ToolboxType string                 `json:"toolboxType,omitempty"`
	Work        string                 `json:"work,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

func (this_ *ToolboxApi) work(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &WorkRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.Work(request.ToolboxId, request.ToolboxType, request.Work, request.Data)
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
