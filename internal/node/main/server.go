package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"regexp"
	"teamide/pkg/util"
	"time"
)

type Server struct {
	ServerUrl  string
	ServerHost string
	ServerPort int
}

func (this_ *Server) Start() (serverUrl string, err error) {

	router := gin.Default()

	router.MaxMultipartMemory = (1024 * 50) << 20 // 设置最大上传大小为50G

	routerGroup := router.Group("/")
	routerGroup.GET("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")

	})

	httpServer := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
	s := &http.Server{
		Addr:    httpServer,
		Handler: router,
	}
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return
	}
	go func() {
		err = s.Serve(ln)
		if err != nil {
			util.Logger.Error("Web启动失败", zap.Error(err))
			panic(err)
		}
	}()
	this_.ServerUrl = fmt.Sprintf("http://127.0.0.1:%d/", this_.ServerPort)

	for {
		time.Sleep(time.Millisecond * 100)
		checkURL := this_.ServerUrl
		util.Logger.Info("监听服务是否启动成功", zap.Any("checkURL", checkURL))
		res, _ := http.Get(checkURL)
		if res.StatusCode == 200 {
			_ = res.Body.Close()
			util.Logger.Info("服务启动成功", zap.Any("ServerUrl", this_.ServerUrl))
			break
		}
		util.Logger.Info("服务未启动完成", zap.Any("StatusCode", res.StatusCode))
	}
	serverUrl = this_.ServerUrl
	return
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
