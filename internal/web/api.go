package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"teamide/pkg/base"
)

func (this_ *Server) bindApi(routerGroup *gin.RouterGroup) (err error) {

	routerGroup.POST("*path", func(c *gin.Context) {
		var path string
		defer func() {
			if e := recover(); e != nil {
				runtimeErr := errors.New(fmt.Sprint(e))
				this_.Logger.Error("post "+path+" runtime error", zap.Error(runtimeErr))
				response := base.HttpResponse{
					Code: "-1",
					Msg:  runtimeErr.Error(),
				}
				c.JSON(http.StatusOK, response)
			}
		}()
		re, _ := regexp.Compile("/+")
		path = c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if this_.api.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
	return
}
