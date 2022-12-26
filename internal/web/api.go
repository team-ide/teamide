package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (this_ *Server) bindApi(routerGroup *gin.RouterGroup) (err error) {

	routerGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if this_.api.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
	return
}
