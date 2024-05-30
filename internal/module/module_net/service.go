package module_net

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"sync"
	"time"
)

type Config struct {
	Address   string      `json:"address"`
	Timeout   int64       `json:"timeout"`
	SSHClient *ssh.Client `json:"-"`
}

func createService(config *Config) (service *Service, err error) {
	service = &Service{
		Config: config,
	}
	return
}

var (
	serviceCache     = map[string]*Service{}
	serviceCacheLock = &sync.Mutex{}
)

func getService(key string) (service *Service) {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()
	service = serviceCache[key]
	return
}

func removeService(key string) {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()
	delete(serviceCache, key)
	return
}
func setService(key string, service *Service) {
	serviceCacheLock.Lock()
	defer serviceCacheLock.Unlock()
	serviceCache[key] = service
	return
}

type Service struct {
	Key string
	*Config
	conn      net.Conn
	ws        *websocket.Conn
	isStopped bool
}

func (this_ *Service) init() (err error) {
	var network = "tcp"
	var timeout = this_.Config.Timeout
	if this_.SSHClient != nil {
		this_.conn, err = this_.SSHClient.Dial(network, this_.Address)
	} else {
		if timeout <= 0 {
			this_.conn, err = net.Dial(network, this_.Address)
		} else {
			this_.conn, err = net.DialTimeout(network, this_.Address, time.Millisecond*time.Duration(timeout))
		}
	}
	return
}

func (this_ *Service) start(ws *websocket.Conn) (err error) {

	_ = ws.WriteMessage(websocket.BinaryMessage, []byte("conn to "+this_.Address+"\n"))
	err = this_.init()
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("conn to "+this_.Address+" error:"+err.Error()+"\n"))
		return
	}
	_ = ws.WriteMessage(websocket.BinaryMessage, []byte("conn to "+this_.Address+" success"+"\n"))
	this_.ws = ws

	go this_.startReadWS()
	go this_.startReadService()
	return
}

func (this_ *Service) stop() {
	this_.isStopped = true
	removeService(this_.Key)
	if this_.conn != nil {
		_ = this_.conn.Close()
	}
	if this_.SSHClient != nil {
		_ = this_.SSHClient.Close()
	}
	if this_.ws != nil {
		_ = this_.ws.Close()
	}
	return
}

func (this_ *Service) startReadWS() {

	defer func() {
		if e := recover(); e != nil {
			err := errors.New(fmt.Sprint(e))
			util.Logger.Error("startReadWS panic error:", zap.Error(err))
		}
	}()

	defer func() { this_.stop() }()
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
		_, writeErr = this_.conn.Write(buf)

		if writeErr != nil {
			break
		}
		if readErr == io.EOF {
			readErr = nil
			break
		}
	}

	if this_.isStopped {
		return
	}

	if readErr != nil && !isClosed {
		util.Logger.Error("ws read error", zap.Error(readErr))
	}

	if writeErr != nil && !isClosed {
		util.Logger.Error("service write error", zap.Error(writeErr))
	}

	util.Logger.Info("ws read is end")

	return
}

func (this_ *Service) startReadService() {

	defer func() {
		if e := recover(); e != nil {
			err := errors.New(fmt.Sprint(e))
			util.Logger.Error("startReadService panic error", zap.Error(err))
		}
		util.Logger.Info("service read end", zap.Any("key", this_.Key))
	}()

	defer func() {
		this_.stop()
	}()

	var n int
	var buf = make([]byte, 1024*32)
	var readErr error
	var writeErr error
	for {
		n, readErr = this_.conn.Read(buf)
		if readErr != nil && readErr != io.EOF {
			break
		}
		//this_.Logger.Info("service on read", zap.Any("bs", string(buf[:n])))

		//n, readErr, writeErr = this_.doSZ(key, n, buf, service)
		//if readErr != nil || writeErr != nil {
		//	break
		//}

		if n > 0 {
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

	if this_.isStopped {
		return
	}

	if readErr != nil {
		util.Logger.Error("service read error", zap.Error(readErr))
	}

	if writeErr != nil {
		util.Logger.Error("ws write error", zap.Error(writeErr))
	}

	util.Logger.Info("service read is end")

	return
}
