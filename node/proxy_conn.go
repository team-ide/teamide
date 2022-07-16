package node

import (
	"go.uber.org/zap"
	"io"
	"net"
)

type ProxyConnListener struct {
	conn            net.Conn
	isStop          bool
	isClose         bool
	trackId         string
	lineNodeIdList  []string
	messageListener *MessageListener
	*Worker
}

func (this_ *ProxyConnListener) stop() {
	this_.isStop = true
	_ = this_.conn.Close()
}

func (this_ *ProxyConnListener) listen(onClose func()) {
	var err error
	this_.isClose = false
	go func() {
		defer func() {
			this_.isClose = true
			if x := recover(); x != nil {
				Logger.Error("proxy conn listen error", zap.Error(err))
				return
			}
			_ = this_.conn.Close()
			onClose()
		}()

		for {
			if this_.isStop {
				return
			}
			var n int
			var bytes = make([]byte, 1024*8)
			n, err = this_.conn.Read(bytes)

			if err != nil {
				if this_.isStop {
					return
				}
				if err == io.EOF {
					return
				}
				Logger.Error("proxy conn read error", zap.Error(err))
				return
			}
			bytes = bytes[0:n]

			_ = this_.messageListener.Send(&Message{
				Method:         methodProxySend,
				TrackId:        this_.trackId,
				LineNodeIdList: this_.lineNodeIdList,
				Bytes:          bytes,
			})

		}
	}()
}
