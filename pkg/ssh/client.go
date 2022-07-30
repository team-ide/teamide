package ssh

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"sync"
	"teamide/pkg/util"
	"time"
)

var (
	TokenCache = map[string]*Config{}
)

func WSSSHConnection(token string, ws *websocket.Conn) (err error) {
	var sshConfig = TokenCache[token]
	client := Client{
		Token:  token,
		Config: *sshConfig,
		ws:     ws,
	}
	shellClient := &ShellClient{
		Client: client,
	}
	shellClient.start()

	return
}

func WSSFPTConnection(token string, ws *websocket.Conn) (err error) {
	var sshConfig = TokenCache[token]
	client := Client{
		Token:  token,
		Config: *sshConfig,
		ws:     ws,
	}
	sftpClient := &SftpClient{
		Client: client,
	}
	sftpClient.start()

	return
}

var (
	SftpCache  = map[string]*SftpClient{}
	ShellCache = map[string]*ShellClient{}
)

type Config struct {
	Type      string `json:"type"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"publicKey"`
}

type Client struct {
	Token              string
	Config             Config
	sshClient          *ssh.Client
	ws                 *websocket.Conn
	isClosedClient     bool
	isClosedWS         bool
	wsWriteLock        sync.Locker
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
	delete(SftpCache, this_.Token)
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

func (this_ *Client) initClient() (err error) {

	if this_.isClosedWS || this_.sshClient == nil {
		err = this_.createClient()
	}
	return
}

func (this_ *Client) createClient() (err error) {

	if this_.Token == "" || this_.Config.Address == "" {
		err = errors.New("令牌会话丢失")
		util.Logger.Error("令牌验证失败", zap.Error(err))
		this_.WSWriteError(err.Error())
		return
	}
	if this_.sshClient, err = NewClient(this_.Config); err != nil {
		util.Logger.Error("createClient error", zap.Error(err))
		this_.WSWriteError("连接失败:" + err.Error())
		return
	}
	go func() {
		err = this_.sshClient.Wait()
		this_.CloseClient()
	}()
	return
}

func (this_ *Client) ListenWS(onEvent func(event string), onMessage func(bs []byte), onClose func()) {
	defer func() {
		if x := recover(); x != nil {
			util.Logger.Error("WebSocket信息监听异常", zap.Any("err", x))
			this_.CloseWS()
			return
		}
	}()
	defer onClose()
	// 第一个协程获取用户的输入
	for {
		if this_.isClosedWS {
			return
		}
		messageType, bs, err := this_.ws.ReadMessage()
		if err != nil {
			this_.CloseWS()
			return
		}
		if messageType == websocket.TextMessage {
			if len(bs) > TeamIDEEventByteLength {
				msg := string(bs[0:TeamIDEEventByteLength])
				if strings.EqualFold(msg, TeamIDEEvent) {
					onEvent(string(bs[TeamIDEEventByteLength:]))
					continue
				}
			}
		}
		onMessage(bs)
	}
}

const (
	TeamIDEEvent       = "^^^^--Team--IDE--^^^^:event:"
	TeamIDEMessage     = "^^^^--Team--IDE--^^^^:TeamIDE:message:"
	TeamIDEError       = "^^^^--Team--IDE--^^^^:TeamIDE:error:"
	TeamIDEAlert       = "^^^^--Team--IDE--^^^^:TeamIDE:alert:"
	TeamIDEConsole     = "^^^^--Team--IDE--^^^^:TeamIDE:console:"
	TeamIDEStdout      = "^^^^--Team--IDE--^^^^:TeamIDE:stdout:"
	TeamIDEBinaryStart = "^^^^--Team--IDE--^^^^:TeamIDE:binary:"
)

var (
	TeamIDEBinaryStartBytes = []byte(TeamIDEBinaryStart)
)

var (
	TeamIDEEventByteLength = len([]byte(TeamIDEEvent))
)

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

func (this_ *Client) WSWriteError(message string) {
	this_.WSWriteText([]byte(TeamIDEError + message))
	return
}

func (this_ *Client) WSWriteMessage(message string) {
	this_.WSWriteText([]byte(TeamIDEMessage + message))
	return
}

func (this_ *Client) WSWriteEvent(event string) {
	this_.WSWriteText([]byte(TeamIDEEvent + event))
	return
}

func (this_ *Client) WSWriteAlert(alert string) {
	this_.WSWriteText([]byte(TeamIDEAlert + alert))
	return
}

func (this_ *Client) WSWriteConsole(console string) {
	this_.WSWriteText([]byte(TeamIDEConsole + console))
	return
}

func (this_ *Client) WSWriteStdout(stdout string) {
	this_.WSWriteText([]byte(TeamIDEStdout + stdout))
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
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:            config.Username,
		Auth:            auth,
		Timeout:         5 * time.Second,
		Config:          sshConfig,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
	}
	client, err = ssh.Dial(config.Type, config.Address, clientConfig)
	if err != nil {
		return
	}
	return
}
