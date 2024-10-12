package module_serial

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"sync"
)

type Config struct {
	PortName                string `json:"portName"`
	BaudRate                uint   `json:"baudRate"` // 波特率
	DataBits                uint   `json:"dataBits"` // 数据位
	StopBits                uint   `json:"stopBits"` // 停止位
	ParityMode              int    `json:"parityMode"`
	InterCharacterTimeout   uint   `json:"interCharacterTimeout"`
	MinimumReadSize         uint   `json:"minimumReadSize"` //
	RTSCTSFlowControl       bool   `json:"rtsCTSFlowControl"`
	Rs485Enable             bool   `json:"rs485Enable"`
	Rs485RtsHighDuringSend  bool   `json:"rs485RtsHighDuringSend"`
	Rs485RtsHighAfterSend   bool   `json:"rs485RtsHighAfterSend"`
	Rs485RxDuringTx         bool   `json:"rs485RxDuringTx"`
	Rs485DelayRtsBeforeSend int    `json:"rs485DelayRtsBeforeSend"`
	Rs485DelayRtsAfterSend  int    `json:"rs485DelayRtsAfterSend"`
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
	conn      io.ReadWriteCloser
	ws        *websocket.Conn
	isStopped bool
}

func (this_ *Service) init() (err error) {

	this_.conn, err = serial.Open(serial.OpenOptions{
		PortName:                this_.Config.PortName,
		BaudRate:                this_.Config.BaudRate,
		DataBits:                this_.Config.DataBits,
		StopBits:                this_.Config.StopBits,
		ParityMode:              serial.ParityMode(this_.Config.ParityMode),
		RTSCTSFlowControl:       this_.Config.RTSCTSFlowControl,
		InterCharacterTimeout:   this_.Config.InterCharacterTimeout,
		MinimumReadSize:         this_.Config.MinimumReadSize,
		Rs485Enable:             this_.Config.Rs485Enable,
		Rs485RtsHighDuringSend:  this_.Config.Rs485RtsHighDuringSend,
		Rs485RtsHighAfterSend:   this_.Config.Rs485RtsHighAfterSend,
		Rs485RxDuringTx:         this_.Config.Rs485RxDuringTx,
		Rs485DelayRtsAfterSend:  this_.Config.Rs485DelayRtsAfterSend,
		Rs485DelayRtsBeforeSend: this_.Config.Rs485DelayRtsBeforeSend,
	})
	if err != nil {
		util.Logger.Error("serial open error", zap.Error(err))
		return
	}
	return
}

func (this_ *Service) start(ws *websocket.Conn) (err error) {

	_ = ws.WriteMessage(websocket.BinaryMessage, []byte("open to "+util.GetStringValue(this_.Config)+"\n"))
	err = this_.init()
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("conn error:"+err.Error()+"\n"))
		return
	}
	_ = ws.WriteMessage(websocket.BinaryMessage, []byte("conn success"+"\n"))
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
