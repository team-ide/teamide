package node

import (
	"net"
	"sync"
	"teamide/pkg/util"
)

type connCache struct {
	_connCache          map[string]net.Conn
	_connCacheLock      sync.Mutex
	_connWriteLockCache map[string]sync.Locker
	MonitorData         *MonitorData
}

func newConnCache(MonitorData *MonitorData) *connCache {
	return &connCache{
		_connCache:          make(map[string]net.Conn),
		_connWriteLockCache: make(map[string]sync.Locker),
		MonitorData:         MonitorData,
	}
}

func (this_ *connCache) clean() {
	this_._connCacheLock.Lock()
	defer this_._connCacheLock.Unlock()

	for _, conn := range this_._connCache {
		_ = conn.Close()
	}
	this_._connCache = make(map[string]net.Conn)
	this_._connWriteLockCache = make(map[string]sync.Locker)
	return
}

func (this_ *connCache) setConn(connId string, conn net.Conn) {
	this_._connCacheLock.Lock()
	defer this_._connCacheLock.Unlock()

	this_._connCache[connId] = conn
	this_._connWriteLockCache[connId] = &sync.Mutex{}
	return
}
func (this_ *connCache) getConn(connId string) (conn net.Conn, writeLock sync.Locker) {
	this_._connCacheLock.Lock()
	defer this_._connCacheLock.Unlock()

	conn, _ = this_._connCache[connId]
	writeLock, _ = this_._connWriteLockCache[connId]
	return
}

func (this_ *connCache) closeConn(connId string) (err error) {
	this_._connCacheLock.Lock()
	defer this_._connCacheLock.Unlock()

	conn, ok := this_._connCache[connId]
	if ok {
		delete(this_._connCache, connId)
		delete(this_._connWriteLockCache, connId)
		_ = conn.Close()
	}
	return
}

func (this_ *connCache) send(connId string, bytes []byte) (err error) {
	conn, writeLock := this_.getConn(connId)
	if conn != nil {
		writeLock.Lock()
		defer writeLock.Unlock()

		start := util.Now().UnixNano()

		_, err = conn.Write(bytes)

		end := util.Now().UnixNano()
		this_.MonitorData.monitorWrite(int64(len(bytes)), end-start)
		//Logger.Info(this_.server.GetServerInfo() + " 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 连接 [" + connId + "] 发送 [" + fmt.Sprint(len(bytes)) + "]")
	} else {
		//Logger.Warn(this_.server.GetServerInfo() + " 代理服务 " + this_.netProxy.Inner.GetInfoStr() + " 连接 [" + connId + "] 不存在")
	}
	return
}
