package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"teamide/internal/base"
	"teamide/internal/static"
	"teamide/pkg/util"
)

func (this_ *Server) bindGet(gouterGroup *gin.RouterGroup) {
	gouterGroup.GET("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")

		if this_.api.DoApi(path, c) {
			return
		}
		if this_.toStatic(path, c) {
			return
		}
		if this_.toFiles(path, c) {
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
		c.Header("Content-Type", "application/javascript")
	}
	c.Writer.Write(bytes)
	c.Status(http.StatusOK)
	return true
}

func (this_ *Server) toFiles(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "files/")
	if index < 0 {
		return false
	}
	name := path[index+len("files/"):]

	filePath := this_.GetFilesFile(name)

	exist, err := util.PathExists(filePath)
	if err != nil {
		base.ResponseJSON(nil, err, c)
		return true
	}
	if !exist {
		err = errors.New("文件[" + name + "]不存在")
		base.ResponseJSON(nil, err, c)
		return true
	}
	fileInfo, err := os.Open(filePath)
	if fileInfo == nil {
		base.ResponseJSON(nil, err, c)
		return true
	}
	defer fileInfo.Close()
	_, err = io.Copy(c.Writer, fileInfo)
	if err != nil {
		base.ResponseJSON(nil, err, c)
		return true
	}
	c.Status(http.StatusOK)
	return true
}
