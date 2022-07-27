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
	methodGetVersion = 2

	methodGetNode = 11

	methodNetProxyNewConn   = 21
	methodNetProxyCloseConn = 22
	methodNetProxySend      = 23
)

func (this_ *Worker) onMessage(msg *Message) {
	if msg == nil {
		return
	}
	callback, ok := this_.cache.getCallback(msg.Id)
	if ok {
		callback(msg)
	} else {
		res, err := this_.doMethod(msg.Method, msg)
		if msg.Id != "" {
			if err != nil {
				err = msg.ReturnError(err.Error(), this_.MonitorData)
				if err != nil {
					Logger.Error("message return error", zap.Error(err))
					return
				}
				return
			}
			err = msg.Return(res, this_.MonitorData)
			if err != nil {
				Logger.Error("message return error", zap.Error(err))
				return
			}
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
	err = listener.Send(msg, this_.MonitorData)
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
	if msg.NotifyAll {
		this_.notifyAll(msg)
	} else {
		if msg.NotifyChildren {
			this_.notifyChildren(msg)
		}
		if msg.NotifyParent {
			this_.notifyParent(msg)
		}
	}
	switch method {
	case methodOK:
		res.Ok = true
		return
	case methodGetVersion:
		res.Version = this_.getVersion(msg.NodeId, msg.NotifiedNodeIdList)
		return
	case methodGetNode:
		res.Node = this_.getNode(msg.NodeId, msg.NotifiedNodeIdList)
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
