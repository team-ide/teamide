package ssh

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"teamide/pkg/terminal"
	"time"
)

func NewTerminalService(config *Config, lastUser string, lastDir string) (res *terminalService) {
	res = &terminalService{
		config:   config,
		lastUser: lastUser,
		lastDir:  lastDir,
	}
	return
}

type terminalService struct {
	config       *Config
	sshClient    *ssh.Client
	sshSession   *ssh.Session
	stdout       io.Reader
	stdin        io.Writer
	onClose      func()
	readeLock    sync.Mutex
	readeErrLock sync.Mutex
	writeLock    sync.Mutex
	isStopped    bool
	lastActive   time.Time
	lastUser     string
	lastDir      string
}

func (this_ *terminalService) IsWindows() (isWindows bool, err error) {
	isWindows = false
	return
}

func (this_ *terminalService) Stop() {
	this_.isStopped = true
	if this_.sshSession != nil {
		_ = this_.sshSession.Close()
	}
	if this_.sshClient != nil {
		_ = this_.sshClient.Close()
	}
	if this_.stdout != nil {
		if readerCloser, ok := this_.stdout.(io.ReadCloser); ok {
			_ = readerCloser.Close()
		}
	}
	if this_.stdin != nil {
		if writeCloser, ok := this_.stdin.(io.WriteCloser); ok {
			_ = writeCloser.Close()
		}
	}
}

func (this_ *terminalService) ChangeSize(size *terminal.Size) (err error) {

	if this_.sshSession == nil {
		return
	}
	if size.Cols > 0 && size.Rows > 0 {
		err = this_.sshSession.WindowChange(size.Rows, size.Cols)
		if err != nil {
			util.Logger.Error("SSH Session Window Change error", zap.Error(err))
			return
		}
	}
	return
}

func (this_ *terminalService) TestClient() (err error) {

	sshClient, err := NewClient(*this_.config)
	if err != nil {
		util.Logger.Error("SSH NewClient error", zap.Error(err))
		return
	}
	defer func() {
		_ = sshClient.Close()
	}()
	return
}
func (this_ *terminalService) Start(size *terminal.Size) (err error) {

	this_.sshClient, err = NewClient(*this_.config)
	if err != nil {
		util.Logger.Error("SSH NewClient error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH NewClient success", zap.Any("address", this_.config.Address))
	go func() {
		err = this_.sshClient.Wait()
		this_.Stop()
		util.Logger.Info("SSH Client end", zap.Any("address", this_.config.Address))
	}()

	this_.sshSession, err = this_.sshClient.NewSession()
	if err != nil {
		util.Logger.Error("SSH NewSession Error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH NewSession success", zap.Any("address", this_.config.Address))

	err = NewSSHShell(size, this_.sshSession)
	if err != nil {
		util.Logger.Error("Create SSH Shell Error", zap.Error(err))
		return
	}
	util.Logger.Info("SSH Format Shell Session success", zap.Any("address", this_.config.Address))
	this_.stdout, err = this_.sshSession.StdoutPipe()
	if err != nil {
		util.Logger.Error("ssh session StdoutPipe error", zap.Error(err))
		return
	}
	this_.stdin, err = this_.sshSession.StdinPipe()
	if err != nil {
		util.Logger.Error("ssh session StdinPipe error", zap.Error(err))
		return
	}
	this_.lastActive = time.Now()
	//if this_.lastUser != "" {
	//	if this_.lastUser != this_.config.Username {
	//		_, _ = this_.Write([]byte("sudo -i\n"))
	//		_, _ = this_.Write([]byte("su " + this_.lastUser + "\n"))
	//	}
	//}
	//if this_.lastDir != "" {
	//	_, _ = this_.Write([]byte("cd " + this_.lastDir + "\n"))
	//}
	// 开启空闲自动发送
	if this_.config.IdleSendOpen {
		idleSendTime := int64(this_.config.IdleSendTime)
		idleSendChar := this_.config.IdleSendChar
		util.Logger.Info("SSH Session open idle send", zap.Any("idleSendTime", idleSendTime), zap.Any("idleSendChar", idleSendChar))
		if idleSendTime > 0 && idleSendChar != "" {
			go func() {
				for !this_.isStopped {
					time.Sleep(time.Second)
					if this_.isStopped {
						break
					}
					nowTime := time.Now()
					idleTime := nowTime.Unix() - this_.lastActive.Unix()
					if idleTime < idleSendTime {
						continue
					}
					this_.lastActive = nowTime
					this_.idleSend(idleSendChar)

				}
			}()
		}
	}
	return
}

var (
	idleSendCharMap = map[string]byte{
		`^@`: '\x00',
		`\0`: '\x00',
		`^A`: '\x01',
		`^B`: '\x02',
		`^C`: '\x03',
		`^D`: '\x04',
		`\n`: '\x0a',
	}
)

func (this_ *terminalService) idleSend(idleSendChar string) {

	sendBytes := []byte(idleSendChar)
	f, ok := idleSendCharMap[idleSendChar]
	if ok {
		sendBytes = []byte{f}
	}

	util.Logger.Info("idleSend", zap.Any("idleSendChar", idleSendChar), zap.ByteString("sendBytes", sendBytes))
	_, err := this_.Write(sendBytes)
	if err != nil {
		util.Logger.Error("idleSend error", zap.Any("idleSendChar", idleSendChar), zap.ByteString("sendBytes", sendBytes), zap.Error(err))
	}
	return
}
func (this_ *terminalService) Write(buf []byte) (n int, err error) {
	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	n = len(buf)
	err = util.Write(this_.stdin, buf, nil)
	return
}

func (this_ *terminalService) Read(buf []byte) (n int, err error) {
	this_.readeLock.Lock()
	defer this_.readeLock.Unlock()

	this_.lastActive = time.Now()
	n, err = this_.stdout.Read(buf)
	return
}
