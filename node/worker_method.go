package node

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

var (
	methodOK         = 1
	methodInitialize = 2

	methodNodeAdd    = 11
	methodNodeRemove = 12

	methodNetProxyAdd       = 21
	methodNetProxyRemove    = 22
	methodNetProxyNewConn   = 23
	methodNetProxyCloseConn = 24
	methodNetProxySend      = 25
)

func (this_ *Worker) onMessage(msg *Message) {
	if msg == nil {
		return
	}
	callback, ok := this_.cache.getCallback(msg.Id)
	if ok {
		callback(msg)
	} else {
		if msg.Method == 0 {
			return
		}
		res, err := this_.doMethod(msg.Method, msg)
		if err != nil {
			err = msg.ReturnError(err.Error())
			if err != nil {
				Logger.Error("message return error", zap.Error(err))
				return
			}
			return
		}
		err = msg.Return(res)
		if err != nil {
			Logger.Error("message return error", zap.Error(err))
			return
		}
	}

}

func (this_ *Worker) Call(listener *MessageListener, method int, msg *Message) (result *Message, err error) {
	msg.Id = uuid.NewString()
	msg.Method = method

	defer func() {
		this_.cache.removeCallback(msg.Id)
	}()

	waitResult := make(chan *Message, 1)
	this_.cache.setCallback(msg.Id, func(msg *Message) {
		waitResult <- msg
	})
	err = listener.Send(msg)
	if err != nil {
		return
	}
	var isEnd bool
	go func() {
		time.Sleep(time.Second * 5)
		if isEnd {
			return
		}
		waitResult <- &Message{
			Error: fmt.Sprintf("请求超时，超时时间%d秒", 5),
		}
	}()
	res := <-waitResult
	isEnd = true
	if res.Error != "" {
		err = errors.New(res.Error)
		return
	}
	result = res
	return
}

func (this_ *Worker) doMethod(method int, msg *Message) (res *Message, err error) {
	if msg == nil {
		return
	}
	res = &Message{}

	switch method {
	case methodOK:
		res.Ok = true
		return
	case methodInitialize:
		this_.initialize(msg.NodeList, msg.NetProxyList)
		return
	case methodNodeAdd:
		if msg.Node != nil {
			_ = this_.AddNode(msg.Node)
		}
		return
	case methodNodeRemove:
		if msg.Node != nil {
			_ = this_.RemoveNode(msg.NodeId)
		}
		return
	case methodNetProxyAdd:
		if msg.NetProxy != nil {
			err = this_.AddNetProxy(msg.NetProxy)
		}
		return
	case methodNetProxyRemove:
		if msg.NetProxy != nil {
			err = this_.RemoveNetProxy(msg.NetProxyId)
		}
		return
	case methodNetProxyNewConn:
		err = this_.netProxyNewConn(msg.LineNodeIdList, msg.NetProxyId, msg.ConnId)
		return
	case methodNetProxyCloseConn:
		err = this_.netProxyCloseConn(msg.IsReverse, msg.LineNodeIdList, msg.NetProxyId, msg.ConnId)
		return
	case methodNetProxySend:
		err = this_.netProxySend(msg.IsReverse, msg.LineNodeIdList, msg.NetProxyId, msg.ConnId, msg.Bytes)
		return
	}

	return
}
