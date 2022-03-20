package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"teamide/internal/config"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	"teamide/pkg/db"
	"teamide/pkg/util"
)

type Server struct {
	ServerContext  string
	ServerHost     string
	ServerPort     int
	ServerUrl      string
	DatabaseWorker db.DatabaseWorker
}

func NewWebServer(serverConfig *config.ServerConfig) (webServer *Server, err error) {
	webServer = &Server{}
	err = webServer.init(serverConfig)
	if err != nil {
		return nil, err
	}
	return
}

func (this_ *Server) init(serverConfig *config.ServerConfig) (err error) {
	if serverConfig.Server != nil {
		this_.ServerContext = serverConfig.Server.Context
		if this_.ServerContext == "" || !strings.HasSuffix(this_.ServerContext, "/") {
			this_.ServerContext = this_.ServerContext + "/"
		}
		this_.ServerHost = serverConfig.Server.Host
		this_.ServerPort = serverConfig.Server.Port
	}

	if this_.ServerHost == "" {
		this_.ServerHost = "127.0.0.1"
	}

	if this_.ServerPort == 0 {
		var listener net.Listener
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			component.Logger.Error(component.LogStr("随机端口获取失败:", err))
			return
		}
		this_.ServerPort = listener.Addr().(*net.TCPAddr).Port
	}

	if this_.ServerHost == "0.0.0.0" || this_.ServerHost == ":" || this_.ServerHost == "::" {
		this_.ServerUrl = fmt.Sprint("http://127.0.0.1:", this_.ServerPort)
	} else {
		this_.ServerUrl = fmt.Sprint("http://", this_.ServerHost, ":", this_.ServerPort)
	}
	this_.DatabaseWorker, err = db.NewDatabaseWorker(*serverConfig.DatabaseConfig)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) Start() (serverUrl string, err error) {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group(this_.ServerContext)

	this_.bindGet(routerGroup)
	this_.bindApi(routerGroup)

	var ints []net.Interface
	ints, err = net.Interfaces()
	if err != nil {
		return
	}
	if !base2.IsStandAlone {
		println("服务启动，访问地址:")
		if this_.ServerHost == "0.0.0.0" || this_.ServerHost == "::" {
			httpServer := fmt.Sprint("127.0.0.1", ":", this_.ServerPort)
			println("\t", "http://"+httpServer+this_.ServerContext)
			for _, iface := range ints {
				if iface.Flags&net.FlagUp == 0 {
					continue // interface down
				}
				if iface.Flags&net.FlagLoopback != 0 {
					continue // loopback interface
				}
				var addrs []net.Addr
				addrs, err = iface.Addrs()
				if err != nil {
					return
				}
				for _, addr := range addrs {
					ip := util.GetIpFromAddr(addr)
					if ip == nil {
						continue
					}
					httpServer := fmt.Sprint(ip, ":", this_.ServerPort)
					println("\t", "http://"+httpServer+this_.ServerContext)
				}
			}
		} else {
			httpServer := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
			println("\t", "http://"+httpServer+this_.ServerContext)
		}
	}

	go func() {
		httpServer := fmt.Sprint(this_.ServerHost, ":", this_.ServerPort)
		err = http.ListenAndServe(httpServer, router)
		if err != nil {
			component.Logger.Error(component.LogStr("Web启动失败:", err))
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
