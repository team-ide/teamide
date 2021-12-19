package service

import (
	"regexp"
	"server/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type DataRequest struct {
	Origin   string `json:"origin"`
	Pathname string `json:"pathname"`
}

type DataResponse struct {
	Url string `json:"url"`
	Api string `json:"api"`
}

func apiData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
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

	res = response
	return
}
