package main

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url string `json:"url,omitempty"`
	Api string `json:"api,omitempty"`
}

func doApiData(path string, c *gin.Context) (res interface{}, err error) {
	path = path[0:strings.LastIndex(path, "api/")]
	request := &DataRequest{}
	if !RequestJSON(request, c) {
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
