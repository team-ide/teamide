package ssh

import (
	"errors"
	"github.com/pkg/sftp"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"teamide/pkg/base"
	"teamide/pkg/filework"
)

func newFileService(config *Config) *fileService {
	return &fileService{
		config: config,
	}
}

var (
	fileServiceCache     = make(map[string]*fileService)
	fileServiceCacheLock = &sync.Mutex{}
)

func CreateOrGetClient(key string, config *Config) (res *fileService) {
	fileServiceCacheLock.Lock()
	defer fileServiceCacheLock.Unlock()
	res, ok := fileServiceCache[key]
	if !ok {
		res = newFileService(config)
		fileServiceCache[key] = res
	}
	return
}

func CloseFileService(key string) {
	fileServiceCacheLock.Lock()
	defer fileServiceCacheLock.Unlock()
	res, ok := fileServiceCache[key]
	if ok {
		delete(fileServiceCache, key)
		res.Close()
	}
	return
}

type fileService struct {
	config      *Config
	sshClient   *ssh.Client
	newSftpLock sync.Mutex

	sftpClient *sftp.Client
}

func (this_ *fileService) getSftp() (sftpClient *sftp.Client, err error) {
	this_.newSftpLock.Lock()
	defer this_.newSftpLock.Unlock()

	if this_.sshClient == nil {
		this_.sftpClient = nil
		err = this_.createClient()
		if err != nil {
			this_.sshClient = nil
			this_.sftpClient = nil
			return
		}
	}
	if this_.sftpClient == nil {
		this_.sftpClient, err = sftp.NewClient(this_.sshClient)
		if err != nil {
			this_.sshClient = nil
			this_.sftpClient = nil
			return
		}
	}

	sftpClient = this_.sftpClient

	return
}

func (this_ *fileService) Close() {
	this_.closeClient()
	return
}

func (this_ *fileService) closeClient() {
	if this_.sshClient != nil {
		_ = this_.sshClient.Close()
		this_.sshClient = nil
	}
	return
}

func (this_ *fileService) createClient() (err error) {

	if this_.sshClient, err = NewClient(*this_.config); err != nil {
		util.Logger.Error("createClient error", zap.Error(err))
		return
	}
	go func() {
		err = this_.sshClient.Wait()
		this_.Close()
	}()
	return
}

func (this_ *fileService) Exist(path string) (exist bool, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	f, err := sftpClient.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	if f != nil {
		exist = true
	}
	return
}

func (this_ *fileService) Create(path string, isDir bool) (err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	path, err = sftpClient.RealPath(path)
	if err != nil {
		return
	}

	exist, err := this_.Exist(path)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("路径[" + path + "]已存在")
		return
	}

	if isDir {
		err = sftpClient.MkdirAll(path)
		if err != nil {
			return
		}
	} else {
		var f *sftp.File
		f, err = sftpClient.Create(path)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
	}
	return
}

func (this_ *fileService) Write(path string, reader io.Reader, onDo func(readSize int64, writeSize int64), callStop *bool) (err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	pathDir := path[0:strings.LastIndex(path, "/")]
	exist, err := this_.Exist(pathDir)
	if err != nil {
		util.Logger.Error("Write Exist path error", zap.Any("path", pathDir), zap.Error(err))
		return
	}
	if !exist {
		err = sftpClient.MkdirAll(pathDir)
		if err != nil {
			return
		}
	}

	f, err := sftpClient.Create(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	buf := make([]byte, 32*1024)
	var readSize int64
	var writeSize int64

	err = util.Read(reader, buf, func(n int) (e error) {
		if *callStop {
			e = base.ProgressCallStoppedError
			return
		}
		readSize += int64(n)
		onDo(readSize, writeSize)
		e = util.Write(f, buf[:n], func(n int) (e error) {
			writeSize += int64(n)
			onDo(readSize, writeSize)
			return
		})
		return
	})

	if err != nil {
		return
	}
	return
}

func (this_ *fileService) Read(path string, writer io.Writer, onDo func(readSize int64, writeSize int64), callStop *bool) (err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	path, err = sftpClient.RealPath(path)
	if err != nil {
		return
	}

	exist, err := this_.Exist(path)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("路径[" + path + "]不存在")
		return
	}

	f, err := sftpClient.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	buf := make([]byte, 32*1024)

	var readSize int64
	var writeSize int64

	err = util.Read(f, buf, func(n int) (e error) {
		if *callStop {
			e = base.ProgressCallStoppedError
			return
		}
		readSize += int64(n)
		onDo(readSize, writeSize)
		e = util.Write(writer, buf[:n], func(n int) (e error) {
			writeSize += int64(n)
			onDo(readSize, writeSize)
			return
		})
		return
	})

	if err != nil {
		return
	}

	return
}

func (this_ *fileService) Rename(oldPath string, newPath string) (err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	oldPath, err = sftpClient.RealPath(oldPath)
	if err != nil {
		return
	}

	newPath, err = sftpClient.RealPath(newPath)
	if err != nil {
		return
	}

	err = sftpClient.Rename(oldPath, newPath)
	if err != nil {
		return
	}

	return
}

func (this_ *fileService) Move(oldPath string, newPath string) (err error) {
	err = this_.Rename(oldPath, newPath)
	return
}

func (this_ *fileService) Remove(path string, onDo func(fileCount int, removeCount int)) (err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	var fileCount int
	var removeCount int

	err = removeFile(sftpClient, path, func() {
		fileCount++
		onDo(fileCount, removeCount)
	}, func() {
		removeCount++
		onDo(fileCount, removeCount)
	})

	return
}

func removeFile(sftpClient *sftp.Client, path string, onLoad func(), onRemove func()) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = sftpClient.Stat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	onLoad()
	if isDir {
		var ds []os.FileInfo
		ds, err = sftpClient.ReadDir(path)
		if err != nil {
			return
		}

		for _, d := range ds {
			err = removeFile(sftpClient, path+"/"+d.Name(), onLoad, onRemove)
			if err != nil {
				return
			}
		}
	}
	err = sftpClient.Remove(path)
	if err != nil {
		return
	}
	onRemove()
	return
}

func (this_ *fileService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *fileService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *fileService) Files(dir string) (parentPath string, files []*filework.FileInfo, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	parentPath = dir
	if parentPath == "" {
		parentPath, err = sftpClient.Getwd()
		if err != nil {
			return
		}
	}
	parentPath, err = sftpClient.RealPath(parentPath)
	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}

	files = []*filework.FileInfo{
		{
			Name:  "..",
			Path:  parentPath + "..",
			IsDir: true,
		},
	}

	fileInfo, err := sftpClient.Stat(parentPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("路径[" + parentPath + "]不存在")
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + parentPath + "]不是目录")
		return
	}

	fs, err := sftpClient.ReadDir(parentPath)
	if err != nil {
		return
	}
	var dirNames []string
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

	sort.Slice(dirNames, func(i, j int) bool {
		return strings.ToLower(dirNames[i]) < strings.ToLower(dirNames[j]) //升序  即前面的值比后面的小 忽略大小写排序
	})
	sort.Slice(fileNames, func(i, j int) bool {
		return strings.ToLower(fileNames[i]) < strings.ToLower(fileNames[j]) //升序  即前面的值比后面的小 忽略大小写排序
	})

	for _, one := range dirNames {
		fileOne := getFileInfoByStat(parentPath+one, fMap[one])
		files = append(files, fileOne)
	}
	for _, one := range fileNames {
		fileOne := getFileInfoByStat(parentPath+one, fMap[one])
		files = append(files, fileOne)
	}

	return
}

func (this_ *fileService) File(path string) (file *filework.FileInfo, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	path, err = sftpClient.RealPath(path)
	if err != nil {
		return
	}

	stat, err := sftpClient.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}

	file = getFileInfoByStat(path, stat)

	return
}

func getFileInfoByStat(path string, stat os.FileInfo) (fileInfo *filework.FileInfo) {

	fileInfo = &filework.FileInfo{
		Name:     stat.Name(),
		Path:     path,
		IsDir:    stat.IsDir(),
		ModTime:  util.GetMilliByTime(stat.ModTime()),
		FileMode: stat.Mode().String(),
		Size:     stat.Size(),
	}
	return
}

func (this_ *fileService) OpenReader(path string) (reader io.ReadCloser, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	reader, err = sftpClient.Open(path)
	if err != nil {
		return
	}
	return
}
func (this_ *fileService) OpenWriter(path string) (writer io.WriteCloser, err error) {
	var sftpClient *sftp.Client
	sftpClient, err = this_.getSftp()
	if err != nil {
		return
	}

	writer, err = sftpClient.Create(path)
	if err != nil {
		return
	}
	return
}
