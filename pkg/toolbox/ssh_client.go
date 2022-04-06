package toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"mime/multipart"
	"strings"
	"sync"
	"time"
)

var (
	SSHSftpCache = map[string]*SSHSftpClient{}
)

type SSHClient struct {
	Token          string
	Config         *SSHConfig
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

func (this_ *SSHClient) createClient() (err error) {

	if this_.Token == "" || this_.Config == nil {
		err = errors.New("令牌会话丢失")
		this_.WSWriteError("令牌会话丢失")
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
	if this_.sshClient, err = ssh.Dial(this_.Config.Type, this_.Config.Address, clientConfig); err != nil {
		this_.WSWriteError("连接失败:" + err.Error())
		return
	}
	return
}

func (this_ *SSHClient) ListenWS(onEvent func(event string), onMessage func(bs []byte), onClose func()) {
	defer func() {
		if x := recover(); x != nil {
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
		_, bs, err := this_.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				this_.CloseWS()
				return
			}
			continue
		}
		if len(bs) > TeamIDEEventByteLength {
			msg := string(bs[0:TeamIDEEventByteLength])
			if strings.EqualFold(msg, TeamIDEEvent) {
				onEvent(string(bs[TeamIDEEventByteLength:]))
				continue
			}
		}
		onMessage(bs)
	}
}

const TeamIDEEvent = "TeamIDE:event:"

var (
	TeamIDEEventByteLength = len([]byte(TeamIDEEvent))
)

func (this_ *SSHClient) WSWrite(bs []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			this_.CloseWS()
			return
		}
	}()

	this_.wsWriteLock.Lock()
	defer this_.wsWriteLock.Unlock()
	err = this_.ws.WriteMessage(websocket.TextMessage, bs)

	if err != nil {
		if websocket.IsCloseError(err) {
			this_.CloseClient()
			this_.CloseWS()
		}
	}
	return
}

func (this_ *SSHClient) WSWriteData(obj interface{}) (err error) {

	res := map[string]interface{}{}
	res["value"] = obj
	res["code"] = 0
	bs, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = this_.WSWrite(bs)
	return
}

func (this_ *SSHClient) WSWriteError(message string) {
	res := map[string]interface{}{}
	res["msg"] = message
	res["code"] = -1
	bs, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = this_.WSWrite(bs)
	return
}

func (this_ *SSHClient) WSWriteMessage(message string) (err error) {
	res := map[string]interface{}{}
	res["msg"] = message
	res["code"] = 0
	bs, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = this_.WSWrite(bs)
	return
}

func (this_ *SSHClient) WSWriteEvent(event string) {
	res := map[string]interface{}{}
	res["event"] = event
	bs, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = this_.WSWrite(bs)
	return
}
