package node

import (
	"go.uber.org/zap"
	"net"
)

type InnerServer struct {
	inner        *NetConfig
	toNodeIdList []string
	isStop       bool
}

func (this_ *InnerServer) Start() (err error) {

	go this_.serverListenerKeepAlive()

	return
}

func (this_ *InnerServer) isStopped() bool {
	return this_.isStop
}

func (this_ *InnerServer) serverListenerKeepAlive() {
	if this_.isStopped() {
		return
	}
	defer func() {
		go this_.serverListenerKeepAlive()
	}()
	listener, err := net.Listen(this_.inner.GetNetwork(), this_.inner.GetAddress())
	if err != nil {
		Logger.Error(this_.inner.GetInfoStr()+" listen error", zap.Error(err))
		return
	}
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			Logger.Error(this_.inner.GetInfoStr()+" listen accept error", zap.Error(err))
			continue
		}
		go this_.onConn(conn)
	}

}

func (this_ *InnerServer) onConn(conn net.Conn) {

}
