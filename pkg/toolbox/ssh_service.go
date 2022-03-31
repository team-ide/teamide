package toolbox

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type SSHClient struct {
	Config   *SSHConfig
	client   *ssh.Client
	channel  ssh.Channel
	session  *ssh.Session
	ws       *websocket.Conn
	isClosed bool
}

func (this_ *SSHClient) Close() {
	if this_.isClosed {
		return
	}
	this_.isClosed = true
	var err error
	if this_.session != nil {
		err = this_.session.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	if this_.channel != nil {
		err = this_.channel.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	if this_.client != nil {
		err = this_.client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	if this_.ws != nil {
		err = this_.ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (this_ *SSHClient) Start(token string, ws *websocket.Conn) (err error) {
	if this_.Config == nil {
		err = errors.New("令牌会话丢失")
		return
	}
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		config       ssh.Config
	)
	auth = []ssh.AuthMethod{}
	if this_.Config.Password != "" {
		auth = append(auth, ssh.Password(this_.Config.Password))
	}
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:            this_.Config.User,
		Auth:            auth,
		Timeout:         5 * time.Second,
		Config:          config,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
	}
	if this_.client, err = ssh.Dial(this_.Config.Type, this_.Config.Address, clientConfig); err != nil {
		fmt.Println(err)
		this_.Close()
		return
	}

	this_.channel, _, err = this_.client.OpenChannel("session", nil)
	if err != nil {
		fmt.Println(err)
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

	ok, err := this_.channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		fmt.Println(err)
		this_.Close()
		return
	}
	ok, err = this_.channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		fmt.Println(err)
		this_.Close()
		return
	}
	// 第一个协程获取用户的输入
	go func() {
		for {
			if this_.isClosed {
				return
			}
			_, p, err := ws.ReadMessage()
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
				_, err = this_.channel.Write(p)
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
			n, err := this_.channel.Read(bs)

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
				err = ws.WriteMessage(websocket.BinaryMessage, bs)
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
		Config: sshConfig,
	}
	err = SSHClient.Start(token, ws)

	return
}
