package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	"teamide/internal/server/config"
)

var (
	ServerContext string = "/"
	ServerHost    string = "127.0.0.1"
	ServerPort    int    = 0
	ServerUrl     string
)

func init() {
	if config.Config.Server != nil {
		ServerContext = config.Config.Server.Context
		if ServerContext == "" || !strings.HasSuffix(ServerContext, "/") {
			ServerContext = ServerContext + "/"
		}
		ServerHost = config.Config.Server.Host
		ServerPort = config.Config.Server.Port
	}

	if ServerHost == "" {
		ServerHost = "127.0.0.1"
	}

	if ServerPort == 0 {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			component.Logger.Error(component.LogStr("随机端口获取失败:", err))
			panic(err)
		}
		ServerPort = listener.Addr().(*net.TCPAddr).Port
	}

	if ServerHost == "0.0.0.0" || ServerHost == ":" || ServerHost == "::" {
		ServerUrl = fmt.Sprint("http://127.0.0.1:", ServerPort)
	} else {
		ServerUrl = fmt.Sprint("http://", ServerHost, ":", ServerPort)
	}
}

func StartServer() (serverUrl string, err error) {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group(ServerContext)

	bindGet(routerGroup)
	bindApi(routerGroup)

	var ints []net.Interface
	ints, err = net.Interfaces()
	if err != nil {
		return
	}
	if !base2.IS_STAND_ALONE {
		println("服务启动，访问地址:")
		if ServerHost == "0.0.0.0" || ServerHost == "::" {
			httpServer := fmt.Sprint("127.0.0.1", ":", ServerPort)
			println("\t", "http://"+httpServer+ServerContext)
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
					ip := base2.GetIpFromAddr(addr)
					if ip == nil {
						continue
					}
					httpServer := fmt.Sprint(ip, ":", ServerPort)
					println("\t", "http://"+httpServer+ServerContext)
				}
			}
		} else {
			httpServer := fmt.Sprint(ServerHost, ":", ServerPort)
			println("\t", "http://"+httpServer+ServerContext)
		}
	}

	go func() {
		httpServer := fmt.Sprint(ServerHost, ":", ServerPort)
		err = http.ListenAndServe(httpServer, router)
		if err != nil {
			component.Logger.Error(component.LogStr("Web启动失败:", err))
			panic(err)
		}
	}()
	serverUrl = ServerUrl
	return serverUrl, err
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
