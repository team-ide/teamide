package module_node

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type WSConn struct {
	id       string
	conn     *websocket.Conn
	connLock sync.Mutex
}

func (this_ *WSConn) Close() {
	_ = this_.conn.Close()
	return
}

func (this_ *WSConn) WriteMessage(bytes []byte) (err error) {
	this_.connLock.Lock()
	defer this_.connLock.Unlock()

	err = this_.conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		return
	}
	return
}

func (this_ *NodeService) addWS(id string, ws *websocket.Conn) (err error) {
	wsConn := &WSConn{
		id:   id,
		conn: ws,
	}
	this_.nodeContext.setWS(wsConn)
	return
}

type Message struct {
	Method       string          `json:"method,omitempty"`
	NodeList     []*NodeInfo     `json:"nodeList,omitempty"`
	NetProxyList []*NetProxyInfo `json:"netProxyList,omitempty"`
}

func (this_ *NodeContext) getWS(id string) (ws *WSConn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	ws = this_.wsCache[id]
	return
}

func (this_ *NodeContext) setWS(ws *WSConn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	this_.wsCache[ws.id] = ws
}

func (this_ *NodeContext) removeWS(id string) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()

	find, ok := this_.wsCache[id]
	if ok {
		find.Close()
	}
	delete(this_.wsCache, id)

}

func (this_ *NodeContext) getWSList() (list []*WSConn) {
	this_.wsCacheLock.Lock()
	defer this_.wsCacheLock.Unlock()
	for _, one := range this_.wsCache {
		list = append(list, one)
	}
	return
}

func (this_ *NodeContext) callMessage(msg *Message) {
	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}
	var list = this_.getWSList()
	for _, one := range list {
		err = one.WriteMessage(bs)
		if err != nil {
			this_.removeWS(one.id)
		}
	}
}
