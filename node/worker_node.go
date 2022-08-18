package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"teamide/pkg/util"
	"time"
)

func (this_ *Worker) doAddToNodeList(toNodeList []*ToNode) (err error) {
	if len(toNodeList) == 0 {
		return
	}
	this_.toNodeListLock.Lock()
	defer this_.toNodeListLock.Unlock()

	Logger.Info(this_.server.GetServerInfo()+" 新增节点 ", zap.Any("toNodeList", toNodeList))

	for _, toNode := range toNodeList {
		var find = this_.findToNode(toNode.Id)

		if find == nil {
			Logger.Info(this_.server.GetServerInfo()+" 添加节点 ", zap.Any("toNode", toNode))
			this_.toNodeList = append(this_.toNodeList, toNode)

			this_.toNodeListenerKeepAlive(toNode.Id, toNode.ConnAddress, toNode.ConnToken, toNode.ConnSize)
		} else {
			var hasChange bool
			Logger.Info(this_.server.GetServerInfo()+" 更新节点 ", zap.Any("toNode", toNode))
			if toNode.Enabled != 0 {
				if toNode.IsEnabled() != find.IsEnabled() {
					hasChange = true
				}
				find.Enabled = toNode.Enabled
			}
			if toNode.ConnAddress != find.ConnAddress {
				find.ConnAddress = toNode.ConnAddress
				hasChange = true
			}
			if toNode.ConnToken != find.ConnToken {
				find.ConnToken = toNode.ConnToken
				hasChange = true
			}
			if toNode.ConnSize != 0 && toNode.ConnSize != find.ConnSize {
				find.ConnSize = toNode.ConnSize
				hasChange = true
			}
			if hasChange {
				this_.removeToNodeListenerPool(toNode.Id)
				if find.IsEnabled() {
					this_.toNodeListenerKeepAlive(find.Id, find.ConnAddress, find.ConnToken, find.ConnSize)
				}
			}
		}
	}

	return
}

func (this_ *Worker) doRemoveToNodeList(removeToNodeIdList []string) (err error) {
	if len(removeToNodeIdList) == 0 {
		return
	}
	this_.toNodeListLock.Lock()
	defer this_.toNodeListLock.Unlock()

	Logger.Info(this_.server.GetServerInfo()+" 移除节点 ", zap.Any("removeToNodeIdList", removeToNodeIdList))

	var list = this_.toNodeList
	var newList []*ToNode
	for _, nodeId := range removeToNodeIdList {
		this_.removeToNodeListenerPool(nodeId)
	}
	for _, one := range list {
		if util.ContainsString(removeToNodeIdList, one.Id) >= 0 {
		} else {
			newList = append(newList, one)
		}
	}
	this_.toNodeList = newList

	return
}

func (this_ *Worker) toNodeListenerKeepAlive(toNodeId string, connAddress, connToken string, connSize int) {
	if connAddress == "" {
		Logger.Warn("连接 [" + toNodeId + "] [" + connAddress + "] 连接地址为空")
		return
	}

	this_.toNodeListenerKeepAliveLock.Lock()
	defer this_.toNodeListenerKeepAliveLock.Unlock()

	var pool = this_.getToNodeListenerPool(toNodeId)
	if pool != nil {
		return
	}
	pool = this_.getToNodeListenerPoolIfAbsentCreate(toNodeId)
	if connSize <= 0 {
		connSize = 5
	}
	for connIndex := 0; connIndex < connSize; connIndex++ {
		go this_.connNodeListener(pool, connAddress, connToken, connIndex)
	}
	return
}

func (this_ *Worker) connNodeListener(pool *MessageListenerPool, connAddress, connToken string, connIndex int) {
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
		go this_.connNodeListener(pool, connAddress, connToken, connIndex)
	}()
	var err error
	var conn net.Conn
	Logger.Info("连接 [" + connAddress + "] 开始")
	conn, err = net.Dial("tcp", GetAddress(connAddress))
	if err != nil {
		Logger.Warn("连接 ["+connAddress+"] 异常", zap.Any("error", err.Error()))
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
		Method: methodOK,
		ConnData: &ConnData{
			ConnIndex: connIndex,
			NodeId:    this_.server.Id,
		},
	}

	err = WriteMessage(conn, msg, this_.MonitorData)
	if err != nil {
		Logger.Warn("连接 [" + connAddress + "] 异常")
		_ = conn.Close()
		return
	}
	msg, err = ReadMessage(conn, this_.MonitorData)
	if err != nil {
		Logger.Warn("连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	if msg.ConnData == nil || msg.ConnData.NodeId == "" {
		Logger.Warn("连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	toNodeId := msg.ConnData.NodeId
	pool = this_.getToNodeListenerPoolIfAbsentCreate(toNodeId)
	Logger.Info("连接 [" + toNodeId + "] [" + connAddress + "] 成功")

	messageListener = &MessageListener{
		conn:      conn,
		onMessage: this_.onMessage,
	}

	messageListener.listen(func() {
		messageListener.stop()
		pool.Remove(messageListener)
		Logger.Info("移除 连接至 [" + toNodeId + "] [" + connAddress + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))

		if !pool.isStop {
			time.Sleep(5 * time.Second)
			go this_.connNodeListener(pool, connAddress, connToken, connIndex)
		}
	}, this_.MonitorData)
	size := pool.Put(messageListener)
	Logger.Info("连接 [" + toNodeId + "] [" + connAddress + "] 成功 现有连接 " + fmt.Sprint(size))

	return
}
