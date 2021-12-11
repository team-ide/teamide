package web

import (
	"fmt"
	"net"
	"server/base"
	"server/config"
	"server/enterpriseService"
	"server/groupService"
	"server/idService"
	"server/jobService"
	"server/logService"
	"server/messageService"
	"server/powerService"
	"server/settingService"
	"server/spaceService"
	"server/systemService"
	"server/userService"
	"server/wbsService"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func handleIndex(rw http.ResponseWriter, r *http.Request) {
	context := config.Config.Server.Context
	if context == "" || context == "/" {
		context = ""
	}
	rw.Header().Set("refresh", "0.1;"+context+"/")
	rw.WriteHeader(200)
}

func handleResource(name, path string, r *mux.Router) {
	r.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(name, ".html") {
			rw.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(name, ".css") {
			rw.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(name, ".js") {
			rw.Header().Set("Content-Type", "application/javascript")
		}
		bytes, _ := Asset(name)
		_, _ = rw.Write(bytes)
	})
}

func bindMapping() (r *mux.Router) {
	r = mux.NewRouter()

	r.HandleFunc(serverContext[0:len(serverContext)-1], handleIndex)
	r.HandleFunc(serverContext+"user/register", userRegister)
	r.HandleFunc(serverContext+"user/search", userSearch)
	r.HandleFunc(serverContext+"user/active", userActive)
	r.HandleFunc(serverContext+"user/cancel", userCancel)

	r.HandleFunc(serverContext+"manage/user/query", manageUserQuery)
	r.HandleFunc(serverContext+"manage/user/insert", manageUserInsert)
	r.HandleFunc(serverContext+"manage/user/update", manageUserUpdate)
	r.HandleFunc(serverContext+"manage/user/lock", manageUserLock)
	r.HandleFunc(serverContext+"manage/user/unlock", manageUserUnlock)
	r.HandleFunc(serverContext+"manage/user/delete", manageUserDelete)

	//静态请求，由AssetFS统一处理。
	for k := range _bintree.Children {
		name := k
		path := serverContext + name
		handleResource(name, path, r)
		if name == "index.html" {
			handleResource(name, serverContext, r)
		}
	}
	return
}

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

	r := bindMapping()
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
		err = http.ListenAndServe(httpServer, r)
		if err != nil {
			panic(err)
		}
	}()

	ginR := gin.Default()
	BindApi(serverContext, ginR)
}

func BindApi(root string, r *gin.Engine) {

	idService.BindApi(root, r)
	userService.BindApi(root, r)
	wbsService.BindApi(root, r)
	logService.BindApi(root, r)
	enterpriseService.BindApi(root, r)
	jobService.BindApi(root, r)
	powerService.BindApi(root, r)
	settingService.BindApi(root, r)
	spaceService.BindApi(root, r)
	systemService.BindApi(root, r)
	messageService.BindApi(root, r)
	groupService.BindApi(root, r)
}
