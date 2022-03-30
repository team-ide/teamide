package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandle(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			c.Writer.Write([]byte(err.Error()))
			break
		}
		fmt.Println("client message " + string(message))
		//写入ws数据
		err = ws.WriteMessage(mt, []byte(time.Now().String()))
		if err != nil {
			break
		}
		fmt.Println("system message " + time.Now().String())
	}
}
