package module_file_manager

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"teamide/internal/base"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/ssh"
	"teamide/pkg/vitess/bytes2"
)

type api struct {
	*worker
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *api {
	return &api{
		worker: NewWorker(toolboxService_, nodeService_),
	}
}

var (
	// 文件管理器 权限

	// Power 文件管理器 基本 权限
	Power           = base.AppendPower(&base.PowerAction{Action: "file_manager", Text: "文件管理器", ShouldLogin: true, StandAlone: true})
	PowerCreate     = base.AppendPower(&base.PowerAction{Action: "file_manager_create", Text: "新建文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerFile       = base.AppendPower(&base.PowerAction{Action: "file_manager_file", Text: "文件信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerFiles      = base.AppendPower(&base.PowerAction{Action: "file_manager_files", Text: "文件列表", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerRead       = base.AppendPower(&base.PowerAction{Action: "file_manager_read", Text: "读取文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerWrite      = base.AppendPower(&base.PowerAction{Action: "file_manager_write", Text: "写入文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerRename     = base.AppendPower(&base.PowerAction{Action: "file_manager_rename", Text: "重命名文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerRemove     = base.AppendPower(&base.PowerAction{Action: "file_manager_remove", Text: "删除文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerCopy       = base.AppendPower(&base.PowerAction{Action: "file_manager_copy", Text: "复制文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerMove       = base.AppendPower(&base.PowerAction{Action: "file_manager_move", Text: "移动文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerUpload     = base.AppendPower(&base.PowerAction{Action: "file_manager_upload", Text: "上传文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerDownload   = base.AppendPower(&base.PowerAction{Action: "file_manager_download", Text: "下载文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerCallAction = base.AppendPower(&base.PowerAction{Action: "file_manager_call_action", Text: "文件操作动作", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerCallStop   = base.AppendPower(&base.PowerAction{Action: "file_manager_call_stop", Text: "文件操作停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerClose      = base.AppendPower(&base.PowerAction{Action: "file_manager_close", Text: "文件管理器关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
	PowerOpen       = base.AppendPower(&base.PowerAction{Action: "file_manager_open", Text: "打开文件", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
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
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/upload"}, Power: PowerUpload, Do: this_.upload, IsUpload: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/download"}, Power: PowerDownload, Do: this_.download, IsGet: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/callAction"}, Power: PowerCallAction, Do: this_.callAction})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/callStop"}, Power: PowerCallStop, Do: this_.callStop})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/close"}, Power: PowerClose, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Apis: []string{"file_manager/open"}, Power: PowerOpen, Do: this_.open, IsGet: true})
	return
}

type FileRequest struct {
	WorkerId      string `json:"workerId,omitempty"`
	FileWorkerKey string `json:"fileWorkerKey,omitempty"`
	Dir           string `json:"dir,omitempty"`
	Place         string `json:"place,omitempty"`
	PlaceId       string `json:"placeId,omitempty"`
	Path          string `json:"path,omitempty"`
	OldPath       string `json:"oldPath,omitempty"`
	NewPath       string `json:"newPath,omitempty"`
	IsDir         bool   `json:"isDir,omitempty"`
	FromPlace     string `json:"fromPlace,omitempty"`
	FromPlaceId   string `json:"fromPlaceId,omitempty"`
	FromPath      string `json:"fromPath,omitempty"`
	Text          string `json:"text,omitempty"`
	ProgressId    string `json:"progressId,omitempty"`
	Action        string `json:"action,omitempty"`
	Force         bool   `json:"force,omitempty"`
}

func (this_ *api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *api) close(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	this_.Close(request.WorkerId)
	ssh.CloseFileService(request.WorkerId)
	return
}

func (this_ *api) create(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = this_.Create(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path, request.IsDir)
	return
}

func (this_ *api) file(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = this_.File(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path)
	return
}

func (this_ *api) files(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	var data = map[string]interface{}{}
	data["dir"], data["files"], err = this_.Files(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Dir)

	res = data
	return
}

func (this_ *api) read(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	response := map[string]interface{}{}
	res = response

	fileInfo, err := this_.File(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path)
	if err != nil {
		return
	}
	response["path"] = request.Path
	response["file"] = fileInfo
	if fileInfo.IsDir {
		err = errors.New("路径[" + request.Path + "]为目录，无法打开!")
		return
	}
	if !request.Force {
		if fileInfo.Size > 10*1024*1024 {
			err = base.NewBaseError(base.FileSizeOversizeErrCode, "文件过大[", fileInfo.Size, "]，无法打开!")
			return
		}
	}

	writer := &bytes2.Buffer{}
	_, err = this_.Read(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path, writer)
	if err != nil {
		return
	}
	if writer.Len() > 0 {
		bytes := writer.Bytes()
		response["text"] = string(bytes)
	}
	return
}

func (this_ *api) write(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	reader := strings.NewReader(request.Text)
	res, err = this_.Write(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path, reader, reader.Len())
	if err != nil {
		return
	}
	return
}

func (this_ *api) rename(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	res, err = this_.Rename(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.OldPath, request.NewPath)
	return
}

func (this_ *api) remove(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = this_.Remove(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path)
	return
}

func (this_ *api) move(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = this_.Move(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.OldPath, request.NewPath)
	return
}

func (this_ *api) copy(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	go this_.Copy(request.WorkerId, request.FileWorkerKey, request.Place, request.PlaceId, request.Path, request.FromPlace, request.FromPlaceId, request.FromPath)
	return
}

func (this_ *api) callAction(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = this_.CallAction(request.ProgressId, request.Action)
	return
}

func (this_ *api) callStop(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &FileRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	err = this_.CallStop(request.ProgressId)
	return
}

func (this_ *api) upload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	workerId := c.PostForm("workerId")
	if workerId == "" {
		err = errors.New("workerId获取失败")
		return
	}
	fileWorkerKey := c.PostForm("fileWorkerKey")
	if fileWorkerKey == "" {
		err = errors.New("fileWorkerKey获取失败")
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

	res, err = this_.Upload(workerId, fileWorkerKey, place, placeId, dir, fullPath, fileList)
	return
}

func (this_ *api) download(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	data := map[string]string{}

	err = c.Bind(&data)
	if err != nil {
		return
	}

	workerId := data["workerId"]
	fileWorkerKey := data["fileWorkerKey"]
	place := data["place"]
	placeId := data["placeId"]
	path := data["path"]

	fileInfo, err := this_.File(workerId, fileWorkerKey, place, placeId, path)
	if err != nil {
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileInfo.Name)))

	// 此处不设置 文件大小，如果设置文件大小，将无法终止下载
	//c.Header("Content-Length", fmt.Sprint(fileInfo.Size))
	c.Header("download-file-name", fileInfo.Name)

	_, err = this_.Read(workerId, fileWorkerKey, place, placeId, path, &cWriter{
		c: c,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusOK)
		err = nil
		this_.Logger.Warn("file manager download file error", zap.Error(err))
		return
	}
	c.Status(http.StatusOK)
	return
}

type cWriter struct {
	c *gin.Context
}

func (this_ *cWriter) Write(buf []byte) (n int, err error) {
	n, err = this_.c.Writer.Write(buf)
	return
}

func (this_ *api) open(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	data := map[string]string{}

	err = c.Bind(&data)
	if err != nil {
		return
	}

	workerId := data["workerId"]
	fileWorkerKey := data["fileWorkerKey"]
	place := data["place"]
	placeId := data["placeId"]
	path := data["path"]

	_, err = this_.Read(workerId, fileWorkerKey, place, placeId, path, &cWriter{
		c: c,
	})
	if err != nil {
		return
	}
	c.Status(http.StatusOK)
	return
}
