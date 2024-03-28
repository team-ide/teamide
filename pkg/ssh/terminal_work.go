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
	"regexp"
	"strings"
	"sync"
	"teamide/pkg/system"
	"teamide/pkg/terminal"
	"time"
)

func NewTerminalService(config *Config, lastUser string, lastDir string) (res *terminalService) {
	res = &terminalService{
		config:     config,
		lastUser:   lastUser,
		lastDir:    lastDir,
		sshClient2: config.SSHClient,
	}
	return
}

type terminalService struct {
	config       *Config
	sshClient2   *ssh.Client
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
		this_.sftpClient = nil
	}
	if this_.sshSession != nil {
		_ = this_.sshSession.Close()
		this_.sshSession = nil
	}
	if this_.sshClient != nil {
		_ = this_.sshClient.Close()
		this_.sshClient = nil
	}
	if this_.sshClient2 != nil {
		_ = this_.sshClient2.Close()
		this_.sshClient2 = nil
	}
	if this_.stdout != nil {
		if readerCloser, ok := this_.stdout.(io.ReadCloser); ok {
			_ = readerCloser.Close()
		}
		this_.stdout = nil
	}
	if this_.stdin != nil {
		if writeCloser, ok := this_.stdin.(io.WriteCloser); ok {
			_ = writeCloser.Close()
		}
		this_.stdin = nil
	}
}

func (this_ *terminalService) ChangeSize(size *terminal.Size) (err error) {

	if this_.sshSession == nil || this_.isStopped {
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
		if this_.config.SSHClient != nil {
			_ = this_.config.SSHClient.Close()
		}
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

func (this_ *terminalService) runCmd(cmd string) (text string, err error) {
	if this_.sshClient == nil {
		err = errors.New("ssh client is null")
		return
	}
	s, err := this_.sshClient.NewSession()
	if err != nil {
		return
	}
	defer func() { _ = s.Close() }()
	bs, err := s.Output(cmd)
	if err != nil {
		return
	}
	text = string(bs)
	return
}

func (this_ *terminalService) fileExist(filepath string) (exist bool) {
	var err error
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
	f, err := this_.sftpClient.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	if f != nil {
		exist = true
	}
	return
}
func (this_ *terminalService) Start(size *terminal.Size) (err error) {

	this_.sshClient, err = NewClient(*this_.config)
	if err != nil {
		util.Logger.Error("SSH NewClient error", zap.Error(err))
		this_.Stop()
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

func (this_ *terminalService) SystemInfo() (res *system.Info, err error) {
	return this_.Info()
}

func (this_ *terminalService) SystemMonitorData() (res *system.MonitorData, err error) {
	return system.GetCacheOrNew()
}

func (this_ *terminalService) Info() (res *system.Info, err error) {
	res = &system.Info{}
	hostname, err := this_.readSSHFile("/etc/hostname")
	if err != nil {
		return
	}
	hostname = strings.ReplaceAll(hostname, "\n", "")

	version, err := this_.readSSHFile("/proc/version")
	if err != nil {
		return
	}
	version = strings.ReplaceAll(version, "\n", "")
	res.HostInfoStat = &system.HostInfoStat{
		Hostname:   hostname,
		KernelArch: version,
	}
	res.HostInfoStat.Platform, res.HostInfoStat.PlatformFamily, res.HostInfoStat.PlatformVersion, err = this_.PlatformInformation()

	cpuInfos, err := this_.GetCpuInfo()
	if err != nil {
		return
	}
	for _, one := range cpuInfos {
		res.CpuInfoStats = append(res.CpuInfoStats, &system.CpuInfoStat{
			CPU:        one.CPU,
			VendorID:   one.VendorID,
			Family:     one.Family,
			Model:      one.Model,
			Stepping:   one.Stepping,
			PhysicalID: one.PhysicalID,
			CoreID:     one.CoreID,
			Cores:      one.Cores,
			ModelName:  one.ModelName,
			Mhz:        one.Mhz,
			CacheSize:  one.CacheSize,
			//Flags:      one.Flags,
			Microcode: one.Microcode,
		})
	}

	text, _ := this_.runCmd("free -b")
	lines := strings.Split(text, "\n")
	r, _ := regexp.Compile("\\s+")
	if len(lines) > 1 {
		line := lines[1]
		line = string(r.ReplaceAll([]byte(line), []byte(" ")))
		vs := strings.Split(line, " ")
		//fmt.Println(line)
		if len(vs) == 7 {
			res.Memory = &system.VirtualMemoryStat{
				Total:     util.StringToUint64(vs[1]),
				Available: util.StringToUint64(vs[6]),
				Used:      util.StringToUint64(vs[2]),
				Free:      util.StringToUint64(vs[3]),
				Cached:    util.StringToUint64(vs[5]),
				Shared:    util.StringToUint64(vs[4]),
			}
		}
	}

	text, _ = this_.runCmd("df -B1")
	lines = strings.Split(text, "\n")
	if len(lines) > 1 {
		for _, line := range lines[1:] {
			line = string(r.ReplaceAll([]byte(line), []byte(" ")))
			//fmt.Println(line)
			vs := strings.Split(line, " ")
			if len(vs) == 6 {
				d := &system.DiskUsageStat{
					Path:  vs[5],
					Total: util.StringToUint64(vs[3]),
					Used:  util.StringToUint64(vs[2]),
				}
				d.Free = d.Total - d.Used
				d.UsedPercent = float64(d.Used) / float64(d.Total) * 100
				res.Disks = append(res.Disks, d)
			}
		}
	}
	return
}

func (this_ *terminalService) PlatformInformation() (platform string, family string, version string, err error) {
	lsb, err := this_.getlsbStruct()
	if err != nil {
		lsb = &lsbStruct{}
	}

	var text string

	if this_.fileExist("/etc/oracle-release") {
		platform = "oracle"
		text, err = this_.readSSHFile("/etc/oracle-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/enterprise-release") {
		platform = "oracle"
		text, err = this_.readSSHFile("/etc/enterprise-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/slackware-version") {
		platform = "slackware"
		text, err = this_.readSSHFile("/etc/slackware-version")
		if err == nil {
			version = getSlackwareVersion(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/debian_version") {
		if lsb.ID == "Ubuntu" {
			platform = "ubuntu"
			version = lsb.Release
		} else if lsb.ID == "LinuxMint" {
			platform = "linuxmint"
			version = lsb.Release
		} else if lsb.ID == "Kylin" {
			platform = "Kylin"
			version = lsb.Release
		} else if lsb.ID == `"Cumulus Linux"` {
			platform = "cumuluslinux"
			version = lsb.Release
		} else {
			if this_.fileExist("/usr/bin/raspi-config") {
				platform = "raspbian"
			} else {
				platform = "debian"
			}
			text, err = this_.readSSHFile("/etc/debian_version")
			contents := strings.Split(text, "\n")
			if err == nil && len(contents) > 0 && contents[0] != "" {
				version = contents[0]
			}
		}
	} else if this_.fileExist("/etc/neokylin-release") {
		text, err = this_.readSSHFile("/etc/neokylin-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
			platform = getRedhatishPlatform(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/redhat-release") {
		text, err = this_.readSSHFile("/etc/redhat-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
			platform = getRedhatishPlatform(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/system-release") {
		text, err = this_.readSSHFile("/etc/system-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
			platform = getRedhatishPlatform(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/gentoo-release") {
		platform = "gentoo"
		text, err = this_.readSSHFile("/etc/gentoo-release")
		if err == nil {
			version = getRedhatishVersion(strings.Split(text, "\n"))
		}
	} else if this_.fileExist("/etc/SuSE-release") {
		text, err = this_.readSSHFile("/etc/SuSE-release")
		if err == nil {
			version = getSuseVersion(strings.Split(text, "\n"))
			platform = getSusePlatform(strings.Split(text, "\n"))
		}
		// TODO: slackware detecion
	} else if this_.fileExist("/etc/arch-release") {
		platform = "arch"
		version = lsb.Release
	} else if this_.fileExist("/etc/alpine-release") {
		platform = "alpine"
		text, err = this_.readSSHFile("/etc/alpine-release")
		contents := strings.Split(text, "\n")
		if err == nil && len(contents) > 0 && contents[0] != "" {
			version = contents[0]
		}
	} else if this_.fileExist("/etc/os-release") {
		p, v, err := this_.GetOSRelease()
		if err == nil {
			platform = p
			version = v
		}
	} else if lsb.ID == "RedHat" {
		platform = "redhat"
		version = lsb.Release
	} else if lsb.ID == "Amazon" {
		platform = "amazon"
		version = lsb.Release
	} else if lsb.ID == "ScientificSL" {
		platform = "scientific"
		version = lsb.Release
	} else if lsb.ID == "XenServer" {
		platform = "xenserver"
		version = lsb.Release
	} else if lsb.ID != "" {
		platform = strings.ToLower(lsb.ID)
		version = lsb.Release
	}

	platform = strings.Trim(platform, `"`)

	switch platform {
	case "debian", "ubuntu", "linuxmint", "raspbian", "Kylin", "cumuluslinux":
		family = "debian"
	case "fedora":
		family = "fedora"
	case "oracle", "centos", "redhat", "scientific", "enterpriseenterprise", "amazon", "xenserver", "cloudlinux", "ibm_powerkvm", "rocky", "almalinux":
		family = "rhel"
	case "suse", "opensuse", "opensuse-leap", "opensuse-tumbleweed", "opensuse-tumbleweed-kubic", "sles", "sled", "caasp":
		family = "suse"
	case "gentoo":
		family = "gentoo"
	case "slackware":
		family = "slackware"
	case "arch":
		family = "arch"
	case "exherbo":
		family = "exherbo"
	case "alpine":
		family = "alpine"
	case "coreos":
		family = "coreos"
	case "solus":
		family = "solus"
	case "neokylin":
		family = "neokylin"
	}

	return platform, family, version, nil
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
