package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	"teamide/internal/server/factory"
)

func getRequestBean(c *gin.Context) (request *base2.RequestBean) {
	request = &base2.RequestBean{}
	JWT := getJWT(c)
	request.JWT = JWT
	return
}

func DoApi(path string, c *gin.Context) bool {

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
		base2.ResponseJSON(res, err, c)
	}
	if api.DoOther != nil {
		api.DoOther(requestBean, c)
	}
	return true
}

func appendApi(apis ...*base2.ApiWorker) {
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

		if base2.IS_STAND_ALONE {
			if !api.Power.AllowNative {
				continue
			}
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
func CacheApi() {
	apiCache = make(map[string]*base2.ApiWorker)
	bindApi(appendApi)

	var apiPowerMap = make(map[string]bool)
	for _, api := range apiCache {
		apiPowerMap[api.Power.Action] = true
	}
	ps := base2.GetPowers()
	for _, one := range ps {
		if base2.IS_STAND_ALONE {
			if !one.AllowNative {
				continue
			}
		}
		_, ok := apiPowerMap[one.Action]
		if !ok {
			component.Logger.Warn(component.LogStr("权限[", one.Action, "]未配置动作"))
		}
	}
}

func bindApi(appendApi func(apis ...*base2.ApiWorker)) {
	appendApi(&base2.ApiWorker{Apis: []string{"", "/", "data"}, Power: base2.PowerData, Do: apiData})
	appendApi(&base2.ApiWorker{Apis: []string{"login"}, Power: base2.PowerLogin, Do: apiLogin})
	appendApi(&base2.ApiWorker{Apis: []string{"autoLogin"}, Power: base2.PowerAutoLogin, Do: apiLogin})
	appendApi(&base2.ApiWorker{Apis: []string{"logout"}, Power: base2.PowerLogout, Do: apiLogout})
	appendApi(&base2.ApiWorker{Apis: []string{"register"}, Power: base2.PowerRegister, Do: apiRegister})
	appendApi(&base2.ApiWorker{Apis: []string{"session"}, Power: base2.PowerSession, Do: apiSession})

	for _, one := range factory.Apis {
		one.BindApi(appendApi)
	}
}
