package toolbox

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
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

type safeWrite struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *safeWrite) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}
func (w *safeWrite) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}
func (w *safeWrite) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
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
	if this_.client, err = ssh.Dial("tcp", this_.Config.Address, clientConfig); err != nil {
		this_.Close()
		return
	}

	this_.channel, _, err = this_.client.OpenChannel("session", []byte(token))
	if err != nil {
		this_.Close()
		return
	}

	ok, err := this_.channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		this_.Close()
		return
	}
	// 第一个协程获取用户的输入
	go func() {
		bufWriter := bufio.NewWriter(this_.channel)
		for {
			if this_.isClosed {
				return
			}
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil && err != io.EOF {
				fmt.Println(err)
				this_.Close()
				return
			}
			fmt.Println("ws read:" + string(p))
			if this_.isClosed {
				return
			}
			_, err = bufWriter.Write(p)
			if err != nil {
				fmt.Println(err)
				this_.Close()
				return
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		bufReader := bufio.NewReader(this_.channel)
		for {
			if this_.isClosed {
				return
			}
			var bs []byte = make([]byte, 1024)
			n, err := bufReader.Read(bs)

			if err != nil && err != io.EOF {
				fmt.Println(err)
				this_.Close()
				return
			}
			bs = bs[0:n]
			fmt.Print("ssh read:" + string(bs))
			if this_.isClosed {
				return
			}
			err = ws.WriteMessage(websocket.BinaryMessage, bs)
			if err != nil {
				fmt.Println(err)
				this_.Close()
				return
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
