package module_toolbox

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/internal/base"
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

func (this_ *ToolboxApi) databaseExportDownload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}
	//err = task.DatabaseExportDownload(data, c)
	//if err != nil {
	//	return
	//}
	res = base.HttpNotResponse
	return
}
