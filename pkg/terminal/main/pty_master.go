package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/creack/pty"
)

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
	this_.ptyFile, err = pty.Start(this_.command, nil, nil)

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

func (this_ *ptyMaster) Refresh() {
	// We wanna force the app to re-draw itself, but there doesn't seem to be a way to do that
	// so we fake it by resizing the window quickly, making it smaller and then back big
	cols, rows, err := this_.GetWinSize()

	if err != nil {
		return
	}

	this_.SetWinSize(rows-1, cols)

	go func() {
		time.Sleep(time.Millisecond * 50)
		this_.SetWinSize(rows, cols)
	}()
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
