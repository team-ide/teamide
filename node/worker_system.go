package node

import "teamide/pkg/system"

func (this_ *Worker) systemQueryMonitorData(lineNodeIdList []string, nodeId string, request *SystemData) (response *SystemData) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemQueryMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			SystemData: &SystemData{
				NodeId:       nodeId,
				QueryRequest: request.QueryRequest,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil {
			response = resMsg.SystemData
		}
		return
	}

	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		response = &SystemData{
			QueryResponse: system.QueryMonitorData(request.QueryRequest),
		}
		return
	}

	return
}

func (this_ *Worker) systemCleanMonitorData(lineNodeIdList []string, nodeId string) (response *SystemData) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemCleanMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			SystemData: &SystemData{
				NodeId: nodeId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil {
			response = resMsg.SystemData
		}
		return
	}

	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		system.CleanMonitorData()
		response = &SystemData{}
		return
	}

	return
}

func (this_ *Worker) systemGetInfo(lineNodeIdList []string, nodeId string) (response *SystemData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemGetInfo, &Message{
			LineNodeIdList: lineNodeIdList,
			SystemData: &SystemData{
				NodeId: nodeId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		if resMsg != nil {
			response = resMsg.SystemData
		}
		return
	}

	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		response = &SystemData{
			Info: system.GetInfo(),
		}
		return
	}

	return
}
