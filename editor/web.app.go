package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/application/base"

	"github.com/gin-gonic/gin"
)

type AppListRequest struct {
}

type AppListResponse struct {
	List []*AppInfo `json:"list,omitempty"`
}

type AppInfo struct {
	Name string `json:"name,omitempty"`
}

func doApiAppList(path string, c *gin.Context) (res interface{}, err error) {

	var abs string
	abs, err = filepath.Abs(appsDir)
	if err != nil {
		return
	}
	appsDirAbsolutePath := filepath.ToSlash(abs)

	var fileList []string
	err = loadDirFiles(appsDirAbsolutePath, &fileList)
	if err != nil {
		return
	}

	sort.Strings(fileList)

	response := &AppListResponse{
		List: []*AppInfo{},
	}
	for _, one := range fileList {
		appInfo := &AppInfo{
			Name: one,
		}
		response.List = append(response.List, appInfo)
	}

	res = response
	return
}

type AppInsertRequest struct {
	Name string `json:"name,omitempty"`
}

func doApiAppInsert(path string, c *gin.Context) (res interface{}, err error) {
	request := &AppInsertRequest{}
	if !RequestJSON(request, c) {
		return
	}
	filename := appsDir + "/" + request.Name

	var exists bool
	exists, err = base.PathExists(filename)
	if err != nil {
		return
	}
	if exists {
		err = errors.New(fmt.Sprint("应用[", request.Name, "]已存在"))
		return
	}

	err = WriteFile(filename, []byte(""))
	if err != nil {
		return
	}
	res = true
	return
}

type AppRenameRequest struct {
	Name   string `json:"name,omitempty"`
	Rename string `json:"rename,omitempty"`
}

func doApiAppRename(path string, c *gin.Context) (res interface{}, err error) {
	request := &AppRenameRequest{}
	if !RequestJSON(request, c) {
		return
	}
	filename := appsDir + "/" + request.Name
	refilename := appsDir + "/" + request.Rename

	var exists bool
	exists, err = base.PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("应用[", request.Name, "]不存在"))
		return
	}
	exists, err = base.PathExists(refilename)
	if err != nil {
		return
	}
	if exists {
		err = errors.New(fmt.Sprint("应用[", request.Name, "]已存在"))
		return
	}
	err = os.Rename(filename, refilename)
	if err != nil {
		return
	}
	res = true
	return
}

func loadDirFiles(dir string, fileList *[]string) (err error) {
	//获取当前目录下的所有文件或目录信息
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Walk [", dir, "] error", err)
			return err
		}
		if info.IsDir() {

		} else {
			var abs string
			abs, err = filepath.Abs(path)
			if err != nil {
				fmt.Println("Abs [", path, "] error", err)
				return err
			}
			fileAbsolutePath := filepath.ToSlash(abs)
			name := strings.TrimPrefix(fileAbsolutePath, dir)
			name = strings.TrimPrefix(name, "/")
			*fileList = append(*fileList, name)
		}
		return nil
	})
	return
}
