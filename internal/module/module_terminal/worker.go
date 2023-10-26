package module_terminal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"teamide/internal/context"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/ssh"
	"teamide/pkg/terminal"
	"time"
)

func NewWorkerFactory(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *WorkerFactory {
	return &WorkerFactory{
		ServerContext:      toolboxService_.ServerContext,
		toolboxService:     toolboxService_,
		nodeService:        nodeService_,
		terminalLogService: NewTerminalLogService(toolboxService_.ServerContext),
		workerCache:        make(map[string]*Worker),
	}
}

type WorkerFactory struct {
	*context.ServerContext
	toolboxService     *module_toolbox.ToolboxService
	nodeService        *module_node.NodeService
	terminalLogService *TerminalLogService
	workerCache        map[string]*Worker
	workerCacheLock    sync.Mutex
}

func (this_ *WorkerFactory) GetService(key string) (res *Worker) {
	this_.workerCacheLock.Lock()
	defer this_.workerCacheLock.Unlock()

	res = this_.workerCache[key]
	return
}

type CreateParam struct {
	place    string
	placeId  string
	workerId string
	lastUser string
	lastDir  string
}

func (this_ *WorkerFactory) createService(param *CreateParam) (worker *Worker, command string, err error) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("createService error", zap.Any("error", e))
		}
	}()
	var service terminal.Service
	switch param.place {
	case "local":
		service = terminal.NewLocalService()
	case "ssh":
		if param.placeId == "" {
			err = errors.New("SSH配置不能为空")
			return
		}
		var id int64
		id, err = strconv.ParseInt(param.placeId, 10, 64)
		if err != nil {
			return
		}
		var tD *module_toolbox.ToolboxModel
		tD, err = this_.toolboxService.Get(id)
		if err != nil {
			return
		}
		if tD == nil || tD.Option == "" {
			err = errors.New("SSH[" + param.placeId + "]配置不存在")
			return
		}

		var config *ssh.Config
		config, err = this_.toolboxService.GetSSHConfig(tD.Option)
		if err != nil {
			return
		}
		if config != nil {
			command = config.Command
		}

		service = ssh.NewTerminalService(config, param.lastUser, param.lastDir)
	case "node":
		if param.placeId == "" {
			err = errors.New("node配置不能为空")
			return
		}
		service = module_node.NewTerminalService(param.placeId, this_.nodeService)
	}
	if service == nil {
		err = errors.New("[" + param.place + "]终端服务不存在")
		return
	}

	worker = &Worker{
		place:         param.place,
		placeId:       param.placeId,
		workerId:      param.workerId,
		service:       service,
		WorkerFactory: this_,
	}
	worker.init()
	return
}

func (this_ *WorkerFactory) Start(key string, param *CreateParam, size *terminal.Size, ws *websocket.Conn, baseLog *TerminalLogModel) (err error) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("Start error", zap.Any("error", e))
		}
	}()

	this_.workerCacheLock.Lock()
	defer this_.workerCacheLock.Unlock()

	worker := this_.workerCache[key]
	if worker != nil {
		err = errors.New("会话服务[" + key + "]已存在")
		return
	}
	var command string
	worker, command, err = this_.createService(param)
	if err != nil {
		return
	}
	// 执行配置的命令
	worker.ws = ws
	isWindow, err := worker.service.IsWindows()
	if err != nil {
		return
	}
	err = worker.service.Start(size)
	if err != nil {
		return
	}
	if command != "" {
		go func() {
			command = strings.ReplaceAll(command, "\n\r", "\n")
			this_.Logger.Info("SSH start run command", zap.Any("command", command))
			commands := strings.Split(command, "\n")
			for _, c := range commands {
				c = strings.TrimSpace(c)
				if c == "" {
					continue
				}
				this_.Logger.Info("SSH start run line command", zap.Any("line", c))
				if strings.HasPrefix(c, "sleep ") {
					t, _ := strconv.Atoi(strings.TrimSpace(c[len("sleep "):]))
					this_.Logger.Info("SSH start run line sleep", zap.Any("time", t))
					if t > 0 {
						time.Sleep(time.Second * time.Duration(t))

					}
					continue
				}
				_, err = worker.service.Write([]byte(c + "\n"))
				if err != nil {
					this_.Logger.Error("SSH start run line error", zap.Error(err))
				}
			}
		}()

	}

	go worker.startReadWS(isWindow, baseLog)
	go worker.startReadService(isWindow)

	this_.workerCache[key] = worker
	return
}

func (this_ *WorkerFactory) getParentDir(place string, placeId string) (dir string) {
	dir = this_.GetFilesDir()
	dir += fmt.Sprintf("toolbox-workers/toolbox-%s-%s/", place, placeId)
	return
}

type LogInfo struct {
	PlaceId  string `json:"placeId"`
	WorkerId string `json:"workerId"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	ModTime  int64  `json:"modTime,omitempty"`
}

func (this_ *WorkerFactory) getLogs(place string, placeId string) (logs []*LogInfo, err error) {
	parentDir := this_.getParentDir(place, placeId)

	ex, _ := util.PathExists(parentDir)
	if !ex {
		return
	}

	fileList, err := os.ReadDir(parentDir)
	if err != nil {
		return
	}
	var names []string
	for _, f := range fileList {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}

	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j]) //升序  即前面的值比后面的小 忽略大小写排序
	})
	size := len(names)
	for i := size - 1; i >= 0; i-- {
		log, _ := this_.getLog(place, placeId, names[i])
		if log == nil {
			continue
		}
		logs = append(logs, log)
	}

	return
}

func (this_ *WorkerFactory) getLog(place string, placeId string, workerId string) (log *LogInfo, err error) {
	path := this_.getLogPath(place, placeId, workerId)

	if ex, _ := util.PathExists(path); !ex {
		return
	}
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	log = &LogInfo{
		PlaceId:  placeId,
		WorkerId: workerId,
		Size:     stat.Size(),
		Path:     path,
		ModTime:  util.GetMilliByTime(stat.ModTime()),
	}
	return
}

func (this_ *WorkerFactory) getLogPath(place string, placeId string, workerId string) (path string) {
	parentDir := this_.getParentDir(place, placeId)

	path = parentDir + workerId + "/command.log"
	return
}

type Worker struct {
	key      string
	place    string
	placeId  string
	workerId string
	dir      string
	*WorkerFactory
	service        terminal.Service
	ws             *websocket.Conn
	commandLogFile *os.File
	isRz           bool
	isSz           bool
	isLastSzEnd    bool
}

func (this_ *Worker) init() {
	dir := this_.getParentDir(this_.place, this_.placeId)
	dir += this_.workerId + "/"

	ex, err := util.PathExists(dir)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(dir, fs.ModePerm)
	}
	if err != nil {
		return
	}
	this_.dir = dir
	return
}

type logContext struct {
	commandBytes []byte
}

func (this_ *Worker) onCommand(logContext *logContext, commandBytes []byte, log TerminalLogModel) {

	logContext.commandBytes = append(logContext.commandBytes, commandBytes...)

	str := string(commandBytes)
	if strings.Contains(str, "\n") ||
		strings.Contains(str, "\r") {
		command := string(logContext.commandBytes)
		log.Command = command
		_ = this_.terminalLogService.Insert(&log)
		logContext.commandBytes = []byte{}
	}

}

var (
	// sshSZStart sz fmt.Sprintf("%+q", "rz\r**\x18B00000000000000\r\x8a\x11")
	//sshSZStart = []byte{13, 42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}
	sshSZStart = []byte{42, 42, 24, 66, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 13, 138, 17}
	// sshSZEnd sz 结束 fmt.Sprintf("%+q", "\r**\x18B0800000000022d\r\x8a")
	//sshSZEnd = []byte{13, 42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}
	sshSZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}
	// sshSZEndOO sz 结束后可能还会发送两个 OO，但是经过测试发现不一定每次都会发送 fmt.Sprintf("%+q", "OO")
	sshSZEndOO = []byte{79, 79}

	// rz fmt.Sprintf("%+q", "**\x18B0100000023be50\r\x8a\x11")
	sshRZStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 50, 51, 98, 101, 53, 48, 13, 138, 17}
	// rz -e fmt.Sprintf("%+q", "**\x18B0100000063f694\r\x8a\x11")
	sshRZEStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 48, 54, 51, 102, 54, 57, 52, 13, 138, 17}
	// rz -S fmt.Sprintf("%+q", "**\x18B0100000223d832\r\x8a\x11")
	sshRZSStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 50, 51, 100, 56, 51, 50, 13, 138, 17}
	// rz -e -S fmt.Sprintf("%+q", "**\x18B010000026390f6\r\x8a\x11")
	sshRZESStart = []byte{42, 42, 24, 66, 48, 49, 48, 48, 48, 48, 48, 50, 54, 51, 57, 48, 102, 54, 13, 138, 17}
	// rz 结束 fmt.Sprintf("%+q", "**\x18B0800000000022d\r\x8a")
	sshRZEnd = []byte{42, 42, 24, 66, 48, 56, 48, 48, 48, 48, 48, 48, 48, 48, 48, 50, 50, 100, 13, 138}

	// **\x18B0
	sshRZCtrlStart = []byte{42, 42, 24, 66, 48}
	// \r\x8a\x11
	sshRZCtrlEnd1 = []byte{13, 138, 17}
	// \r\x8a
	sshRZCtrlEnd2 = []byte{13, 138}

	// zmodem 取消 \x18\x18\x18\x18\x18\x08\x08\x08\x08\x08
	sshCancel = []byte{24, 24, 24, 24, 24, 8, 8, 8, 8, 8}
)

func (this_ *Worker) onServiceRead(bs []byte) {
	if this_.dir == "" {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("onServiceRead error", zap.Any("err", e))
		}
	}()

	if !this_.isRz && (bytes.Contains(bs, sshRZStart) || bytes.Contains(bs, sshRZSStart) || bytes.Contains(bs, sshRZEStart) || bytes.Contains(bs, sshRZESStart)) {
		this_.isRz = true
		bs = []byte(fmt.Sprintf("\n开始上传文件:%s\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000")))

	} else if this_.isRz && bytes.Contains(bs, sshRZEnd) {
		this_.isRz = false
		bs = []byte(fmt.Sprintf("\n结束上传文件:%s\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000")))
	} else if !this_.isSz && bytes.Contains(bs, sshSZStart) {
		this_.isSz = true
		bs = []byte(fmt.Sprintf("\n开始下载文件:%s\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000")))
	} else if this_.isSz && bytes.Contains(bs, sshSZEnd) {
		this_.isSz = false
		this_.isLastSzEnd = true
		bs = []byte(fmt.Sprintf("\n结束下载文件:%s\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000")))
	} else {
		if this_.isSz || this_.isRz {
			return
		}
		if this_.isLastSzEnd && bytes.Equal(bs, sshSZEndOO) {
			return
		}
		this_.isLastSzEnd = false
	}

	if this_.commandLogFile == nil {
		ex, err := util.PathExists(this_.dir)
		if err != nil {
			return
		}
		if !ex {
			err = os.MkdirAll(this_.dir, fs.ModePerm)
		}
		path := this_.getLogPath(this_.place, this_.placeId, this_.workerId)
		if ex, _ = util.PathExists(path); ex {
			this_.commandLogFile, _ = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
		} else {
			this_.commandLogFile, _ = os.Create(path)
		}
	}
	if this_.commandLogFile == nil {
		return
	}
	str := string(bs)
	// 配色
	re, err := regexp.Compile("\u001B\\[[0-9]+[;0-9]*m")
	if err != nil {
		this_.Logger.Error("regexp.Compile error", zap.Error(err))
	} else {
		str = re.ReplaceAllString(str, "")
	}
	re, err = regexp.Compile("\u001B\\[[0-9]+[XABD]+")
	if err != nil {
		this_.Logger.Error("regexp.Compile error", zap.Error(err))
	} else {
		str = re.ReplaceAllString(str, "")
	}
	str = strings.ReplaceAll(str, "\u001B]0;", "")
	str = strings.ReplaceAll(str, "\a", "")
	str = strings.ReplaceAll(str, "\u001B[C", "")
	str = strings.ReplaceAll(str, "\u001B(B", "")
	str = strings.ReplaceAll(str, "\u001B[K", "")
	str = strings.ReplaceAll(str, "\u001B[m", "")
	if err != nil {
		this_.Logger.Error("regexp.Compile error", zap.Error(err))
	} else {
		str = re.ReplaceAllString(str, "")
	}
	writer := bufio.NewWriter(this_.commandLogFile)
	_, err = writer.WriteString(str)

	if err != nil {
		if this_.commandLogFile != nil {
			_ = this_.commandLogFile.Close()
		}
		this_.commandLogFile = nil
		return
	}
	_ = writer.Flush()
}

func (this_ *Worker) startReadWS(isWindow bool, baseLog *TerminalLogModel) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("startReadWS error", zap.Any("error", e))
		}
	}()

	defer func() { this_.stopAll() }()
	logContext_ := &logContext{}
	var buf []byte
	var readErr error
	var writeErr error

	var isClosed bool
	this_.ws.SetCloseHandler(func(code int, text string) error {
		isClosed = true
		return nil
	})
	for !isClosed {
		_, buf, readErr = this_.ws.ReadMessage()
		if readErr != nil && readErr != io.EOF {
			break
		}
		//this_.Logger.Info("ws on read", zap.Any("bs", string(buf)))
		this_.onCommand(logContext_, buf, *baseLog)
		_, writeErr = this_.service.Write(buf)

		if writeErr != nil {
			break
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.GetService(this_.key) == nil {
		return
	}

	if readErr != nil && !isClosed {
		this_.Logger.Error("ws read error", zap.Error(readErr))
	}

	if writeErr != nil && !isClosed {
		this_.Logger.Error("service write error", zap.Error(writeErr))
	}

	this_.Logger.Info("ws read is end")

	return
}

func (this_ *Worker) startReadService(isWindow bool) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("startReadService error", zap.Any("error", e))
		}
		this_.Logger.Info("service read end", zap.Any("key", this_.key))
	}()

	defer func() {
		this_.onServiceRead([]byte(fmt.Sprintf("\n\n结束时间:%s\n\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000"))))
		this_.stopAll()
	}()

	var n int
	var buf = make([]byte, 1024*32)
	var readErr error
	var writeErr error
	this_.onServiceRead([]byte(fmt.Sprintf("\n\n开始时间:%s\n\n", util.TimeFormat(time.Now(), "2006-01-02 15:04:05.000"))))
	for {
		n, readErr = this_.service.Read(buf)
		if readErr != nil && readErr != io.EOF {
			break
		}
		//this_.Logger.Info("service on read", zap.Any("bs", string(buf[:n])))

		//n, readErr, writeErr = this_.doSZ(key, n, buf, service)
		//if readErr != nil || writeErr != nil {
		//	break
		//}

		if n > 0 {
			this_.onServiceRead(buf[:n])
			writeErr = this_.ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if writeErr != nil {
				break
			}
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.GetService(this_.key) == nil {
		return
	}

	if readErr != nil {
		this_.Logger.Error("service read error", zap.Error(readErr))
	}

	if writeErr != nil {
		this_.Logger.Error("ws write error", zap.Error(writeErr))
	}

	this_.Logger.Info("service read is end")

	return
}

func (this_ *WorkerFactory) stopService(key string) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("stopService error", zap.Any("error", e))
		}
	}()

	this_.workerCacheLock.Lock()
	defer this_.workerCacheLock.Unlock()

	find := this_.workerCache[key]
	if find == nil {
		return
	}
	delete(this_.workerCache, key)
	this_.Logger.Info("stop service", zap.Any("key", key))
	find.service.Stop()
}

func (this_ *Worker) stopAll() {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("stopAll error", zap.Any("error", e))
		}
	}()

	this_.Logger.Info("stopAll", zap.Any("key", this_.key))
	this_.stopService(this_.key)
	if this_ != nil {
		this_.service.Stop()
	}
	if this_.ws != nil {
		_ = this_.ws.Close()
	}
	if this_.commandLogFile != nil {
		_ = this_.commandLogFile.Close()
	}
}
