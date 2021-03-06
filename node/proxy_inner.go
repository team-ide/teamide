package node

import (
	"go.uber.org/zap"
	"io"
	"net"
	"teamide/pkg/util"
	"time"
)

type InnerServer struct {
	netProxy *NetProxy
	isStop   bool
	*Worker
	*connCache
	serverListener net.Listener
	MonitorData    *MonitorData
}

func (this_ *InnerServer) Start() {
	this_.MonitorData = &MonitorData{}
	this_.connCache = newConnCache(this_.MonitorData)

	go this_.serverListenerKeepAlive()

	return
}

func (this_ *InnerServer) Stop() {
	this_.isStop = true
	_ = this_.serverListener.Close()
	this_.connCache.clean()
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
		this_.notifyAll(&Message{
			NetProxyInnerStatusChangeList: []*StatusChange{
				{
					Id:          this_.netProxy.Id,
					Status:      StatusStopped,
					StatusError: "",
				},
			},
		})
		if !this_.isStopped() {
			return
		}
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive()
	}()
	var err error
	Logger.Info(this_.server.GetServerInfo() + " 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 启动")

	this_.serverListener, err = net.Listen(this_.netProxy.Inner.GetType(), this_.netProxy.Inner.GetAddress())
	if err != nil {
		Logger.Error(this_.server.GetServerInfo()+" 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 监听异常", zap.Error(err))

		this_.notifyAll(&Message{
			NetProxyInnerStatusChangeList: []*StatusChange{
				{
					Id:          this_.netProxy.Id,
					Status:      StatusError,
					StatusError: err.Error(),
				},
			},
		})
		return
	}
	Logger.Info(this_.server.GetServerInfo() + " 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 启动成功")
	this_.notifyAll(&Message{
		NetProxyInnerStatusChangeList: []*StatusChange{
			{
				Id:          this_.netProxy.Id,
				Status:      StatusStarted,
				StatusError: "",
			},
		},
	})
	for {
		if this_.isStopped() {
			break
		}
		var conn net.Conn
		conn, err = this_.serverListener.Accept()
		if err != nil {
			if this_.isStopped() {
				break
			}
			Logger.Error(this_.netProxy.Inner.GetInfoStr()+" listen accept error", zap.Error(err))
			break
		}
		go this_.onConn(conn)
	}

}

func (this_ *InnerServer) onConn(conn net.Conn) {
	if this_.isStopped() {
		_ = conn.Close()
		return
	}
	//Logger.Info(this_.server.GetServerInfo() + " 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 新连接")
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
		Logger.Error(this_.server.GetServerInfo()+" 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 节点线连接创建异常", zap.Error(err))
		return
	}
	for {
		if this_.isStopped() {
			break
		}
		start := util.Now().UnixNano()

		var n int
		var bytes = make([]byte, 1024*8)
		n, err = conn.Read(bytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			//Logger.Error(this_.server.GetServerInfo()+" 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 读取异常", zap.Error(err))
			break
		}
		bytes = bytes[:n]

		end := util.Now().UnixNano()
		this_.MonitorData.monitorRead(int64(n), end-start)

		err = this_.netProxySend(false, this_.netProxy.LineNodeIdList, netProxyId, connId, bytes)
		if err != nil {
			Logger.Error(this_.server.GetServerInfo()+" 代理服务 "+this_.netProxy.Inner.GetInfoStr()+" 节点线流发送异常", zap.Error(err))
			break
		}
	}

}
