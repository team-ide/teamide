package main

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func bindGet(routerGroup *gin.RouterGroup) {
	routerGroup.GET("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if toStatic(path, c) {
			return
		}
		toIndex(c)
	})
}

func toIndex(c *gin.Context) bool {

	bytes, _ := Asset("index.html")
	if bytes == nil {
		return false
	}

	c.Header("Content-Type", "text/html")
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}

func toStatic(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "static/")
	if index < 0 {
		return false
	}
	name := path[index:]

	bytes, _ := Asset(name)
	if bytes == nil {
		return false
	}

	if strings.HasSuffix(name, ".html") {
		c.Header("Content-Type", "text/html")
	} else if strings.HasSuffix(name, ".css") {
		c.Header("Content-Type", "text/css")
	} else if strings.HasSuffix(name, ".js") {
		c.Header("Content-Type", "application/javascript")
	}
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}

func bindPost(routerGroup *gin.RouterGroup) {
	routerGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if doApi(path, c) {
			return
		}
		toIndex(c)
	})
}

func doApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	name := path[index+len("api/"):]

	var res interface{}
	var err error
	if name == "" || name == "/" || name == "/data" {
		res, err = doApiData(path, c)
	} else if name == "session" || name == "/session" {
		res, err = doApiSession(path, c)
	} else if name == "app/list" || name == "/app/list" {
		res, err = doApiAppList(path, c)
	} else if name == "app/insert" || name == "/app/insert" {
		res, err = doApiAppInsert(path, c)
	} else if name == "app/rename" || name == "/app/rename" {
		res, err = doApiAppRename(path, c)
	} else if name == "context/get" || name == "/context/get" {
		res, err = doApiContextGet(path, c)
	} else if name == "context/save" || name == "/context/save" {
		res, err = doApiContextSave(path, c)
	}
	if res != nil || err != nil {
		c.JSON(200, toWebResponseJSON(res, err))
	}
	return true
}

func bindWebServerApis(routerGroup *gin.RouterGroup) (err error) {
	bindGet(routerGroup)

	bindPost(routerGroup)
	return
}

var (
	webServerHost    = "127.0.0.1"
	webServerPort    = 21010
	webServerAddress = ""
	webServerUrl     = ""
	htmlServerUrl    = "http://127.0.0.1:21011/"
)

func init() {
	if !isDev {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		webServerPort = listener.Addr().(*net.TCPAddr).Port
	}
	webServerAddress = fmt.Sprint(webServerHost, ":", webServerPort)
	webServerUrl = "http://" + webServerAddress
	if !isDev {
		htmlServerUrl = webServerUrl
	}

}

func startWebServer() (err error) {

	gin.DefaultWriter = &nullWriter{}

	router := gin.Default()

	routerGroup := router.Group("")

	err = bindWebServerApis(routerGroup)
	if err != nil {
		return
	}
	go func() {
		err = http.ListenAndServe(webServerAddress, router)
		if err != nil {
			panic(err)
		}
	}()
	// time.Sleep(1000 * 3)
	return
}

type nullWriter struct{}

func (*nullWriter) Write(bs []byte) (int, error) {

	return 0, nil
}
