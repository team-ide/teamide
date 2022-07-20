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
	Id                   string
	BindAddress          string
	BindToken            string
	ConnAddress          string
	ConnToken            string
	ConnSize             int
	serverListener       net.Listener
	cache                *Cache
	worker               *Worker
	OnNodeListChange     func([]*Info)
	OnNetProxyListChange func([]*NetProxy)
	rootNode             *Info
}

func (this_ *Server) AddNodeList(nodeList []*Info) (err error) {
	err = this_.worker.addNodeList(nodeList, []string{})
	return
}

func (this_ *Server) RemoveNodeList(nodeIdList []string) (err error) {
	err = this_.worker.removeNodeList(nodeIdList, []string{})
	return
}

func (this_ *Server) GetNodeLineByFromTo(fromNodeId, toNodeId string) (lineIdList []string) {

	return this_.worker.getNodeLineByFromTo(fromNodeId, toNodeId)
}

func (this_ *Server) AddNetProxyList(netProxyList []*NetProxy) (err error) {
	err = this_.worker.addNetProxyList(netProxyList, []string{})
	return
}

func (this_ *Server) RemoveNetProxyList(netProxyIdList []string) (err error) {
	err = this_.worker.removeNetProxyList(netProxyIdList, []string{})
	return
}

func (this_ *Server) Start() (err error) {
	this_.rootNode = &Info{
		Id:          this_.Id,
		BindAddress: this_.BindAddress,
		BindToken:   this_.BindToken,
		ConnAddress: this_.ConnAddress,
		ConnToken:   this_.ConnToken,
	}
	this_.cache = newCache()
	this_.worker = &Worker{
		server: this_,
		cache:  this_.cache,
	}

	_ = this_.worker.doAddNodeList([]*Info{this_.rootNode})

	if this_.BindAddress != "" {
		go this_.serverListenerKeepAlive()
	}
	if this_.ConnAddress != "" {
		this_.connNodeListenerKeepAlive(this_.ConnAddress, this_.ConnToken, this_.ConnSize)
	}
	return
}

func (this_ *Server) Stop() {
	if this_.serverListener != nil {
		_ = this_.serverListener.Close()
	}
	if this_.worker != nil {
		this_.worker.Stop()
	}

}

func (this_ *Server) GetServerInfo() (str string) {
	return fmt.Sprintf("节点服务[%s][%s]", this_.Id, this_.BindAddress)
}

func (this_ *Server) serverListenerKeepAlive() {

	defer func() {
		// 删除所有连接
		var list = this_.cache.getNodeListenerPoolListByToNodeId(this_.Id)
		for _, one := range list {
			this_.cache.removeNodeListenerPool(one.fromNodeId, one.toNodeId)
		}
		time.Sleep(5 * time.Second)
		go this_.serverListenerKeepAlive()
	}()
	var err error
	Logger.Info(this_.GetServerInfo() + " 服务启动 开始")
	this_.serverListener, err = net.Listen("tcp", GetAddress(this_.BindAddress))
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
			if this_.BindToken != token {
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
			if msg.Method != methodOK || msg.Node == nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var clientNode = msg.Node
			var clientIndex = msg.ClientIndex

			// 发送当前节点ID
			err = WriteMessage(conn, &Message{
				Node: this_.rootNode,
			})
			if err != nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var fromNodeId = clientNode.Id

			pool := this_.cache.newIfAbsentNodeListenerPool(fromNodeId, this_.Id)

			messageListener := &MessageListener{
				conn:      conn,
				onMessage: this_.worker.onMessage,
			}
			messageListener.listen(func() {
				messageListener.stop()
				pool.Remove(messageListener)
				Logger.Info(this_.GetServerInfo() + " 移除 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))
				if len(pool.listeners) == 0 {
					this_.cache.removeNodeListenerPool(fromNodeId, this_.Id)
				}
			})
			pool.Put(messageListener)
			Logger.Info(this_.GetServerInfo() + " 添加 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(len(pool.listeners)))

			if clientIndex == 0 {
				clientNode.addConnNodeId(this_.Id)
				//_ = this_.worker.addNodeList([]*Info{clientNode})

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
		Node:        this_.rootNode,
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
	if msg.Node == nil {
		Logger.Error(this_.GetServerInfo() + " 连接 [" + connAddress + "] 接口异常")
		_ = conn.Close()
		return
	}
	serverNode := msg.Node
	Logger.Info(this_.GetServerInfo() + " 连接 [" + connAddress + "] 成功")

	toNodeId := serverNode.Id
	var fromNodeId = this_.Id

	pool = this_.cache.newIfAbsentNodeListenerPool(fromNodeId, toNodeId)

	messageListener := &MessageListener{
		conn:      conn,
		onMessage: this_.worker.onMessage,
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
		this_.rootNode.addConnNodeId(this_.Id)
		//_ = this_.worker.addNodeList([]*Info{serverNode})

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
