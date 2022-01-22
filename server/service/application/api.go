package applicationService

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"teamide/application/model"
	"teamide/server/base"
	"teamide/server/component"
	"teamide/server/config"
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
}

func GetAppsDir(requestBean *base.RequestBean) string {
	dir := config.Config.Server.Data + "/apps/"
	if config.Config.IsNative {
		dir = config.Config.Server.Data + "/native/apps/"
	}
	dir = util.GetAbsolutePath(dir)

	var exist bool
	exist, _ = util.PathExists(dir)
	if !exist {
		os.MkdirAll(dir, 0777)
	}
	return dir
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
	var files []string
	err = util.LoadDirFilenames(&files, GetAppsDir(requestBean))
	if err != nil {
		return
	}
	for _, one := range files {
		appInfo := &AppInfo{
			Name: one,
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
	var f *os.File
	f, err = os.Create(appPath)
	if err != nil {
		return
	}
	defer f.Close()

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
	if !exist {
		err = errors.New(fmt.Sprint("应用", request.Name, "不存在"))
		return
	}
	err = os.Remove(appPath)
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

	var bs []byte
	bs, err = util.ReadFile(appPath)
	if err != nil {
		return
	}
	if len(bs) > 0 {
		bs, err = base.AesCBCDecrypt(bs, []byte(component.GetKey()))
		if err != nil {
			return
		}
		err = base.ToBean(bs, context)
		if err != nil {
			return
		}
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

	err = base.ToBean([]byte(request.Content), context)
	if err != nil {
		return
	}
	var bs []byte
	bs, err = json.Marshal(context)
	if err != nil {
		return
	}
	bs, err = base.AesCBCEncrypt(bs, []byte(component.GetKey()))
	if err != nil {
		return
	}
	err = util.WriteFile(appPath, bs)
	if err != nil {
		return
	}
	res = context

	return
}
