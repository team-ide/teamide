package application

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"teamide/application/base"
	"teamide/application/common"
	"teamide/application/invoke"
	"teamide/application/model"
)

type ContextDoc struct {
	Dir             string
	dirAbsolutePath string
	App             common.IApplication
}

func (this_ *ContextDoc) Out() (err error) {
	var path string
	path, err = os.Getwd()
	if err != nil {
		return
	}
	dirPath := path + "/" + this_.Dir
	var exists bool
	exists, err = base.PathExists(dirPath)
	if err != nil {
		return
	}
	if !exists {
		os.MkdirAll(dirPath, 0777)
	}
	var abs string
	abs, err = filepath.Abs(dirPath)
	if err != nil {
		return
	}
	this_.dirAbsolutePath = filepath.ToSlash(abs)

	for _, serverWeb := range this_.App.GetContext().ServerWebs {
		err = this_.outWebApi(serverWeb)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ContextDoc) outWebApi(serverWeb *model.ServerWebModel) (err error) {
	var fileName = "web.api.md"
	if base.IsNotEmpty(serverWeb.Name) {
		fileName = serverWeb.Name + "." + fileName
	}

	var exists bool
	indexPath := this_.dirAbsolutePath + "/" + fileName
	exists, err = base.PathExists(indexPath)
	if err != nil {
		return
	}
	if exists {
		os.RemoveAll(indexPath)
	}
	var abs string
	abs, err = filepath.Abs(indexPath)
	if err != nil {
		return
	}
	indexPath = filepath.ToSlash(abs)

	if this_.App.GetLogger() != nil {
		this_.App.GetLogger().Info("out web api doc to [", indexPath, "]")
	}

	var content string
	content, err = this_.getApiWebContent(serverWeb)
	if err != nil {
		return
	}
	var file *os.File
	file, err = os.Create(indexPath)
	defer func() {
		e := file.Close()
		if e != nil {
			fmt.Println("file close error:", e)
		}
	}()
	file.WriteString(content)
	return
}

func (this_ *ContextDoc) getApiWebContent(serverWeb *model.ServerWebModel) (content string, err error) {
	appendLine(&content, "## Web Server Api")
	appendLine(&content, "")
	contextPath := serverWeb.ContextPath
	if !strings.HasPrefix(contextPath, "/") {
		contextPath = "/" + contextPath
	}
	if !strings.HasSuffix(contextPath, "/") {
		contextPath = contextPath + "/"
	}
	var url string = "http://"
	if serverWeb.Host == ":" || serverWeb.Host == "::" || serverWeb.Host == "0.0.0.0" {
		url += "127.0.0.1"
	}
	url += ":" + fmt.Sprint(serverWeb.Port)
	url += contextPath
	appendLine(&content, "### 访问路径")
	appendLine(&content, "> "+url)
	appendLine(&content, "")
	err = this_.appendServiceApiWebTokenContent(serverWeb, &content)
	if err != nil {
		return
	}
	for _, one := range this_.App.GetContext().Services {
		if one.Api == nil || one.Api.Request == nil {
			continue
		}
		err = this_.appendServiceApiWebContent(one, &content)
		if err != nil {
			fmt.Println("service [", one.Name, "] api error:", err)
			return
		}
	}
	return
}

func (this_ *ContextDoc) appendServiceApiWebTokenContent(serverWeb *model.ServerWebModel, content *string) (err error) {
	if serverWeb.Token == nil {
		return
	}
	tokenValidateService := this_.App.GetContext().GetService(serverWeb.Token.ValidateService)
	if tokenValidateService == nil {
		err = base.NewError("", "token validata service [", serverWeb.Token.ValidateService, "] not defind")
		return
	}

	appendLine(content, "### Token 验证")
	appendLine(content, "")
	appendLine(content, "#### 验证路径")
	appendLine(content, "> "+strings.ReplaceAll(serverWeb.Token.Include, ",", " , "))
	appendLine(content, "")
	appendLine(content, "#### 忽略路径")
	appendLine(content, "> "+strings.ReplaceAll(serverWeb.Token.Exclude, ",", " , "))
	appendLine(content, "")

	paths := getRequestPath(tokenValidateService)
	if len(paths) > 0 {
		appendLine(content, "#### 路径参数（Path）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range paths {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")

	}

	params := getRequestParam(tokenValidateService)
	if len(params) > 0 {
		appendLine(content, "#### 请求参数（URL）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range params {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	headers := getRequestHeader(tokenValidateService)
	if len(headers) > 0 {
		appendLine(content, "#### 请求头（Header）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range headers {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	tokenCreateService := this_.App.GetContext().GetService(serverWeb.Token.CreateService)
	if tokenCreateService == nil {
		err = base.NewError("", "token create service [", serverWeb.Token.CreateService, "] not defind")
		return
	}
	err = this_.appendServiceApiWebContent(tokenCreateService, content)
	if err != nil {
		return
	}
	return
}
func getRequestHeader(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_HEADER {
			res = append(res, one)
		}
	}
	return
}
func getRequestBody(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if base.IsEmpty(one.DataPlace) || model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_BODY {
			res = append(res, one)
		}
	}
	return
}
func getRequestFile(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_FILE {
			res = append(res, one)
		}
	}
	return
}
func getRequestForm(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_FORM {
			res = append(res, one)
		}
	}
	return
}

func getRequestParam(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_PARAM {
			res = append(res, one)
		}
	}
	return
}

func getRequestPath(service *model.ServiceModel) (res []*model.VariableModel) {
	for _, one := range service.InVariables {
		if model.GetVariableDataPlace(one.DataPlace) == model.DATA_PLACE_PATH {
			res = append(res, one)
		}
	}
	return
}

func (this_ *ContextDoc) appendServiceApiWebContent(service *model.ServiceModel, content *string) (err error) {

	var webApiJavascript string
	webApiJavascript, err = common.GetWebApiJavascriptByService(this_.App, service, true)
	if err != nil {
		return
	}
	functionParser := invoke.NewFunctionParser(webApiJavascript)

	var invokeNamespace *common.InvokeNamespace
	invokeNamespace, err = common.NewInvokeNamespace(this_.App)
	if err != nil {
		return
	}
	parseInfo := &invoke.ParseInfo{
		App:             this_.App,
		InvokeNamespace: invokeNamespace,
	}
	err = functionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	if base.IsEmpty(service.Comment) {
		appendLine(content, fmt.Sprint("### ", service.Name, ""))
	} else {
		appendLine(content, fmt.Sprint("### ", service.Comment, ""))
	}
	appendLine(content, "")

	if base.IsNotEmpty(service.Comment) || base.IsNotEmpty(service.Description) {
		appendLine(content, "#### 接口功能")
		appendLine(content, "")
		if base.IsEmpty(service.Description) {
			appendLine(content, "> "+service.Comment)
		} else {
			appendLine(content, "> "+service.Description)
		}
		appendLine(content, "")
	}

	appendLine(content, "#### 请求地址")
	appendLine(content, "")
	requestPath := service.Api.Request.Path
	if !strings.HasPrefix(requestPath, "/") {
		requestPath = "/" + requestPath
	}
	appendLine(content, "> "+requestPath)
	appendLine(content, "")

	appendLine(content, "#### 请求方式")
	appendLine(content, "")
	if base.IsEmpty(service.Api.Request.Method) {
		appendLine(content, "> POST")
	} else {
		appendLine(content, "> "+service.Api.Request.Method)
	}
	appendLine(content, "")

	err = this_.appendInVariable(service, content, parseInfo)
	if err != nil {
		return
	}

	err = this_.appendOutVariable(service, content, parseInfo)
	if err != nil {
		return
	}

	appendLine(content, "")

	return
}

func (this_ *ContextDoc) appendInVariable(service *model.ServiceModel, content *string, parseInfo *invoke.ParseInfo) (err error) {

	paths := getRequestPath(service)
	if len(paths) > 0 {
		appendLine(content, "#### 路径参数（Path）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range paths {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")

	}

	params := getRequestParam(service)
	if len(params) > 0 {
		appendLine(content, "#### 请求参数（URL）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range params {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	headers := getRequestHeader(service)
	if len(headers) > 0 {
		appendLine(content, "#### 请求头（Header）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range headers {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	forms := getRequestForm(service)
	if len(forms) > 0 {
		appendLine(content, "#### 请求表单（Form）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range forms {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	files := getRequestFile(service)
	if len(files) > 0 {
		appendLine(content, "#### 表单文件（Form）")
		appendLine(content, "")
		appendLine(content, "|字段名称    |字段说明   |必填   |")
		appendLine(content, "| ----------|:--------:|:--------:|")
		for _, one := range files {
			appendLine(content, "|"+one.Name+"|"+one.Comment+"|Y|")
		}
		appendLine(content, "")
	}

	bodys := getRequestBody(service)
	if len(bodys) > 0 {
		var dataInfos []*common.InvokeDataInfo
		for _, one := range parseInfo.InvokeNamespace.DataInfos {
			for _, inV := range bodys {
				if one.Name == inV.Name {
					dataInfos = append(dataInfos, one)
				}
			}
		}

		var json string = ""
		err = this_.appendInJSONDataInfos(service, &json, dataInfos, 0)
		if err != nil {
			return
		}
		if base.IsNotEmpty(json) {
			appendLine(content, "#### 请求数据（JSON）")
			appendLine(content, "")
			appendLine(content, "```json")
			appendLine(content, "")
			appendLine(content, json)
			appendLine(content, "```")
			appendLine(content, "")
		}
	}

	return
}

func (this_ *ContextDoc) appendOutVariable(service *model.ServiceModel, content *string, parseInfo *invoke.ParseInfo) (err error) {

	// 输出文件
	if service.Api.Response != nil && (service.Api.Response.DownloadFile || service.Api.Response.OpenFile) {
		if service.Api.Response.DownloadFile {
			appendLine(content, "#### 下载文件")
			appendLine(content, "")
			appendLine(content, "```")
			appendLine(content, "文件流")
			appendLine(content, "```")
			appendLine(content, "")
		} else if service.Api.Response.OpenFile {
			appendLine(content, "#### 打开文件")
			appendLine(content, "")
			appendLine(content, "```")
			appendLine(content, "文件流")
			appendLine(content, "```")
			appendLine(content, "")
		}
		return
	}
	var json string = ""
	json += "{\n"
	json += `  "code": "0",` + " // 错误码\n"
	json += `  "msg": "", // 错误信息` + "\n"
	if service.OutVariable != nil {
		var outDataInfo *common.InvokeDataInfo
		for _, one := range parseInfo.InvokeNamespace.DataInfos {
			if one.Name == service.OutVariable.Name {
				outDataInfo = one
				break
			}
		}
		err = this_.appendJSONDataInfo(service, &json, "value", outDataInfo, true, true, 1)
		if err != nil {
			return
		}
	}

	json += "}"
	if base.IsNotEmpty(json) {
		appendLine(content, "#### 返回数据（JSON）")
		appendLine(content, "")
		appendLine(content, "```json")
		appendLine(content, "")
		appendLine(content, json)
		appendLine(content, "```")
		appendLine(content, "")
	}
	return
}

func (this_ *ContextDoc) appendInJSONDataInfos(service *model.ServiceModel, content *string, dataInfos []*common.InvokeDataInfo, tab int) (err error) {
	for i := 0; i < tab; i++ {
		*content += "  "
	}
	*content += "{\n"
	for index, one := range dataInfos {
		isEnd := len(dataInfos) == index+1
		name := one.Name
		if base.IsNotEmpty(one.Value) {
			name = one.Value
		}
		if len(one.DataInfos) > 0 && base.IsEmpty(one.Value) {
			if !one.IsList {
				name = ""
			}
		}
		err = this_.appendJSONDataInfo(service, content, name, one, false, isEnd, tab+1)
		if err != nil {
			return
		}
	}
	for i := 0; i < tab; i++ {
		*content += "  "
	}
	*content += "}\n"
	return
}

func (this_ *ContextDoc) appendJSONDataInfo(service *model.ServiceModel, content *string, name string, dataInfo *common.InvokeDataInfo, isOut bool, isEnd bool, tab int) (err error) {

	if len(dataInfo.DataInfos) > 0 {
		var list []*common.InvokeDataInfo

		if isOut {
			list = append(list, dataInfo.DataInfos...)
		} else {

			for _, one := range dataInfo.DataInfos {
				if !one.IsSetValue && one.IsUse {
					list = append(list, one)
				}
			}
		}
		if base.IsNotEmpty(name) {
			for i := 0; i < tab; i++ {
				*content += "  "
			}
			*content += `"` + name + `":`
			if dataInfo.IsList {
				*content += "[{"
			} else {
				*content += "{"
			}
			comment := dataInfo.Comment
			if dataInfo.DataType != nil {
				if base.IsNotEmpty(comment) {
					comment += ", "
				}
				comment += dataInfo.DataType.Text

			}
			if base.IsNotEmpty(comment) {
				*content += " // " + comment
			}
			*content += "\n"
		}
		for index, one := range list {
			isEnd := len(list) == index+1
			tab_ := tab
			if base.IsNotEmpty(name) {
				tab_++
			}
			err = this_.appendJSONDataInfo(service, content, one.Name, one, isOut, isEnd, tab_)
			if err != nil {
				return
			}
		}

		if base.IsNotEmpty(name) {
			for i := 0; i < tab; i++ {
				*content += "  "
			}
			if dataInfo.IsList {
				*content += "}]"
			} else {
				*content += "}"
			}
			if !isEnd {
				*content += ","
			}
			*content += "\n"
		}
	} else {
		for i := 0; i < tab; i++ {
			*content += "  "
		}
		*content += `"` + name + `":`
		if dataInfo.IsList {
			*content += "[" + name + "]"
		} else {
			*content += name
		}
		if !isEnd {
			*content += ","
		}
		comment := dataInfo.Comment
		if dataInfo.DataType != nil {
			if base.IsNotEmpty(comment) {
				comment += ", "
			}
			comment += dataInfo.DataType.Text

		}
		if base.IsNotEmpty(comment) {
			*content += " // " + comment
		}
		*content += "\n"
	}
	return
}

func appendLine(content *string, line string) {
	*content += line + "\n"
}
