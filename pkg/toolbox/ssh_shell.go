package toolbox

import (
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"time"
)

type SSHShellClient struct {
	SSHClient
	sessionChannel   ssh.Channel
	startReadChannel bool
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

func (this_ *SSHShellClient) startShell(terminalSize TerminalSize) (err error) {
	defer func() {
		if x := recover(); x != nil {
			this_.Logger.Error("SSH Shell Start Error", zap.Any("err", x))
			return
		}
		this_.sessionChannel = nil
	}()
	if this_.sessionChannel != nil {
		err = this_.sessionChannel.Close()
		if err != nil {
			this_.Logger.Error("SSH Shell Shell Session Close Error", zap.Error(err))
		}
		this_.sessionChannel = nil
	}
	err = this_.initClient()
	if err != nil {
		this_.Logger.Error("createShell initClient error", zap.Error(err))
		this_.WSWriteError("SSH客户端创建失败:" + err.Error())
		return
	}
	this_.sessionChannel, _, err = this_.sshClient.OpenChannel("session", nil)
	if err != nil {
		this_.Logger.Error("createShell OpenChannel error", zap.Error(err))
		this_.WSWriteError("SSH会话创建失败:" + err.Error())
		return
	}
	defer this_.sessionChannel.Close()

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
	ok, err := this_.sessionChannel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		this_.Logger.Error("createShell SendRequest pty-req error", zap.Error(err))
		this_.WSWriteError("SSH会话请求失败:" + err.Error())
		return
	}
	ok, err = this_.sessionChannel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		this_.Logger.Error("createShell SendRequest shell error", zap.Error(err))
		this_.WSWriteError("SSH Shell创建失败:" + err.Error())
		return
	}

	for {
		if !this_.startReadChannel {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		var bs = make([]byte, 1024)
		var n int
		n, err = this_.sessionChannel.Read(bs)
		if err != nil {
			this_.Logger.Error("SSH Shell 消息读取异常", zap.Error(err))
			if err == io.EOF {
				return
			}
			//this_.WSWriteError("SSH Shell 消息读取失败:" + err.Error())
			continue
		}
		bs = bs[0:n]
		this_.WSWrite(bs)
	}
	return
}

func (this_ *SSHShellClient) start() {
	go this_.ListenWS(this_.onEvent, this_.onMessage, this_.CloseClient)
	this_.WSWriteEvent("shell ready")
}

func (this_ *SSHShellClient) onEvent(event string) {
	var err error
	this_.Logger.Info("SSH Shell On Event:", zap.Any("event", event))

	if strings.HasPrefix(event, "shell start") {
		jsonStr := event[len("shell start"):]
		var terminalSize *TerminalSize
		if jsonStr != "" {
			_ = json.Unmarshal([]byte(jsonStr), &terminalSize)
		}
		go func() {
			err = this_.startShell(*terminalSize)
		}()
		go func() {
			for {
				if
			}
		}()
		time.Sleep(2000 * time.Millisecond)
		if err != nil {
			return
		}
		this_.WSWriteEvent("shell created")
		time.Sleep(1000 * time.Millisecond)
		this_.startReadChannel = true
		return
	}
	switch strings.ToLower(event) {
	}
}

func (this_ *SSHShellClient) onMessage(bs []byte) {
	defer func() {
		if x := recover(); x != nil {
			this_.Logger.Error("SSH Shell Write Error", zap.Any("err", x))
			return
		}
	}()
	if this_.sessionChannel == nil {
		return
	}
	var err error

	_, err = this_.sessionChannel.Write(bs)
	if err != nil {
		this_.WSWriteError("SSH Shell Write失败:" + err.Error())
	}
}
