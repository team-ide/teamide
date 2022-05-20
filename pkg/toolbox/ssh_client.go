package toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/wxnacy/wgo/file"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"mime/multipart"
	"strings"
	"sync"
	"time"
)

var (
	SSHSftpCache  = map[string]*SSHSftpClient{}
	SSHShellCache = map[string]*SSHShellClient{}
)

type SSHClient struct {
	Token          string
	Config         SSHConfig
	sshClient      *ssh.Client
	ws             *websocket.Conn
	isClosedClient bool
	isClosedWS     bool
	wsWriteLock    sync.RWMutex
	Logger         *zap.Logger
}

type ConfirmInfo struct {
	ConfirmId   string `json:"confirmId,omitempty"`
	IsConfirm   bool   `json:"isConfirm,omitempty"`
	Confirm     string `json:"confirm,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	IsFileExist bool   `json:"isFileExist,omitempty"`
	IsOk        bool   `json:"isOk,omitempty"`
	IsCancel    bool   `json:"isCancel,omitempty"`
}

type UploadFile struct {
	Dir      string
	Place    string
	WorkId   string
	FullPath string
	File     *multipart.FileHeader
}

func (this_ *SSHClient) CloseClient() {
	this_.isClosedClient = true
	if this_.sshClient != nil {
		err := this_.sshClient.Close()
		if err != nil {
			this_.Logger.Error("SSH Client close error", zap.Error(err))
		}
	}
	this_.sshClient = nil
}

func (this_ *SSHClient) CloseWS() {
	delete(SSHSftpCache, this_.Token)
	delete(SSHShellCache, this_.Token)
	this_.isClosedWS = true
	if this_.ws != nil {
		err := this_.ws.Close()
		if err != nil {
			this_.Logger.Error("WebSocket close error", zap.Error(err))
		}
	}
	this_.ws = nil
}

func (this_ *SSHClient) initClient() (err error) {

	if this_.isClosedWS || this_.sshClient == nil {
		err = this_.createClient()
	}
	return
}

func NewSSHClient(config SSHConfig) (client *ssh.Client, err error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshConfig    ssh.Config
	)
	auth = []ssh.AuthMethod{}

	if config.PublicKey != "" {
		var publicKeyContent string
		publicKeyContent, err = file.ReadFile(config.PublicKey)
		if err != nil {
			return
		}
		publicKeyBytes := []byte(publicKeyContent)
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

func (this_ *SSHClient) createClient() (err error) {

	if this_.Token == "" || this_.Config.Address == "" {
		err = errors.New("令牌会话丢失")
		this_.Logger.Error("令牌验证失败", zap.Error(err))
		this_.WSWriteError(err.Error())
		return
	}
	if this_.sshClient, err = NewSSHClient(this_.Config); err != nil {
		this_.Logger.Error("createClient error", zap.Error(err))
		this_.WSWriteError("连接失败:" + err.Error())
		return
	}
	return
}

func WSIsCloseError(err error) bool {
	_, ok := err.(*websocket.CloseError)
	return ok
}

func (this_ *SSHClient) ListenWS(onEvent func(event string), onMessage func(bs []byte), onClose func()) {
	defer func() {
		if x := recover(); x != nil {
			this_.Logger.Error("WebSocket信息监听异常", zap.Any("err", x))
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
			if WSIsCloseError(err) {
				this_.CloseWS()
				return
			}
			this_.Logger.Error("WebSocket信息读取异常", zap.Error(err))
			continue
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
	TeamIDEEvent   = "TeamIDE:event:"
	TeamIDEMessage = "TeamIDE:message:"
	TeamIDEError   = "TeamIDE:error:"
)

var (
	TeamIDEEventByteLength = len([]byte(TeamIDEEvent))
)

func (this_ *SSHClient) WSWrite(bs []byte) {
	defer func() {
		if x := recover(); x != nil {
			this_.Logger.Error("WebSocket信息写入异常", zap.Any("err", x))
			this_.CloseWS()
			return
		}
	}()
	if this_.isClosedWS {
		return
	}

	this_.wsWriteLock.Lock()
	defer this_.wsWriteLock.Unlock()
	err := this_.ws.WriteMessage(websocket.TextMessage, bs)

	if err != nil {
		if WSIsCloseError(err) {
			this_.CloseWS()
			return
		}
		this_.Logger.Error("WebSocket信息写入异常", zap.Error(err))
	}
	return
}

func (this_ *SSHClient) WSWriteBinary(bs []byte) {
	defer func() {
		if x := recover(); x != nil {
			this_.Logger.Error("WebSocket信息写入异常", zap.Any("err", x))
			this_.CloseWS()
			return
		}
	}()
	if this_.isClosedWS {
		return
	}

	this_.wsWriteLock.Lock()
	defer this_.wsWriteLock.Unlock()
	err := this_.ws.WriteMessage(websocket.BinaryMessage, bs)

	if err != nil {
		if WSIsCloseError(err) {
			this_.CloseWS()
			return
		}
		this_.Logger.Error("WebSocket信息写入异常", zap.Error(err))
	}
	return
}

func (this_ *SSHClient) WSWriteData(obj interface{}) {

	bs, err := json.Marshal(obj)
	if err != nil {
		this_.Logger.Error("WSWriteData转换JSON异常", zap.Error(err))
		return
	}
	this_.WSWrite(bs)
	return
}

func (this_ *SSHClient) WSWriteError(message string) {
	this_.WSWrite([]byte(TeamIDEError + message))
	return
}

func (this_ *SSHClient) WSWriteMessage(message string) {
	this_.WSWrite([]byte(TeamIDEMessage + message))
	return
}

func (this_ *SSHClient) WSWriteEvent(event string) {
	this_.WSWrite([]byte(TeamIDEEvent + event))
	return
}
