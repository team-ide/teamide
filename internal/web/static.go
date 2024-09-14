package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"teamide/internal/static"
	"teamide/pkg/base"
)

func (this_ *Server) bindGet(routerGroup *gin.RouterGroup) {
	routerGroup.GET("*path", func(c *gin.Context) {
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
	return this_.toStaticByName("index.html", c)
}

func (this_ *Server) toStatic(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "static/")
	if index < 0 {
		return false
	}
	name := path[index:]

	return this_.toStaticByName(name, c)
}

func (this_ *Server) setHeaderByName(name string, c *gin.Context) {
	if strings.HasSuffix(name, ".html") {
		c.Header("Content-Type", "text/html")
		c.Header("Cache-Control", "no-cache")
	} else if strings.HasSuffix(name, ".css") {
		c.Header("Content-Type", "text/css")
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".js") {
		c.Header("Content-Type", "application/javascript")
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	} else if strings.HasSuffix(name, ".woff") ||
		strings.HasSuffix(name, ".ttf") ||
		strings.HasSuffix(name, ".woff2") ||
		strings.HasSuffix(name, ".eot") {
		// max-age 缓存 过期时间 秒为单位
		c.Header("Cache-Control", "max-age=31536000")
	}
}

func (this_ *Server) toStaticByName(name string, c *gin.Context) bool {

	// 查看本地文件是否有静态文件

	filePath := this_.RootDir + "statics/" + name

	if is, e := util.IsSubPath(this_.RootDir+"statics/", filePath); !is || e != nil {
		return false
	}
	exist, _ := util.PathExists(filePath)
	if exist {
		fileInfo, err := os.Open(filePath)
		if err != nil {
			base.ResponseJSON(nil, err, c)
			return true
		}
		defer func() { _ = fileInfo.Close() }()
		this_.setHeaderByName(name, c)
		_, err = io.Copy(c.Writer, fileInfo)
		if err != nil {
			base.ResponseJSON(nil, err, c)
			return true
		}
	} else {
		bytes := static.Asset(name)
		if bytes == nil {
			return false
		}
		this_.setHeaderByName(name, c)
		_, _ = c.Writer.Write(bytes)
	}

	c.Status(http.StatusOK)
	return true
}

func (this_ *Server) toFiles(path string, c *gin.Context) bool {

	if c.Query("isDownload") == "true" || c.Query("isDownload") == "1" {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Transfer-Encoding", "binary")
	}

	index := strings.LastIndex(path, "files/")
	if index < 0 {
		return false
	}
	name := path[index+len("files/"):]

	filePath := this_.GetFilesFile(name)

	if is, e := util.IsSubPath(this_.GetFilesDir(), filePath); !is || e != nil {
		return false
	}
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
	if err != nil {
		base.ResponseJSON(nil, err, c)
		return true
	}
	defer func() { _ = fileInfo.Close() }()
	_, err = io.Copy(c.Writer, fileInfo)
	if err != nil {
		base.ResponseJSON(nil, err, c)
		return true
	}
	c.Status(http.StatusOK)
	return true
}
