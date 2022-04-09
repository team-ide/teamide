package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func (this_ *Server) bindApi(gouterGroup *gin.RouterGroup) (err error) {

	gouterGroup.POST("*path", func(c *gin.Context) {
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
