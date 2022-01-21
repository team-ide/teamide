package web

import (
	"fmt"
	"net"
	"strings"
	"teamide/server/base"
	"teamide/server/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ServerContext string
	ServerHost    string
	ServerPort    int
	ServerUrl     string
)

func init() {
	ServerContext = config.Config.Server.Context
	if ServerContext == "" || !strings.HasSuffix(ServerContext, "/") {
		ServerContext = ServerContext + "/"
	}

	ServerHost = config.Config.Server.Host
	ServerPort = config.Config.Server.Port

	if base.IsLocalStartup {
		if ServerHost == "" {
			ServerHost = "127.0.0.1"
		}
		if ServerPort == 0 {
			listener, err := net.Listen("tcp", ":0")
			if err != nil {
				panic(err)
			}
			ServerPort = listener.Addr().(*net.TCPAddr).Port
		}
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
	if !base.IsLocalStartup {
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
					ip := base.GetIpFromAddr(addr)
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
