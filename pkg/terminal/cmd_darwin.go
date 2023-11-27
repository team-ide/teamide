package terminal

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func IsWindows() bool {
	return false
}

func start(size *Size) (starter *terminalStart, err error) {
	obj := ptyMasterNew()

	command := "bash"
	_, err = os.Stat("/bin/bash")
	if os.IsNotExist(err) {
		command = "sh"
	}
	err = obj.Start(command, nil, nil, size.Cols, size.Rows)
	if err != nil {
		return
	}
	starter = &terminalStart{
		Stop: func() {
			_ = obj.Stop()
		},
		Write_: obj.Write,
		Read_:  obj.Read,
	}
	return
}

// This defines a PTY Master whih will encapsulate the command we want to run, and provide simple
// access to the command, to write and read IO, but also to control the window size.
type ptyMaster struct {
	ptyFile *os.File
	command *exec.Cmd
}

func ptyMasterNew() *ptyMaster {
	return &ptyMaster{}
}

func (this_ *ptyMaster) Start(command string, args []string, envVars []string, cols int, rows int) (err error) {
	this_.command = exec.Command(command, args...)
	this_.command.Env = envVars
	this_.ptyFile, err = pty.Start(this_.command)

	if err != nil {
		return
	}

	this_.SetWinSize(rows, cols)
	return
}

func (this_ *ptyMaster) Write(b []byte) (int, error) {
	return this_.ptyFile.Write(b)
}

func (this_ *ptyMaster) Read(b []byte) (int, error) {
	return this_.ptyFile.Read(b)
}

func (this_ *ptyMaster) SetWinSize(rows, cols int) {
	pty.Setsize(this_.ptyFile, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

func (this_ *ptyMaster) Wait() (err error) {
	err = this_.command.Wait()
	return
}

func (this_ *ptyMaster) Stop() (err error) {
	signal.Ignore(syscall.SIGWINCH)

	this_.command.Process.Signal(syscall.SIGTERM)
	// TODO: Find a proper wai to close the running command. Perhaps have a timeout after which,
	// if the command hasn't reacted to SIGTERM, then send a SIGKILL
	// (bash for example doesn't finish if only a SIGTERM has been sent)
	this_.command.Process.Signal(syscall.SIGKILL)
	return
}
