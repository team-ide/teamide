package ssh

import (
	"errors"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"teamide/pkg/terminal"
)

type ShellClient struct {
	Client
	shellSession                     *ssh.Session
	startReadChannel                 bool
	shellOK                          bool
	DisableZModemSZ, DisableZModemRZ bool
	ZModemSZ, ZModemRZ, ZModemSZOO   bool
	rzFileSize                       int64
	rzFileUploadSize                 int64
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type TerminalSize struct {
	Cols   int `json:"cols"`
	Rows   int `json:"rows"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (this_ *ShellClient) changeSize(terminalSize TerminalSize) (err error) {

	if this_.shellSession == nil {
		return
	}
	if terminalSize.Cols > 0 && terminalSize.Rows > 0 {
		err = this_.shellSession.WindowChange(terminalSize.Rows, terminalSize.Cols)
		if err != nil {
			util.Logger.Error("SSH Shell Session Window Change error", zap.Error(err))
			return
		}
	}
	if terminalSize.Width > 0 && terminalSize.Height > 0 {
		err = this_.shellSession.WindowChange(terminalSize.Height, terminalSize.Width)
		if err != nil {
			util.Logger.Error("SSH Shell Session Window Change error", zap.Error(err))
			return
		}
	}
	return
}

func (this_ *ShellClient) closeSession(session *ssh.Session) {
	if session == nil {
		return
	}
	err := session.Close()
	if err != nil {
		if err == io.EOF {
			return
		}
		util.Logger.Error("SSH Shell Session Close Error", zap.Error(err))
		return
	}
}

func NewSSHShell(terminalSize *terminal.Size, sshSession *ssh.Session) (err error) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Modelist: string(modeList),
	}
	if terminalSize.Cols > 0 && terminalSize.Rows > 0 {
		req.Columns = uint32(terminalSize.Cols)
		req.Rows = uint32(terminalSize.Rows)
	}
	_, err = sshSession.SendRequest("pty-req", true, ssh.Marshal(&req))
	if err != nil {
		return
	}

	var ok bool
	ok, err = sshSession.SendRequest("shell", true, nil)
	if !ok || err != nil {
		if err != nil {
			err = errors.New("SSH Shell Send Request Fail")
		}
		return
	}
	return
}
