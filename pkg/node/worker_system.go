package node

import "teamide/pkg/system"

func (this_ *Worker) systemQueryMonitorData(lineNodeIdList []string, request *SystemData) (response *SystemData) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemQueryMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
			SystemData: &SystemData{
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

	response = &SystemData{
		QueryResponse: system.QueryMonitorData(request.QueryRequest),
	}
	return
}

func (this_ *Worker) systemCleanMonitorData(lineNodeIdList []string) (response *SystemData) {

	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemCleanMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
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

	system.CleanMonitorData()
	response = &SystemData{}

	return
}

func (this_ *Worker) systemGetInfo(lineNodeIdList []string) (response *SystemData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemGetInfo, &Message{
			LineNodeIdList: lineNodeIdList,
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

	response = &SystemData{
		Info: system.GetInfo(),
	}
	return

}

func (this_ *Worker) systemMonitorData(lineNodeIdList []string) (response *SystemData) {
	var resMsg *Message
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		resMsg, e = this_.Call(listener, methodSystemMonitorData, &Message{
			LineNodeIdList: lineNodeIdList,
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

	response = &SystemData{}
	response.MonitorData, err = system.GetCacheOrNew()
	return

}
