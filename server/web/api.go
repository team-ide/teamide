package web

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"server/base"
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
	apiCache map[string]*base.ApiWorker
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

func getRequestBean(c *gin.Context) (request *base.RequestBean) {
	request = &base.RequestBean{}
	sessionBean := GetSessionBean(c.Request)
	request.Session = sessionBean
	if sessionBean != nil {
		request.User = sessionBean.User
	}
	return
}

func doApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	requestBean := getRequestBean(c)
	name := path[index+len("api/"):]

	if name == "" || name == "/" {
		apiData(path[0:index], requestBean, c)
		return true
	}

	api := apiCache[name]
	if api == nil {
		return false
	}
	if api.Do != nil {
		res, err := api.Do(requestBean, c)
		base.ResponseJSON(res, err, c)
	}
	if api.DoOther != nil {
		api.DoOther(requestBean, c)
	}
	return true
}
func appendApi(apis []*base.ApiWorker) {
	if len(apis) == 0 {
		return
	}
	for _, api := range apis {
		if api.Api == "" {
			panic(errors.New(fmt.Sprint("api is null.", api)))
		}
		_, find := apiCache[api.Api]
		if find {
			panic(errors.New(fmt.Sprint("api [", api.Api, "] already exists.", api)))
		}
		apiCache[api.Api] = api
	}
}
func cacheApi() {
	apiCache = make(map[string]*base.ApiWorker)

	apiCache["login"] = &base.ApiWorker{
		Do: apiLogin,
	}
	apiCache["logout"] = &base.ApiWorker{
		Do: apiLogout,
	}
	apiCache["register"] = &base.ApiWorker{
		Do: apiRegister,
	}
	apiCache["session"] = &base.ApiWorker{
		Do: apiSession,
	}

	appendApi(idService.BindApi())
	appendApi(userService.BindApi())
	appendApi(wbsService.BindApi())
	appendApi(logService.BindApi())
	appendApi(enterpriseService.BindApi())
	appendApi(jobService.BindApi())
	appendApi(powerService.BindApi())
	appendApi(settingService.BindApi())
	appendApi(spaceService.BindApi())
	appendApi(systemService.BindApi())
	appendApi(messageService.BindApi())
	appendApi(groupService.BindApi())
}
