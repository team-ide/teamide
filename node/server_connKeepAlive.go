package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

func (this_ *Server) connNodeListenerKeepAlive(toNode *Info, connAddress, connToken string, connSize int) {
	if connAddress == "" {
		Logger.Warn(this_.GetServerInfo() + " 连接 [" + connAddress + "] 连接地址为空")
		return
	}
	this_.connNodeListenerKeepAliveLock.Lock()
	defer this_.connNodeListenerKeepAliveLock.Unlock()

	var pool *MessageListenerPool
	if toNode != nil {
		pool = this_.cache.getNodeListenerPool(this_.Id, toNode.Id)
		if pool != nil {
			return
		}
		if this_.Id == toNode.Id {
			Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 不可连接两个节点ID相同的节点")
			return
		}
		pool = this_.cache.newIfAbsentNodeListenerPool(this_.Id, toNode.Id)
	}
	if connSize <= 0 {
		connSize = 5
	}
	for i := 0; i < connSize; i++ {
		go this_.connNodeListener(pool, connAddress, connToken, i)
	}
	return
}

func (this_ *Server) connNodeListener(pool *MessageListenerPool, connAddress, connToken string, clientIndex int) {
	if pool != nil && pool.isStop {
		return
	}
	var messageListener *MessageListener
	defer func() {
		if messageListener != nil {
			return
		}
		if pool != nil && pool.isStop {
			return
		}
		time.Sleep(5 * time.Second)
		go this_.connNodeListener(pool, connAddress, connToken, clientIndex)
	}()
	var err error
	var conn net.Conn
	Logger.Info(this_.GetServerInfo() + " 连接 [" + connAddress + "] 开始")
	conn, err = net.Dial("tcp", GetAddress(connAddress))
	if err != nil {
		if pool != nil && len(pool.listeners) == 0 {
			this_.worker.notifyAll(&Message{
				NodeId:          pool.toNodeId,
				NodeStatus:      StatusError,
				NodeStatusError: err.Error(),
			})
		}
		Logger.Warn(this_.GetServerInfo()+" 连接 ["+connAddress+"] 异常", zap.Any("error", err.Error()))
		return
	}

	var tokenBytes = []byte(connToken)
	if len(tokenBytes) > tokenByteSize {
		tokenBytes = tokenBytes[:tokenByteSize]
	} else if len(tokenBytes) < tokenByteSize {
		for i := len(tokenBytes); i < tokenByteSize; i++ {
			tokenBytes = append(tokenBytes, []byte(" ")...)
		}
	}

	_, err = conn.Write(tokenBytes)
	if err != nil {
		_ = conn.Close()
		return
	}

	var msg = &Message{
		Method:      methodOK,
		Node:        this_.rootNode,
		Ok:          true,
		ClientIndex: clientIndex,
	}

	if clientIndex == 0 {
		msg.NodeList = this_.cache.nodeList
		msg.NetProxyList = this_.cache.netProxyList
	}
	err = WriteMessage(conn, msg)
	if err != nil {
		Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	msg, err = ReadMessage(conn)
	if err != nil {
		Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	if msg.Node == nil {
		Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	serverNode := msg.Node
	Logger.Info(this_.GetServerInfo() + " 连接 [" + connAddress + "] 成功")

	toNodeId := serverNode.Id
	var fromNodeId = this_.Id

	if pool != nil {
		if pool.toNodeId != toNodeId {
			Logger.Error(this_.GetServerInfo()+" 连接 ["+connAddress+"] 节点ID不一致", zap.Any("toNodeId", pool.toNodeId), zap.Any("serverNodeId", toNodeId))
			this_.cache.removeNodeListenerPool(pool.fromNodeId, pool.toNodeId)
			return
		}
	}
	if pool == nil {
		pool = this_.cache.newIfAbsentNodeListenerPool(fromNodeId, toNodeId)
	}

	messageListener = &MessageListener{
		conn:      conn,
		onMessage: this_.worker.onMessage,
	}

	messageListener.listen(func() {
		messageListener.stop()
		pool.Remove(messageListener)
		Logger.Info(this_.GetServerInfo() + " 移除 连接至 [" + toNodeId + "][" + connAddress + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))

		if clientIndex == 0 {
			var notifyMsg = &Message{
				NodeId: serverNode.Id,
			}
			if pool.isStop {
				notifyMsg.NodeStatus = StatusStopped
				notifyMsg.NodeStatusError = ""
			} else {
				notifyMsg.NodeStatus = StatusError
				notifyMsg.NodeStatusError = "连接异常"
			}

			this_.worker.notifyAll(notifyMsg)
		}
		if !pool.isStop {
			time.Sleep(5 * time.Second)
			go this_.connNodeListener(pool, connAddress, connToken, clientIndex)
		}
	})
	pool.Put(messageListener)
	Logger.Info(this_.GetServerInfo() + " 连接 [" + toNodeId + "][" + connAddress + "] 成功 现有连接 " + fmt.Sprint(len(pool.listeners)))

	return
}
