package terminal

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"os/exec"
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
	cmd          *exec.Cmd
}

func (this_ *localService) IsWindows() (isWindows bool, err error) {
	return IsWindows(), nil
}

func (this_ *localService) Stop() {
	if this_.stdout != nil {
		_ = this_.stdout.Close()
		this_.stdout = nil
	}
	if this_.stderr != nil {
		_ = this_.stderr.Close()
		this_.stderr = nil
	}
	if this_.stdin != nil {
		_ = this_.stdin.Close()
		this_.stdin = nil
	}
}

func (this_ *localService) ChangeSize(size *Size) (err error) {
	return
}

func (this_ *localService) Start(size *Size) (err error) {

	this_.cmd = getCmd()

	this_.stdout, err = this_.cmd.StdoutPipe()
	if err != nil {
		util.Logger.Error("cmd StdoutPipe error", zap.Error(err))
		return
	}

	this_.stderr, err = this_.cmd.StderrPipe()
	if err != nil {
		util.Logger.Error("cmd StderrPipe error", zap.Error(err))
		return
	}

	this_.stdin, err = this_.cmd.StdinPipe()
	if err != nil {
		util.Logger.Error("cmd StdinPipe error", zap.Error(err))
		return
	}

	err = this_.cmd.Start()
	if err != nil {
		util.Logger.Error("cmd Start error", zap.Error(err))
		return
	}

	util.Logger.Info("terminal local start success", zap.Any("path", this_.cmd.Path), zap.Any("args", this_.cmd.Args), zap.Any("env", this_.cmd.Env))

	return
}

func (this_ *localService) Write(buf []byte) (n int, err error) {

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Write err", zap.Any("err", e))
		}
	}()
	if this_.stdin == nil {
		err = errors.New("stdin is close")
		return
	}

	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	//util.Logger.Info("local terminal write start")

	n = len(buf)
	err = util.Write(this_.stdin, buf, nil)
	return
}

func (this_ *localService) Read(buf []byte) (n int, err error) {

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("Read err", zap.Any("err", e))
		}
	}()
	if this_.stdout == nil {
		err = errors.New("stdout is close")
		return
	}

	this_.readeLock.Lock()
	defer this_.readeLock.Unlock()

	//util.Logger.Info("local terminal read start")
	n, err = this_.stdout.Read(buf)
	return
}

func (this_ *localService) ReadError(buf []byte) (n int, err error) {

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("ReadError err", zap.Any("err", e))
		}
	}()
	if this_.stderr == nil {
		err = errors.New("stderr is close")
		return
	}

	this_.readeErrLock.Lock()
	defer this_.readeErrLock.Unlock()

	//util.Logger.Info("local terminal read error start")
	n, err = this_.stderr.Read(buf)
	return
}
