package node

import (
	"go.uber.org/zap"
	"io"
	"net"
	"teamide/pkg/util"
)

type OuterListener struct {
	netProxy *NetProxy
	isStop   bool
	*Worker
	*connCache
	MonitorData *MonitorData
}

func (this_ *OuterListener) Start() {
	this_.MonitorData = &MonitorData{}
	this_.connCache = newConnCache(this_.MonitorData)

	this_.notifyAll(&Message{
		NetProxyOuterStatusChangeList: []*StatusChange{
			{
				Id:          this_.netProxy.Id,
				Status:      StatusStarted,
				StatusError: "",
			},
		},
	})
	return
}

func (this_ *OuterListener) Stop() {
	this_.isStop = true

	this_.connCache.clean()

	this_.notifyAll(&Message{
		NetProxyOuterStatusChangeList: []*StatusChange{
			{
				Id:          this_.netProxy.Id,
				Status:      StatusStopped,
				StatusError: "",
			},
		},
	})
	return
}

func (this_ *OuterListener) isStopped() bool {
	return this_.isStop
}

func (this_ *OuterListener) newConn(connId string) (err error) {
	if this_.isStopped() {
		return
	}

	//Logger.Info(this_.server.GetServerInfo() + " 新建至 " + this_.netProxy.Outer.GetInfoStr() + " 的连接 [" + connId + "]")

	conn, err := net.Dial(this_.netProxy.Outer.GetType(), this_.netProxy.Outer.GetAddress())
	if err != nil {
		Logger.Error(this_.server.GetServerInfo()+" 至 "+this_.netProxy.Outer.GetInfoStr()+" 连接 ["+connId+"] 异常", zap.Error(err))
		return
	}
	//Logger.Info(this_.server.GetServerInfo() + " 至 " + this_.netProxy.Outer.GetInfoStr() + " 连接 [" + connId + "] 成功")
	this_.setConn(connId, conn)
	go func() {
		var netProxyId = this_.netProxy.Id
		defer func() {
			_ = this_.closeConn(connId)
			_ = this_.netProxyCloseConn(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId)
		}()

		for {
			if this_.isStopped() {
				return
			}

			start := util.Now().UnixNano()

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

			end := util.Now().UnixNano()
			this_.MonitorData.monitorRead(int64(n), end-start)

			err = this_.netProxySend(true, this_.netProxy.ReverseLineNodeIdList, netProxyId, connId, bytes)
			if err != nil {
				Logger.Error(this_.server.GetServerInfo()+" 至 "+this_.netProxy.Outer.GetInfoStr()+" 连接 发送异常", zap.Error(err))
				break
			}
		}
	}()
	return
}
