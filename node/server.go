package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

var tokenByteSize = 128

type Server struct {
	Id             string
	Address        string
	Token          string
	ConnAddress    string
	ConnToken      string
	ConnSize       int
	serverListener net.Listener
	cache          *Cache
	*Worker
	OnNodeListChange     func([]*Info)
	OnNetProxyListChange func([]*NetProxy)
}

func (this_ *Server) Start() (err error) {
	this_.cache = newCache()
	this_.Worker = &Worker{
		server: this_,
		cache:  this_.cache,
	}

	go this_.serverListenerKeepAlive()
	if this_.ConnAddress != "" {
		this_.connNodeListenerKeepAlive(this_.ConnAddress, this_.ConnToken, this_.ConnSize)
	}
	return
}

func (this_ *Server) Stop() {
	if this_.serverListener != nil {
		_ = this_.serverListener.Close()
	}
	if this_.Worker != nil {
		this_.Worker.Stop()
	}

}

func (this_ *Server) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.Address)
}

func (this_ *Server) serverListenerKeepAlive() {

	defer func() {
		// 删除所有连接
		var list = this_.cache.fromNodeIdList
		for _, fromNodeId := range list {
			this_.cache.removeFromNodeId(fromNodeId)
			this_.cache.removeNodeListenerPool(fromNodeId, this_.Id)
		}
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive()
	}()
	var err error
	Logger.Info(this_.GetServerInfo() + " 服务启动 开始")
	this_.serverListener, err = net.Listen("tcp", GetAddress(this_.Address))
	if err != nil {
		Logger.Error(this_.GetServerInfo()+" 服务启动 异常", zap.Error(err))
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
		go func(conn_ net.Conn) {
			var bytes = make([]byte, tokenByteSize)
			_, err = conn_.Read(bytes)
			if err != nil {
				_ = conn_.Close()
				return
			}
			token := strings.TrimSpace(string(bytes))
			if this_.Token != token {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 Token验证异常")
				_ = conn_.Close()
				return
			}

			var msg *Message
			msg, err = ReadMessage(conn_)
			if err != nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			if msg.Method != methodOK || msg.FromNodeId == "" {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var clientIndex = msg.ClientIndex

			// 发送当前节点ID
			err = WriteMessage(conn, &Message{
				FromNodeId: this_.Id,
			})
			if err != nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var fromNodeId = msg.FromNodeId

			this_.cache.addFromNodeId(fromNodeId)
			pool := this_.cache.newIfAbsentNodeListenerPool(fromNodeId, this_.Id)

			messageListener := &MessageListener{
				conn:      conn,
				onMessage: this_.onMessage,
			}
			messageListener.listen(func() {
				messageListener.stop()
				pool.Remove(messageListener)
				Logger.Info(this_.GetServerInfo() + " 移除 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
				if len(pool.listeners) == 0 {
					this_.cache.removeFromNodeId(fromNodeId)
					this_.cache.removeNodeListenerPool(fromNodeId, this_.Id)
				}
			})
			pool.Put(messageListener)
			Logger.Info(this_.GetServerInfo() + " 添加 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))

			if clientIndex == 0 {
				err = messageListener.Send(&Message{
					Method:       methodNotifyParentRefresh,
					NodeList:     this_.cache.nodeList,
					NetProxyList: this_.cache.netProxyList,
				})
				if err != nil {
					Logger.Error(this_.GetServerInfo()+" 通知 ["+fromNodeId+"] 刷新 异常 ", zap.Error(err))
				}
			}

		}(conn)
	}
	return
}

func (this_ *Server) connNodeListenerKeepAlive(connAddress, connToken string, connSize int) {
	if connAddress == "" {
		Logger.Warn(this_.GetServerInfo() + " 连接 [" + connAddress + "] 连接地址为空")
		return
	}

	if connSize <= 0 {
		connSize = 5
	}

	for i := 0; i < connSize; i++ {
		_ = this_.connNodeListener(connAddress, connToken, i)
	}
	return
}

func (this_ *Server) connNodeListener(connAddress, connToken string, clientIndex int) (pool *MessageListenerPool) {
	defer func() {
		if pool != nil {
			return
		}
		time.Sleep(5 * time.Second)
		go this_.connNodeListener(connAddress, connToken, clientIndex)
	}()
	var err error
	var conn net.Conn
	Logger.Info(this_.GetServerInfo() + " 连接 [" + connAddress + "] 开始")
	conn, err = net.Dial("tcp", GetAddress(connAddress))
	if err != nil {
		Logger.Error(this_.GetServerInfo()+" 连接 ["+connAddress+"] 异常", zap.Error(err))
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
		FromNodeId:  this_.Id,
		Ok:          true,
		ClientIndex: clientIndex,
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
	if msg.FromNodeId == "" {
		Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}

	Logger.Info(this_.GetServerInfo() + " 连接 [" + connAddress + "] 成功")

	toNodeId := msg.FromNodeId
	var fromNodeId = this_.Id

	pool = this_.cache.newIfAbsentNodeListenerPool(fromNodeId, toNodeId)

	messageListener := &MessageListener{
		conn:      conn,
		onMessage: this_.onMessage,
	}

	messageListener.listen(func() {
		messageListener.stop()
		pool.Remove(messageListener)
		Logger.Info(this_.GetServerInfo() + " 移除 连接至 [" + toNodeId + "][" + connAddress + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
		if !pool.isStop {
			time.Sleep(5 * time.Second)
			go this_.connNodeListener(connAddress, connToken, clientIndex)
		}
	})
	pool.Put(messageListener)
	Logger.Info(this_.GetServerInfo() + " 连接 [" + toNodeId + "][" + connAddress + "] 成功 现有连接 " + fmt.Sprint(len(pool.listeners)))

	if clientIndex == 0 {
		var toNode = &Info{
			Id:          toNodeId,
			ParentId:    fromNodeId,
			ConnToken:   connToken,
			ConnAddress: connAddress,
		}
		_ = this_.addNodeList([]*Info{toNode})

		err = messageListener.Send(&Message{
			Method:       methodNotifyParentRefresh,
			NodeList:     this_.cache.nodeList,
			NetProxyList: this_.cache.netProxyList,
		})
		if err != nil {
			Logger.Error(this_.GetServerInfo()+" 通知 ["+toNodeId+"]["+connAddress+"] 刷新 异常 ", zap.Error(err))
		}
	}
	return
}
