package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"teamide/internal/server/service"
)

func bindApi(gouterGroup *gin.RouterGroup) {
	service.CacheApi()

	gouterGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if service.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
}
