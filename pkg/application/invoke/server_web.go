package invoke

import (
	"fmt"
	"net"
	"net/http"
	"teamide/pkg/application/base"
	"teamide/pkg/application/common"
	"teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

func startServerWeb(app common.IApplication, serverWeb *model.ServerWebModel) (err error) {
	contextPath := serverWeb.ContextPath
	if contextPath == "" || contextPath[len(contextPath)-1:] != "/" {
		contextPath = contextPath + "/"
	}

	serverHost := serverWeb.Host
	serverPort := serverWeb.Port

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group(contextPath)

	err = ServerWebBindApis(app, serverWeb.Token, routerGroup)
	if err != nil {
		return
	}
	var ints []net.Interface
	ints, err = net.Interfaces()
	if err != nil {
		return
	}

	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("服务启动，访问地址:")
	}
	if serverHost == "0.0.0.0" || serverHost == "::" {

		httpServer := fmt.Sprint("127.0.0.1", ":", serverPort)

		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("\t", "http://"+httpServer+contextPath)
		}
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
				httpServer := fmt.Sprint(ip, ":", serverPort)

				if app.GetLogger() != nil && app.GetLogger().OutDebug() {
					app.GetLogger().Debug("\t", "http://"+httpServer+contextPath)
				}
			}
		}
	} else {
		httpServer := fmt.Sprint(serverHost, ":", serverPort)

		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("\t", "http://"+httpServer+contextPath)
		}
	}
	go func() {
		httpServer := fmt.Sprint(serverHost, ":", serverPort)
		err = http.ListenAndServe(httpServer, router)
		if err != nil {
			panic(err)
		}
	}()
	return
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
