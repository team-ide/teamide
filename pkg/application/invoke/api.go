package invoke

import (
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"teamide/pkg/application/base"
	"teamide/pkg/application/common"
	"teamide/pkg/application/model"

	"github.com/gin-gonic/gin"
)

func ServerWebBindApis(app common.IApplication, serverWebToken *model.ServerWebToken, gouterGroup *gin.RouterGroup) (err error) {

	if len(app.GetContext().Actions) == 0 {
		return
	}
	for _, one := range app.GetContext().Actions {
		err = serverWebBindActionApi(app, serverWebToken, gouterGroup, one)
		if err != nil {
			return
		}
	}
	return
}

type webResponseJSON struct {
	Code  string      `json:"code,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func toWebResponseJSON(value interface{}, err error) *webResponseJSON {
	response := &webResponseJSON{
		Code: "0",
		Msg:  "成功",
	}
	if err != nil {
		response.Msg = err.Error()
		baseErr, baseErrOk := err.(*base.ErrorBase)
		if baseErrOk {
			response.Code = baseErr.Code
			response.Msg = baseErr.Msg
		} else {
			response.Code = "-1"
		}
	} else {
		response.Value = value
	}
	return response
}

func serverWebBindActionApi(app common.IApplication, serverWebToken *model.ServerWebToken, gouterGroup *gin.RouterGroup, action *model.ActionModel) (err error) {
	if action.Api == nil || action.Api.Request == nil {
		return
	}
	requestPath := action.Api.Request.Path
	if requestPath != "" {
		if strings.Index(requestPath, "/") != 0 {
			requestPath = "/" + requestPath
		}
	}
	requestMethod := strings.ToUpper(action.Api.Request.Method)
	var bindFunc func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	switch requestMethod {
	case "GET":
		bindFunc = gouterGroup.GET
	case "DELETE":
		bindFunc = gouterGroup.DELETE
	case "PUT":
		bindFunc = gouterGroup.PUT
	case "HEAD":
		bindFunc = gouterGroup.HEAD
	case "PATCH":
		bindFunc = gouterGroup.PATCH
	case "OPTIONS":
		bindFunc = gouterGroup.OPTIONS
	default:
		bindFunc = gouterGroup.POST
	}
	if strings.Contains(requestPath, "{") && strings.Contains(requestPath, "}") {
		requestPath_ := ""
		var re *regexp.Regexp
		re, err = regexp.Compile(`{(.+?)}`)
		if err != nil {
			return
		}
		indexsList := re.FindAllIndex([]byte(requestPath), -1)
		var lastIndex int = 0
		for _, indexs := range indexsList {
			requestPath_ += requestPath[lastIndex:indexs[0]]
			lastIndex = indexs[1]

			name := requestPath[indexs[0]+1 : indexs[1]-1]
			if strings.Index(name, "*") == 0 {
				requestPath_ += name
			} else {
				requestPath_ += ":" + name
			}
		}
		requestPath_ += requestPath[lastIndex:]

		requestPath = requestPath_
	}

	var shouldValidataToken bool

	// 访问路径 匹配忽略路径
	if serverWebToken.Exclude != "" && matchRequestTokenRule(requestPath, serverWebToken.Exclude) {

	} else if serverWebToken.Include != "" && matchRequestTokenRule(requestPath, serverWebToken.Exclude) {
		shouldValidataToken = true
	}

	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("web api bind [", requestPath, "] should token [", shouldValidataToken, "] for action [", action.Name, "]")
	}

	bindFunc(requestPath, func(c *gin.Context) {
		var err error
		var res interface{}

		defer func() {
			if err == nil {
				err = base.PanicError(recover())
			}
			if err != nil {
				if app.GetLogger() != nil {
					app.GetLogger().Error("web api request [", requestPath, "] action [", action.Name, "] error:", err)
				}
			}
			if action.Api.Response != nil {
				if action.Api.Response.DownloadFile || action.Api.Response.OpenFile {
					if err != nil {
						response := toWebResponseJSON(res, err)
						c.JSON(200, response)
						return
					}
					var reader io.Reader = nil
					var fileName string
					if res != nil {
						switch v := res.(type) {
						case io.Reader:
							reader = v
						case *model.FileInfo:
							fileName = v.Name
							reader, err = v.GetReader()
						case map[string]interface{}:
							fileInfo := &model.FileInfo{}
							if v["name"] != nil {
								fileInfo.Name = v["name"].(string)
							}
							if v["path"] != nil {
								fileInfo.Path = v["path"].(string)
							}
							if v["absolutePath"] != nil {
								fileInfo.AbsolutePath = v["absolutePath"].(string)
							}
							if v["fileHeader"] != nil {
								fileInfo.FileHeader = v["fileHeader"].(*multipart.FileHeader)
							}
							if v["file"] != nil {
								fileInfo.File = v["file"].(*os.File)
							}
							fileName = fileInfo.Name
							reader, err = fileInfo.GetReader()
						}
						if err != nil {
							response := toWebResponseJSON(nil, err)
							c.JSON(200, response)
							return
						}
					}
					if reader == nil {
						c.Status(404)
						return
					}
					if action.Api.Response.DownloadFile {
						c.Header("contentType", "multipart/form-data")
						c.Header("Content-Disposition", "attachment;fileName="+fileName)
					}
					io.Copy(c.Writer, reader)
					return
				}
			}
			response := toWebResponseJSON(res, err)
			c.JSON(200, response)
		}()

		res, err = InvokeWebApi(app, serverWebToken, shouldValidataToken, c, action)
		if err != nil {
			return
		}
	})
	return
}

func matchRequestTokenRule(path string, rule_ string) bool {
	if path == rule_ {
		return true
	}
	rules := strings.Split(rule_, ",")
	for _, rule := range rules {
		rule = strings.ReplaceAll(rule, "**", "(.+?)")
		re := regexp.MustCompile(rule)
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func InvokeWebApi(app common.IApplication, serverWebToken *model.ServerWebToken, shouldValidataToken bool, c *gin.Context, action *model.ActionModel) (res interface{}, err error) {
	var invokeNamespace *common.InvokeNamespace
	invokeNamespace, err = common.NewInvokeNamespace(app)
	if err != nil {
		return
	}
	invokeNamespace.RequestContext = c
	invokeNamespace.ServerWebToken = serverWebToken
	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("invoke web api [", action.Api.Request, "] start")
	}

	startTime := base.GetNowTime()
	defer func() {
		endTime := base.GetNowTime()
		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("invoke web api [", action.Api.Request, "] end, use:", (endTime - startTime), "ms")
		}
	}()

	if base.IsEmpty(action.WebApiJavascript) {
		action.WebApiJavascript, err = common.GetWebApiJavascriptByAction(app, action, shouldValidataToken)
		if err != nil {
			return
		}
	}
	res, err = invokeJavascript(app, invokeNamespace, action.WebApiJavascript)
	if err != nil {
		return
	}
	return
}

func invokeWebApiValidateToken(app common.IApplication, serverWebToken *model.ServerWebToken, c *gin.Context) (res interface{}, err error) {
	if serverWebToken.ValidateAction == "" {
		err = base.NewError("", "request token validata action not defind")
		return
	}
	action := app.GetContext().GetAction(serverWebToken.ValidateAction)
	if action == nil {
		err = base.NewError("", "request token validata action not defind")
		return
	}
	var invokeNamespace *common.InvokeNamespace
	invokeNamespace, err = common.NewInvokeNamespace(app)
	if err != nil {
		return
	}
	invokeNamespace.RequestContext = c
	invokeNamespace.ServerWebToken = serverWebToken

	if app.GetLogger() != nil && app.GetLogger().OutDebug() {
		app.GetLogger().Debug("invoke web api validata token [", action.Name, "] start")
	}

	startTime := base.GetNowTime()
	defer func() {
		endTime := base.GetNowTime()
		if app.GetLogger() != nil && app.GetLogger().OutDebug() {
			app.GetLogger().Debug("invoke web api validata token[", action.Name, "] end, use:", (endTime - startTime), "ms")
		}
	}()

	if base.IsEmpty(action.WebApiJavascript) {
		action.WebApiJavascript, err = common.GetWebApiJavascriptByAction(app, action, false)
		if err != nil {
			return
		}
	}
	res, err = invokeJavascript(app, invokeNamespace, action.WebApiJavascript)
	if err != nil {
		return
	}
	return
}
