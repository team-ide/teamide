package module

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"teamide/internal/base"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url      string `json:"url,omitempty"`
	Api      string `json:"api,omitempty"`
	FilesUrl string `json:"filesUrl,omitempty"`
}

func (this_ *Api) apiData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	path := requestBean.Path[0:strings.LastIndex(requestBean.Path, "api/")]
	request := &DataRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DataResponse{}

	pathname := request.Pathname
	re, _ := regexp.Compile("/+")
	pathname = re.ReplaceAllLiteralString(pathname, "/")

	path = strings.TrimSuffix(path, "/")
	pathname = strings.TrimSuffix(pathname, "/")
	pathname = strings.TrimSuffix(pathname, path)

	if !strings.HasSuffix(pathname, "/") {
		pathname += "/"
	}

	response.Url = request.Origin + pathname
	response.Api = response.Url + "api/"
	response.FilesUrl = response.Api + "files/"

	res = response
	return
}
