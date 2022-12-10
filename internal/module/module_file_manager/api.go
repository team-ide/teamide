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
	Power           = base.AppendPower(&base.PowerAction{Action: "fileManager", Text: "文件管理器", ShouldLogin: true, StandAlone: true})
	createPower     = base.AppendPower(&base.PowerAction{Action: "create", Text: "新建文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	filePower       = base.AppendPower(&base.PowerAction{Action: "file", Text: "文件信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	filesPower      = base.AppendPower(&base.PowerAction{Action: "files", Text: "文件列表", ShouldLogin: true, StandAlone: true, Parent: Power})
	readPower       = base.AppendPower(&base.PowerAction{Action: "read", Text: "读取文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	writePower      = base.AppendPower(&base.PowerAction{Action: "write", Text: "写入文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	renamePower     = base.AppendPower(&base.PowerAction{Action: "rename", Text: "重命名文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	removePower     = base.AppendPower(&base.PowerAction{Action: "remove", Text: "删除文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	copyPower       = base.AppendPower(&base.PowerAction{Action: "copy", Text: "复制文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	movePower       = base.AppendPower(&base.PowerAction{Action: "move", Text: "移动文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	uploadPower     = base.AppendPower(&base.PowerAction{Action: "upload", Text: "上传文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	downloadPower   = base.AppendPower(&base.PowerAction{Action: "download", Text: "下载文件", ShouldLogin: true, StandAlone: true, Parent: Power})
	callActionPower = base.AppendPower(&base.PowerAction{Action: "callAction", Text: "文件操作动作", ShouldLogin: true, StandAlone: true, Parent: Power})
	callStopPower   = base.AppendPower(&base.PowerAction{Action: "callStop", Text: "文件操作停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower      = base.AppendPower(&base.PowerAction{Action: "close", Text: "文件管理器关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
	openPower       = base.AppendPower(&base.PowerAction{Action: "open", Text: "打开文件", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: Power, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Power: createPower, Do: this_.create})
	apis = append(apis, &base.ApiWorker{Power: filePower, Do: this_.file})
	apis = append(apis, &base.ApiWorker{Power: filesPower, Do: this_.files})
	apis = append(apis, &base.ApiWorker{Power: readPower, Do: this_.read})
	apis = append(apis, &base.ApiWorker{Power: writePower, Do: this_.write})
	apis = append(apis, &base.ApiWorker{Power: renamePower, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Power: removePower, Do: this_.remove})
	apis = append(apis, &base.ApiWorker{Power: copyPower, Do: this_.move})
	apis = append(apis, &base.ApiWorker{Power: movePower, Do: this_.copy})
	apis = append(apis, &base.ApiWorker{Power: uploadPower, Do: this_.upload, IsUpload: true})
	apis = append(apis, &base.ApiWorker{Power: downloadPower, Do: this_.download, IsGet: true})
	apis = append(apis, &base.ApiWorker{Power: callActionPower, Do: this_.callAction})
	apis = append(apis, &base.ApiWorker{Power: callStopPower, Do: this_.callStop})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Power: openPower, Do: this_.open, IsGet: true})
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
