package module_toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/pkg/toolbox"
)

// Work 执行
func (this_ *ToolboxService) Work(toolboxId int64, work string, data map[string]interface{}) (res interface{}, err error) {

	toolboxData, err := this_.Get(toolboxId)
	if err != nil {
		return
	}
	if toolboxData == nil {
		err = errors.New("工具配置丢失")
		return
	}

	option := map[string]interface{}{}
	if toolboxData.Option != "" {
		err = json.Unmarshal([]byte(toolboxData.Option), &option)
		if err != nil {
			return
		}
	}

	if len(option) == 0 {
		err = errors.New("工具未设置配置")
		return
	}

	toolboxWorker := toolbox.GetWorker(toolboxData.ToolboxType)
	if toolboxWorker == nil {
		err = errors.New("不支持的工具类型[" + toolboxData.ToolboxType + "]")
		return
	}

	res, err = toolboxWorker.Work(work, option, data)
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

func SSHConnection(c *gin.Context) (err error) {

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
	toolbox.WSSSHConnection(token, ws)
	return
}

func SFTPConnection(c *gin.Context) (err error) {

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
	err = toolbox.WSSFPTConnection(token, ws)
	if err != nil {
		ws.Close()
		return
	}
	return
}

func SFTPUpload(c *gin.Context) (res interface{}, err error) {

	res, err = toolbox.SFTPUpload(c)
	return
}

func SFTPDownload(data map[string]string, c *gin.Context) (err error) {
	err = toolbox.SFTPDownload(data, c)
	return
}
