package module_file_manager

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/filework"
)

type Api struct {
	*context.ServerContext
	ToolboxService *module_toolbox.ToolboxService
}

func NewApi(ToolboxService *module_toolbox.ToolboxService) *Api {
	return &Api{
		ServerContext:  ToolboxService.ServerContext,
		ToolboxService: ToolboxService,
	}
}

var (
	// 文件管理器 权限

	// Power 文件管理器 基本 权限
	Power       = base.AppendPower(&base.PowerAction{Action: "file_manager", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerFile   = base.AppendPower(&base.PowerAction{Action: "file_manager_file", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerFiles  = base.AppendPower(&base.PowerAction{Action: "file_manager_files", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerRead   = base.AppendPower(&base.PowerAction{Action: "file_manager_read", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerRename = base.AppendPower(&base.PowerAction{Action: "file_manager_rename", Text: "工具", ShouldLogin: false, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager"}, Power: Power, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/file"}, Power: PowerFile, Do: this_.file})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/files"}, Power: PowerFiles, Do: this_.files})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/read"}, Power: PowerRead, Do: this_.read})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/rename"}, Power: PowerRename, Do: this_.rename})
	return
}

type FileRequest struct {
	Dir           string `json:"dir,omitempty"`
	Place         string `json:"place,omitempty"`
	PlaceId       string `json:"placeId,omitempty"`
	Path          string `json:"path,omitempty"`
	FullPath      string `json:"fullPath,omitempty"`
	Name          string `json:"name,omitempty"`
	OldPath       string `json:"oldPath,omitempty"`
	NewPath       string `json:"newPath,omitempty"`
	FromFilePlace string `json:"fromFilePlace,omitempty"`
	ToFilePlace   string `json:"toFilePlace,omitempty"`
	ConfirmId     string `json:"confirmId,omitempty"`
	IsDir         bool   `json:"isDir,omitempty"`
	IsNew         bool   `json:"isNew,omitempty"`
}

func (this_ *Api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *Api) file(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = filework.File(request.Place, request.PlaceId, request.Path)
	return
}

func (this_ *Api) files(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = filework.Files(request.Place, request.PlaceId, request.Dir)
	return
}

func (this_ *Api) read(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var bytes []byte
	bytes, err = filework.Read(request.Place, request.PlaceId, request.Path)
	if err != nil {
		return
	}
	if len(bytes) > 0 {
		res = string(bytes)
	}
	return
}

func (this_ *Api) rename(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = filework.Rename(request.Place, request.PlaceId, request.OldPath, request.NewPath)
	return
}
