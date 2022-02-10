package applicationService

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/application"
	"teamide/application/model"
	"teamide/server/base"
	"teamide/server/component"
	"teamide/util"

	"github.com/gin-gonic/gin"
)

func (this_ *Service) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	appendApi(&base.ApiWorker{Apis: []string{"application/list"}, Power: base.PowerApplicationList, Do: apiList})
	appendApi(&base.ApiWorker{Apis: []string{"application/insert"}, Power: base.PowerApplicationInsert, Do: apiInsert})
	appendApi(&base.ApiWorker{Apis: []string{"application/rename"}, Power: base.PowerApplicationRename, Do: apiRename})
	appendApi(&base.ApiWorker{Apis: []string{"application/delete"}, Power: base.PowerApplicationDelete, Do: apiDelete})

	appendApi(&base.ApiWorker{Apis: []string{"application/context/load"}, Power: base.PowerApplicationContextLoad, Do: apiContextLoad})
	appendApi(&base.ApiWorker{Apis: []string{"application/context/save"}, Power: base.PowerApplicationContextSave, Do: apiContextSave})

	appendApi(&base.ApiWorker{Apis: []string{"application/model/insert"}, Power: base.PowerApplicationContextSave, Do: apiModelInsert})
	appendApi(&base.ApiWorker{Apis: []string{"application/model/save"}, Power: base.PowerApplicationContextSave, Do: apiModelSave})
	appendApi(&base.ApiWorker{Apis: []string{"application/model/delete"}, Power: base.PowerApplicationContextSave, Do: apiModelDelete})
	appendApi(&base.ApiWorker{Apis: []string{"application/model/rename"}, Power: base.PowerApplicationContextSave, Do: apiModelRename})
}

func GetAppsDir(requestBean *base.RequestBean) string {
	return component.GetUserApps(requestBean.JWT)
}

func GetAppPath(requestBean *base.RequestBean, name string) string {
	appPath := GetAppsDir(requestBean) + "/" + name
	appPath = util.GetAbsolutePath(appPath)
	return appPath
}

type ListRequest struct {
}

type ListResponse struct {
	List []*AppInfo `json:"list,omitempty"`
}

type AppInfo struct {
	Name string `json:"name,omitempty"`
}

func apiList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ListResponse{}

	appsDir := GetAppsDir(requestBean)

	//获取当前目录下的文件或目录名(包含路径)
	filePaths, err := filepath.Glob(filepath.Join(appsDir, "*"))
	if err != nil {
		return
	}

	sort.Strings(filePaths)

	for _, filePath := range filePaths {

		var abs string
		abs, err = filepath.Abs(filePath)
		if err != nil {
			return
		}
		fileAbsolutePath := filepath.ToSlash(abs)
		name := strings.TrimPrefix(fileAbsolutePath, appsDir)
		name = strings.TrimPrefix(name, "/")

		appInfo := &AppInfo{
			Name: name,
		}
		response.List = append(response.List, appInfo)
	}

	res = response
	return
}

type InsertRequest struct {
	Name string `json:"name,omitempty"`
}

type InsertResponse struct {
}

func apiInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	if request.Name == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appPath := GetAppPath(requestBean, request.Name)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New(fmt.Sprint("应用", request.Name, "已存在"))
		return
	}
	err = os.Mkdir(appPath, os.ModeDir)
	if err != nil {
		return
	}

	response := &InsertResponse{}
	res = response
	return
}

type RenameRequest struct {
	Name   string `json:"name,omitempty"`
	Rename string `json:"rename,omitempty"`
}

type RenameResponse struct {
}

func apiRename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Name == "" {
		err = errors.New("原应用名称不能为空")
		return
	}
	if request.Rename == "" {
		err = errors.New("新应用名称不能为空")
		return
	}
	appPath := GetAppPath(requestBean, request.Name)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.Name, "不存在"))
		return
	}
	newAppPath := GetAppPath(requestBean, request.Rename)
	exist, err = util.PathExists(newAppPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New(fmt.Sprint("应用", request.Rename, "已存在"))
		return
	}
	err = os.Rename(appPath, newAppPath)
	if err != nil {
		return
	}
	response := &RenameResponse{}

	res = response
	return
}

type DeleteRequest struct {
	Name string `json:"name,omitempty"`
}

type DeleteResponse struct {
}

func apiDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.Name)
	if err != nil {
		return
	}
	err = os.Remove(appDir)
	if err != nil {
		return
	}
	response := &DeleteResponse{}

	res = response
	return
}

type ContextLoadRequest struct {
	AppName string `json:"appName,omitempty"`
}

func apiContextLoad(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextLoadRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var context *model.ModelContext
	context, err = application.LoadContext(appDir)
	if err != nil {
		return
	}
	res = context

	return
}

type ContextSaveRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func apiContextSave(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextSaveRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var context = &model.ModelContext{}
	context, err = application.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application.SaveContext(appDir, context)
	if err != nil {
		return
	}
	res = context

	return
}

func checkApp(requestBean *base.RequestBean, appName string) (appDir string, err error) {
	if appName == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appDir = GetAppPath(requestBean, appName)
	var exist bool
	exist, err = util.PathExists(appDir)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", appName, "不存在"))
		return
	}
	return
}

type ModelInsertRequest struct {
	AppName   string `json:"appName,omitempty"`
	ModelType string `json:"modelType,omitempty"`
	ModelName string `json:"modelName,omitempty"`
}

func apiModelInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelInsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model.GetModelType(request.ModelType)
	modelPath := appDir + "/" + modelType.Dir + "/" + request.ModelName + ".yaml"
	var exist bool
	exist, err = util.PathExists(modelPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New(fmt.Sprint("应用模型", request.ModelName, "已存在"))
		return
	}
	var file *os.File
	file, err = os.Create(modelPath)
	if err != nil {
		return
	}
	defer file.Close()
	return
}

type ModelSaveRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func apiModelSave(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelSaveRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	if request.AppName == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appPath := GetAppPath(requestBean, request.AppName)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.AppName, "不存在"))
		return
	}

	var context = &model.ModelContext{}
	context, err = application.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application.SaveContext(appPath, context)
	if err != nil {
		return
	}
	res = context

	return
}

type ModelRenameRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func apiModelRename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelRenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	if request.AppName == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appPath := GetAppPath(requestBean, request.AppName)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.AppName, "不存在"))
		return
	}

	var context = &model.ModelContext{}
	context, err = application.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application.SaveContext(appPath, context)
	if err != nil {
		return
	}
	res = context

	return
}

type ModelDeleteRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func apiModelDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelDeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	if request.AppName == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appPath := GetAppPath(requestBean, request.AppName)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.AppName, "不存在"))
		return
	}

	var context = &model.ModelContext{}
	context, err = application.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application.SaveContext(appPath, context)
	if err != nil {
		return
	}
	res = context

	return
}
