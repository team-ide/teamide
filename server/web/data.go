package web

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

func apiData(path string, c *gin.Context) {
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
	base.ResponseJSON(response, nil, c)
}
