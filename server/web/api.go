package web

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"server/base"
	certificateService "server/service/certificate"
	enterpriseService "server/service/enterprise"
	groupService "server/service/group"
	idService "server/service/id"
	jobService "server/service/job"
	logService "server/service/log"
	loginService "server/service/login"
	messageService "server/service/message"
	organizationService "server/service/organization"
	powerService "server/service/power"
	spaceService "server/service/space"
	systemService "server/service/system"
	userService "server/service/user"
	wbsService "server/service/wbs"
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
			// println("add api path :" + apiName + ",action:" + api.Power.Action)
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

	idService.BindApi(appendApi)
	userService.BindApi(appendApi)
	wbsService.BindApi(appendApi)
	logService.BindApi(appendApi)
	enterpriseService.BindApi(appendApi)
	organizationService.BindApi(appendApi)
	jobService.BindApi(appendApi)
	powerService.BindApi(appendApi)
	spaceService.BindApi(appendApi)
	systemService.BindApi(appendApi)
	loginService.BindApi(appendApi)
	messageService.BindApi(appendApi)
	groupService.BindApi(appendApi)
	certificateService.BindApi(appendApi)

	var apiPowerMap = make(map[string]bool)
	for _, api := range apiCache {
		apiPowerMap[api.Power.Action] = true
	}
	ps := base.GetPowers()
	for _, one := range ps {
		_, ok := apiPowerMap[one.Action]
		if !ok {
			panic(errors.New("权限[" + one.Action + "]未配置动作!"))
		}
	}
}
