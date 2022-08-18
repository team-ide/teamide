package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

func (this_ *Server) serverListenerKeepAlive() {

	defer func() {
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive()
	}()
	var err error
	Logger.Info(this_.GetServerInfo() + " 服务启动 开始")
	this_.serverListener, err = net.Listen("tcp", GetAddress(this_.BindAddress))
	if err != nil {
		Logger.Error(this_.GetServerInfo()+" 服务启动 异常", zap.Any("error", err.Error()))
		return
	}
	Logger.Info(this_.GetServerInfo() + " 服务启动 成功")
	for {
		var conn net.Conn
		conn, err = this_.serverListener.Accept()
		if err != nil {
			Logger.Error(this_.GetServerInfo()+" 服务监听 异常", zap.Error(err))
			break
		}
		_ = this_.onServerConn(conn)
	}
	return
}

func (this_ *Server) onServerConn(conn net.Conn) (err error) {

	var bytes = make([]byte, tokenByteSize)
	_, err = conn.Read(bytes)
	if err != nil {
		_ = conn.Close()
		return
	}
	token := strings.TrimSpace(string(bytes))
	if this_.BindToken != token {
		Logger.Error(this_.GetServerInfo() + " 来之客户端连接 Token验证异常")
		_ = conn.Close()
		return
	}

	var clientMsg *Message
	clientMsg, err = ReadMessage(conn, this_.MonitorData)
	if err != nil {
		Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
	if clientMsg.ConnData == nil || clientMsg.ConnData.NodeId == "" {
		Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
	var fromNodeId = clientMsg.ConnData.NodeId

	// 发送当前节点ID
	err = WriteMessage(conn, &Message{
		ConnData: &ConnData{
			NodeId:    this_.Id,
			NodeToken: this_.BindToken,
		},
	}, this_.MonitorData)
	if err != nil {
		Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
		_ = conn.Close()
		return
	}
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
		Logger.Info(this_.GetServerInfo() + " 移除 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
		if len(pool.listeners) == 0 {
			this_.removeFromNodeListenerPool(fromNodeId)
		}
	}, this_.MonitorData)
	size := pool.Put(messageListener)
	Logger.Info(this_.GetServerInfo() + " 添加 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(size))

	return
}
