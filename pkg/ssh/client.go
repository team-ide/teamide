package ssh

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"time"
)

var (
	ShellCache = map[string]*ShellClient{}
)

type Config struct {
	Type         string `json:"type"`
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PublicKey    string `json:"publicKey"`
	Command      string `json:"command"`
	Timeout      int    `json:"timeout"`
	IdleSendOpen bool   `json:"idleSendOpen"`
	IdleSendTime int    `json:"idleSendTime"`
	IdleSendChar string `json:"idleSendChar"`
}

type Client struct {
	Token              string
	Config             Config
	sshClient          *ssh.Client
	ws                 *websocket.Conn
	isClosedClient     bool
	isClosedWS         bool
	wsWriteLock        sync.Mutex
	writeWSMessageList chan *writeWSMessage
}

type writeWSMessage struct {
	Type int
	Data []byte
}

func (this_ *Client) CloseClient() {
	this_.isClosedClient = true
	if this_.sshClient != nil {
		err := this_.sshClient.Close()
		if err != nil {
			util.Logger.Error("SSH Client close error", zap.Error(err))
		}
	}
	this_.sshClient = nil
}

func (this_ *Client) CloseWS() {
	delete(ShellCache, this_.Token)
	this_.isClosedWS = true
	if this_.ws != nil {
		err := this_.ws.Close()
		if err != nil {
			util.Logger.Error("WebSocket close error", zap.Error(err))
		}
	}
	this_.ws = nil
}

func (this_ *Client) WSWriteText(bs []byte) {
	this_.WSWriteByType(websocket.TextMessage, bs)
	return
}

func (this_ *Client) WSWriteBinary(bs []byte) {
	this_.WSWriteByType(websocket.BinaryMessage, bs)
	return
}

func (this_ *Client) WSWriteByType(messageType int, bs []byte) {
	defer func() {
		if x := recover(); x != nil {
			util.Logger.Error("WebSocket信息写入异常", zap.Any("err", x))
			this_.CloseWS()
			return
		}
	}()
	if this_.isClosedWS {
		return
	}

	this_.wsWriteLock.Lock()
	defer this_.wsWriteLock.Unlock()

	if this_.writeWSMessageList == nil {
		this_.writeWSMessageList = make(chan *writeWSMessage, 1)
		go func() {
			for {
				select {
				case msg := <-this_.writeWSMessageList:

					if this_.isClosedWS {
						close(this_.writeWSMessageList)
						return
					}
					//fmt.Println("write message:", string(msg.Data))
					err := this_.ws.WriteMessage(msg.Type, msg.Data)

					if err != nil {
						this_.CloseWS()
						return
					}
				}
			}
		}()
	}

	this_.writeWSMessageList <- &writeWSMessage{
		Type: messageType,
		Data: bs,
	}
}

func (this_ *Client) WSWriteData(obj interface{}) {

	bs, err := json.Marshal(obj)
	if err != nil {
		util.Logger.Error("WSWriteData转换JSON异常", zap.Error(err))
		return
	}
	this_.WSWriteText(bs)
	return
}

func NewClient(config Config) (client *ssh.Client, err error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshConfig    ssh.Config
	)
	auth = []ssh.AuthMethod{}

	if config.PublicKey != "" {
		var publicKeyBytes []byte
		publicKeyBytes, err = os.ReadFile(config.PublicKey)
		if err != nil {
			return
		}
		var publicKeySigner ssh.Signer
		if config.Password != "" {
			publicKeySigner, err = ssh.ParsePrivateKeyWithPassphrase(publicKeyBytes, []byte(config.Password))
		} else {
			publicKeySigner, err = ssh.ParsePrivateKey(publicKeyBytes)
		}
		if err != nil {
			return
		}
		auth = append(auth, ssh.PublicKeys(publicKeySigner))

	} else if config.Password != "" {
		auth = append(auth, ssh.Password(config.Password))
	}

	sshConfig = ssh.Config{
		Ciphers: Ciphers,
	}
	var timeout = 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}
	clientConfig = &ssh.ClientConfig{
		User:            config.Username,
		Auth:            auth,
		Timeout:         timeout,
		Config:          sshConfig,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
	}
	if config.Type == "" {
		config.Type = "tcp"
	}
	client, err = ssh.Dial(config.Type, config.Address, clientConfig)
	if err != nil {
		return
	}
	return
}

var (
	Ciphers = []string{
		"aes128-ctr",
		"aes192-ctr",
		"aes256-ctr",
		"aes128-gcm@openssh.com",
		"aes256-gcm@openssh.com",
		"chacha20-poly1305@openssh.com",
		"arcfour256",
		"arcfour128",
		"arcfour",
		"aes128-cbc",
		"3des-cbc",
		"aes192-cbc",
		"aes256-cbc",
	}
)
