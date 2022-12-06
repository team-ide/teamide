package node

import (
	"errors"
	"go.uber.org/zap"
	"net"
	"teamide/pkg/util"
)

type OuterListener struct {
	netProxy *NetProxyOuter
	isStop   bool
	worker   *Worker
	*connCache
	MonitorData *MonitorData
}

func (this_ *OuterListener) Start() {
	this_.MonitorData = &MonitorData{}
	this_.connCache = newConnCache(this_.MonitorData)

	return
}

func (this_ *OuterListener) Stop() {
	this_.isStop = true

	this_.connCache.clean()

	return
}

func (this_ *OuterListener) isStopped() bool {
	return this_.isStop
}

func (this_ *OuterListener) newConn(connId string) (err error) {
	if this_.isStopped() {
		return
	}

	//Logger.Info(" OuterListener newConn [" + connId + "]")

	conn, err := net.Dial(this_.netProxy.GetType(), this_.netProxy.GetAddress())
	if err != nil {
		Logger.Error(this_.netProxy.GetInfoStr()+" 连接 ["+connId+"] 异常", zap.Error(err))
		return
	}
	//Logger.Info(this_.server.GetServerInfo() + " 至 " + this_.netProxy.Outer.GetInfoStr() + " 连接 [" + connId + "] 成功")
	this_.setConn(connId, conn)
	go func() {
		var netProxyId = this_.netProxy.Id
		defer func() {
			_ = this_.closeConn(connId)
			_ = this_.worker.netProxyCloseConn(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId)
		}()

		var buf = make([]byte, 1024*32)

		start := util.Now().UnixNano()
		err = util.Read(conn, buf, func(n int) (e error) {
			if this_.isStopped() {
				e = errors.New("proxy outer is stopped")
				return
			}

			end := util.Now().UnixNano()
			this_.MonitorData.monitorRead(int64(n), end-start)

			e = this_.worker.netProxySend(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId, buf[:n])
			if e != nil {
				Logger.Error(this_.netProxy.GetInfoStr()+" 连接 发送异常", zap.Error(e))
				return
			}
			start = util.Now().UnixNano()
			return
		})

	}()
	return
}
