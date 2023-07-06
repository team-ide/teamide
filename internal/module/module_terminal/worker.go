package module_terminal

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"regexp"
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

func (this_ *WorkerFactory) createService(place string, placeId string, workerId string) (worker *Worker, command string, err error) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("createService error", zap.Any("error", e))
		}
	}()
	var service terminal.Service
	switch place {
	case "local":
		service = terminal.NewLocalService()
	case "ssh":
		if placeId == "" {
			err = errors.New("SSH配置不能为空")
			return
		}
		var id int64
		id, err = strconv.ParseInt(placeId, 10, 64)
		if err != nil {
			return
		}
		var tD *module_toolbox.ToolboxModel
		tD, err = this_.toolboxService.Get(id)
		if err != nil {
			return
		}
		if tD == nil || tD.Option == "" {
			err = errors.New("SSH[" + placeId + "]配置不存在")
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

		service = ssh.NewTerminalService(config)
	case "node":
		if placeId == "" {
			err = errors.New("node配置不能为空")
			return
		}
		service = module_node.NewTerminalService(placeId, this_.nodeService)
	}
	if service == nil {
		err = errors.New("[" + place + "]文件服务不存在")
		return
	}

	worker = &Worker{
		place:         place,
		placeId:       placeId,
		workerId:      workerId,
		service:       service,
		WorkerFactory: this_,
	}
	worker.init()
	return
}

func (this_ *WorkerFactory) Start(key string, place string, placeId string, workerId string, size *terminal.Size, ws *websocket.Conn, baseLog *TerminalLogModel) (err error) {

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
	worker, command, err = this_.createService(place, placeId, workerId)
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
}

func (this_ *Worker) init() {
	dir := this_.GetFilesDir()
	dir += fmt.Sprintf("%s/toolbox-%s/%s", "toolbox-workers", this_.placeId, this_.workerId) + "/"

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

func (this_ *Worker) onServiceRead(bs []byte) {
	if this_.dir == "" {
		return
	}
	if this_.commandLogFile == nil {
		ex, err := util.PathExists(this_.dir)
		if err != nil {
			return
		}
		if !ex {
			err = os.MkdirAll(this_.dir, fs.ModePerm)
		}
		if ex, _ = util.PathExists(this_.dir + "/command.log"); ex {
			this_.commandLogFile, _ = os.OpenFile(this_.dir+"/command.log", os.O_WRONLY|os.O_APPEND, 0666)
		} else {
			this_.commandLogFile, _ = os.Create(this_.dir + "/command.log")
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
	str = strings.ReplaceAll(str, "\u001B]0;", "")
	str = strings.ReplaceAll(str, "\a", "")
	str = strings.ReplaceAll(str, "\u001B[C", "")
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
	for {
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

	if readErr != nil {
		this_.Logger.Error("ws read error", zap.Error(readErr))
	}

	if writeErr != nil {
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

	defer func() { this_.stopAll() }()

	var n int
	var buf = make([]byte, 1024*32)
	var readErr error
	var writeErr error
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
