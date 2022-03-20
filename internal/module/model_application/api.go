package model_application

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"sort"
	"strings"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	application2 "teamide/pkg/application"
	model2 "teamide/pkg/application/model"
	"teamide/pkg/util"
)

type Service struct {
}

func (this_ *Service) BindApi(appendApi func(apis ...*base2.ApiWorker)) {
	appendApi(&base2.ApiWorker{Apis: []string{"application/list"}, Power: base2.PowerApplicationList, Do: apiList})
	appendApi(&base2.ApiWorker{Apis: []string{"application/insert"}, Power: base2.PowerApplicationInsert, Do: apiInsert})
	appendApi(&base2.ApiWorker{Apis: []string{"application/rename"}, Power: base2.PowerApplicationRename, Do: apiRename})
	appendApi(&base2.ApiWorker{Apis: []string{"application/delete"}, Power: base2.PowerApplicationDelete, Do: apiDelete})

	appendApi(&base2.ApiWorker{Apis: []string{"application/context/load"}, Power: base2.PowerApplicationContextLoad, Do: apiContextLoad})
	appendApi(&base2.ApiWorker{Apis: []string{"application/context/save"}, Power: base2.PowerApplicationContextSave, Do: apiContextSave})

	appendApi(&base2.ApiWorker{Apis: []string{"application/model/insert"}, Power: base2.PowerApplicationContextSave, Do: apiModelInsert})
	appendApi(&base2.ApiWorker{Apis: []string{"application/model/save"}, Power: base2.PowerApplicationContextSave, Do: apiModelSave})
	appendApi(&base2.ApiWorker{Apis: []string{"application/model/delete"}, Power: base2.PowerApplicationContextSave, Do: apiModelDelete})
	appendApi(&base2.ApiWorker{Apis: []string{"application/model/rename"}, Power: base2.PowerApplicationContextSave, Do: apiModelRename})
}

func GetAppsDir(requestBean *base2.RequestBean) string {
	return component.GetUserApps(requestBean.JWT)
}

func GetAppPath(requestBean *base2.RequestBean, name string) string {
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

func apiList(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base2.RequestJSON(request, c) {
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

func apiInsert(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base2.RequestJSON(request, c) {
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

func apiRename(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &RenameRequest{}
	if !base2.RequestJSON(request, c) {
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

func apiDelete(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base2.RequestJSON(request, c) {
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

func apiContextLoad(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextLoadRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var context *model2.ModelContext
	context, err = application2.LoadContext(appDir)
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

func apiContextSave(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextSaveRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var context = &model2.ModelContext{}
	context, err = application2.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application2.SaveContext(appDir, context)
	if err != nil {
		return
	}
	res = context

	return
}

func checkApp(requestBean *base2.RequestBean, appName string) (appDir string, err error) {
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

func apiModelInsert(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelInsertRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model2.GetModelType(request.ModelType)
	err = application2.NewWorker(appDir).ModelInsert(modelType, request.ModelName)

	if err != nil {
		return
	}
	return
}

type ModelSaveRequest struct {
	AppName   string `json:"appName,omitempty"`
	ModelType string `json:"modelType,omitempty"`
	ModelName string `json:"modelName,omitempty"`
	ModelText string `json:"modelText,omitempty"`
}

func apiModelSave(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelSaveRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model2.GetModelType(request.ModelType)
	err = application2.NewWorker(appDir).ModelInsert(modelType, request.ModelName)

	if err != nil {
		return
	}

	return
}

type ModelRenameRequest struct {
	AppName     string `json:"appName,omitempty"`
	ModelType   string `json:"modelType,omitempty"`
	ModelName   string `json:"modelName,omitempty"`
	ModelRename string `json:"modelRename,omitempty"`
}

func apiModelRename(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelRenameRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model2.GetModelType(request.ModelType)
	err = application2.NewWorker(appDir).ModelRename(modelType, request.ModelName, request.ModelRename)

	if err != nil {
		return
	}

	return
}

type ModelDeleteRequest struct {
	AppName   string `json:"appName,omitempty"`
	ModelType string `json:"modelType,omitempty"`
	ModelName string `json:"modelName,omitempty"`
}

func apiModelDelete(requestBean *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelDeleteRequest{}
	if !base2.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model2.GetModelType(request.ModelType)
	err = application2.NewWorker(appDir).ModelDelete(modelType, request.ModelName)

	if err != nil {
		return
	}

	return
}
