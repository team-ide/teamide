package node

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

func (this_ *Worker) addChildrenNode(childrenNode *Info) {
	var find = this_.findChildrenNode(childrenNode.Id)
	if find == nil {
		_ = this_.getChildrenNodeListenerPool(childrenNode)
		this_.childrenNodeList = append(this_.childrenNodeList, childrenNode)
	} else {
		copyNode(childrenNode, find)
	}
}

func (this_ *Worker) getChildrenNodeListenerPool(childrenNode *Info) (pool *MessageListenerPool) {
	this_.childrenNodeListenerLock.Lock()
	defer this_.childrenNodeListenerLock.Unlock()

	pool, ok := this_.childrenNodeListenerPoolCache[childrenNode.Id]
	if !ok {
		pool = &MessageListenerPool{}
		this_.childrenNodeListenerPoolCache[childrenNode.Id] = pool

		for i := 0; i < childrenNode.GetConnSize(); i++ {
			listener := &MessageListener{
				onMessage: this_.onMessage,
				isClose:   true,
			}
			this_.messageListenerKeepAlive(childrenNode, listener, i == 0)
			pool.Put(listener)
		}

	}
	return
}

func (this_ *Worker) removeChildrenNodeListener(childrenNode *Info) {
	this_.childrenNodeListenerLock.Lock()
	defer this_.childrenNodeListenerLock.Unlock()

	pool, ok := this_.childrenNodeListenerPoolCache[childrenNode.Id]
	if ok {
		pool.Stop()
		delete(this_.childrenNodeListenerPoolCache, childrenNode.Id)
	}
	return
}

func (this_ *Worker) messageListenerKeepAlive(node *Info, listener *MessageListener, isFirst bool) {

	if listener.isStop {
		return
	}
	if !listener.isClose {
		return
	}
	conn, err := net.Dial(node.GetNetwork(), node.GetAddress())
	if err != nil {
		node.Status = StatusError
		node.StatusError = err.Error()
		Logger.Error(this_.Node.GetNodeStr()+" dial "+node.Address+" error", zap.Error(err))
		go func() {
			time.Sleep(5 * time.Second)
			this_.messageListenerKeepAlive(node, listener, isFirst)
		}()
		return
	} else {
		var isEnd = false
		listener.conn = conn
		listener.listen(func() {
			isEnd = true
			if listener.isStop {
				return
			}
			this_.messageListenerKeepAlive(node, listener, isFirst)
		})

		_ = listener.Send(&Message{
			Token:      node.Token,
			Method:     methodOK,
			Ok:         true,
			FromNodeId: this_.Node.Id,
		})

		if isFirst {

			_, _ = this_.Call(node, listener, methodInitialize, &Message{
				NodeList:     this_.nodeList,
				NetProxyList: this_.netProxyList,
			})

			go func() {
				for {
					if isEnd || listener.isStop {
						return
					}
					if listener.isClose {
						node.Status = StatusStopped
						node.StatusError = ""
					} else {
						res, err := this_.Call(node, listener, methodOK, &Message{
							Ok: true,
						})
						if err != nil {
							node.Status = StatusError
							node.StatusError = err.Error()
							return
						} else {
							if !res.Ok {
								node.Status = StatusError
								node.StatusError = "服务节点验证失败"
							} else {
								node.Status = StatusStarted
								node.StatusError = ""
							}
						}
					}
					switch node.Status {
					case StatusStarted:
						//Logger.Info(fmt.Sprintf(this_.Node.GetNodeStr() + " 子节点 " + node.GetNodeStr() + " 验证成功"))
						break
					case StatusStopped:
						Logger.Info(fmt.Sprintf(this_.Node.GetNodeStr() + " 子节点 " + node.GetNodeStr() + " 未启动"))
						break
					case StatusError:
						Logger.Info(fmt.Sprintf(this_.Node.GetNodeStr()+" 子节点 "+node.GetNodeStr()+" 验证异常 [%s]", node.StatusError))
						break
					}
					if isEnd || listener.isStop {
						return
					}
					time.Sleep(5 * time.Second)
				}
			}()
		}

	}
}
