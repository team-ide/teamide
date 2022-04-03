package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"sort"
	"strings"
	"teamide/pkg/util"
)

func (this_ *SSHClient) StartSftp(ws *websocket.Conn) (err error) {
	this_.ws = ws
	err = this_.initClient()
	if err != nil {
		fmt.Println("StartSftp error", err)
		this_.Close()
		return
	}

	this_.sftpClient, err = sftp.NewClient(this_.sshClient)
	if err != nil {
		fmt.Println("NewClient error", err)
		this_.Close()
		return
	}
	if this_.UploadFile == nil {
		this_.UploadFile = make(chan *UploadFile, 10)
	}
	go func() {
		for {
			select {
			case uploadFile := <-this_.UploadFile:
				this_.work(&SFTPRequest{
					Work:     "upload",
					WorkId:   uploadFile.WorkId,
					Dir:      uploadFile.Dir,
					Place:    uploadFile.Place,
					File:     uploadFile.File,
					FullPath: uploadFile.FullPath,
				})
			}
		}

	}()
	// 第一个协程获取用户的输入
	go func() {
		for {
			if this_.isClosed {
				return
			}
			_, p, err := this_.ws.ReadMessage()
			if err != nil && err != io.EOF {
				fmt.Println("sftp ws read err:", err)
				this_.Close()
				return
			}
			//fmt.Println("ws read:" + string(p))
			if len(p) > 0 {
				if this_.isClosed {
					return
				}
				var request *SFTPRequest
				err = json.Unmarshal(p, &request)
				if err != nil {
					fmt.Println("sftp ws message to struct err:", err)
					continue
				}
				this_.work(request)
			}
		}
	}()
	return
}

type SFTPRequest struct {
	Work     string                `json:"work,omitempty"`
	WorkId   string                `json:"workId,omitempty"`
	Dir      string                `json:"dir,omitempty"`
	Place    string                `json:"place,omitempty"`
	Path     string                `json:"path,omitempty"`
	FullPath string                `json:"fullPath,omitempty"`
	Name     string                `json:"name,omitempty"`
	OldPath  string                `json:"oldPath,omitempty"`
	NewPath  string                `json:"newPath,omitempty"`
	File     *multipart.FileHeader `json:"-`
	FromFile *SFTPFile             `json:"fromFile,omitempty"`
	ToFile   *SFTPFile             `json:"toFile,omitempty"`
}
type SFTPResponse struct {
	Work   string      `json:"work,omitempty"`
	WorkId string      `json:"workId,omitempty"`
	Dir    string      `json:"dir,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Files  []*SFTPFile `json:"files,omitempty"`
	Place  string      `json:"place,omitempty"`
	Path   string      `json:"path,omitempty"`
	Name   string      `json:"name,omitempty"`
}
type SFTPFile struct {
	Name  string `json:"name,omitempty"`
	IsDir bool   `json:"isDir,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Place string `json:"place,omitempty"`
	Path  string `json:"path,omitempty"`
}

func (this_ *SSHClient) work(request *SFTPRequest) {
	if this_.isClosed {
		return
	}
	response := &SFTPResponse{}
	var err error
	switch request.Work {
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
		if request.Place == "local" {
			response, err = this_.localUpdate(request)
		} else if request.Place == "remote" {
			response, err = this_.remoteUpdate(request)
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
		response, err = this_.copy(request)
	case "remove":
		if request.Place == "local" {
			response, err = this_.localRemove(request)
		} else if request.Place == "remote" {
			response, err = this_.remoteRemove(request)
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
		response.Msg = err.Error()
	}
	response.Work = request.Work
	response.WorkId = request.WorkId
	response.Place = request.Place
	bs, err := json.Marshal(response)
	if err != nil {
		fmt.Println("sftp message to json err:", err)
		return
	}
	err = this_.ws.WriteMessage(websocket.TextMessage, bs)
	if err != nil {
		fmt.Println("sftp ws write err:", err)
		this_.Close()
		return
	}

	return
}

func (this_ *SSHClient) localUpdate(request *SFTPRequest) (response *SFTPResponse, err error) {

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

	fileInfo, err := os.Create(path)
	if err != nil {
		return
	}
	defer fileInfo.Close()

	uploadF, err := request.File.Open()
	if err != nil {
		return
	}
	defer uploadF.Close()
	_, err = io.Copy(fileInfo, uploadF)
	if err != nil {
		return
	}
	return
}

func (this_ *SSHClient) remoteUpdate(request *SFTPRequest) (response *SFTPResponse, err error) {
	path := request.Dir + "/" + request.File.Filename
	if request.FullPath != "" {
		path = request.Dir + "/" + strings.TrimPrefix(request.FullPath, "/")
	}
	response = &SFTPResponse{
		Path: path,
		Dir:  request.Dir,
	}

	pathDir := path[0:strings.LastIndex(path, "/")]
	_, err = this_.sftpClient.Lstat(pathDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = this_.sftpClient.MkdirAll(pathDir)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	fileInfo, err := this_.sftpClient.Create(path)
	if err != nil {
		return
	}
	defer fileInfo.Close()

	uploadF, err := request.File.Open()
	if err != nil {
		return
	}
	defer uploadF.Close()
	_, err = io.Copy(fileInfo, uploadF)
	if err != nil {
		return
	}

	return
}

func (this_ *SSHClient) copy(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}
	err = this_.copyAll(request.FromFile.Place, request.FromFile.Path, request.ToFile.Place, request.ToFile.Path)
	if err != nil {
		return
	}
	return
}

func (this_ *SSHClient) copyAll(fromPlace string, fromPath string, toPlace string, toPath string) (err error) {

	var isDir bool

	if fromPlace == "local" {
		var info fs.FileInfo
		info, err = os.Lstat(fromPath)
		if err != nil {
			return
		}
		isDir = info.IsDir()
	} else if fromPlace == "remote" {
		var info os.FileInfo
		info, err = this_.sftpClient.Lstat(fromPath)
		if err != nil {
			return
		}
		isDir = info.IsDir()
	}

	if isDir {
		if fromPlace == "local" {
			var fs []os.DirEntry
			fs, err = os.ReadDir(fromPath)
			if err != nil {
				return
			}

			for _, f := range fs {
				err = this_.copyAll(fromPlace, fromPath+"/"+f.Name(), toPlace, toPath+"/"+f.Name())
				if err != nil {
					return
				}
			}
		} else if fromPlace == "remote" {
			var fs []os.FileInfo
			fs, err = this_.sftpClient.ReadDir(fromPath)
			if err != nil {
				return
			}
			for _, f := range fs {
				err = this_.copyAll(fromPlace, fromPath+"/"+f.Name(), toPlace, toPath+"/"+f.Name())
				if err != nil {
					return
				}
			}
		}

	} else {
		var fromReader io.Reader
		if fromPlace == "local" {
			fromReader, err = os.Open(fromPath)
		} else if fromPlace == "remote" {
			fromReader, err = this_.sftpClient.Open(fromPath)
		}
		if err != nil {
			return
		}
		defer fromReader.(io.Closer).Close()
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
			_, err = this_.sftpClient.Lstat(pathDir)
			if err != nil {
				if os.IsNotExist(err) {
					err = this_.sftpClient.MkdirAll(pathDir)
					if err != nil {
						return
					}
				} else {
					return
				}
			}

			toWriter, err = this_.sftpClient.Create(toPath)
		}
		if err != nil {
			return
		}

		defer toWriter.(io.Closer).Close()
		_, err = io.Copy(toWriter, fromReader)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *SSHClient) localRemove(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}
	info, err := os.Lstat(request.Path)
	if err != nil {
		return
	}
	if info.IsDir() {
		err = os.RemoveAll(request.Path)
	} else {
		err = os.Remove(request.Path)
	}
	if err != nil {
		return
	}
	return
}

func (this_ *SSHClient) remoteRemove(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.Path,
		Dir:  request.Dir,
	}

	err = this_.remoteRemoveAll(request.Path)
	if err != nil {
		return
	}

	return
}

func (this_ *SSHClient) remoteRemoveAll(path string) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = this_.sftpClient.Lstat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	if isDir {
		var fs []os.FileInfo
		fs, err = this_.sftpClient.ReadDir(path)
		if err != nil {
			return
		}
		for _, f := range fs {
			err = this_.remoteRemoveAll(path + "/" + f.Name())
			if err != nil {
				return
			}
		}

	}
	err = this_.sftpClient.Remove(path)

	return
}
func (this_ *SSHClient) localRename(request *SFTPRequest) (response *SFTPResponse, err error) {
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

func (this_ *SSHClient) remoteRename(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Path: request.OldPath,
		Dir:  request.Dir,
	}

	_, err = this_.sftpClient.Lstat(request.OldPath)
	if err != nil {
		return
	}

	err = this_.sftpClient.Rename(request.OldPath, request.NewPath)
	if err != nil {
		return
	}

	return
}

func (this_ *SSHClient) localFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Files: []*SFTPFile{},
	}
	dir := request.Dir
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return
		}
	}

	fileInfo, err := os.Lstat(dir)
	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dir = util.FormatPath(dir)
	if err != nil {
		return
	}
	response.Dir = dir
	fs, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	dirNames := []string{".."}
	var fileNames []string

	fMap := map[string]os.DirEntry{}
	for _, f := range fs {
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
		response.Files = append(response.Files, &SFTPFile{
			Name:  one,
			IsDir: true,
			Place: "local",
		})
	}
	for _, one := range fileNames {
		f := fMap[one]
		var size int64
		if !f.IsDir() {
			var fi os.FileInfo
			fi, err = f.Info()
			if err != nil {
				return
			}
			size = fi.Size()
		}
		response.Files = append(response.Files, &SFTPFile{
			Name:  one,
			Size:  size,
			Place: "local",
		})
	}

	return
}

func (this_ *SSHClient) remoteFiles(request *SFTPRequest) (response *SFTPResponse, err error) {
	response = &SFTPResponse{
		Files: []*SFTPFile{},
	}
	dir := request.Dir
	if dir == "" {
		dir, err = this_.sftpClient.Getwd()
		if err != nil {
			return
		}
	}

	fileInfo, err := this_.sftpClient.Lstat(dir)
	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dir, err = this_.sftpClient.RealPath(dir)
	if err != nil {
		return
	}
	response.Dir = dir
	fs, err := this_.sftpClient.ReadDir(dir)
	if err != nil {
		return
	}
	dirNames := []string{".."}
	var fileNames []string

	fMap := map[string]os.FileInfo{}
	for _, f := range fs {
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
		response.Files = append(response.Files, &SFTPFile{
			Name:  one,
			IsDir: true,
			Place: "remote",
		})
	}
	for _, one := range fileNames {
		response.Files = append(response.Files, &SFTPFile{
			Name:  one,
			Size:  fMap[one].Size(),
			Place: "remote",
		})
	}

	return
}

func SFTPUpload(c *gin.Context) (res error, err error) {
	token := c.PostForm("token")
	//fmt.Println("token=" + token)
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	dir := c.PostForm("dir")
	//fmt.Println("token=" + token)
	if dir == "" {
		err = errors.New("dir获取失败")
		return
	}
	place := c.PostForm("place")
	//fmt.Println("token=" + token)
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	workId := c.PostForm("workId")
	//fmt.Println("token=" + token)
	if workId == "" {
		err = errors.New("workId获取失败")
		return
	}
	sshClient := SSHClientCache[token]
	if sshClient == nil {
		err = errors.New("SSH会话丢失")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	uploadFile := &UploadFile{
		Dir:      dir,
		Place:    place,
		WorkId:   workId,
		File:     file,
		FullPath: c.PostForm("fullPath"),
	}
	sshClient.UploadFile <- uploadFile

	return
}
