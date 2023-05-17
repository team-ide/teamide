package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"net"
	"net/http"
	"teamide/internal/context"
	"teamide/internal/module"
	"teamide/internal/module/module_toolbox"
	"time"
)

type Server struct {
	*context.ServerContext
	api            *module.Api
	toolboxService *module_toolbox.ToolboxService
}

func NewWebServer(ServerContext *context.ServerContext, api *module.Api) (webServer *Server) {
	webServer = &Server{
		api:            api,
		ServerContext:  ServerContext,
		toolboxService: module_toolbox.NewToolboxService(ServerContext),
	}
	return
}

func (this_ *Server) Start() (serverUrl string, err error) {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	router.MaxMultipartMemory = (1024 * 50) << 20 // 设置最大上传大小为50G

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
		s := "http"
		if this_.ServerHost == "0.0.0.0" || this_.ServerHost == "::" {
			fmt.Printf("\t%s://127.0.0.1:%d%s\n", s, this_.ServerPort, this_.ServerContext.ServerContext)
			for _, in := range ins {
				if in.Flags&net.FlagUp == 0 {
					continue
				}
				if in.Flags&net.FlagLoopback != 0 {
					continue
				}
				var adders []net.Addr
				adders, err = in.Addrs()
				if err != nil {
					return
				}
				for _, addr := range adders {
					ip := util.GetIpFromAddr(addr)
					if ip == nil {
						continue
					}
					fmt.Printf("\t%s://%s:%d%s\n", s, ip.String(), this_.ServerPort, this_.ServerContext.ServerContext)
				}
			}
		} else {
			fmt.Printf("\t%s://%s:%d%s\n", s, this_.ServerHost, this_.ServerPort, this_.ServerContext.ServerContext)
		}
	}
	addr := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
	this_.Logger.Info("http server start", zap.Any("addr", addr))
	s := &http.Server{
		Addr:    addr,
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
	var checkStartTime = util.GetNowMilli()
	for {
		var newTime = util.GetNowMilli()
		if (newTime - checkStartTime) > 1000*5 {
			this_.Logger.Warn("服务启动检查超过5秒，不再检测")
			break
		}
		time.Sleep(time.Millisecond * 100)
		checkURL := this_.ServerUrl + this_.ServerContext.ServerContext
		this_.Logger.Info("监听服务是否启动成功", zap.Any("checkURL", checkURL))
		res, e := http.Get(checkURL)
		if e != nil {
			this_.Logger.Warn("监听服务连接失败，将继续监听", zap.Any("checkURL", checkURL), zap.Any("error", e.Error()))
			continue
		}
		if res == nil {
			this_.Logger.Warn("监听服务连接无返回，不再监听", zap.Any("checkURL", checkURL))
			continue
		}
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
