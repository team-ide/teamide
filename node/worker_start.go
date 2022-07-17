package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"teamide/pkg/util"
	"time"
)

func (this_ *Worker) Start() (err error) {
	this_.callbackCache = make(map[string]func(msg *Message))
	this_.childrenNodeListenerPoolCache = make(map[string]*MessageListenerPool)
	this_.fromNodeListenerPoolCache = make(map[string]*MessageListenerPool)
	this_.netProxyInnerCache = make(map[string]*InnerServer)
	this_.netProxyOuterCache = make(map[string]*OuterListener)

	err = this_.AddNode(this_.Node)

	if err != nil {
		return
	}

	go this_.serverListenerKeepAlive()

	return
}

func (this_ *Worker) serverListenerKeepAlive() {
	if this_.isStopped() {
		return
	}
	defer func() {
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive()
	}()
	Logger.Info(this_.Node.GetNodeStr() + " 服务启动 开始")
	listener, err := net.Listen(this_.Node.GetNetwork(), this_.Node.GetAddress())
	if err != nil {
		Logger.Error(this_.Node.GetNodeStr()+" 服务启动 异常", zap.Error(err))
		return
	}
	Logger.Info(this_.Node.GetNodeStr() + " 服务启动 成功")
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			Logger.Error(this_.Node.GetNodeStr()+" 服务监听 异常", zap.Error(err))
			continue
		}
		go func(conn_ net.Conn) {
			for {
				var msg *Message
				msg, err = ReadMessage(conn_)
				if err != nil {
					_ = conn_.Close()
					return
				}
				if msg.Method != methodOK {
					_ = conn_.Close()
					return
				}
				if msg.FromNodeId != "" {
					pool := this_.getFromNodeListenerPool(msg.FromNodeId)
					var find = this_.findNode(msg.FromNodeId)
					messageListener := &MessageListener{
						conn:      conn,
						onMessage: this_.onMessage,
						id:        util.UUID(),
					}
					messageListener.listen(func() {
						messageListener.stop()
						pool.Remove(messageListener)
						if find != nil {
							Logger.Info(this_.Node.GetNodeStr() + " 移除至 " + find.GetNodeStr() + " 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
						} else {
							Logger.Info(this_.Node.GetNodeStr() + " 移除至 " + msg.FromNodeId + " 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
						}
					})
					pool.Put(messageListener)

					this_.Node.ParentId = msg.FromNodeId
					var findThisNode = this_.findNode(this_.Node.Id)
					if findThisNode != nil {
						if findThisNode.ParentId != msg.FromNodeId {
							findThisNode.ParentId = msg.FromNodeId
							this_.refreshNodeList()
						}
					}
					if find != nil {
						Logger.Info(this_.Node.GetNodeStr() + " 添加至 " + find.GetNodeStr() + " 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
					} else {
						Logger.Info(this_.Node.GetNodeStr() + " 添加至 " + msg.FromNodeId + " 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
					}
				}
				return
			}
		}(conn)
	}

}

func (this_ *Worker) getFromNodeListenerPool(fromNodeId string) (pool *MessageListenerPool) {
	this_.fromNodeListenerLock.Lock()
	defer this_.fromNodeListenerLock.Unlock()

	pool, ok := this_.fromNodeListenerPoolCache[fromNodeId]
	if !ok {
		pool = &MessageListenerPool{}
		this_.fromNodeListenerPoolCache[fromNodeId] = pool
	}
	return
}
