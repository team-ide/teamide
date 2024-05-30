package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"strings"
	"sync"
	"time"
)

func (this_ *Server) serverListenerKeepAlive(localNode *LocalNode) {
	if localNode.IsStop {
		Logger.Info("本地节点停止", zap.Any("localNode", localNode))
		return
	}
	defer func() {
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive(localNode)
	}()
	var err error
	Logger.Info("本地节点 启动 开始", zap.Any("localNode", localNode))
	localNode.serverListener, err = net.Listen("tcp", GetAddress(localNode.BindAddress))
	if err != nil {
		Logger.Error("本地节点 启动 异常", zap.Any("localNode", localNode), zap.Error(err))
		return
	}
	Logger.Info("本地节点 启动 成功", zap.Any("localNode", localNode))

	var locker = &sync.Mutex{}
	for {
		var conn net.Conn
		conn, err = localNode.serverListener.Accept()
		if err != nil {
			Logger.Error("本地节点 监听 异常", zap.Any("localNode", localNode), zap.Error(err))
			break
		}
		_ = this_.onServerConn(locker, localNode, conn)
	}
	return
}

func (this_ *Server) onServerConn(locker sync.Locker, localNode *LocalNode, conn net.Conn) (err error) {
	locker.Lock()
	defer locker.Unlock()
	var bytes = make([]byte, tokenByteSize)
	_, err = conn.Read(bytes)
	if err != nil {
		_ = conn.Close()
		return
	}
	token := strings.TrimSpace(string(bytes))
	if localNode.BindToken != token {
		Logger.Error(localNode.GetServerInfo() + " 来之客户端连接 Token验证异常")
		_ = conn.Close()
		return
	}

	var clientMsg *Message
	clientMsg, err = ReadMessage(conn, this_.MonitorData)
	if err != nil {
		Logger.Error(localNode.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
	if clientMsg.ConnData == nil {
		Logger.Error(localNode.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
	var fromNodeIdList []string
	if clientMsg.ConnData.NodeId != "" {
		fromNodeIdList = append(fromNodeIdList, clientMsg.ConnData.NodeId)
	}
	for _, id := range clientMsg.ConnData.NodeIdList {
		fromNodeIdList = append(fromNodeIdList, id)
	}

	// 发送当前节点ID
	err = WriteMessage(conn, &Message{
		ConnData: &ConnData{
			NodeId:    localNode.Id,
			NodeToken: localNode.BindToken,
		},
	}, this_.MonitorData)
	if err != nil {
		Logger.Error(localNode.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
	for _, fromNodeId := range fromNodeIdList {
		pool := this_.getFromNodeListenerPoolIfAbsentCreate(fromNodeId)

		if pool != nil && pool.isStop {
			return
		}
		messageListener := &MessageListener{
			conn:      conn,
			onMessage: this_.onMessage,
		}
		messageListener.listen(func() {
			messageListener.stop()
			pool.Remove(messageListener)
			Logger.Info(localNode.GetServerInfo() + " 移除 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
			if len(pool.listeners) == 0 {
				this_.removeFromNodeListenerPool(fromNodeId)
			}
		}, this_.MonitorData)
		size := pool.Put(messageListener)
		Logger.Info(localNode.GetServerInfo() + " 添加 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(size))
	}

	return
}
