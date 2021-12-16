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
	JWT := getJWT(c)
	request.JWT = JWT
	return
}

func doApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	requestBean := getRequestBean(c)
	requestBean.Path = path
	name := path[index+len("api/"):]

	api := apiCache[name]
	if api == nil {
		return false
	}
	if !checkPower(api, requestBean.JWT, c) {
		return true
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

func appendApi(apis ...*base.ApiWorker) {
	if len(apis) == 0 {
		return
	}
	for _, api := range apis {
		if api.Power == nil {
			panic(errors.New(fmt.Sprint("API未设置权限!", api)))
		}
		if len(api.Apis) == 0 {
			panic(errors.New(fmt.Sprint("API未设置映射路径!", api)))
		}
		for _, apiName := range api.Apis {

			_, find := apiCache[apiName]
			if find {
				panic(errors.New(fmt.Sprint("API映射路径[", apiName, "]已存在!", api)))
			}
			apiCache[apiName] = api
		}
	}
}
func cacheApi() {
	apiCache = make(map[string]*base.ApiWorker)

	appendApi(&base.ApiWorker{Apis: []string{"", "/", "data"}, Power: base.PowerData, Do: apiData})
	appendApi(&base.ApiWorker{Apis: []string{"login"}, Power: base.PowerLogin, Do: apiLogin})
	appendApi(&base.ApiWorker{Apis: []string{"autoLogin"}, Power: base.PowerAutoLogin, Do: apiLogin})
	appendApi(&base.ApiWorker{Apis: []string{"logout"}, Power: base.PowerLogout, Do: apiLogout})
	appendApi(&base.ApiWorker{Apis: []string{"register"}, Power: base.PowerRegister, Do: apiRegister})
	appendApi(&base.ApiWorker{Apis: []string{"session"}, Power: base.PowerSession, Do: apiSession})

	appendApi(idService.BindApi()...)
	appendApi(userService.BindApi()...)
	appendApi(wbsService.BindApi()...)
	appendApi(logService.BindApi()...)
	appendApi(enterpriseService.BindApi()...)
	appendApi(jobService.BindApi()...)
	appendApi(powerService.BindApi()...)
	appendApi(settingService.BindApi()...)
	appendApi(spaceService.BindApi()...)
	appendApi(systemService.BindApi()...)
	appendApi(messageService.BindApi()...)
	appendApi(groupService.BindApi()...)
}
