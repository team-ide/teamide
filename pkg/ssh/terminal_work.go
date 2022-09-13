package ssh

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"teamide/pkg/terminal"
	"teamide/pkg/util"
)

func NewTerminalService(config *Config) (res *terminalService) {
	res = &terminalService{
		config: config,
	}
	return
}

type terminalService struct {
	config     *Config
	sshClient  *ssh.Client
	sshSession *ssh.Session
	reader     io.Reader
	writer     io.Writer
	stdout     io.Reader
	stderr     io.Reader
	stdin      io.Writer
	onClose    func()
}

func (this_ *terminalService) IsWindows() (isWindows bool, err error) {
	isWindows = false
	return
}

func (this_ *terminalService) Stop() {
	if this_.sshSession != nil {
		_ = this_.sshSession.Close()
	}
	if this_.sshClient != nil {
		_ = this_.sshClient.Close()
	}
}

func (this_ *terminalService) ChangeSize(size *terminal.Size) (err error) {
	return
}

func (this_ *terminalService) Start(size *terminal.Size) (err error) {

	this_.sshClient, err = NewClient(*this_.config)
	if err != nil {
		util.Logger.Error("NewClient error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH NewClient success", zap.Any("address", this_.config.Address))
	go func() {
		err = this_.sshClient.Wait()
		this_.Stop()
	}()

	this_.sshSession, err = this_.sshClient.NewSession()
	if err != nil {
		util.Logger.Error("SSH NewSession Error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH NewSession success", zap.Any("address", this_.config.Address))

	err = NewSSHShell(size, this_.sshSession)
	if err != nil {
		util.Logger.Error("Create SSH Shell Error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH Format Shell Session success", zap.Any("address", this_.config.Address))
	this_.stdout, err = this_.sshSession.StdoutPipe()
	if err != nil {
		util.Logger.Error("ssh session StdoutPipe error", zap.Error(err))
		return
	}
	this_.stderr, err = this_.sshSession.StderrPipe()
	if err != nil {
		util.Logger.Error("ssh session StderrPipe error", zap.Error(err))
		return
	}
	this_.stdin, err = this_.sshSession.StdinPipe()
	if err != nil {
		util.Logger.Error("ssh session StdinPipe error", zap.Error(err))
		return
	}

	return
}

func (this_ *terminalService) Write(buf []byte) (n int, err error) {
	n, err = this_.stdin.Write(buf)
	return
}

func (this_ *terminalService) Read(buf []byte) (n int, err error) {
	n, err = this_.stdout.Read(buf)
	return
}

func (this_ *terminalService) ReadError(buf []byte) (n int, err error) {
	return
}
