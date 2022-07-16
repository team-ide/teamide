package node

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

var (
	methodOK            = 1
	methodAddNode       = 2
	methodRemoveNode    = 3
	methodResetNodeList = 4
	methodProxySend     = 5
)

func (this_ *Worker) onMessage(msg *Message) {
	if msg == nil {
		return
	}
	callback, ok := this_.getCallback(msg.Id)
	if ok {
		callback(msg)
	} else {
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

func (this_ *Worker) getCallback(id string) (callback func(msg *Message), ok bool) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()
	callback, ok = this_.callbackCache[id]
	return
}

func (this_ *Worker) setCallback(id string, callback func(msg *Message)) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	this_.callbackCache[id] = callback
}

func (this_ *Worker) removeCallback(id string) {
	this_.callbackCacheLock.Lock()
	defer this_.callbackCacheLock.Unlock()

	delete(this_.callbackCache, id)
}

func (this_ *Worker) Call(node *Info, listener *MessageListener, method int, msg *Message) (result *Message, err error) {
	msg.Id = uuid.NewString()
	msg.Method = method
	msg.Token = node.Token

	defer func() {
		this_.removeCallback(msg.Id)
	}()

	waitResult := make(chan *Message, 1)
	this_.setCallback(msg.Id, func(msg *Message) {
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
	case methodAddNode:
		if msg.Node != nil {
			this_.AddNode(msg.Node)
		}
		return
	case methodRemoveNode:
		if msg.Node != nil {
			this_.RemoveNode(msg.Node)
		}
		return
	case methodResetNodeList:
		this_.resetNodeList(msg.NodeList)
		return
	}

	return
}
