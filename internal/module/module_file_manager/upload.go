package module_file_manager

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"strings"
	"sync"
	"time"
)

type UploadReader struct {
	dataChan  chan []byte
	dataCache []byte
	close     chan struct{}
	mutex     sync.Mutex
}

func NewUploadReader() *UploadReader {
	return &UploadReader{
		dataChan: make(chan []byte),
		close:    make(chan struct{}),
	}
}

func (ur *UploadReader) Write(data []byte) (err error) {
	ur.mutex.Lock()
	defer ur.mutex.Unlock()
	select {
	case <-ur.close:
		err = errors.New("closed")
	case ur.dataChan <- data:
	}
	return
}

func (ur *UploadReader) Read(p []byte) (n int, err error) {
	var data = ur.dataCache
	var dataLen = len(data)
	if dataLen > 0 {
		n = copy(p, data)
	} else {
		select {
		case <-ur.close:
			err = errors.New("closed")
			return
		case data = <-ur.dataChan:
			n = copy(p, data)
			dataLen = len(data)
		}
	}
	if dataLen == 0 {
		err = io.EOF
		return
	}
	if n < dataLen {
		ur.dataCache = data[n:]
	} else {
		ur.dataCache = nil
	}
	return
}

func (ur *UploadReader) Close() {
	ur.mutex.Lock()
	defer ur.mutex.Unlock()

	defer func() {
		if x := recover(); x != nil {
			util.Logger.Error("UploadReader close panic error:", zap.Any("error", x))
		}
	}()
	close(ur.close)
	close(ur.dataChan)
}

var chunkUploadCache = map[string]*ChunkUpload{}
var chunkUploadCacheLock = &sync.Mutex{}

func getChunkUpload(key string) (res *ChunkUpload) {
	chunkUploadCacheLock.Lock()
	res = chunkUploadCache[key]
	chunkUploadCacheLock.Unlock()
	return
}
func setChunkUpload(key string, chunkUpload *ChunkUpload) (res *ChunkUpload) {
	chunkUploadCacheLock.Lock()
	chunkUploadCache[key] = chunkUpload
	chunkUploadCacheLock.Unlock()
	return
}
func removeChunkUpload(key string) {
	chunkUploadCacheLock.Lock()
	res := chunkUploadCache[key]
	delete(chunkUploadCache, key)
	chunkUploadCacheLock.Unlock()
	if res != nil {
		res.close()
	}
	return
}

type ChunkUpload struct {
	chunkUploadKey string
	*worker
	param         *BaseParam
	fileWorkerKey string
	dir           string
	fullPath      string
	filename      string
	size          int64
	uploadReader  *UploadReader
	closed        bool

	callStop *bool
}

func (this_ *ChunkUpload) close() {
	if this_.closed {
		return
	}
	this_.closed = true
	if this_.uploadReader != nil {
		this_.uploadReader.Close()
	}
}
func (this_ *ChunkUpload) Start() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	service, err := this_.GetService(this_.fileWorkerKey, this_.param)
	if err != nil {
		return
	}

	if !strings.HasSuffix(this_.dir, "/") {
		this_.dir += "/"
	}

	if strings.HasPrefix(this_.fullPath, "/") {
		this_.fullPath = this_.fullPath[1:]
	}

	path := this_.dir + this_.filename
	if len(this_.fullPath) > 0 {
		path = this_.dir + this_.fullPath
	}
	this_.callStop = new(bool)
	*this_.callStop = false

	progress := newProgress(this_.param, "upload", func() {
		this_.close()
		*this_.callStop = true
	})

	progress.Data.FileWorkerKey = this_.fileWorkerKey
	progress.Data.Dir = this_.dir
	progress.Data.FullPath = this_.fullPath
	progress.Data.Filename = this_.filename
	progress.Data.Path = path
	progress.Data.Size = this_.size
	progress.Data.SuccessSize = 0

	var exist bool
	exist, err = service.Exist(path)

	if exist {
		var action string
		action, err = progress.waitAction("文件["+this_.filename+"]已存在，是否覆盖？",
			[]*Action{
				newAction("是", "yes", "color-green"),
				newAction("否", "no", "color-orange"),
			})
		if err != nil {
			return
		}
		if action != "yes" {
			return
		}
	}

	this_.uploadReader = NewUploadReader()

	go func() {

		defer func() {
			removeChunkUpload(this_.chunkUploadKey)
			this_.close()
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprint(e))
			}

			progress.end(err)
		}()
		err = service.Write(path, this_.uploadReader, func(readSize int64, writeSize int64) {
			if progress.Data.FileInfo == nil {
				pathDir := path[0:strings.LastIndex(path, "/")]
				var pathDirExist bool
				pathDirExist, _ = service.Exist(pathDir)
				progress.Data.FileInfo, _ = service.File(path)
				if !pathDirExist {
					progress.Data.FileDir, _ = service.File(pathDir)
				}
			}
			progress.Data.SuccessSize = writeSize
			if progress.Data.FileInfo != nil {
				progress.Data.FileInfo.Size = writeSize
			}
			progress.Data.Timestamp = time.Now().UnixMilli()
		}, this_.callStop)
		if err != nil {
			return
		}

	}()
	return

}

func (this_ *ChunkUpload) Append(bs []byte, isEnd bool) (err error) {
	if this_.closed {
		err = errors.New("closed")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	err = this_.uploadReader.Write(bs)
	if err != nil {
		return
	}
	if isEnd {
		err = this_.uploadReader.Write([]byte{})
		if err != nil {
			return
		}
	}
	return
}
