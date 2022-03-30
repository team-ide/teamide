package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/static"
)

func (this_ *Server) bindGet(gouterGroup *gin.RouterGroup) {
	gouterGroup.GET("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		//fmt.Println("path=" + path)
		if strings.HasSuffix(path, "api/ws/toolbox/ssh/connection") {
			err := module_toolbox.SSHConnection(c)
			if err != nil {
				base.ResponseJSON(nil, err, c)
			}
			return
		}
		if this_.toStatic(path, c) {
			return
		}
		if this_.toUploads(path, c) {
			return
		}
		this_.toIndex(c)
	})
}

func (this_ *Server) toIndex(c *gin.Context) bool {

	bytes := static.Asset("index.html")
	if bytes == nil {
		return false
	}

	c.Header("Content-Type", "text/html")
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}

func (this_ *Server) toStatic(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "static/")
	if index < 0 {
		return false
	}
	name := path[index:]

	bytes := static.Asset(name)
	if bytes == nil {
		return false
	}

	if strings.HasSuffix(name, ".html") {
		c.Header("Content-Type", "text/html")
	} else if strings.HasSuffix(name, ".css") {
		c.Header("Content-Type", "text/css")
	} else if strings.HasSuffix(name, ".js") {
		c.Header("Content-Type", "teamide/application/javascript")
	}
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}

func (this_ *Server) toUploads(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "uploads/")
	if index < 0 {
		return false
	}
	name := path[index:]

	bytes := static.Asset(name)
	if bytes == nil {
		return false
	}
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}
