package toolbox

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (this_ *SSHClient) StartSSH(ws *websocket.Conn) (err error) {
	this_.ws = ws
	err = this_.initClient()
	if err != nil {
		fmt.Println("StartSSH error", err)
		this_.Close()
		return
	}
	this_.sshChannel, _, err = this_.sshClient.OpenChannel("session", nil)
	if err != nil {
		fmt.Println("SSH OpenChannel error", err)
		this_.Close()
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
		fmt.Println("SSH SendRequest error", err)
		this_.Close()
		return
	}
	ok, err = this_.sshChannel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		fmt.Println("SSH SendRequest error", err)
		this_.Close()
		return
	}
	// 第一个协程获取用户的输入
	go func() {
		for {
			if this_.isClosed {
				return
			}
			_, p, err := this_.ws.ReadMessage()
			if err != nil && err != io.EOF {
				fmt.Println("ws read err:", err)
				this_.Close()
				return
			}
			//fmt.Println("ws read:" + string(p))
			if len(p) > 0 {
				if this_.isClosed {
					return
				}
				//fmt.Println("ssh write:", p)
				_, err = this_.sshChannel.Write(p)
				if err != nil {
					fmt.Println("ssh write err:", err)
					this_.Close()
					return
				}
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		for {
			if this_.isClosed {
				return
			}
			var bs []byte = make([]byte, 1024)
			n, err := this_.sshChannel.Read(bs)

			if err != nil && err != io.EOF {
				fmt.Println("ssh read err:", err)
				this_.Close()
				return
			}
			bs = bs[0:n]
			//fmt.Print("ssh read:" + string(bs))
			if len(bs) > 0 {
				if this_.isClosed {
					return
				}
				//fmt.Println("ws write:", bs)
				err = this_.WSWriteMessage(bs)
				if err != nil {
					fmt.Println("ws write err:", err)
					this_.Close()
					return
				}
			}
		}

	}()
	return
}

func WSSSHConnection(token string, ws *websocket.Conn) (err error) {
	var sshConfig *SSHConfig = sshTokenCache[token]
	if sshConfig == nil {
		err = errors.New("令牌会话丢失")
		return
	}
	SSHClient := &SSHClient{
		Token:  token,
		Config: sshConfig,
	}
	err = SSHClient.StartSSH(ws)

	return
}
