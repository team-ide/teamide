package ssh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"teamide/pkg/util"
	"time"
)

type SFTPRequest struct {
	Work          string `json:"work,omitempty"`
	WorkId        string `json:"workId,omitempty"`
	Dir           string `json:"dir,omitempty"`
	Place         string `json:"place,omitempty"`
	Path          string `json:"path,omitempty"`
	FullPath      string `json:"fullPath,omitempty"`
	Name          string `json:"name,omitempty"`
	OldPath       string `json:"oldPath,omitempty"`
	NewPath       string `json:"newPath,omitempty"`
	File          *multipart.FileHeader
	FromFilePlace string         `json:"fromFilePlace,omitempty"`
	FromFile      *util.FileInfo `json:"fromFile,omitempty"`
	ToFilePlace   string         `json:"toFilePlace,omitempty"`
	ToFile        *util.FileInfo `json:"toFile,omitempty"`
	ConfirmId     string         `json:"confirmId,omitempty"`
	IsOk          bool           `json:"isOk,omitempty"`
	IsCancel      bool           `json:"isCancel,omitempty"`
	ScrollTop     int            `json:"scrollTop,omitempty"`
	Text          string         `json:"text,omitempty"`
	Pattern       string         `json:"pattern,omitempty"`
	IsDir         bool           `json:"isDir,omitempty"`
	IsNew         bool           `json:"isNew,omitempty"`
}
type SFTPResponse struct {
	Work      string           `json:"work,omitempty"`
	WorkId    string           `json:"workId,omitempty"`
	Dir       string           `json:"dir,omitempty"`
	Msg       string           `json:"msg,omitempty"`
	Files     []*util.FileInfo `json:"files,omitempty"`
	Place     string           `json:"place,omitempty"`
	Path      string           `json:"path,omitempty"`
	Name      string           `json:"name,omitempty"`
	ScrollTop int              `json:"scrollTop,omitempty"`
	Text      string           `json:"text,omitempty"`
	Pattern   string           `json:"pattern,omitempty"`
	Size      int64            `json:"size,omitempty"`
}

func (this_ *SftpClient) callConfirm(confirmInfo *util.FileConfirmInfo) (res *util.FileConfirmInfo, err error) {

	if this_.confirmMap == nil {
		this_.confirmMap = map[string]chan *util.FileConfirmInfo{}
	}
	confirmInfo.IsConfirm = true
	if confirmInfo.ConfirmId == "" {
		confirmInfo.ConfirmId = util.UUID()
	}
	this_.confirmMap[confirmInfo.ConfirmId] = make(chan *util.FileConfirmInfo, 1)
	bs, err := json.Marshal(confirmInfo)
	if err != nil {
		util.Logger.Error("call confirm to json err error", zap.Error(err))
		return
	}
	this_.WSWriteText(bs)
	res = <-this_.confirmMap[confirmInfo.ConfirmId]

	close(this_.confirmMap[confirmInfo.ConfirmId])
	delete(this_.confirmMap, confirmInfo.ConfirmId)
	return

}

func (this_ *SftpClient) callProgress(request *SFTPRequest, progress interface{}) {
	for {
		time.Sleep(100 * time.Millisecond)

		if this_.isClosedWS {
			return
		}
		var waitCall bool
		var endTime int64 = -1
		UploadProgress, UploadProgressOk := progress.(*util.FileUploadProgress)
		if UploadProgressOk {
			UploadProgress.Timestamp = util.GetNowTime()
			endTime = UploadProgress.EndTime
			waitCall = UploadProgress.WaitCall
		}

		CopyProgress, CopyProgressOK := progress.(*util.FileCopyProgress)
		if CopyProgressOK {
			CopyProgress.Timestamp = util.GetNowTime()
			endTime = CopyProgress.EndTime
			waitCall = CopyProgress.WaitCall
		}

		RemoveProgress, RemoveProgressOk := progress.(*util.FileRemoveProgress)
		if RemoveProgressOk {
			RemoveProgress.Timestamp = util.GetNowTime()
			endTime = RemoveProgress.EndTime
			waitCall = RemoveProgress.WaitCall
		}
		if endTime == -1 {
			return
		}
		if waitCall {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		out := map[string]interface{}{
			"work":       request.Work,
			"workId":     request.WorkId,
			"isProgress": true,
			"progress":   progress,
		}

		bs, err := json.Marshal(out)
		if err != nil {
			util.Logger.Error("sftp upload progress to json err", zap.Error(err))
			continue
		}
		this_.WSWriteText(bs)

		if endTime > 0 {
			break
		}
	}
}

func (this_ *SftpClient) work(request *SFTPRequest) {
	response := &SFTPResponse{}
	var err error
	switch request.Work {
	case "confirmResult":
		if this_.confirmMap == nil {
			return
		}
		this_.confirmMap[request.ConfirmId] <- &util.FileConfirmInfo{
			ConfirmId: request.ConfirmId,
			IsCancel:  request.IsCancel,
			IsOk:      request.IsOk,
		}
		return

	case "files":
		if request.Place == "local" {
			response, err = this_.localFiles(request)
		} else if request.Place == "remote" {
			response, err = this_.remoteFiles(request)
		}
	case "upload":
		if request.File == nil {
			err = errors.New("上传文件丢失")
			break
		}
		progress := &util.FileUploadProgress{
			StartTime: util.GetNowTime(),
		}
		go this_.callProgress(request, progress)
		if request.Place == "local" {
			response, err = this_.localUpdate(request, progress)
		} else if request.Place == "remote" {
			response, err = this_.remoteUpdate(request, progress)
		}
	case "copy":
		if request.FromFile == nil {
			err = errors.New("源文件信息丢失")
			break
		}
		if request.ToFile == nil {
			err = errors.New("目标文件信息丢失")
			break
		}
		progress := &util.FileCopyProgress{
			StartTime: util.GetNowTime(),
		}
		go this_.callProgress(request, progress)
		response, err = this_.copy(request, progress)
	case "remove":
		progress := &util.FileRemoveProgress{
			StartTime: util.GetNowTime(),
		}
		go this_.callProgress(request, progress)
		if request.Place == "local" {
			response, err = this_.localRemove(request, progress)
		} else if request.Place == "remote" {
			response, err = this_.remoteRemove(request, progress)
		}
	case "rename":
		if request.Place == "local" {
			response, err = this_.localRename(request)
		} else if request.Place == "remote" {
			response, err = this_.remoteRename(request)
		}
	}
	if response == nil {
		response = &SFTPResponse{}
	}
	if err != nil {
		util.Logger.Error("ssh ftp work{"+request.Work+"} error", zap.Error(err))
		response.Msg = err.Error()
	}
	response.Work = request.Work
	response.WorkId = request.WorkId
	response.Place = request.Place
	response.ScrollTop = request.ScrollTop

	this_.WSWriteData(response)

	return
}

func (this_ *SftpClient) localUpdate(request *SFTPRequest, progress *util.FileUploadProgress) (response *SFTPResponse, err error) {

	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	progress.Size = request.File.Size
	defer func() {
		progress.EndTime = util.GetNowTime()
	}()

	path := request.Dir + "/" + request.File.Filename
	if request.FullPath != "" {
		path = request.Dir + "/" + strings.TrimPrefix(request.FullPath, "/")
	}
	response = &SFTPResponse{
		Path: path,
		Dir:  request.Dir,
	}

	err = util.FileUpload(os.Lstat,
		func(s string) error {
			return os.MkdirAll(s, os.ModePerm)
		},
		func() (io.Reader, error) {
			return request.File.Open()
		},
		util.LocalFileWrite,
		request.File.Size, request.File.Filename, path, this_.callConfirm, progress)

	if err != nil {
		return
	}
	return
}

func (this_ *SftpClient) remoteUpdate(request *SFTPRequest, progress *util.FileUploadProgress) (response *SFTPResponse, err error) {

	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	progress.Size = request.File.Size
	defer func() {
		progress.EndTime = util.GetNowTime()
	}()

	path := request.Dir + "/" + request.File.Filename
	if request.FullPath != "" {
		path = request.Dir + "/" + strings.TrimPrefix(request.FullPath, "/")
	}
	response = &SFTPResponse{
		Path: path,
		Dir:  request.Dir,
	}

	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	err = util.FileUpload(sftpClient.Lstat, sftpClient.MkdirAll,
		func() (io.Reader, error) {
			return request.File.Open()
		},
		func(s string) (io.Writer, error) {
			return sftpClient.Create(s)
		},
		request.File.Size, request.File.Filename, path, this_.callConfirm, progress)

	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) localDownload(c *gin.Context, path string) (err error) {

	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	var fileInfo *os.File
	fileInfo, err = os.Open(path)
	if err != nil {
		return
	}
	defer closeFile(fileInfo)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	err = util.CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SftpClient) remoteDownload(c *gin.Context, path string) (err error) {

	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	var fileName string
	var fileSize int64
	ff, err := sftpClient.Lstat(path)
	if err != nil {
		return
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	var fileInfo *sftp.File
	fileInfo, err = sftpClient.Open(path)
	if err != nil {
		return
	}
	defer closeFtpFile(fileInfo)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	err = util.CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SftpClient) copy(request *SFTPRequest, progress *util.FileCopyProgress) (response *SFTPResponse, err error) {

	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}
	defer func() {
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()

	var sftpClient *sftp.Client
	if request.FromFilePlace == "remote" || request.ToFilePlace == "remote" {
		sftpClient, err = this_.newSftp()
		if err != nil {
			return
		}
		defer this_.closeSftClient(sftpClient)
	}

	var fromLstat func(string) (os.FileInfo, error)
	var fromLoadFiles func(string) ([]os.FileInfo, error)
	var fromReader func(string) (io.Reader, error)
	var toLstat func(string) (os.FileInfo, error)
	var toMkdirAll func(string) error
	var toWriter func(string) (io.Writer, error)
	var callConfirm func(*util.FileConfirmInfo) (*util.FileConfirmInfo, error)

	callConfirm = this_.callConfirm
	if request.FromFilePlace == "remote" {
		fromLstat = sftpClient.Lstat
		fromLoadFiles = sftpClient.ReadDir
		fromReader = func(s string) (io.Reader, error) {
			return sftpClient.Open(s)
		}
	} else {
		fromLstat = os.Lstat
		fromLoadFiles = util.LocalLoadFiles
		fromReader = util.LocalFileOpen
	}
	if request.ToFilePlace == "remote" {
		toLstat = sftpClient.Lstat
		toMkdirAll = sftpClient.MkdirAll
		toWriter = func(s string) (io.Writer, error) {
			return sftpClient.Create(s)
		}
	} else {
		toLstat = os.Lstat
		toMkdirAll = func(s string) error {
			return os.MkdirAll(s, os.ModePerm)
		}
		toWriter = util.LocalFileWrite
	}

	err = util.FileCopy(fromLstat, fromLoadFiles, fromReader, request.FromFile.Path, toLstat, toMkdirAll, toWriter, request.ToFile.Path, callConfirm, progress)
	if err != nil {
		return
	}
	return
}

func closeFile(obj *os.File) {
	if obj == nil {
		return
	}
	_ = obj.Close()
}

func closeFtpFile(obj *sftp.File) {
	if obj == nil {
		return
	}
	_ = obj.Close()
}
func closeIfCloser(obj interface{}) {
	if obj == nil {
		return
	}
	closer, closerOk := obj.(io.Closer)
	if closerOk {
		_ = closer.Close()
	}
}

func (this_ *SftpClient) localRemove(request *SFTPRequest, progress *util.FileRemoveProgress) (response *SFTPResponse, err error) {

	defer func() {
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()

	err = util.FileRemove(os.Lstat, util.LocalLoadFiles, request.Path, progress)

	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) remoteRemove(request *SFTPRequest, progress *util.FileRemoveProgress) (response *SFTPResponse, err error) {

	defer func() {
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	err = util.FileRemove(sftpClient.Lstat, sftpClient.ReadDir, request.Path, progress)

	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) localRename(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.OldPath,
		Dir:  request.Dir,
	}

	if request.IsNew {
		if request.IsDir {
			err = os.MkdirAll(request.NewPath, os.ModePerm)
		} else {
			var f *os.File
			f, err = os.Create(request.NewPath)
			defer closeFile(f)
		}
		if err != nil {
			return
		}

		return
	}

	_, err = os.Lstat(request.OldPath)
	if err != nil {
		return
	}
	err = os.Rename(request.OldPath, request.NewPath)
	if err != nil {
		return
	}
	return
}

func (this_ *SftpClient) remoteRename(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.OldPath,
		Dir:  request.Dir,
	}
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	if request.IsNew {
		if request.IsDir {
			err = sftpClient.MkdirAll(request.NewPath)
		} else {
			var f *sftp.File
			f, err = sftpClient.Create(request.NewPath)
			defer closeIfCloser(f)
		}
		if err != nil {
			return
		}

		return
	}

	_, err = sftpClient.Lstat(request.OldPath)
	if err != nil {
		return
	}

	err = sftpClient.Rename(request.OldPath, request.NewPath)
	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) localFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	dir := request.Dir
	if dir == "" {
		dir, err = os.UserHomeDir()
		if err != nil {
			return
		}
	}

	dir = util.FormatPath(dir)
	if err != nil {
		return
	}
	response.Dir = dir

	response.Files, err = util.LoadFiles(os.Lstat, util.LocalLoadFiles, dir)
	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) remoteFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	dir := request.Dir
	if dir == "" {
		dir, err = sftpClient.Getwd()
		if err != nil {
			return
		}
	}
	dir, err = sftpClient.RealPath(dir)
	if err != nil {
		return
	}
	response.Dir = dir

	response.Files, err = util.LoadFiles(sftpClient.Lstat, sftpClient.ReadDir, dir)
	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) LocalReadText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	response.Size, response.Text, err = util.FileRead(os.Lstat, util.LocalFileOpen, request.Path, 1024*1024*10)
	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) RemoteReadText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	response.Size, response.Text, err = util.FileRead(sftpClient.Lstat, func(s string) (io.Reader, error) {
		return sftpClient.Open(s)
	}, request.Path, 1024*1024*10)
	if err != nil {
		return
	}
	return
}

func (this_ *SftpClient) LocalSaveText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}

	err = util.FileWrite(os.Lstat, util.LocalFileWrite, request.Path, []byte(request.Text))
	if err != nil {
		return
	}
	return
}

func (this_ *SftpClient) RemoteSaveText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	err = util.FileWrite(sftpClient.Lstat, func(s string) (io.Writer, error) {
		return sftpClient.Create(s)
	}, request.Path, []byte(request.Text))
	if err != nil {
		return
	}
	return
}
