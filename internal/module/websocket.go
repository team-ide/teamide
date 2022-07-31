package module

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"teamide/internal/base"
	"teamide/internal/context"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this_ *Api) apiWebsocket(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	id := c.Query("id")
	if id == "" {
		err = errors.New("id获取失败")
		return
	}
	if request.JWT == nil || request.JWT.UserId == 0 {
		err = errors.New("登录用户获取失败")
		return
	}
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	context.AddServerWebsocket(id, fmt.Sprint(request.JWT.UserId), ws)

	res = base.HttpNotResponse
	return
}
