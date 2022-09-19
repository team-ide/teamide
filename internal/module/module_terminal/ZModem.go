package module_terminal

import (
	"errors"
	"io"
	"mime/multipart"
	"teamide/pkg/util"
)

func (this_ *worker) GetFileHeadersChan(key string) (fileHeaders []*multipart.FileHeader) {
	this_.fileHeadersChanCacheLock.Lock()
	defer this_.fileHeadersChanCacheLock.Unlock()

	if this_.fileHeadersChanCache[key] == nil {
		return
	}
	chan_ := this_.fileHeadersChanCache[key]
	fileHeaders = <-chan_
	delete(this_.fileHeadersChanCache, key)
	close(chan_)
	return
}

func (this_ *worker) NewFileHeadersChan() (key string) {
	this_.fileHeadersChanCacheLock.Lock()
	defer this_.fileHeadersChanCacheLock.Unlock()

	key = util.UUID()
	this_.fileHeadersChanCache[key] = make(chan []*multipart.FileHeader)
	return
}

func (this_ *worker) SetFileHeadersChan(key string, fileHeaders []*multipart.FileHeader) (err error) {
	this_.fileHeadersChanCacheLock.Lock()
	defer this_.fileHeadersChanCacheLock.Unlock()

	if this_.fileHeadersChanCache[key] == nil {
		err = errors.New("fileHeaders chan is null")
		return
	}
	chan_ := this_.fileHeadersChanCache[key]
	chan_ <- fileHeaders
	return
}

func (this_ *worker) GetWriterChan(key string) (writer io.Writer) {
	this_.writerChanCacheLock.Lock()
	defer this_.writerChanCacheLock.Unlock()

	if this_.writerChanCache[key] == nil {
		return
	}
	chan_ := this_.writerChanCache[key]
	delete(this_.writerChanCache, key)
	writer = <-chan_
	close(chan_)
	return
}

func (this_ *worker) NewWriterChan() (key string) {
	this_.writerChanCacheLock.Lock()
	defer this_.writerChanCacheLock.Unlock()

	key = util.UUID()
	this_.writerChanCache[key] = make(chan io.Writer)
	return
}

func (this_ *worker) SetWriterChan(key string, writer io.Writer) (err error) {
	this_.writerChanCacheLock.Lock()
	defer this_.writerChanCacheLock.Unlock()

	chan_ := this_.writerChanCache[key]
	if chan_ == nil {
		err = errors.New("writer chan is null")
		return
	}
	chan_ <- writer
	return
}
