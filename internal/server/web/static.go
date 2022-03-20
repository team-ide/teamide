package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"teamide/internal/server/static"
)

func bindGet(gouterGroup *gin.RouterGroup) {
	gouterGroup.GET("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if toStatic(path, c) {
			return
		}
		if toUploads(path, c) {
			return
		}
		toIndex(c)
	})
}

func toIndex(c *gin.Context) bool {

	bytes := static.Asset("index.html")
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

func toUploads(path string, c *gin.Context) bool {

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
