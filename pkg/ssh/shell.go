package ssh

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"teamide/pkg/util"
	"time"
)

type ShellClient struct {
	Client
	shellSession                     *ssh.Session
	startReadChannel                 bool
	shellOK                          bool
	DisableZModemSZ, DisableZModemRZ bool
	ZModemSZ, ZModemRZ, ZModemSZOO   bool
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

func NewSSHShell(terminalSize TerminalSize, sshSession *ssh.Session) (err error) {
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
	if terminalSize.Width > 0 && terminalSize.Height > 0 {
		req.Width = uint32(terminalSize.Width)
		req.Height = uint32(terminalSize.Height)
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
func (this_ *ShellClient) startShell(terminalSize TerminalSize) (err error) {
	this_.shellOK = false
	this_.startReadChannel = false
	defer func() {
		if x := recover(); x != nil {
			util.Logger.Error("SSH Shell Start Error", zap.Any("err", x))
			return
		}
		this_.shellSession = nil
	}()
	if this_.shellSession != nil {
		err = this_.shellSession.Close()
		if err != nil {
			util.Logger.Error("SSH Shell Shell Session Close Error", zap.Error(err))
		}
		this_.shellSession = nil
	}
	err = this_.initClient()
	if err != nil {
		util.Logger.Error("Create Shell Init Client Error", zap.Error(err))
		this_.WSWriteError("SSH?????????????????????:" + err.Error())
		return
	}
	this_.shellSession, err = this_.sshClient.NewSession()
	if err != nil {
		util.Logger.Error("Create Shell Open Channel Error", zap.Error(err))
		this_.WSWriteError("SSH??????????????????:" + err.Error())
		return
	}
	defer this_.closeSession(this_.shellSession)
	defer this_.WSWriteEvent("ssh session closed")

	err = NewSSHShell(terminalSize, this_.shellSession)
	if err != nil {
		util.Logger.Error("Create Shell Error", zap.Error(err))
		this_.WSWriteError("SSH Shell????????????:" + err.Error())
		return
	}
	this_.shellOK = true
	var errReader io.Reader
	errReader, err = this_.shellSession.StderrPipe()
	go func() {
		err = this_.startRead(errReader, true)
	}()
	var reader io.Reader
	reader, err = this_.shellSession.StdoutPipe()
	err = this_.startRead(reader, false)
	return
}

func (this_ *ShellClient) startRead(reader io.Reader, isError bool) (err error) {

	for {
		if !this_.startReadChannel {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if err != nil {
			util.Logger.Error("SSH Shell Stderr Pipe Error", zap.Error(err))
			continue
		}
		var buffSize = 1024 * 8
		var bs = make([]byte, buffSize)
		var n int
		n, err = reader.Read(bs)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			util.Logger.Error("SSH Shell ??????????????????", zap.Error(err))
			//this_.WSWriteError("SSH Shell ??????????????????:" + err.Error())
			continue
		}
		if isError {

		}

		var isZModem bool
		isZModem, _ = this_.processZModem(bs, n, buffSize)
		if !isZModem {
			out := TeamIDEBinaryStartBytes
			if n == buffSize {
				out = append(out, bs...)
			} else {
				out = append(out, bs[0:n]...)
			}
			this_.WSWriteBinary(out)
		}

	}
}

func (this_ *ShellClient) start() {
	ShellCache[this_.Token] = this_
	go this_.ListenWS(this_.onEvent, this_.ONSSHMessage, this_.CloseClient)
	this_.WSWriteEvent("shell ready")
}

func (this_ *ShellClient) onEvent(event string) {
	var err error
	util.Logger.Info("SSH Shell On Event:", zap.Any("event", event))

	if strings.HasPrefix(event, "shell start") {
		jsonStr := event[len("shell start"):]
		var terminalSize *TerminalSize
		if jsonStr != "" {
			_ = json.Unmarshal([]byte(jsonStr), &terminalSize)
		}
		go func() {
			err = this_.startShell(*terminalSize)
			if err != nil {
				util.Logger.Error("SSH Shell Start Shell error", zap.Error(err))
			}
		}()
		for {
			time.Sleep(100 * time.Millisecond)
			if err != nil || this_.shellOK {
				break
			}
		}
		if err != nil {
			this_.WSWriteEvent("shell create error")
			return
		}
		//time.Sleep(1000 * time.Millisecond)
		this_.WSWriteEvent("shell created")
		this_.startReadChannel = true
		return
	} else if strings.HasPrefix(event, "change size") {
		jsonStr := event[len("change size"):]
		var terminalSize *TerminalSize
		err = json.Unmarshal([]byte(jsonStr), &terminalSize)
		if err != nil {
			return
		}
		err = this_.changeSize(*terminalSize)

	} else if strings.HasPrefix(event, "shell cancel upload file") {
		// ????????????
		this_.SSHWrite(ZModemCancel)
		//this_.WSWrite([]byte("????????????"))
	}

	switch strings.ToLower(event) {
	case "ssh session close":
		this_.closeSession(this_.shellSession)
	}
}

func (this_ *ShellClient) ONSSHMessage(bs []byte) {
	this_.SSHWrite(bs)
}

func (this_ *ShellClient) SSHWrite(bs []byte) {
	defer func() {
		if x := recover(); x != nil {
			util.Logger.Error("SSH Shell Write Error", zap.Any("err", x))
			return
		}
	}()
	if this_.shellSession == nil {
		return
	}
	var err error
	var writer io.Writer
	writer, err = this_.shellSession.StdinPipe()
	if err != nil {
		util.Logger.Error("SSH Shell Stderr Pipe Error", zap.Error(err))
	}

	_, err = writer.Write(bs)
	if err != nil {
		this_.WSWriteError("SSH Shell Write??????:" + err.Error())
	}
}
