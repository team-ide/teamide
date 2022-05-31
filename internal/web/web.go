package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"teamide/internal/context"
	"teamide/internal/module"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/util"
	"time"
)

type Server struct {
	*context.ServerContext
	api            *module.Api
	toolboxService *module_toolbox.ToolboxService
}

func NewWebServer(ServerContext *context.ServerContext) (webServer *Server) {
	webServer = &Server{
		ServerContext:  ServerContext,
		toolboxService: module_toolbox.NewToolboxService(ServerContext),
	}
	return
}

func (this_ *Server) Start() (serverUrl string, err error) {

	this_.api, err = module.NewApi(this_.ServerContext)
	if err != nil {
		return
	}

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
	if this_.IsServer {
		println("服务启动，访问地址:")
		if this_.ServerHost == "0.0.0.0" || this_.ServerHost == "::" {
			httpServer := fmt.Sprint("localhost", ":", this_.ServerPort)
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
			this_.Logger.Error("Web启动失败", zap.Error(err))
			panic(err)
		}
	}()
	for {
		time.Sleep(time.Millisecond * 100)
		checkURL := this_.ServerUrl + this_.ServerContext.ServerContext
		this_.Logger.Info("监听服务是否启动成功", zap.Any("checkURL", checkURL))
		res, _ := http.Get(checkURL)
		if res.StatusCode == 200 {
			_ = res.Body.Close()
			this_.Logger.Info("服务启动成功", zap.Any("ServerUrl", this_.ServerUrl))
			break
		}
		this_.Logger.Info("服务未启动完成", zap.Any("StatusCode", res.StatusCode))
	}
	serverUrl = this_.ServerUrl
	return serverUrl, err
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
