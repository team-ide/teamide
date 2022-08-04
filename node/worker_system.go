package node

import "teamide/pkg/system"

func (this_ *Worker) systemQueryMonitorData(nodeId string, NotifiedNodeIdList []string, request *SystemData) (response *SystemData) {
	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		response = &SystemData{
			QueryResponse: system.QueryMonitorData(request.QueryRequest),
		}
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if response == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				msg := &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					SystemData: &SystemData{
						NodeId:       nodeId,
						QueryRequest: request.QueryRequest,
					},
				}
				res, _ := this_.Call(listener, methodSystemQueryMonitorData, msg)
				if res != nil && res.SystemData != nil {
					response = res.SystemData
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) systemCleanMonitorData(nodeId string, NotifiedNodeIdList []string) (response *SystemData) {
	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		system.CleanMonitorData()
		response = &SystemData{}
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if response == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				msg := &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					SystemData: &SystemData{
						NodeId: nodeId,
					},
				}
				res, _ := this_.Call(listener, methodSystemCleanMonitorData, msg)
				if res != nil && res.SystemData != nil {
					response = res.SystemData
				}
				return
			})
		}
	}
	return
}

func (this_ *Worker) systemGetInfo(nodeId string, NotifiedNodeIdList []string) (response *SystemData) {
	if nodeId == "" {
		return
	}
	if this_.server.Id == nodeId {
		response = &SystemData{
			Info: system.GetInfo(),
		}
		return
	}

	var callPools = this_.getOtherPool(&NotifiedNodeIdList)

	for _, pool := range callPools {
		if response == nil {
			_ = pool.Do("", func(listener *MessageListener) (e error) {
				msg := &Message{
					NotifiedNodeIdList: NotifiedNodeIdList,
					SystemData: &SystemData{
						NodeId: nodeId,
					},
				}
				res, _ := this_.Call(listener, methodSystemGetInfo, msg)
				if res != nil && res.SystemData != nil {
					response = res.SystemData
				}
				return
			})
		}
	}
	return
}
