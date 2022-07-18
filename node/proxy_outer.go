package node

import (
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
)

type OuterListener struct {
	netProxy *NetProxy
	isStop   bool
	*Worker
	connCache     map[string]net.Conn
	connCacheLock sync.Mutex
}

func (this_ *OuterListener) Start() {
	this_.connCache = make(map[string]net.Conn)

	return
}
func (this_ *OuterListener) Stop() {
	this_.isStop = true

	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	for _, conn := range this_.connCache {
		_ = conn.Close()
	}
	this_.connCache = make(map[string]net.Conn)
	return
}

func (this_ *OuterListener) newConn(connId string) (err error) {
	if this_.isStop {
		return
	}
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()

	Logger.Info(this_.server.GetServerInfo() + " 新建至 " + this_.netProxy.Outer.GetInfoStr() + " 的连接 [" + connId + "]")

	conn, err := net.Dial(this_.netProxy.Outer.GetNetwork(), this_.netProxy.Outer.GetAddress())
	if err != nil {
		Logger.Error(this_.server.GetServerInfo()+" 至 "+this_.netProxy.Outer.GetInfoStr()+" 连接 ["+connId+"] 异常", zap.Error(err))
		return
	}
	Logger.Info(this_.server.GetServerInfo() + " 至 " + this_.netProxy.Outer.GetInfoStr() + " 连接 [" + connId + "] 成功")
	this_.connCache[connId] = conn
	go func() {
		var netProxyId = this_.netProxy.Id
		defer func() {
			_ = this_.closeConn(connId)
			_ = this_.netProxyCloseConn(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId)
		}()

		for {
			var n int
			var bytes = make([]byte, 1024*8)
			n, err = conn.Read(bytes)
			if err != nil {
				if err == io.EOF {
					break
				}
				//Logger.Error(this_.server.GetServerInfo()+" 至 "+this_.netProxy.Outer.GetInfoStr()+" 连接 读取异常", zap.Error(err))
				break
			}
			bytes = bytes[:n]
			err = this_.netProxySend(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId, bytes)
			if err != nil {
				Logger.Error(this_.server.GetServerInfo()+" 至 "+this_.netProxy.Outer.GetInfoStr()+" 连接 发送异常", zap.Error(err))
				break
			}
		}
	}()
	return
}

func (this_ *OuterListener) getConn(connId string) (conn net.Conn) {
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	conn, _ = this_.connCache[connId]
	return
}

func (this_ *OuterListener) closeConn(connId string) (err error) {
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	conn, ok := this_.connCache[connId]
	if ok {
		delete(this_.connCache, connId)
		_ = conn.Close()
	}
	return
}

func (this_ *OuterListener) send(connId string, bytes []byte) (err error) {
	conn := this_.getConn(connId)
	if conn != nil {
		_, err = conn.Write(bytes)
		//Logger.Info(this_.server.GetServerInfo() + " 至 " + this_.netProxy.Outer.GetInfoStr() + " 连接 [" + connId + "] 发送 [" + fmt.Sprint(len(bytes)) + "]")
	} else {
		Logger.Warn(this_.server.GetServerInfo() + " 至 " + this_.netProxy.Outer.GetInfoStr() + " 连接 [" + connId + "] 不存在")
	}
	return
}
