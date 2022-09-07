package module_terminal

import (
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"strconv"
	"sync"
	"teamide/internal/context"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/ssh"
	"teamide/pkg/terminal"
	"teamide/pkg/util"
)

func NewWorker(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *worker {
	return &worker{
		ServerContext:  toolboxService_.ServerContext,
		toolboxService: toolboxService_,
		nodeService:    nodeService_,
		serviceCache:   make(map[string]terminal.Service),
	}
}

type worker struct {
	*context.ServerContext
	toolboxService   *module_toolbox.ToolboxService
	nodeService      *module_node.NodeService
	serviceCache     map[string]terminal.Service
	serviceCacheLock sync.Mutex
}

func (this_ *worker) GetService(key string) (service terminal.Service) {
	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	service = this_.serviceCache[key]
	return
}

func (this_ *worker) createService(key string, place string, placeId string) (service terminal.Service, err error) {
	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

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
		service = module_node.NewTerminalService(placeId)
	}
	if service == nil {
		err = errors.New("[" + place + "]文件服务不存在")
		return
	}

	this_.serviceCache[key] = service
	return
}

func (this_ *worker) Start(key string, place string, placeId string, size *terminal.Size, ws *websocket.Conn) (err error) {
	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	service := this_.serviceCache[key]
	if service != nil {
		err = errors.New("会话服务[" + key + "]已存在")
		return
	}

	service, err = this_.createService(key, place, placeId)
	if err != nil {
		return
	}
	err = service.Start(size)
	if err != nil {
		return
	}
	go this_.startReadWS(key, ws, service)
	go this_.startReadService(key, ws, service)

	this_.serviceCache[key] = service
	return
}

func (this_ *worker) startReadWS(key string, ws *websocket.Conn, service terminal.Service) {

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
		if messageType == websocket.TextMessage {

		}
		writeErr = util.Write(service, buf, nil)
		if readErr == io.EOF {
			readErr = nil
			break
		}
		if writeErr != nil {
			break
		}
	}

	if readErr != nil {
		this_.Logger.Error("ws read error", zap.Error(readErr))
	}

	if writeErr != nil {
		this_.Logger.Error("service write error", zap.Error(writeErr))
	}

	return
}

func (this_ *worker) stopService(key string) {

	this_.serviceCacheLock.Lock()
	defer this_.serviceCacheLock.Unlock()

	service := this_.serviceCache[key]
	if service == nil {
		return
	}
	delete(this_.serviceCache, key)
	service.Stop()
}

func (this_ *worker) stopAll(key string, ws *websocket.Conn, service terminal.Service) {

	this_.stopService(key)
	if service != nil {
		service.Stop()
	}
	if ws != nil {
		_ = ws.Close()
	}
}

func (this_ *worker) startReadService(key string, ws *websocket.Conn, service terminal.Service) {

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

		writeErr = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
		if readErr == io.EOF {
			readErr = nil
			break
		}
		if writeErr != nil {
			break
		}
	}

	if readErr != nil {
		this_.Logger.Error("service read error", zap.Error(readErr))
	}

	if writeErr != nil {
		this_.Logger.Error("ws write error", zap.Error(writeErr))
	}

	return
}
