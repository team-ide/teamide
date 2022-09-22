package terminal

import (
	"go.uber.org/zap"
	"io"
	"sync"
	"teamide/pkg/util"
)

func NewLocalService() (res *localService) {
	res = &localService{}
	return
}

type localService struct {
	reader       io.Reader
	writer       io.Writer
	stdout       io.ReadCloser
	stderr       io.ReadCloser
	stdin        io.WriteCloser
	onClose      func()
	readeLock    sync.Mutex
	readeErrLock sync.Mutex
	writeLock    sync.Mutex
}

func (this_ *localService) IsWindows() (isWindows bool, err error) {
	return IsWindows(), nil
}

func (this_ *localService) Stop() {
	if this_.stdout != nil {
		_ = this_.stdout.Close()
	}
	if this_.stderr != nil {
		_ = this_.stderr.Close()
	}
	if this_.stdin != nil {
		_ = this_.stdin.Close()
	}
}

func (this_ *localService) ChangeSize(size *Size) (err error) {
	return
}

func (this_ *localService) Start(size *Size) (err error) {

	cmd := getCmd()

	this_.stdout, err = cmd.StdoutPipe()
	if err != nil {
		util.Logger.Error("cmd StdoutPipe error", zap.Error(err))
		return
	}

	this_.stderr, err = cmd.StderrPipe()
	if err != nil {
		util.Logger.Error("cmd StderrPipe error", zap.Error(err))
		return
	}

	this_.stdin, err = cmd.StdinPipe()
	if err != nil {
		util.Logger.Error("cmd StdinPipe error", zap.Error(err))
		return
	}

	err = cmd.Start()
	if err != nil {
		util.Logger.Error("cmd Start error", zap.Error(err))
		return
	}

	return
}

func (this_ *localService) Write(buf []byte) (n int, err error) {
	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	n = len(buf)
	err = util.Write(this_.stdin, buf, nil)
	return
}

func (this_ *localService) Read(buf []byte) (n int, err error) {
	this_.readeLock.Lock()
	defer this_.readeLock.Unlock()

	n, err = this_.stdout.Read(buf)
	return
}

func (this_ *localService) ReadError(buf []byte) (n int, err error) {
	this_.readeErrLock.Lock()
	defer this_.readeErrLock.Unlock()

	n, err = this_.stderr.Read(buf)
	return
}
