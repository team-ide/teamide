package node

import (
	"go.uber.org/zap"
	"net"
)

func (this_ *Worker) Start() (err error) {
	this_.callbackCache = make(map[string]func(msg *Message))
	this_.childrenNodeListenerCache = make(map[string]*MessageListener)

	go this_.serverListenerKeepAlive()

	return
}

func (this_ *Worker) serverListenerKeepAlive() {
	if this_.isStopped() {
		return
	}
	defer func() {
		go this_.serverListenerKeepAlive()
	}()
	listener, err := net.Listen(this_.Node.GetNetwork(), this_.Node.Address)
	if err != nil {
		Logger.Error(this_.Node.GetNodeStr()+" listen error", zap.Error(err))
		return
	}
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			Logger.Error(this_.Node.GetNodeStr()+" listen accept error", zap.Error(err))
			continue
		}
		messageListener := &MessageListener{
			conn:      conn,
			onMessage: this_.onMessage,
		}
		messageListener.listen(func() {

		})
	}

}
