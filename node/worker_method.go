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

	methodNodeAddToNodeList      = 101
	methodNodeRemoveToNodeList   = 102
	methodNodeGetNodeMonitorData = 103
	methodNodeGetStatus          = 104

	methodNetProxyNewConn                 = 201
	methodNetProxyCloseConn               = 202
	methodNetProxySend                    = 203
	methodNetProxyGetInnerMonitorData     = 204
	methodNetProxyGetOuterMonitorData     = 205
	methodNetProxyAddNetProxyInnerList    = 206
	methodNetProxyAddNetProxyOuterList    = 207
	methodNetProxyRemoveNetProxyInnerList = 208
	methodNetProxyRemoveNetProxyOuterList = 209
	methodNetProxyGetInnerStatus          = 210
	methodNetProxyGetOuterStatus          = 211

	methodFileFiles         = 301
	methodFileCopy          = 302
	methodFileRemove        = 303
	methodFileRename        = 304
	methodFileUpload        = 305
	methodFileConfirmResult = 306
	methodFileProgress      = 307
	methodFileRead          = 308

	methodShellNewConn   = 401
	methodShellCloseConn = 402
	methodShellSend      = 403

	methodSystemGetInfo          = 501
	methodSystemQueryMonitorData = 502
	methodSystemCleanMonitorData = 503
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
		this_.removeCallback(msg.Id)
	}()

	waitResult := make(chan *Message, 1)
	this_.setCallback(msg.Id, func(msg *Message) {
		waitResult <- msg
	})
	err = listener.Send(msg, this_.MonitorData)
	if err != nil {
		return
	}
	var isEnd bool
	go func() {
		time.Sleep(time.Second * 60 * 1)
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
		return
	case methodGetVersion:
		if msg.NodeWorkData != nil {
			version := this_.getVersion(msg.LineNodeIdList)
			if version != "" {
				res.NodeWorkData = &WorkData{
					Version: version,
				}
			}
		}
		return
	case methodSystemGetInfo:
		response := this_.systemGetInfo(msg.LineNodeIdList)
		if response != nil {
			res.SystemData = response
		}
		return
	case methodSystemQueryMonitorData:
		if msg.SystemData != nil {
			response := this_.systemQueryMonitorData(msg.LineNodeIdList, msg.SystemData)
			if response != nil {
				res.SystemData = response
			}
		}
		return
	case methodSystemCleanMonitorData:
		response := this_.systemCleanMonitorData(msg.LineNodeIdList)
		if response != nil {
			res.SystemData = response
		}
		return
	case methodNodeGetNodeMonitorData:
		monitorData := this_.getNodeMonitorData(msg.LineNodeIdList)
		if monitorData != nil {
			res.NodeWorkData = &WorkData{
				MonitorData: monitorData,
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
			monitorData := this_.getNetProxyInnerMonitorData(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId)
			if monitorData != nil {
				res.NetProxyWorkData = &NetProxyWorkData{
					MonitorData: monitorData,
				}
			}
		}
		return
	case methodNetProxyGetOuterMonitorData:
		if msg.NetProxyWorkData != nil {
			monitorData := this_.getNetProxyOuterMonitorData(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId)
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
