package web

import (
	"fmt"
	"net"
	"server/base"
	"server/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	serverContext string
	serverHost    string
	serverPort    int
)

func init() {
	serverContext = config.Config.Server.Context
	if serverContext == "" || serverContext[len(serverContext)-1:] != "/" {
		serverContext = serverContext + "/"
	}

	serverHost = config.Config.Server.Host
	serverPort = config.Config.Server.Port
}

func StartServer() {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group(serverContext)

	bindGet(routerGroup)
	bindApi(routerGroup)

	ints, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	println("服务启动，访问地址:")
	if serverHost == "0.0.0.0" || serverHost == "::" {
		httpServer := fmt.Sprint("127.0.0.1", ":", serverPort)
		println("\t", "http://"+httpServer+serverContext)
		for _, iface := range ints {
			if iface.Flags&net.FlagUp == 0 {
				continue // interface down
			}
			if iface.Flags&net.FlagLoopback != 0 {
				continue // loopback interface
			}
			addrs, err := iface.Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				ip := base.GetIpFromAddr(addr)
				if ip == nil {
					continue
				}
				httpServer := fmt.Sprint(ip, ":", serverPort)
				println("\t", "http://"+httpServer+serverContext)
			}
		}
	} else {
		httpServer := fmt.Sprint(serverHost, ":", serverPort)
		println("\t", "http://"+httpServer+serverContext)
	}
	go func() {
		httpServer := fmt.Sprint(serverHost, ":", serverPort)
		err = http.ListenAndServe(httpServer, router)
		if err != nil {
			panic(err)
		}
	}()

}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
