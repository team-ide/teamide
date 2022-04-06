package toolbox

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type SSHShellClient struct {
	SSHClient
	sshChannel    ssh.Channel
	isClosedShell bool
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (this_ *SSHShellClient) CloseShell() {

	this_.isClosedShell = true
	if this_.sshChannel != nil {
		err := this_.sshChannel.Close()
		if err != nil {
			this_.Logger.Error("SSH Session close error", zap.Error(err))
		}
	}
	this_.sshChannel = nil
	this_.CloseClient()
}

func (this_ *SSHShellClient) initShell() (err error) {
	if this_.isClosedShell || this_.sshChannel == nil {
		err = this_.createShell()
	}
	return
}

func (this_ *SSHShellClient) createShell() (err error) {
	err = this_.initClient()
	if err != nil {
		return
	}

	this_.sshChannel, _, err = this_.sshClient.OpenChannel("session", nil)
	if err != nil {
		this_.WSWriteError("SSH会话创建失败:" + err.Error())
		this_.CloseShell()
		return
	}

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
		Term: "xterm",
		//Columns:  100,
		//Rows:     40,
		//Width:    uint32(100 * 8),
		//Height:   uint32(40 * 8),
		Modelist: string(modeList),
	}

	ok, err := this_.sshChannel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		this_.WSWriteError("SSH会话请求失败:" + err.Error())
		this_.CloseShell()
		return
	}
	ok, err = this_.sshChannel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		this_.WSWriteError("SSH Shell创建失败:" + err.Error())
		this_.CloseShell()
		return
	}

	return
}

func (this_ *SSHShellClient) start() (err error) {
	err = this_.initShell()
	if err != nil {
		return
	}

	return
}
