package toolbox

import (
	"bufio"
	"errors"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

type SSHClient struct {
	Config  *SSHConfig
	Session *ssh.Session
	Client  *ssh.Client
	channel ssh.Channel
}

func (this_ *SSHClient) GenerateClient(ws *websocket.Conn) (err error) {
	if this_.Config == nil {
		err = errors.New("令牌会话丢失")
		return
	}
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
	)
	auth = make([]ssh.AuthMethod, 0)
	if this_.Config.Password != "" {
		auth = append(auth, ssh.Password(this_.Config.Password))
	}
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    this_.Config.User,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	if client, err = ssh.Dial("tcp", this_.Config.Address, clientConfig); err != nil {
		return err
	}
	this_.Client = client

	channel, _, err := this_.Client.OpenChannel("session", nil)
	if err != nil {
		return nil
	}
	this_.channel = channel

	//ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&inRequests))
	//if !ok || err != nil {
	//	return
	//}
	ok, err := channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return
	}
	//这里第一个协程获取用户的输入
	go func() {
		for {
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil {
				return
			}
			_, err = this_.channel.Write(p)
			if err != nil {
				return
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		br := bufio.NewReader(this_.channel)
		var buf []byte
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			defer this_.Client.Close()
			defer this_.Session.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					log.Println(err)
					ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					ws.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		// 主循环
		for {
			select {
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
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
	err = SSHClient.GenerateClient(ws)

	return
}
