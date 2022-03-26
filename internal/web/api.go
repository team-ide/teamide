package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"teamide/internal/module"
)

func (this_ *Server) bindApi(gouterGroup *gin.RouterGroup) (err error) {
	api, err := module.NewApi(this_.ServerContext)

	if err != nil {
		return
	}

	gouterGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if api.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
	return
}
