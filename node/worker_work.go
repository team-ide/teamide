package node

func (this_ *Worker) netProxyNewConn(lineNodeIdList []string, netProxyId string, connId string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, connId, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyNewConn, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyId: netProxyId,
				ConnId:     connId,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	outer := this_.getNetProxyOuter(netProxyId)
	if outer != nil {
		err = outer.newConn(connId)
	}
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxyCloseConn(isReverse bool, lineNodeIdList []string, netProxyId string, connId string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, connId, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxyCloseConn, &Message{
			LineNodeIdList: lineNodeIdList,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyId: netProxyId,
				ConnId:     connId,
				IsReverse:  isReverse,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	if isReverse {
		inner := this_.getNetProxyInner(netProxyId)
		if inner != nil {
			err = inner.closeConn(connId)
		}
	} else {
		outer := this_.getNetProxyOuter(netProxyId)
		if outer != nil {
			err = outer.closeConn(connId)
		}
	}
	if err != nil {
		return
	}
	return
}

func (this_ *Worker) netProxySend(isReverse bool, lineNodeIdList []string, netProxyId string, connId string, bytes []byte) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, connId, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodNetProxySend, &Message{
			LineNodeIdList: lineNodeIdList,
			HasBytes:       true,
			Bytes:          bytes,
			NetProxyWorkData: &NetProxyWorkData{
				NetProxyId: netProxyId,
				ConnId:     connId,
				IsReverse:  isReverse,
			},
		})
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}
	if isReverse {
		inner := this_.getNetProxyInner(netProxyId)
		if inner != nil {
			err = inner.send(connId, bytes)
		}
	} else {
		outer := this_.getNetProxyOuter(netProxyId)
		if outer != nil {
			err = outer.send(connId, bytes)
		}
	}
	if err != nil {
		return
	}
	return
}
