package node

import (
	"github.com/team-ide/go-tool/util"
	"sync"
)

var (
	StatusStarted int8 = 1
	StatusStopped int8 = 2
	StatusError   int8 = 3
)

type MonitorData struct {
	ReadSize           int64 `json:"readSize,omitempty"`
	ReadTime           int64 `json:"readTime,omitempty"`
	ReadLastSize       int64 `json:"readLastSize,omitempty"`
	ReadLastTime       int64 `json:"readLastTime,omitempty"`
	ReadLastTimestamp  int64 `json:"readLastTimestamp,omitempty"`
	readLock           sync.Mutex
	WriteSize          int64 `json:"writeSize,omitempty"`
	WriteTime          int64 `json:"writeTime,omitempty"`
	WriteLastSize      int64 `json:"writeLastSize,omitempty"`
	WriteLastTime      int64 `json:"writeLastTime,omitempty"`
	WriteLastTimestamp int64 `json:"writeLastTimestamp,omitempty"`
	writeLock          sync.Mutex
}

func (this_ *MonitorData) monitorRead(bytesSize int64, useTime int64) {
	this_.readLock.Lock()
	defer this_.readLock.Unlock()

	var nowTime = util.GetNowTime()
	if this_.ReadLastTimestamp == 0 {
		this_.ReadLastTimestamp = nowTime
	}
	if (nowTime - this_.ReadLastTimestamp) <= 1000 {
		this_.ReadLastSize += bytesSize
		this_.ReadLastTime += useTime
	} else {
		this_.ReadLastTimestamp = nowTime
		this_.ReadLastSize = bytesSize
		this_.ReadLastTime = useTime
	}

	this_.ReadSize += bytesSize
	this_.ReadTime += useTime
}

func (this_ *MonitorData) monitorWrite(bytesSize int64, useTime int64) {
	this_.writeLock.Lock()
	defer this_.writeLock.Unlock()

	var nowTime = util.GetNowTime()
	if this_.WriteLastTimestamp == 0 {
		this_.WriteLastTimestamp = nowTime
	}
	if (nowTime - this_.WriteLastTimestamp) <= 1000 {
		this_.WriteLastSize += bytesSize
		this_.WriteLastTime += useTime
	} else {
		this_.WriteLastTimestamp = nowTime
		this_.WriteLastSize = bytesSize
		this_.WriteLastTime = useTime
	}

	this_.WriteSize += bytesSize
	this_.WriteTime += useTime
}
