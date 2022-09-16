package module_terminal

import (
	"errors"
	"io"
	"teamide/pkg/util"
)

func (this_ *worker) GetReaderChan(key string) (reader io.Reader) {
	this_.readerChanCacheLock.Lock()
	defer this_.readerChanCacheLock.Unlock()

	if this_.readerChanCache[key] == nil {
		return
	}
	chan_ := this_.readerChanCache[key]
	reader = <-chan_
	delete(this_.readerChanCache, key)
	close(chan_)
	return
}

func (this_ *worker) NewReaderChan() (key string) {
	this_.readerChanCacheLock.Lock()
	defer this_.readerChanCacheLock.Unlock()

	key = util.UUID()
	this_.readerChanCache[key] = make(chan io.Reader)
	return
}

func (this_ *worker) SetReaderChan(key string, reader io.Reader) (err error) {
	this_.readerChanCacheLock.Lock()
	defer this_.readerChanCacheLock.Unlock()

	if this_.readerChanCache[key] == nil {
		err = errors.New("reader chan is null")
		return
	}
	chan_ := this_.readerChanCache[key]
	chan_ <- reader
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
