package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"teamide/internal/module"
)

func (this_ *Server) bindApi(gouterGroup *gin.RouterGroup) {
	api := module.CacheApi(this_.DatabaseWorker)

	gouterGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if api.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
}
