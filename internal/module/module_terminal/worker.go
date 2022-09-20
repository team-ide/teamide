package module_terminal

import (
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"strconv"
	"sync"
	"teamide/internal/context"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/ssh"
	"teamide/pkg/terminal"
)

func NewWorker(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *worker {
	return &worker{
		ServerContext:        toolboxService_.ServerContext,
		toolboxService:       toolboxService_,
		nodeService:          nodeService_,
		serviceCache:         make(map[string]terminal.Service),
		fileHeadersChanCache: make(map[string]chan []*multipart.FileHeader),
		writerChanCache:      make(map[string]chan io.Writer),
	}
}

type worker struct {
	*context.ServerContext
	toolboxService           *module_toolbox.ToolboxService
	nodeService              *module_node.NodeService
	serviceCache             map[string]terminal.Service
	serviceCacheLock         sync.Mutex
	fileHeadersChanCache     map[string]chan []*multipart.FileHeader
	fileHeadersChanCacheLock sync.Mutex
	writerChanCache          map[string]chan io.Writer
	writerChanCacheLock      sync.Mutex
}

func (this_ *worker) GetService(key string) (service terminal.Service) {
	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	service = this_.serviceCache[key]
	return
}

func (this_ *worker) createService(place string, placeId string) (service terminal.Service, err error) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("createService error", zap.Any("error", e))
		}
	}()

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

	return
}

func (this_ *worker) Start(key string, place string, placeId string, size *terminal.Size, ws *websocket.Conn) (err error) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("Start error", zap.Any("error", e))
		}
	}()

	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	service := this_.serviceCache[key]
	if service != nil {
		err = errors.New("会话服务[" + key + "]已存在")
		return
	}

	service, err = this_.createService(place, placeId)
	if err != nil {
		return
	}
	isWindow, err := service.IsWindows()
	if err != nil {
		return
	}
	err = service.Start(size)
	if err != nil {
		return
	}
	go this_.startReadWS(key, isWindow, ws, service)
	go this_.startReadService(key, isWindow, ws, service)
	go this_.startReadErrService(key, isWindow, ws, service)

	this_.serviceCache[key] = service
	return
}

func (this_ *worker) startReadWS(key string, isWindow bool, ws *websocket.Conn, service terminal.Service) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("startReadWS error", zap.Any("error", e))
		}
	}()

	defer func() { this_.stopAll(key, ws, service) }()

	var messageType int
	var buf []byte
	var readErr error
	var writeErr error
	for {
		messageType, buf, readErr = ws.ReadMessage()
		if readErr != nil && readErr != io.EOF {
			break
		}
		if messageType == websocket.BinaryMessage {

		}
		//this_.Logger.Info("ws on read", zap.Any("bs", string(buf)))
		_, writeErr = service.Write(buf)
		if writeErr != nil {
			break
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.GetService(key) == nil {
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

func (this_ *worker) startRZ(key string, service terminal.Service) (err error) {
	progress := newProgress(key)
	progress.Data["readSize"] = 0
	progress.Data["writeSize"] = 0
	progress.Data["successSize"] = 0
	defer func() { progress.end(err) }()

	action, err := progress.waitAction("upload")
	if err != nil {
		return
	}

	if action == nil {
		return
	}
	fileHeaders, ok := action.([]*multipart.FileHeader)
	if !ok {
		err = errors.New("action can not to []*multipart.FileHeader")
		return
	}
	if len(fileHeaders) == 0 {
		return
	}

	return
}

func (this_ *worker) startReadService(key string, isWindow bool, ws *websocket.Conn, service terminal.Service) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("startReadService error", zap.Any("error", e))
		}
	}()

	defer func() { this_.stopAll(key, ws, service) }()

	var n int
	var buf = make([]byte, 1024*32)
	var readErr error
	var writeErr error
	for {
		n, readErr = service.Read(buf)
		if readErr != nil && readErr != io.EOF {
			break
		}
		//this_.Logger.Info("service on read", zap.Any("bs", string(buf[:n])))

		if x, ok := ByteContains(buf[:n], ZModemSZStart); ok {
			var rsErr error
			n, rsErr = this_.startSZ(key, x, buf, service)
			if rsErr != nil {
				writeErr = rsErr
				break
			}
		}

		if n > 0 {
			writeErr = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if writeErr != nil {
				break
			}
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.GetService(key) == nil {
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

func (this_ *worker) startReadErrService(key string, isWindow bool, ws *websocket.Conn, service terminal.Service) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("startReadErrService error", zap.Any("error", e))
		}
	}()

	var n int
	var buf = make([]byte, 1024*32)
	var readErr error
	var writeErr error
	for {
		n, readErr = service.ReadError(buf)
		if readErr != nil && readErr != io.EOF {
			break
		}
		//this_.Logger.Info("service on read err", zap.Any("bs", string(buf[:n])))

		if n > 0 {
			writeErr = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if writeErr != nil {
				break
			}
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.GetService(key) == nil {
		return
	}

	if readErr != nil {
		this_.Logger.Error("service read err error", zap.Error(readErr))
	}

	if writeErr != nil {
		this_.Logger.Error("ws write err error", zap.Error(writeErr))
	}

	return
}

func (this_ *worker) stopService(key string) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("stopService error", zap.Any("error", e))
		}
	}()

	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	progressList := getProgressList(key)
	for _, one := range progressList {
		one.closeCallAction()
	}

	service := this_.serviceCache[key]
	if service == nil {
		return
	}
	delete(this_.serviceCache, key)
	service.Stop()
}

func (this_ *worker) stopAll(key string, ws *websocket.Conn, service terminal.Service) {

	defer func() {
		if e := recover(); e != nil {
			this_.Logger.Error("stopAll error", zap.Any("error", e))
		}
	}()

	this_.stopService(key)
	if service != nil {
		service.Stop()
	}
	if ws != nil {
		_ = ws.Close()
	}
}

func (this_ *worker) CallAction(progressId string, action interface{}) (err error) {
	progress := getProgress(progressId)
	if progress != nil {
		progress.callAction(action)
	}
	return
}
