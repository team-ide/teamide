package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"teamide/pkg/util"
	"time"
)

type SFTPRequest struct {
	Work      string `json:"work,omitempty"`
	WorkId    string `json:"workId,omitempty"`
	Dir       string `json:"dir,omitempty"`
	Place     string `json:"place,omitempty"`
	Path      string `json:"path,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
	Name      string `json:"name,omitempty"`
	OldPath   string `json:"oldPath,omitempty"`
	NewPath   string `json:"newPath,omitempty"`
	File      *multipart.FileHeader
	FromFile  *SFTPFile `json:"fromFile,omitempty"`
	ToFile    *SFTPFile `json:"toFile,omitempty"`
	ConfirmId string    `json:"confirmId,omitempty"`
	IsOk      bool      `json:"isOk,omitempty"`
	IsCancel  bool      `json:"isCancel,omitempty"`
	ScrollTop int       `json:"scrollTop,omitempty"`
	Text      string    `json:"text,omitempty"`
}
type SFTPResponse struct {
	Work      string      `json:"work,omitempty"`
	WorkId    string      `json:"workId,omitempty"`
	Dir       string      `json:"dir,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Files     []*SFTPFile `json:"files,omitempty"`
	Place     string      `json:"place,omitempty"`
	Path      string      `json:"path,omitempty"`
	Name      string      `json:"name,omitempty"`
	ScrollTop int         `json:"scrollTop,omitempty"`
	Text      string      `json:"text,omitempty"`
}
type SFTPFile struct {
	Name     string     `json:"name,omitempty"`
	IsDir    bool       `json:"isDir,omitempty"`
	Size     int64      `json:"size,omitempty"`
	Place    string     `json:"place,omitempty"`
	Path     string     `json:"path,omitempty"`
	ModTime  *time.Time `json:"modTime,omitempty"`
	FileMode int64      `json:"fileMode,omitempty"`
}

type RemoveProgress struct {
	WaitCall     bool  `json:"-"`
	StartTime    int64 `json:"startTime"`
	EndTime      int64 `json:"endTime"`
	Timestamp    int64 `json:"timestamp"`
	Count        int64 `json:"count"`
	Size         int64 `json:"size"`
	SuccessCount int64 `json:"successCount"`
}

type CopyProgress struct {
	WaitCall     bool     `json:"-"`
	StartTime    int64    `json:"startTime"`
	EndTime      int64    `json:"endTime"`
	Timestamp    int64    `json:"timestamp"`
	Size         int64    `json:"size"`
	SuccessSize  int64    `json:"successSize"`
	Count        int64    `json:"count"`
	SuccessCount int64    `json:"successCount"`
	Copying      *Copying `json:"copying,omitempty"`
}

type Copying struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	SuccessSize int64  `json:"successSize"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Timestamp   int64  `json:"timestamp"`
}

type UploadProgress struct {
	WaitCall     bool       `json:"-"`
	StartTime    int64      `json:"startTime"`
	EndTime      int64      `json:"endTime"`
	Timestamp    int64      `json:"timestamp"`
	Size         int64      `json:"size"`
	SuccessSize  int64      `json:"successSize"`
	Count        int64      `json:"count"`
	SuccessCount int64      `json:"successCount"`
	Uploading    *Uploading `json:"uploading,omitempty"`
}

type Uploading struct {
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	SuccessSize int64  `json:"successSize"`
}

func (this_ *SSHSftpClient) callConfirm(confirmInfo *ConfirmInfo) (res *ConfirmInfo, err error) {

	if this_.confirmMap == nil {
		this_.confirmMap = map[string]chan *ConfirmInfo{}
	}
	confirmInfo.IsConfirm = true
	if confirmInfo.ConfirmId == "" {
		confirmInfo.ConfirmId = util.GenerateUUID()
	}
	this_.confirmMap[confirmInfo.ConfirmId] = make(chan *ConfirmInfo, 1)
	bs, err := json.Marshal(confirmInfo)
	if err != nil {
		fmt.Println("call confirm to json err:", err)
		return
	}
	this_.WSWrite(bs)
	//fmt.Println("等待[", confirmInfo.ConfirmId, "]结果")
	res = <-this_.confirmMap[confirmInfo.ConfirmId]
	//fmt.Println("接收[", confirmInfo.ConfirmId, "]结果", res)

	close(this_.confirmMap[confirmInfo.ConfirmId])
	delete(this_.confirmMap, confirmInfo.ConfirmId)
	return

}

func (this_ *SSHSftpClient) callProgress(request *SFTPRequest, progress interface{}) {
	for {
		time.Sleep(100 * time.Millisecond)

		if this_.isClosedWS {
			return
		}
		var waitCall bool
		var endTime int64 = -1
		UploadProgress, UploadProgressOk := progress.(*UploadProgress)
		if UploadProgressOk {
			UploadProgress.Timestamp = util.GetNowTime()
			endTime = UploadProgress.EndTime
			waitCall = UploadProgress.WaitCall
		}

		CopyProgress, CopyProgressOK := progress.(*CopyProgress)
		if CopyProgressOK {
			CopyProgress.Timestamp = util.GetNowTime()
			endTime = CopyProgress.EndTime
			waitCall = CopyProgress.WaitCall
		}

		RemoveProgress, RemoveProgressOk := progress.(*RemoveProgress)
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
			fmt.Println("sftp upload progress to json err:", err)
			continue
		}
		this_.WSWrite(bs)

		if endTime > 0 {
			break
		}
	}
}

func (this_ *SSHSftpClient) work(request *SFTPRequest) {
	response := &SFTPResponse{}
	var err error
	switch request.Work {
	case "confirmResult":
		//fmt.Println("WS 收到[", request.ConfirmId, "]结果", request)
		if this_.confirmMap == nil {
			return
		}
		//fmt.Println("WS 通知[", request.ConfirmId, "]携程", request)
		this_.confirmMap[request.ConfirmId] <- &ConfirmInfo{
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
		progress := &UploadProgress{
			StartTime: GetNowTime(),
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
		progress := &CopyProgress{
			StartTime: GetNowTime(),
		}
		go this_.callProgress(request, progress)
		response, err = this_.copy(request, progress)
	case "remove":
		progress := &RemoveProgress{
			StartTime: GetNowTime(),
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
		this_.Logger.Error("ssh ftp work{"+request.Work+"} error", zap.Error(err))
		response.Msg = err.Error()
	}
	response.Work = request.Work
	response.WorkId = request.WorkId
	response.Place = request.Place
	response.ScrollTop = request.ScrollTop

	this_.WSWriteData(response)

	return
}

func CopyBytes(dst io.Writer, src io.Reader, call func(readSize int64, writeSize int64)) (err error) {
	var buf = make([]byte, 32*1024)
	var errInvalidWrite = errors.New("invalid write result")
	var ErrShortWrite = errors.New("short write")
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			call(int64(nr), 0)
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			call(0, int64(nw))
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return
}

func (this_ *SSHSftpClient) localUpdate(request *SFTPRequest, progress *UploadProgress) (response *SFTPResponse, err error) {

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

	pathDir := path[0:strings.LastIndex(path, "/")]

	_, err = os.Lstat(pathDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(pathDir, 0777)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	_, err = os.Lstat(path)
	if err == nil {
		progress.WaitCall = true
		defer func() {
			progress.WaitCall = false
		}()
		confirmInfo := &ConfirmInfo{
			IsFileExist: true,
			Path:        path,
			Name:        request.File.Filename,
		}
		var res *ConfirmInfo
		res, err = this_.callConfirm(confirmInfo)
		if err != nil {
			return
		}
		if res.IsCancel {
			progress.SuccessCount++
			progress.Size -= request.File.Size
			return
		}
		progress.WaitCall = false
	}

	var fileInfo *os.File
	fileInfo, err = os.Create(path)
	if err != nil {
		return
	}
	defer closeFile(fileInfo)

	uploadF, err := request.File.Open()
	if err != nil {
		return
	}
	defer closeUploadFile(uploadF)

	err = CopyBytes(fileInfo, uploadF, func(readSize int64, writeSize int64) {
		progress.SuccessSize += writeSize
	})
	if err != nil {
		return
	}

	progress.SuccessCount++
	return
}

func (this_ *SSHSftpClient) remoteUpdate(request *SFTPRequest, progress *UploadProgress) (response *SFTPResponse, err error) {

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

	pathDir := path[0:strings.LastIndex(path, "/")]

	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	_, err = sftpClient.Lstat(pathDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = sftpClient.MkdirAll(pathDir)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	_, err = sftpClient.Lstat(path)
	if err == nil {
		progress.WaitCall = true
		defer func() {
			progress.WaitCall = false
		}()
		confirmInfo := &ConfirmInfo{
			IsFileExist: true,
			Path:        path,
			Name:        request.File.Filename,
		}
		var res *ConfirmInfo
		res, err = this_.callConfirm(confirmInfo)
		if err != nil {
			return
		}
		if res.IsCancel {
			progress.SuccessCount++
			progress.Size -= request.File.Size
			return
		}
		progress.WaitCall = false
	}

	fileInfo, err := sftpClient.Create(path)
	if err != nil {
		return
	}
	defer closeFtpFile(fileInfo)

	uploadF, err := request.File.Open()
	if err != nil {
		return
	}
	defer closeUploadFile(uploadF)

	err = CopyBytes(fileInfo, uploadF, func(readSize int64, writeSize int64) {
		progress.SuccessSize += writeSize
	})
	if err != nil {
		return
	}

	progress.SuccessCount++

	return
}

func (this_ *SSHSftpClient) localDownload(c *gin.Context, path string) (err error) {

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

	err = CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SSHSftpClient) remoteDownload(c *gin.Context, path string) (err error) {

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

	err = CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *SSHSftpClient) copy(request *SFTPRequest, progress *CopyProgress) (response *SFTPResponse, err error) {
	defer func() {
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()

	var sftpClient *sftp.Client
	if request.FromFile.Place == "remote" || request.ToFile.Place == "remote" {
		sftpClient, err = this_.newSftp()
		if err != nil {
			return
		}
		defer this_.closeSftClient(sftpClient)
	}

	progress.Count, progress.Size, err = this_.fileCount(request.FromFile.Place, request.FromFile.Path, sftpClient)
	if err != nil {
		return
	}
	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}
	err = this_.copyAll(request.FromFile.Place, request.FromFile.Path, request.ToFile.Place, request.ToFile.Path, sftpClient, progress)
	if err != nil {
		return
	}
	return
}

func closeUploadFile(obj multipart.File) {
	if obj == nil {
		return
	}
	_ = obj.Close()
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

func (this_ *SSHSftpClient) copyAll(fromPlace string, fromPath string, toPlace string, toPath string, sftpClient *sftp.Client, progress *CopyProgress) (err error) {

	var isDir bool
	var fileName string
	var fileSize int64

	if fromPlace == "local" {
		var info fs.FileInfo
		info, err = os.Lstat(fromPath)
		if err != nil {
			return
		}
		isDir = info.IsDir()
		fileName = info.Name()
		if !isDir {
			fileSize = info.Size()
		}
	} else if fromPlace == "remote" {
		var info os.FileInfo
		info, err = sftpClient.Lstat(fromPath)
		if err != nil {
			return
		}
		isDir = info.IsDir()
		fileName = info.Name()
		if !isDir {
			fileSize = info.Size()
		}
	}

	if isDir {
		progress.SuccessCount++
		if fromPlace == "local" {
			var dirFiles []os.DirEntry
			dirFiles, err = os.ReadDir(fromPath)
			if err != nil {
				return
			}

			for _, f := range dirFiles {
				err = this_.copyAll(fromPlace, fromPath+"/"+f.Name(), toPlace, toPath+"/"+f.Name(), sftpClient, progress)
				if err != nil {
					return
				}
			}
		} else if fromPlace == "remote" {
			var dirFiles []os.FileInfo
			dirFiles, err = sftpClient.ReadDir(fromPath)
			if err != nil {
				return
			}
			for _, f := range dirFiles {
				err = this_.copyAll(fromPlace, fromPath+"/"+f.Name(), toPlace, toPath+"/"+f.Name(), sftpClient, progress)
				if err != nil {
					return
				}
			}
		}

	} else {

		var isExist bool

		if toPlace == "local" {
			_, err = os.Lstat(toPath)
			if err == nil {
				isExist = true
			}
		} else if toPlace == "remote" {
			_, err = sftpClient.Lstat(toPath)
			if err == nil {
				isExist = true
			}
		}
		if isExist {
			progress.WaitCall = true
			defer func() {
				progress.WaitCall = false
			}()
			confirmInfo := &ConfirmInfo{
				IsFileExist: true,
				Path:        toPath,
			}
			var res *ConfirmInfo
			res, err = this_.callConfirm(confirmInfo)
			if err != nil {
				return
			}
			if res.IsCancel {
				progress.Size -= fileSize
				progress.SuccessCount++
				return
			}
			progress.WaitCall = false
		}

		var fromReader io.Reader
		if fromPlace == "local" {
			fromReader, err = os.Open(fromPath)
		} else if fromPlace == "remote" {
			fromReader, err = sftpClient.Open(fromPath)
		}

		if err != nil {
			return
		}

		defer closeIfCloser(fromReader)
		var toWriter io.Writer

		if toPlace == "local" {

			pathDir := toPath[0:strings.LastIndex(toPath, "/")]
			_, err = os.Lstat(pathDir)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(pathDir, 0777)
					if err != nil {
						return
					}
				} else {
					return
				}
			}

			toWriter, err = os.Create(toPath)
		} else if toPlace == "remote" {

			pathDir := toPath[0:strings.LastIndex(toPath, "/")]
			_, err = sftpClient.Lstat(pathDir)
			if err != nil {
				if os.IsNotExist(err) {
					err = sftpClient.MkdirAll(pathDir)
					if err != nil {
						return
					}
				} else {
					return
				}
			}

			toWriter, err = sftpClient.Create(toPath)
		}
		if err != nil {
			return
		}

		defer closeIfCloser(toWriter)

		Copying := &Copying{}
		Copying.Name = fileName
		Copying.StartTime = util.GetNowTime()
		Copying.Size = fileSize
		progress.Copying = Copying
		err = CopyBytes(toWriter, fromReader, func(readSize int64, writeSize int64) {
			progress.SuccessSize += writeSize
			Copying.SuccessSize += writeSize
		})
		if err != nil {
			return
		}

		progress.SuccessCount++
	}

	return
}

func (this_ *SSHSftpClient) localRemove(request *SFTPRequest, progress *RemoveProgress) (response *SFTPResponse, err error) {
	defer func() {
		progress.EndTime = util.GetNowTime()
	}()
	progress.StartTime = util.GetNowTime()

	progress.Count, progress.Size, err = this_.fileCount("local", request.Path, nil)
	if err != nil {
		return
	}

	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}

	err = this_.localRemoveAll(request.Path, progress)
	if err != nil {
		return
	}
	return
}

func (this_ *SSHSftpClient) localRemoveAll(path string, progress *RemoveProgress) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = os.Lstat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	if isDir {
		var dirFiles []os.DirEntry
		dirFiles, err = os.ReadDir(path)
		if err != nil {
			return
		}
		for _, f := range dirFiles {
			err = this_.localRemoveAll(path+"/"+f.Name(), progress)
			if err != nil {
				return
			}
		}

	}
	err = os.Remove(path)
	if err != nil {
		return
	}
	progress.SuccessCount++
	return
}

func (this_ *SSHSftpClient) fileCount(place string, path string, sftpClient *sftp.Client) (fileCount int64, fileSize int64, err error) {
	var isDir bool

	var thisFileSize int64
	if place == "local" {
		var info fs.FileInfo
		info, err = os.Lstat(path)
		if err != nil {
			return
		}
		isDir = info.IsDir()
		if !isDir {
			thisFileSize = info.Size()
		}
	} else if place == "remote" {
		var info os.FileInfo
		info, err = sftpClient.Lstat(path)
		if err != nil {
			return
		}
		isDir = info.IsDir()
		if !isDir {
			thisFileSize = info.Size()
		}
	}

	fileCount++
	fileSize += thisFileSize
	if isDir {
		if place == "local" {
			var dirFiles []os.DirEntry
			dirFiles, err = os.ReadDir(path)
			if err != nil {
				return
			}

			for _, f := range dirFiles {
				var fileCount_ int64
				var fileSize_ int64
				fileCount_, fileSize_, err = this_.fileCount(place, path+"/"+f.Name(), sftpClient)
				if err != nil {
					return
				}
				fileCount += fileCount_
				fileSize += fileSize_
			}
		} else if place == "remote" {
			var dirFiles []os.FileInfo
			dirFiles, err = sftpClient.ReadDir(path)
			if err != nil {
				return
			}
			for _, f := range dirFiles {
				var fileCount_ int64
				var fileSize_ int64
				fileCount_, fileSize_, err = this_.fileCount(place, path+"/"+f.Name(), sftpClient)
				if err != nil {
					return
				}
				fileCount += fileCount_
				fileSize += fileSize_
			}
		}

	}
	return
}

func (this_ *SSHSftpClient) remoteRemove(request *SFTPRequest, progress *RemoveProgress) (response *SFTPResponse, err error) {
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

	progress.Count, progress.Size, err = this_.fileCount("remote", request.Path, sftpClient)
	if err != nil {
		return
	}
	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}

	err = this_.remoteRemoveAll(request.Path, sftpClient, progress)
	if err != nil {
		return
	}

	return
}

func (this_ *SSHSftpClient) remoteRemoveAll(path string, sftpClient *sftp.Client, progress *RemoveProgress) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = sftpClient.Lstat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	if isDir {
		var dirFiles []os.FileInfo
		dirFiles, err = sftpClient.ReadDir(path)
		if err != nil {
			return
		}
		for _, f := range dirFiles {
			err = this_.remoteRemoveAll(path+"/"+f.Name(), sftpClient, progress)
			if err != nil {
				return
			}
		}

	}
	err = sftpClient.Remove(path)
	if err != nil {
		return
	}
	progress.SuccessCount++
	return
}

func (this_ *SSHSftpClient) localRename(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.OldPath,
		Dir:  request.Dir,
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

func (this_ *SSHSftpClient) remoteRename(request *SFTPRequest) (response *SFTPResponse, err error) {
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

func (this_ *SSHSftpClient) localFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Files: []*SFTPFile{
			{
				Name:  "..",
				IsDir: true,
				Place: "local",
			},
		},
	}
	dir := request.Dir
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return
		}
	}

	dir = util.FormatPath(dir)
	if err != nil {
		return
	}
	response.Dir = dir

	fileInfo, err := os.Lstat(dir)
	if err != nil {
		if err == os.ErrNotExist {
			err = nil
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var dirNames []string
	var fileNames []string

	fMap := map[string]os.DirEntry{}
	for _, f := range dirFiles {
		fName := f.Name()
		fMap[fName] = f
		if f.IsDir() {
			dirNames = append(dirNames, fName)
		} else {
			fileNames = append(fileNames, fName)
		}
	}

	sort.Strings(dirNames)
	sort.Strings(fileNames)

	for _, one := range dirNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		response.Files = append(response.Files, &SFTPFile{
			Name:     one,
			IsDir:    true,
			Place:    "local",
			ModTime:  &ModTime,
			FileMode: int64(fi.Mode()),
		})
	}
	for _, one := range fileNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		response.Files = append(response.Files, &SFTPFile{
			Name:     one,
			Size:     fi.Size(),
			Place:    "local",
			ModTime:  &ModTime,
			FileMode: int64(fi.Mode()),
		})
	}

	return
}

func (this_ *SSHSftpClient) remoteFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Files: []*SFTPFile{
			{
				Name:  "..",
				IsDir: true,
				Place: "remote",
			},
		},
	}
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

	fileInfo, err := sftpClient.Lstat(dir)
	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dirFiles, err := sftpClient.ReadDir(dir)
	if err != nil {
		return
	}
	var dirNames []string
	var fileNames []string

	fMap := map[string]os.FileInfo{}
	for _, f := range dirFiles {
		fName := f.Name()
		fMap[fName] = f
		if f.IsDir() {
			dirNames = append(dirNames, fName)
		} else {
			fileNames = append(fileNames, fName)
		}
	}

	sort.Strings(dirNames)
	sort.Strings(fileNames)

	for _, one := range dirNames {
		ModTime := fMap[one].ModTime()
		response.Files = append(response.Files, &SFTPFile{
			Name:     one,
			IsDir:    true,
			Place:    "remote",
			ModTime:  &ModTime,
			FileMode: int64(fMap[one].Mode()),
		})
	}
	for _, one := range fileNames {
		ModTime := fMap[one].ModTime()

		response.Files = append(response.Files, &SFTPFile{
			Name:     one,
			Size:     fMap[one].Size(),
			Place:    "remote",
			ModTime:  &ModTime,
			FileMode: int64(fMap[one].Mode()),
		})
	}

	return
}

func (this_ *SSHSftpClient) localReadText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	fileInfo, err := os.Lstat(request.Path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + request.Path + "]是目录")
		return
	}
	if fileInfo.Size() > 1024*1024*10 {
		err = errors.New("只支持打开10M以内的文件在线查看")
		return
	}
	f, err := os.Open(request.Path)
	if err != nil {
		return
	}
	defer closeFile(f)
	bs, err := io.ReadAll(f)
	if err != nil {
		return
	}
	if len(bs) > 0 {
		response.Text = string(bs)
	}

	return
}

func (this_ *SSHSftpClient) remoteReadText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	fileInfo, err := sftpClient.Lstat(request.Path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + request.Path + "]是目录")
		return
	}
	if fileInfo.Size() > 1024*1024*10 {
		err = errors.New("只支持打开10M以内的文件在线查看")
		return
	}
	f, err := sftpClient.Open(request.Path)
	if err != nil {
		return
	}
	defer closeIfCloser(f)
	bs, err := io.ReadAll(f)
	if err != nil {
		return
	}
	if len(bs) > 0 {
		response.Text = string(bs)
	}
	return
}

func (this_ *SSHSftpClient) localSaveText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	fileInfo, err := os.Lstat(request.Path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + request.Path + "]是目录")
		return
	}

	f, err := os.Create(request.Path)
	if err != nil {
		return
	}
	defer closeFile(f)

	_, err = f.Write([]byte(request.Text))
	if err != nil {
		return
	}

	return
}

func (this_ *SSHSftpClient) remoteSaveText(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
	}
	var sftpClient *sftp.Client
	sftpClient, err = this_.newSftp()
	if err != nil {
		return
	}
	defer this_.closeSftClient(sftpClient)

	fileInfo, err := sftpClient.Lstat(request.Path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + request.Path + "]是目录")
		return
	}

	f, err := sftpClient.Create(request.Path)
	if err != nil {
		return
	}
	defer closeIfCloser(f)

	_, err = f.Write([]byte(request.Text))
	if err != nil {
		fmt.Println("文件:"+request.Path+",写入异常", err)
		return
	}
	return
}
