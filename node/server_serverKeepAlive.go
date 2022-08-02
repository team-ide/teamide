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
		this_.worker.notifyAll(&Message{
			NotifyChange: &NotifyChange{
				NodeStatusChangeList: []*StatusChange{
					{
						Id:          this_.rootNode.Id,
						Status:      StatusStopped,
						StatusError: "",
					},
				},
			},
		})
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
		this_.worker.notifyAll(&Message{
			NotifyChange: &NotifyChange{
				NodeStatusChangeList: []*StatusChange{
					{
						Id:          this_.rootNode.Id,
						Status:      StatusError,
						StatusError: err.Error(),
					},
				},
			},
		})
		Logger.Error(this_.GetServerInfo()+" 服务启动 异常", zap.Any("error", err.Error()))
		return
	}
	Logger.Info(this_.GetServerInfo() + " 服务启动 成功")

	_ = this_.worker.doChangeNodeStatus([]*StatusChange{
		{
			Id:          this_.rootNode.Id,
			Status:      StatusStarted,
			StatusError: "",
		},
	})
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

			var clientMsg *Message
			clientMsg, err = ReadMessage(conn_, this_.worker.MonitorData)
			if err != nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			if clientMsg.Method != methodOK || clientMsg.ClientData == nil || clientMsg.ClientData.Node == nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var clientNode = clientMsg.ClientData.Node
			var clientIndex = clientMsg.ClientData.Index

			// 发送当前节点ID
			err = WriteMessage(conn, &Message{
				ClientData: &ClientData{
					Node: this_.rootNode,
				},
			}, this_.worker.MonitorData)
			if err != nil {
				Logger.Error(this_.GetServerInfo() + " 来之客户端连接 接口异常")
				_ = conn_.Close()
				return
			}
			var fromNodeId = clientNode.Id

			pool := this_.cache.newIfAbsentNodeListenerPool(fromNodeId, this_.Id)
			if pool != nil && pool.isStop {
				return
			}
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

					var notifyMsg = &Message{
						NotifyChange: &NotifyChange{
							NodeStatusChangeList: []*StatusChange{
								{
									Id:          clientNode.Id,
									Status:      StatusStopped,
									StatusError: "",
								},
							},
						},
					}
					this_.worker.notifyAll(notifyMsg)
				}
			}, this_.worker.MonitorData)
			size := pool.Put(messageListener)
			Logger.Info(this_.GetServerInfo() + " 添加 来至 [" + fromNodeId + "] 节点的连接 现有连接 " + fmt.Sprint(size))

			if clientIndex == 0 {
				if clientMsg.ClientData != nil {
					for _, one := range clientMsg.ClientData.NodeList {
						if one.Id == clientNode.Id {
							one.addConnNodeId(this_.Id)
							one.Status = StatusStarted
							one.StatusError = ""
						} else if one.Id == this_.rootNode.Id {
							one.Status = StatusStarted
							one.StatusError = ""
						}
					}
					_ = this_.worker.doAddNodeList(clientMsg.ClientData.NodeList)
					_ = this_.worker.doAddNetProxyList(clientMsg.ClientData.NetProxyList)
				}

				this_.worker.notifyAll(&Message{
					NotifyChange: &NotifyChange{
						NodeStatusChangeList: []*StatusChange{
							{
								Id:          clientNode.Id,
								Status:      StatusStarted,
								StatusError: "",
							},
						},
						NodeList:     this_.cache.nodeList,
						NetProxyList: this_.cache.netProxyList,
					},
				})
			}

		}(conn)
	}
	return
}
