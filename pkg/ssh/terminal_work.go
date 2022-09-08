package ssh

import (
	"go.uber.org/zap"
	"io"
	"os/exec"
	"runtime"
	"syscall"
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
	config  *Config
	reader  io.Reader
	writer  io.Writer
	stdout  io.ReadCloser
	stdin   io.WriteCloser
	onClose func()
}

func (this_ *terminalService) IsWindows() (isWindows bool, err error) {
	isWindows = runtime.GOOS == "windows"
	return
}

func (this_ *terminalService) getCmd() (cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		cmd = exec.Command("bash", "-c")
	}
	return
}

func (this_ *terminalService) Stop() {
	if this_.stdout != nil {
		_ = this_.stdout.Close()
	}
	if this_.stdin != nil {
		_ = this_.stdin.Close()
	}
}

func (this_ *terminalService) ChangeSize(size *terminal.Size) (err error) {
	return
}

func (this_ *terminalService) Start(size *terminal.Size) (err error) {

	cmd := this_.getCmd()

	this_.stdout, err = cmd.StdoutPipe()
	if err != nil {
		util.Logger.Error("cmd StdoutPipe error", zap.Error(err))
		return
	}

	go func() { _ = this_.readStdout() }()

	this_.stdin, err = cmd.StdinPipe()
	if err != nil {
		util.Logger.Error("cmd StdinPipe error", zap.Error(err))
		return
	}

	go func() { _ = this_.writeStdin() }()

	err = cmd.Start()
	if err != nil {
		util.Logger.Error("cmd Start error", zap.Error(err))
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

func (this_ *terminalService) readStdout() (err error) {
	defer func() { this_.onClose() }()
	defer func() { _ = this_.stdout.Close() }() // 保证关闭输出流

	var buf = make([]byte, 1024*32)
	err = util.Read(this_.stdout, buf, func(n int) (err error) {
		_, err = this_.writer.Write(buf[:n])

		return
	})

	return
}

func (this_ *terminalService) writeStdin() (err error) {
	defer func() { _ = this_.stdin.Close() }() // 保证关闭输出流

	var buf = make([]byte, 1024*32)
	err = util.Read(this_.reader, buf, func(n int) (err error) {
		_, err = this_.stdin.Write(buf[:n])
		return
	})

	return
}
