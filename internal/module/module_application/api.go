package module_application

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/internal/base"
	"teamide/pkg/application"
	"teamide/pkg/application/model"
	"teamide/pkg/util"
)

type ApplicationApi struct {
	ApplicationService *ApplicationService
}

func NewApplicationApi(ApplicationService *ApplicationService) *ApplicationApi {

	return &ApplicationApi{
		ApplicationService: ApplicationService,
	}
}

var (
	// 应用 权限

	// PowerApplication 应用基本 权限
	PowerApplication        = base.AppendPower(&base.PowerAction{Action: "application", Text: "应用", ShouldLogin: true, StandAlone: false})
	PowerApplicationPage    = base.AppendPower(&base.PowerAction{Action: "application_page", Text: "应用页面", Parent: PowerApplication, ShouldLogin: true, StandAlone: false})
	PowerApplicationList    = base.AppendPower(&base.PowerAction{Action: "application_list", Text: "应用列表", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationInsert  = base.AppendPower(&base.PowerAction{Action: "application_insert", Text: "应用新增", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationUpdate  = base.AppendPower(&base.PowerAction{Action: "application_update", Text: "应用修改", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationRename  = base.AppendPower(&base.PowerAction{Action: "application_rename", Text: "应用重命名", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationDelete  = base.AppendPower(&base.PowerAction{Action: "application_delete", Text: "应用删除", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationStart   = base.AppendPower(&base.PowerAction{Action: "application_start", Text: "应用启动", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationStop    = base.AppendPower(&base.PowerAction{Action: "application_stop", Text: "应用停止", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})
	PowerApplicationRestart = base.AppendPower(&base.PowerAction{Action: "application_restart", Text: "应用重启", Parent: PowerApplicationPage, ShouldLogin: true, StandAlone: false})

	// PowerApplicationContext 应用Context 权限
	PowerApplicationContext     = base.AppendPower(&base.PowerAction{Action: "application_context", Text: "应用Context", ShouldLogin: true, StandAlone: false})
	PowerApplicationContextLoad = base.AppendPower(&base.PowerAction{Action: "application_context_load", Text: "应用Context加载", Parent: PowerApplicationContext, ShouldLogin: true, StandAlone: false})
	PowerApplicationContextSave = base.AppendPower(&base.PowerAction{Action: "application_context_save", Text: "应用Context保存", Parent: PowerApplicationContext, ShouldLogin: true, StandAlone: false})
)

func (this_ *ApplicationApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/list"}, Power: PowerApplicationList, Do: this_.apiList})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/insert"}, Power: PowerApplicationInsert, Do: this_.apiInsert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/rename"}, Power: PowerApplicationRename, Do: this_.apiRename})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/delete"}, Power: PowerApplicationDelete, Do: this_.apiDelete})

	apis = append(apis, &base.ApiWorker{Apis: []string{"application/context/load"}, Power: PowerApplicationContextLoad, Do: this_.apiContextLoad})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/context/save"}, Power: PowerApplicationContextSave, Do: this_.apiContextSave})

	apis = append(apis, &base.ApiWorker{Apis: []string{"application/model/insert"}, Power: PowerApplicationContextSave, Do: this_.apiModelInsert})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/model/save"}, Power: PowerApplicationContextSave, Do: this_.apiModelSave})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/model/delete"}, Power: PowerApplicationContextSave, Do: this_.apiModelDelete})
	apis = append(apis, &base.ApiWorker{Apis: []string{"application/model/rename"}, Power: PowerApplicationContextSave, Do: this_.apiModelRename})

	return
}

type ListRequest struct {
}

type ListResponse struct {
	List []*AppInfo `json:"list,omitempty"`
}

type AppInfo struct {
	Name string `json:"name,omitempty"`
}

func (this_ *ApplicationApi) apiList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ListRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &ListResponse{}

	appsDir := this_.ApplicationService.GetAppsDir(requestBean)

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

func (this_ *ApplicationApi) apiInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &InsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	if request.Name == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appPath := this_.ApplicationService.GetAppPath(requestBean, request.Name)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New(fmt.Sprint("应用", request.Name, "已存在"))
		return
	}
	err = os.Mkdir(appPath, 0777)
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

func (this_ *ApplicationApi) apiRename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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
	appPath := this_.ApplicationService.GetAppPath(requestBean, request.Name)
	var exist bool
	exist, err = util.PathExists(appPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.Name, "不存在"))
		return
	}
	newAppPath := this_.ApplicationService.GetAppPath(requestBean, request.Rename)
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

func (this_ *ApplicationApi) apiDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &DeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.Name)
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

func (this_ *ApplicationApi) apiContextLoad(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextLoadRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var modelContext *model.ModelContext
	modelContext, err = application.LoadContext(appDir)
	if err != nil {
		return
	}
	res = modelContext

	return
}

type ContextSaveRequest struct {
	AppName string `json:"appName,omitempty"`
	Content string `json:"content,omitempty"`
}

func (this_ *ApplicationApi) apiContextSave(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ContextSaveRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}

	var modelContext = &model.ModelContext{}
	modelContext, err = application.GetContextByText(request.Content)
	if err != nil {
		return
	}

	err = application.SaveContext(appDir, modelContext)
	if err != nil {
		return
	}
	res = modelContext

	return
}

func (this_ *ApplicationApi) checkApp(requestBean *base.RequestBean, appName string) (appDir string, err error) {
	if appName == "" {
		err = errors.New("应用名称不能为空")
		return
	}
	appDir = this_.ApplicationService.GetAppPath(requestBean, appName)
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

func (this_ *ApplicationApi) apiModelInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelInsertRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model.GetModelType(request.ModelType)
	err = application.NewWorker(appDir).ModelInsert(modelType, request.ModelName)

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

func (this_ *ApplicationApi) apiModelSave(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelSaveRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model.GetModelType(request.ModelType)
	err = application.NewWorker(appDir).ModelInsert(modelType, request.ModelName)

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

func (this_ *ApplicationApi) apiModelRename(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelRenameRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model.GetModelType(request.ModelType)
	err = application.NewWorker(appDir).ModelRename(modelType, request.ModelName, request.ModelRename)

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

func (this_ *ApplicationApi) apiModelDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ModelDeleteRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var appDir string
	appDir, err = this_.checkApp(requestBean, request.AppName)
	if err != nil {
		return
	}
	modelType := model.GetModelType(request.ModelType)
	err = application.NewWorker(appDir).ModelDelete(modelType, request.ModelName)

	if err != nil {
		return
	}

	return
}
