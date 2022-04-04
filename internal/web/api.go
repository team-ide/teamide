package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"teamide/internal/base"
	"teamide/internal/module"
	"teamide/internal/module/module_toolbox"
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
		if strings.HasSuffix(path, "api/upload") {
			res, err := upload(c)
			base.ResponseJSON(res, err, c)
			return
		} else if strings.HasSuffix(path, "api/download") {
			data := map[string]string{}
			err = c.BindJSON(&data)
			if err != nil {
				base.ResponseJSON(nil, err, c)
				return
			}
			err = download(data, c)
			if err != nil {
				base.ResponseJSON(nil, err, c)
				return
			}
			return
		}

		if api.DoApi(path, c) {
			return
		}
		c.Status(http.StatusNotFound)
	})
	//gouterGroup.GET("ws/toolbox/ssh/connection", func(c *gin.Context) {
	//	err := module_toolbox.SSHConnection(c)
	//	if err != nil {
	//		base.ResponseJSON(nil, err, c)
	//	}
	//})
	return
}

func upload(c *gin.Context) (res interface{}, err error) {

	uploadType := c.PostForm("type")

	switch uploadType {
	case "sftp":
		res, err = module_toolbox.SFTPUpload(c)
		break
	}

	return
}

func download(data map[string]string, c *gin.Context) (err error) {

	uploadType := data["type"]
	switch uploadType {
	case "sftp":
		err = module_toolbox.SFTPDownload(data, c)
		break
	}

	return
}
