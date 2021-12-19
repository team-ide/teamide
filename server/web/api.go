package web

import (
	"net/http"
	"regexp"
	"server/service"

	"github.com/gin-gonic/gin"
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
		println("path:" + path)
		c.Status(http.StatusNotFound)
	})
}
