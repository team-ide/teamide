package web

import (
	"net/http"
	"regexp"
	"server/enterpriseService"
	"server/groupService"
	"server/idService"
	"server/jobService"
	"server/logService"
	"server/messageService"
	"server/powerService"
	"server/settingService"
	"server/spaceService"
	"server/systemService"
	"server/userService"
	"server/wbsService"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	apiCache map[string]func(c *gin.Context)
)

func bindApi(gouterGroup *gin.RouterGroup) {
	cacheApi()

	gouterGroup.POST("*path", func(c *gin.Context) {
		re, _ := regexp.Compile("/+")
		path := c.Params.ByName("path")
		path = re.ReplaceAllLiteralString(path, "/")
		if doApi(path, c) {
			return
		}
		println("path:" + path)
		c.Status(http.StatusNotFound)
	})
}

func doApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	name := path[index+len("api/"):]

	if name == "" || name == "/" {
		apiData(path[0:index], c)
		return true
	}

	api := apiCache[name]
	if api == nil {
		return false
	}
	api(c)
	return true
}

func cacheApi() {
	apiCache = make(map[string]func(c *gin.Context))
	idService.BindApi(apiCache)
	userService.BindApi(apiCache)
	wbsService.BindApi(apiCache)
	logService.BindApi(apiCache)
	enterpriseService.BindApi(apiCache)
	jobService.BindApi(apiCache)
	powerService.BindApi(apiCache)
	settingService.BindApi(apiCache)
	spaceService.BindApi(apiCache)
	systemService.BindApi(apiCache)
	messageService.BindApi(apiCache)
	groupService.BindApi(apiCache)
}
