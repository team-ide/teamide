package ssh

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"teamide/internal/context"
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
	WorkerId  string           `json:"workerId,omitempty"`
	Work      string           `json:"work,omitempty"`
	WorkId    string           `json:"workId,omitempty"`
	Dir       string           `json:"dir,omitempty"`
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
	confirmInfo.WorkerId = this_.WorkerId
	this_.confirmMap[confirmInfo.ConfirmId] = make(chan *util.FileConfirmInfo, 1)

	context.ServerWebsocketOutEvent("ftp-data", confirmInfo)

	res = <-this_.confirmMap[confirmInfo.ConfirmId]

	close(this_.confirmMap[confirmInfo.ConfirmId])
	delete(this_.confirmMap, confirmInfo.ConfirmId)
	return

}

func (this_ *SftpClient) callProgress(request *SFTPRequest, progress interface{}) {
	var needBreak bool
	for {
		time.Sleep(200 * time.Millisecond)

		var waitCall bool
		var endTime int64 = -1
		var err error
		switch p := progress.(type) {
		case *util.FileUploadProgress:
			p.Timestamp = util.GetNowTime()
			endTime = p.EndTime
			waitCall = p.WaitCall
			err = p.Error
			break
		case *util.FileCopyProgress:
			p.Timestamp = util.GetNowTime()
			endTime = p.EndTime
			waitCall = p.WaitCall
			err = p.Error
		case *util.FileRemoveProgress:
			p.Timestamp = util.GetNowTime()
			endTime = p.EndTime
			waitCall = p.WaitCall
			err = p.Error
			break
		case *util.FileRenameProgress:
			p.Timestamp = util.GetNowTime()
			endTime = p.EndTime
			waitCall = p.WaitCall
			err = p.Error
			break
		case *util.FileDownloadProgress:
			p.Timestamp = util.GetNowTime()
			endTime = p.EndTime
			waitCall = p.WaitCall
			err = p.Error
			break
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
			"place":      request.Place,
			"dir":        request.Dir,
			"isProgress": true,
			"progress":   progress,
			"workerId":   this_.WorkerId,
		}
		if err != nil {
			out["error"] = err.Error()
		}
		context.ServerWebsocketOutEvent("ftp-data", out)
		if needBreak {
			break
		}
		if endTime > 0 {
			needBreak = true
		}
	}
}

func (this_ *SftpClient) Work(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{}
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
		if request.Place == "local" {
			go this_.callProgress(request, progress)
			go this_.localUpdate(request, progress)
		} else if request.Place == "remote" {
			go this_.callProgress(request, progress)
			go this_.remoteUpdate(request, progress)
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
		go this_.copy(request, progress)
	case "remove":
		progress := &util.FileRemoveProgress{
			StartTime: util.GetNowTime(),
		}
		if request.Place == "local" {
			go this_.callProgress(request, progress)
			go this_.localRemove(request, progress)
		} else if request.Place == "remote" {
			go this_.callProgress(request, progress)
			go this_.remoteRemove(request, progress)
		}
	case "rename":
		progress := &util.FileRenameProgress{
			StartTime: util.GetNowTime(),
		}
		go this_.callProgress(request, progress)
		if request.Place == "local" {
			err = this_.localRename(request, progress)
		} else if request.Place == "remote" {
			err = this_.remoteRename(request, progress)
		}
	}
	if response == nil {
		response = &SFTPResponse{}
	}
	if err != nil {
		return
	}
	response.Work = request.Work
	response.WorkId = request.WorkId
	response.Place = request.Place
	response.ScrollTop = request.ScrollTop

	return
}

func (this_ *SftpClient) localUpdate(request *SFTPRequest, progress *util.FileUploadProgress) {
	var err error

	defer func() {
		if err != nil {
			progress.Error = err
			progress.EndTime = util.GetNowTime()
		}
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	progress.Size = request.File.Size

	path := request.Dir + "/" + request.File.Filename
	if request.FullPath != "" {
		path = request.Dir + "/" + strings.TrimPrefix(request.FullPath, "/")
	}

	util.FileUpload(os.Lstat,
		func(s string) error {
			return os.MkdirAll(s, os.ModePerm)
		},
		func() (io.Reader, error) {
			return request.File.Open()
		},
		util.LocalFileWrite,
		request.File.Size, request.File.Filename, path, this_.callConfirm, progress)

	return
}

func (this_ *SftpClient) remoteUpdate(request *SFTPRequest, progress *util.FileUploadProgress) {
	var err error

	defer func() {
		if err != nil {
			progress.Error = err
			progress.EndTime = util.GetNowTime()
		}
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	progress.Size = request.File.Size

	path := request.Dir + "/" + request.File.Filename
	if request.FullPath != "" {
		path = request.Dir + "/" + strings.TrimPrefix(request.FullPath, "/")
	}

	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

	util.FileUpload(sftpClient.Lstat, sftpClient.MkdirAll,
		func() (io.Reader, error) {
			return request.File.Open()
		},
		func(s string) (io.Writer, error) {
			return sftpClient.Create(s)
		},
		request.File.Size, request.File.Filename, path, this_.callConfirm, progress)

	return
}

func (this_ *SftpClient) LocalDownload(c *gin.Context, workId string, path string) (err error) {
	progress := &util.FileDownloadProgress{}
	defer func() {
		progress.Error = err
		progress.SuccessCount = 1
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	go this_.callProgress(&SFTPRequest{WorkId: workId, Work: "download"}, progress)

	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	fileName = ff.Name()
	fileSize = ff.Size()
	progress.Size = fileSize

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
		progress.SuccessSize += writeSize
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SftpClient) RemoteDownload(c *gin.Context, workId string, path string) (err error) {
	progress := &util.FileDownloadProgress{}
	defer func() {
		progress.Error = err
		progress.SuccessCount = 1
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1
	go this_.callProgress(&SFTPRequest{WorkId: workId, Work: "download"}, progress)

	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

	var fileName string
	var fileSize int64
	ff, err := sftpClient.Lstat(path)
	if err != nil {
		return
	}
	fileName = ff.Name()
	fileSize = ff.Size()
	progress.Size = fileSize

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
		progress.SuccessSize += writeSize
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SftpClient) copy(request *SFTPRequest, progress *util.FileCopyProgress) {
	var err error

	defer func() {
		if err != nil {
			progress.Error = err
			progress.EndTime = util.GetNowTime()
		}
	}()
	progress.StartTime = util.GetNowTime()

	var sftpClient *sftp.Client
	if request.FromFilePlace == "remote" || request.ToFilePlace == "remote" {
		sftpClient, err = this_.getSftp()
		if err != nil {
			return
		}
		//defer func() { _ = sftpClient.Close() }()
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

	util.FileCopy(fromLstat, fromLoadFiles, fromReader, request.FromFile.Path, toLstat, toMkdirAll, toWriter, request.ToFile.Path, callConfirm, progress)
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

func (this_ *SftpClient) localRemove(request *SFTPRequest, progress *util.FileRemoveProgress) {
	var err error
	defer func() {
		if err != nil {
			progress.Error = err
			progress.EndTime = util.GetNowTime()
		}
	}()
	progress.StartTime = util.GetNowTime()

	util.FileRemove(os.Lstat, util.LocalLoadFiles, os.Remove, request.Path, progress)

	return
}

func (this_ *SftpClient) remoteRemove(request *SFTPRequest, progress *util.FileRemoveProgress) {
	var err error
	defer func() {
		if err != nil {
			progress.Error = err
			progress.EndTime = util.GetNowTime()
		}
	}()
	progress.StartTime = util.GetNowTime()
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

	util.FileRemove(sftpClient.Lstat, sftpClient.ReadDir, sftpClient.Remove, request.Path, progress)

	return
}

func (this_ *SftpClient) localRename(request *SFTPRequest, progress *util.FileRenameProgress) (err error) {

	defer func() {
		if err == nil {
			progress.SuccessCount = 1
		}
		progress.Error = err
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1

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

func (this_ *SftpClient) remoteRename(request *SFTPRequest, progress *util.FileRenameProgress) (err error) {
	defer func() {
		if err == nil {
			progress.SuccessCount = 1
		}
		progress.Error = err
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()
	progress.Count = 1

	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

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
	response = &SFTPResponse{}
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
	response = &SFTPResponse{}
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

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
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

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
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}
	//defer func() { _ = sftpClient.Close() }()

	err = util.FileWrite(sftpClient.Lstat, func(s string) (io.Writer, error) {
		return sftpClient.Create(s)
	}, request.Path, []byte(request.Text))
	if err != nil {
		return
	}
	return
}
