package node

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"teamide/pkg/util"
)

type InnerServer struct {
	netProxy *NetProxy
	isStop   bool
	*Worker
	connCache     map[string]net.Conn
	connCacheLock sync.Mutex
}

func (this_ *InnerServer) Start() {
	this_.connCache = make(map[string]net.Conn)

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
	Logger.Info(this_.Node.GetNodeStr() + " 服务 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 启动")
	listener, err := net.Listen(this_.netProxy.Inner.GetNetwork(), this_.netProxy.Inner.GetAddress())
	if err != nil {
		Logger.Error(this_.Node.GetNodeStr()+" 服务 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 监听异常", zap.Error(err))
		return
	}
	Logger.Info(this_.Node.GetNodeStr() + " 服务 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 启动成功")
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			Logger.Error(this_.netProxy.Inner.GetInfoStr()+" listen accept error", zap.Error(err))
			continue
		}
		go this_.onConn(conn)
	}

}

func (this_ *InnerServer) onConn(conn net.Conn) {
	Logger.Info(this_.Node.GetNodeStr() + " 服务 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 新连接")
	var connId = util.UUID()
	var netProxyId = this_.netProxy.Id
	this_.setConn(connId, conn)

	defer func() {
		_ = this_.closeConn(connId)
		_ = this_.netProxyCloseConn(false, this_.netProxy.LineNodeIdList, netProxyId, connId)
	}()
	var err error

	err = this_.netProxyNewConn(this_.netProxy.LineNodeIdList, netProxyId, connId)

	if err != nil {
		Logger.Error(this_.Node.GetNodeStr()+" 服务 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 节点线连接创建异常", zap.Error(err))
		return
	}
	for {
		var bytes = make([]byte, 1024*8)
		n, err := conn.Read(bytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			Logger.Error(this_.Node.GetNodeStr()+" 服务 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 读取异常", zap.Error(err))
			break
		}
		bytes = bytes[:n]
		err = this_.netProxySend(false, this_.netProxy.LineNodeIdList, netProxyId, connId, bytes)
		if err != nil {
			Logger.Error(this_.Node.GetNodeStr()+" 服务 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 节点线流发送异常", zap.Error(err))
			break
		}
	}

}

func (this_ *InnerServer) setConn(connId string, conn net.Conn) {
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	this_.connCache[connId] = conn
	return
}
func (this_ *InnerServer) getConn(connId string) (conn net.Conn) {
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	conn, _ = this_.connCache[connId]
	return
}

func (this_ *InnerServer) closeConn(connId string) (err error) {
	this_.connCacheLock.Lock()
	defer this_.connCacheLock.Unlock()
	conn, ok := this_.connCache[connId]
	if ok {
		delete(this_.connCache, connId)
		_ = conn.Close()
	}
	return
}

func (this_ *InnerServer) send(connId string, bytes []byte) (err error) {
	conn := this_.getConn(connId)
	if conn != nil {
		_, err = conn.Write(bytes)
		Logger.Info(this_.Node.GetNodeStr() + " 服务 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 连接 [" + connId + "] 发送 [" + fmt.Sprint(len(bytes)) + "]")
	} else {
		Logger.Warn(this_.Node.GetNodeStr() + " 服务 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 连接 [" + connId + "] 不存在")
	}
	return
}
