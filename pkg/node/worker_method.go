package node

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"teamide/pkg/filework"
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

	methodFileExist     MethodType = 301
	methodFileFile      MethodType = 302
	methodFileFiles     MethodType = 303
	methodFileCreate    MethodType = 304
	methodFileRemove    MethodType = 305
	methodFileRename    MethodType = 306
	methodFileMove      MethodType = 307
	methodFileWrite     MethodType = 308
	methodFileRead      MethodType = 309
	methodFileCount     MethodType = 310
	methodFileCountSize MethodType = 311

	methodTerminalStart      MethodType = 401
	methodTerminalWrite      MethodType = 402
	methodTerminalChangeSize MethodType = 403
	methodTerminalStop       MethodType = 404
	methodTerminalIsWindows  MethodType = 405

	methodSystemGetInfo          MethodType = 501
	methodSystemQueryMonitorData MethodType = 502
	methodSystemCleanMonitorData MethodType = 503

	methodSendBytesStart MethodType = 601
	methodSendBytes      MethodType = 602
	methodSendBytesEnd   MethodType = 603
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
		if e := recover(); e != nil {
			Logger.Error("call error", zap.Any("error", e))
		}
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
		defer func() {
			if e := recover(); e != nil {
				Logger.Error("call wait result error", zap.Any("error", e))
			}
		}()
		time.Sleep(time.Second * 60 * 1)
		if isEnd || waitResult == nil {
			return
		}
		waitResult <- &Message{
			Error: fmt.Sprintf("请求超时，超时时间%d秒", 5),
		}
	}()
	res := <-waitResult

	close(waitResult)
	waitResult = nil

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

	case methodFileExist:
		if msg.FileWorkData != nil {
			var exist bool
			exist, err = this_.workExist(msg.LineNodeIdList, msg.FileWorkData.Path)
			if err != nil {
				return
			}
			res.FileWorkData = &FileWorkData{
				Exist: exist,
			}
		}
		return
	case methodFileFile:
		if msg.FileWorkData != nil {
			var file *filework.FileInfo
			file, err = this_.workFile(msg.LineNodeIdList, msg.FileWorkData.Path)
			if err != nil {
				return
			}
			res.FileWorkData = &FileWorkData{
				File: file,
			}
		}
		return
	case methodFileFiles:
		if msg.FileWorkData != nil {
			var fileList []*filework.FileInfo
			var path string
			path, fileList, err = this_.workFiles(msg.LineNodeIdList, msg.FileWorkData.Dir)
			if err != nil {
				return
			}
			res.FileWorkData = &FileWorkData{
				FileList: fileList,
				Path:     path,
			}
		}
		return
	case methodFileCreate:
		if msg.FileWorkData != nil {
			err = this_.workFileCreate(msg.LineNodeIdList, msg.FileWorkData.Path, msg.FileWorkData.IsDir)
			if err != nil {
				return
			}
		}
		return
	case methodFileRename:
		if msg.FileWorkData != nil {
			err = this_.workFileRename(msg.LineNodeIdList, msg.FileWorkData.OldPath, msg.FileWorkData.NewPath)
			if err != nil {
				return
			}
		}
		return
	case methodFileMove:
		if msg.FileWorkData != nil {
			err = this_.workFileMove(msg.LineNodeIdList, msg.FileWorkData.OldPath, msg.FileWorkData.NewPath)
			if err != nil {
				return
			}
		}
		return
	case methodFileRemove:
		if msg.FileWorkData != nil {
			var fileCount int
			var removeCount int
			fileCount, removeCount, err = this_.workFileRemove(msg.LineNodeIdList, msg.FileWorkData.Path)
			if err != nil {
				return
			}
			res.FileWorkData = &FileWorkData{
				FileCount:   fileCount,
				RemoveCount: removeCount,
			}
		}
		return
	case methodFileRead:
		if msg.FileWorkData != nil {
			err = this_.workFileRead(msg.LineNodeIdList, msg.FileWorkData.Path, msg.SendKey)
			if err != nil {
				return
			}
		}
		return
	case methodFileWrite:
		if msg.FileWorkData != nil {
			var sendKey string
			sendKey, err = this_.workFileWrite(msg.LineNodeIdList, msg.FileWorkData.Path)
			if err != nil {
				return
			}
			res.SendKey = sendKey
		}
		return
	case methodFileCount:
		return
	case methodFileCountSize:
		return

	case methodTerminalStart:
		if msg.TerminalWorkData != nil {
			var key string
			key, err = this_.workTerminalStart(msg.LineNodeIdList, msg.TerminalWorkData.Size, msg.TerminalWorkData.ReadKey, msg.TerminalWorkData.ReadErrorKey)
			if err != nil {
				return
			}
			res.TerminalWorkData = &TerminalWorkData{
				Key: key,
			}
		}
		return
	case methodTerminalIsWindows:
		if msg.TerminalWorkData != nil {
			var isWindows bool
			isWindows, err = this_.workTerminalIsWindows(msg.LineNodeIdList)
			if err != nil {
				return
			}
			res.TerminalWorkData = &TerminalWorkData{
				IsWindows: isWindows,
			}
		}
		return
	case methodTerminalStop:
		if msg.TerminalWorkData != nil {
			err = this_.workTerminalStop(msg.LineNodeIdList, msg.TerminalWorkData.Key)
			if err != nil {
				return
			}
		}
		return
	case methodTerminalWrite:
		if msg.TerminalWorkData != nil {
			err = this_.workTerminalWrite(msg.LineNodeIdList, msg.TerminalWorkData.Key, msg.Bytes)
			if err != nil {
				return
			}
		}
		return
	case methodTerminalChangeSize:
		if msg.TerminalWorkData != nil {
			err = this_.workTerminalChangeSize(msg.LineNodeIdList, msg.TerminalWorkData.Key, msg.TerminalWorkData.Size)
			if err != nil {
				return
			}
		}
		return

	case methodSendBytesStart:
		err = this_.workSendBytesStart(msg.LineNodeIdList, msg.SendKey)
		if err != nil {
			return
		}
		return
	case methodSendBytes:
		err = this_.workSendBytes(msg.LineNodeIdList, msg.SendKey, msg.Bytes)
		if err != nil {
			return
		}
		return
	case methodSendBytesEnd:
		err = this_.workSendBytesEnd(msg.LineNodeIdList, msg.SendKey)
		if err != nil {
			return
		}
		return
	}

	return
}
