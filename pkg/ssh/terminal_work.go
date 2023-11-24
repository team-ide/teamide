package ssh

import (
	"bytes"
	"errors"
	"github.com/pkg/sftp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
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
	sftpClient   *sftp.Client
}

func (this_ *terminalService) IsWindows() (isWindows bool, err error) {
	isWindows = false
	return
}

func (this_ *terminalService) Stop() {
	this_.isStopped = true
	if this_.sftpClient != nil {
		_ = this_.sftpClient.Close()
	}
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
func (this_ *terminalService) readSSHFile(filepath string) (text string, err error) {
	if this_.sftpClient == nil {
		if this_.sshClient == nil {
			err = errors.New("ssh client is null")
			return
		}
		this_.sftpClient, err = sftp.NewClient(this_.sshClient)
		if err != nil {
			return
		}
	}
	f, err := this_.sftpClient.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	defer func() { _ = f.Close() }()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, f)
	if err != nil {
		return
	}
	text = buf.String()
	return
}

func (this_ *terminalService) GetCpuStats() (res []cpu.TimesStat, err error) {
	statText, err := this_.readSSHFile("/proc/stat")
	if err != nil {
		return
	}
	res = ParseProcStat(statText)
	return
}

func (this_ *terminalService) GetCpuPercent() (res []float64, err error) {
	// Get CPU usage at the start of the interval.
	cpuTimes1, err := this_.GetCpuStats()
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	cpuTimes2, err := this_.GetCpuStats()
	if err != nil {
		return
	}
	return calculateAllBusy(cpuTimes1, cpuTimes2)
}

func (this_ *terminalService) GetCpuInfo() (res []cpu.InfoStat, err error) {
	cpuInfoText, err := this_.readSSHFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	res, err = ParseProcCpuInfo(cpuInfoText, func(filepath string) []string {
		text, _ := this_.readSSHFile(filepath)
		return strings.Split(text, "\n")
	})
	if err != nil {
		return
	}
	return
}

func (this_ *terminalService) GetMemInfo() (res *mem.VirtualMemoryStat, err error) {

	memInfoText, err := this_.readSSHFile("/proc/meminfo")
	if err != nil {
		return
	}
	zoneInfoText, err := this_.readSSHFile("/proc/zoneinfo")
	if err != nil {
		return
	}
	res, _, err = ParseProcMemInfo(memInfoText, zoneInfoText)
	if err != nil {
		return
	}
	return
}

func (this_ *terminalService) GetDiskStats() (res map[string]disk.IOCountersStat, err error) {
	diskStatsText, err := this_.readSSHFile("/proc/diskstats")
	if err != nil {
		return
	}
	res, err = ParseProcDiskStats(diskStatsText)
	if err != nil {
		return
	}
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
