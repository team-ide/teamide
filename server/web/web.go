package web

import (
	"fmt"
	"server/config"

	"net/http"
	"strings"

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

func StartServer() {

	context := config.Config.Server.Context
	if context == "" || context == "/" {
		context = ""
	}

	r := mux.NewRouter()

	r.HandleFunc(context+"", handleIndex)

	//静态请求，由AssetFS统一处理。
	for k := range _bintree.Children {
		name := k
		path := context + "/" + name
		handleResource(name, path, r)
		if name == "index.html" {
			handleResource(name, context+"/", r)
		}
	}

	host := config.Config.Server.Host
	port := config.Config.Server.Port

	httpServer := fmt.Sprint(host, ":", port)
	if context != "/" {
		context += "/"
	}
	println("url:http://" + httpServer + context)
	err := http.ListenAndServe(httpServer, r)
	if err != nil {
		panic(err)
	}
}
