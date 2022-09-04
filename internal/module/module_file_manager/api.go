package module_file_manager

import (
	"errors"
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
	Power           = base.AppendPower(&base.PowerAction{Action: "file_manager", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerCreate     = base.AppendPower(&base.PowerAction{Action: "file_manager_create", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerFile       = base.AppendPower(&base.PowerAction{Action: "file_manager_file", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerFiles      = base.AppendPower(&base.PowerAction{Action: "file_manager_files", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerRead       = base.AppendPower(&base.PowerAction{Action: "file_manager_read", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerWrite      = base.AppendPower(&base.PowerAction{Action: "file_manager_write", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerRename     = base.AppendPower(&base.PowerAction{Action: "file_manager_rename", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerRemove     = base.AppendPower(&base.PowerAction{Action: "file_manager_remove", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerCopy       = base.AppendPower(&base.PowerAction{Action: "file_manager_copy", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerMove       = base.AppendPower(&base.PowerAction{Action: "file_manager_move", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerUpload     = base.AppendPower(&base.PowerAction{Action: "file_manager_upload", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerCallAction = base.AppendPower(&base.PowerAction{Action: "file_manager_call_action", Text: "工具", ShouldLogin: false, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager"}, Power: Power, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/create"}, Power: PowerCreate, Do: this_.create})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/file"}, Power: PowerFile, Do: this_.file})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/files"}, Power: PowerFiles, Do: this_.files})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/read"}, Power: PowerRead, Do: this_.read})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/write"}, Power: PowerWrite, Do: this_.write})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/rename"}, Power: PowerRename, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/remove"}, Power: PowerRemove, Do: this_.remove})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/move"}, Power: PowerCopy, Do: this_.move})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/copy"}, Power: PowerMove, Do: this_.copy})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/upload"}, Power: PowerUpload, Do: this_.upload})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/callAction"}, Power: PowerCallAction, Do: this_.callAction})
	return
}

type FileRequest struct {
	WorkerId    string `json:"workerId,omitempty"`
	Dir         string `json:"dir,omitempty"`
	Place       string `json:"place,omitempty"`
	PlaceId     string `json:"placeId,omitempty"`
	Path        string `json:"path,omitempty"`
	OldPath     string `json:"oldPath,omitempty"`
	NewPath     string `json:"newPath,omitempty"`
	IsDir       bool   `json:"isDir,omitempty"`
	FromPlace   string `json:"fromPlace,omitempty"`
	FromPlaceId string `json:"fromPlaceId,omitempty"`
	FromPath    string `json:"fromPath,omitempty"`
	Text        string `json:"text,omitempty"`
	ProgressId  string `json:"progressId,omitempty"`
	Action      string `json:"action,omitempty"`
}

func (this_ *Api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *Api) create(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = filework.Create(request.WorkerId, request.Place, request.PlaceId, request.Path, request.IsDir)
	return
}

func (this_ *Api) file(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = filework.File(request.WorkerId, request.Place, request.PlaceId, request.Path)
	return
}

func (this_ *Api) files(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	var data = map[string]interface{}{}
	data["dir"], data["files"], err = filework.Files(request.WorkerId, request.Place, request.PlaceId, request.Dir)

	res = data
	return
}

func (this_ *Api) read(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	var bytes []byte
	bytes, err = filework.Read(request.WorkerId, request.Place, request.PlaceId, request.Path)
	if err != nil {
		return
	}
	if len(bytes) > 0 {
		res = string(bytes)
	}
	return
}

func (this_ *Api) write(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = filework.Write(request.WorkerId, request.Place, request.PlaceId, request.Path, []byte(request.Text))
	if err != nil {
		return
	}
	return
}

func (this_ *Api) rename(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = filework.Rename(request.WorkerId, request.Place, request.PlaceId, request.OldPath, request.NewPath)
	return
}

func (this_ *Api) remove(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = filework.Remove(request.WorkerId, request.Place, request.PlaceId, request.Path)
	return
}

func (this_ *Api) move(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = filework.Move(request.WorkerId, request.Place, request.PlaceId, request.OldPath, request.NewPath)
	return
}

func (this_ *Api) copy(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	go filework.Copy(request.WorkerId, request.Place, request.PlaceId, request.Path, request.FromPlace, request.FromPlaceId, request.FromPath)
	return
}

func (this_ *Api) callAction(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = filework.CallAction(request.ProgressId, request.Action)
	return
}

func (this_ *Api) upload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	workerId := c.PostForm("workerId")
	if workerId == "" {
		err = errors.New("workerId获取失败")
		return
	}
	dir := c.PostForm("dir")
	if dir == "" {
		err = errors.New("dir获取失败")
		return
	}
	place := c.PostForm("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	placeId := c.PostForm("placeId")
	fullPath := c.PostForm("fullPath")
	mF, err := c.MultipartForm()
	if err != nil {
		return
	}
	fileList := mF.File["file"]

	res, err = filework.Upload(workerId, place, placeId, dir, fullPath, fileList)
	return
}
