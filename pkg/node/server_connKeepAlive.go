package node

func (this_ *Server) connNodeListenerKeepAlive(connAddress, connToken string, connSize int) {
	if connAddress == "" {
		Logger.Warn(this_.GetServerInfo() + " 连接 [" + connAddress + "] 连接地址为空")
		return
	}
	this_.connNodeListenerKeepAliveLock.Lock()
	defer this_.connNodeListenerKeepAliveLock.Unlock()

	if connSize <= 0 {
		connSize = 5
	}
	for connIndex := 0; connIndex < connSize; connIndex++ {
		go this_.connNodeListener(nil, connAddress, connToken, connIndex)
	}
	return
}
