package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"teamide/internal/context"
	"teamide/pkg/util"
)

type Server struct {
	*context.ServerContext
}

func NewWebServer(ServerContext *context.ServerContext) (webServer *Server) {
	webServer = &Server{
		ServerContext: ServerContext,
	}
	return
}

func (this_ *Server) Start() (serverUrl string, err error) {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group(this_.ServerContext.ServerContext)

	this_.bindGet(routerGroup)
	err = this_.bindApi(routerGroup)
	if err != nil {
		return
	}

	var ins []net.Interface
	ins, err = net.Interfaces()
	if err != nil {
		return
	}
	if !this_.IsStandAlone {
		println("服务启动，访问地址:")
		if this_.ServerHost == "0.0.0.0" || this_.ServerHost == "::" {
			httpServer := fmt.Sprint("127.0.0.1", ":", this_.ServerPort)
			println("\t", "http://"+httpServer+this_.ServerContext.ServerContext)
			for _, iface := range ins {
				if iface.Flags&net.FlagUp == 0 {
					continue
				}
				if iface.Flags&net.FlagLoopback != 0 {
					continue
				}
				var adders []net.Addr
				adders, err = iface.Addrs()
				if err != nil {
					return
				}
				for _, addr := range adders {
					ip := util.GetIpFromAddr(addr)
					if ip == nil {
						continue
					}
					httpServer := fmt.Sprint(ip, ":", this_.ServerPort)
					println("\t", "http://"+httpServer+this_.ServerContext.ServerContext)
				}
			}
		} else {
			httpServer := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
			println("\t", "http://"+httpServer+this_.ServerContext.ServerContext)
		}
	}

	go func() {
		httpServer := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
		err = http.ListenAndServe(httpServer, router)
		if err != nil {
			this_.Logger.Error("Web启动失败", zap.Error(err))
			panic(err)
		}
	}()
	serverUrl = this_.ServerUrl
	return serverUrl, err
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
