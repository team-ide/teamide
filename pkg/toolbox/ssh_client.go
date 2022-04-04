package toolbox

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"mime/multipart"
	"sync"
	"time"
)

var (
	SSHClientCache = map[string]*SSHClient{}
)

type SSHClient struct {
	Token       string
	Config      *SSHConfig
	sshClient   *ssh.Client
	sshChannel  ssh.Channel
	sshSession  *ssh.Session
	ws          *websocket.Conn
	sftpClient  *sftp.Client
	isClosed    bool
	UploadFile  chan *UploadFile
	wsWriteLock sync.RWMutex
	confirmMap  map[string]chan *ConfirmInfo
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

func (this_ *SSHClient) Close() {
	delete(SSHClientCache, this_.Token)
	if this_.isClosed {
		return
	}
	this_.isClosed = true
	var err error
	if this_.sshChannel != nil {
		err = this_.sshChannel.Close()
		if err != nil {
			fmt.Println("sshChannel close error", err)
		}
	}
	if this_.sshSession != nil {
		err = this_.sshSession.Close()
		if err != nil {
			fmt.Println("sshSession close error", err)
		}
	}
	if this_.sshClient != nil {
		err = this_.sshClient.Close()
		if err != nil {
			fmt.Println("sshClient close error", err)
		}
	}
	if this_.sftpClient != nil {
		err = this_.sftpClient.Close()
		if err != nil {
			fmt.Println("sftpClient close error", err)
		}
	}
	if this_.ws != nil {
		err = this_.ws.Close()
		if err != nil {
			fmt.Println("ws close error", err)
		}
	}
}

func (this_ *SSHClient) initClient() (err error) {

	if !this_.isClosed && this_.sshClient != nil {
		return
	}
	if this_.Token == "" || this_.Config == nil {
		err = errors.New("令牌会话丢失")
		return
	}
	SSHClientCache[this_.Token] = this_
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
		fmt.Println("SSH Client Dial error", err)
		this_.Close()
		return
	}
	return
}
