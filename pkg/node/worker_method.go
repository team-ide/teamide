package node

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

var (
	methodOK         MethodType = 1
	methodGetVersion MethodType = 2

	methodNodeAddToNodeList      MethodType = 101
	methodNodeRemoveToNodeList   MethodType = 102
	methodNodeGetNodeMonitorData MethodType = 103
	methodNodeGetStatus          MethodType = 104

	methodNetProxyNewConn                 MethodType = 201
	methodNetProxyCloseConn               MethodType = 202
	methodNetProxySend                    MethodType = 203
	methodNetProxyGetInnerMonitorData     MethodType = 204
	methodNetProxyGetOuterMonitorData     MethodType = 205
	methodNetProxyAddNetProxyInnerList    MethodType = 206
	methodNetProxyAddNetProxyOuterList    MethodType = 207
	methodNetProxyRemoveNetProxyInnerList MethodType = 208
	methodNetProxyRemoveNetProxyOuterList MethodType = 209
	methodNetProxyGetInnerStatus          MethodType = 210
	methodNetProxyGetOuterStatus          MethodType = 211

	methodFileFiles         MethodType = 301
	methodFileCopy          MethodType = 302
	methodFileRemove        MethodType = 303
	methodFileRename        MethodType = 304
	methodFileUpload        MethodType = 305
	methodFileConfirmResult MethodType = 306
	methodFileProgress      MethodType = 307
	methodFileRead          MethodType = 308

	methodTerminalNewConn   MethodType = 401
	methodTerminalCloseConn MethodType = 402
	methodTerminalSend      MethodType = 403

	methodSystemGetInfo          MethodType = 501
	methodSystemQueryMonitorData MethodType = 502
	methodSystemCleanMonitorData MethodType = 503
)

type MethodType int

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

func (this_ *Worker) Call(listener *MessageListener, method MethodType, msg *Message) (result *Message, err error) {
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

func (this_ *Worker) doMethod(method MethodType, msg *Message) (res *Message, err error) {
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

	case methodNodeGetNodeMonitorData:
		monitorData := this_.getNodeMonitorData(msg.LineNodeIdList)
		if monitorData != nil {
			res.NodeWorkData = &WorkData{
				MonitorData: monitorData,
			}
		}
		return
	case methodNodeGetStatus:
		status := this_.getNodeStatus(msg.LineNodeIdList)
		res.NodeWorkData = &WorkData{
			Status: status,
		}
		return
	case methodNodeAddToNodeList:
		if msg.NodeWorkData != nil {
			this_.addToNodeList(msg.LineNodeIdList, msg.NodeWorkData.ToNodeList)
		}
		return
	case methodNodeRemoveToNodeList:
		if msg.NodeWorkData != nil {
			this_.removeToNodeList(msg.LineNodeIdList, msg.NodeWorkData.ToNodeIdList)
		}
		return

	case methodNetProxyAddNetProxyInnerList:
		if msg.NetProxyWorkData != nil {
			this_.addNetProxyInnerList(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyInnerList)
		}
		return
	case methodNetProxyRemoveNetProxyInnerList:
		if msg.NetProxyWorkData != nil {
			this_.removeNetProxyInnerList(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyIdList)
		}
		return
	case methodNetProxyAddNetProxyOuterList:
		if msg.NetProxyWorkData != nil {
			this_.addNetProxyOuterList(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyOuterList)
		}
		return
	case methodNetProxyRemoveNetProxyOuterList:
		if msg.NetProxyWorkData != nil {
			this_.removeNetProxyOuterList(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyIdList)
		}
		return
	case methodNetProxyGetInnerStatus:
		if msg.NetProxyWorkData != nil {
			status := this_.getNetProxyInnerStatus(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId)
			res.NetProxyWorkData = &NetProxyWorkData{
				Status: status,
			}
		}
		return
	case methodNetProxyGetOuterStatus:
		if msg.NetProxyWorkData != nil {
			status := this_.getNetProxyOuterStatus(msg.LineNodeIdList, msg.NetProxyWorkData.NetProxyId)
			res.NetProxyWorkData = &NetProxyWorkData{
				Status: status,
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

	case methodFileFiles:
		return
	case methodFileCopy:
		return
	case methodFileProgress:
		return
	case methodFileRemove:
		return
	case methodFileRead:
		return
	case methodFileConfirmResult:
		return
	case methodFileRename:
		return
	case methodFileUpload:
		return

	case methodTerminalNewConn:
		return
	case methodTerminalCloseConn:
		return
	case methodTerminalSend:
		return
	}

	return
}
