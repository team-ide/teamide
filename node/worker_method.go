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

	methodGetNode            = 11
	methodGetNodeMonitorData = 12

	methodNetProxyNewConn             = 21
	methodNetProxyCloseConn           = 22
	methodNetProxySend                = 23
	methodNetProxyGetInnerMonitorData = 24
	methodNetProxyGetOuterMonitorData = 25

	methodFileFiles         = 31
	methodFileCopy          = 32
	methodFileRemove        = 33
	methodFileRename        = 34
	methodFileUpload        = 35
	methodFileConfirmResult = 36
	methodFileProgress      = 37
	methodFileRead          = 38

	methodShellNewConn   = 41
	methodShellCloseConn = 42
	methodShellSend      = 43
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
	if msg.NotifyChange != nil {
		if msg.NotifyChange.NotifyAll {
			this_.notifyAll(msg)
		} else {
			if msg.NotifyChange.NotifyChildren {
				this_.notifyChildren(msg)
			}
			if msg.NotifyChange.NotifyParent {
				this_.notifyParent(msg)
			}
		}
	}
	switch method {
	case methodOK:
		return
	case methodGetVersion:
		if msg.NodeWorkData != nil {
			version := this_.getVersion(msg.NodeWorkData.NodeId, msg.NotifiedNodeIdList)
			if version != "" {
				res.NodeWorkData = &NodeWorkData{
					Version: version,
				}
			}
		}
		return
	case methodGetNode:
		if msg.NodeWorkData != nil {
			node := this_.getNode(msg.NodeWorkData.NodeId, msg.NotifiedNodeIdList)
			if node != nil {
				res.NodeWorkData = &NodeWorkData{
					Node: node,
				}
			}
		}
		return
	case methodGetNodeMonitorData:
		if msg.NodeWorkData != nil {
			monitorData := this_.getNodeMonitorData(msg.NodeWorkData.NodeId, msg.NotifiedNodeIdList)
			if monitorData != nil {
				res.NodeWorkData = &NodeWorkData{
					MonitorData: monitorData,
				}
			}
		}
		return
	case methodNetProxyNewConn:
		if msg.NetProxyWorkData != nil {
			err = this_.netProxyNewConn(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId, msg.NetProxyWorkData.ConnId)
		}
		return
	case methodNetProxyCloseConn:
		if msg.NetProxyWorkData != nil {
			err = this_.netProxyCloseConn(msg.NetProxyWorkData.IsReverse, msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId, msg.NetProxyWorkData.ConnId)
		}
		return
	case methodNetProxySend:
		if msg.NetProxyWorkData != nil {
			err = this_.netProxySend(msg.NetProxyWorkData.IsReverse, msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId, msg.NetProxyWorkData.ConnId, msg.Bytes)
		}
	case methodNetProxyGetInnerMonitorData:
		if msg.NetProxyWorkData != nil {
			monitorData := this_.getNetProxyInnerMonitorData(msg.NetProxyWorkData.NetProxyId, msg.NotifiedNodeIdList)
			if monitorData != nil {
				res.NetProxyWorkData = &NetProxyWorkData{
					MonitorData: monitorData,
				}
			}
		}
		return
	case methodNetProxyGetOuterMonitorData:
		if msg.NetProxyWorkData != nil {
			monitorData := this_.getNetProxyOuterMonitorData(msg.NetProxyWorkData.NetProxyId, msg.NotifiedNodeIdList)
			if monitorData != nil {
				res.NetProxyWorkData = &NetProxyWorkData{
					MonitorData: monitorData,
				}
			}
		}
		return
	}

	return
}
