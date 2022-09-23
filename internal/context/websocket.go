package context

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/util"
)

var (
	idServerWebsocket   map[string]*ServerWebsocket
	userServerWebsocket map[string][]*ServerWebsocket
	serverWebsocketLock sync.Mutex

	websocketMessageOut chan *WebsocketMessage

	websocketMessageLock  sync.Mutex
	serverWebsocketInited bool
)

func initServerWebsocket() {
	if serverWebsocketInited {
		return
	}
	serverWebsocketInited = true
	idServerWebsocket = make(map[string]*ServerWebsocket)
	userServerWebsocket = make(map[string][]*ServerWebsocket)

	websocketMessageOut = make(chan *WebsocketMessage, 10)

	go func() {
		for {
			select {
			case msg := <-websocketMessageOut:
				if msg == nil {
					continue
				}
				var list []*ServerWebsocket
				if msg.id != "" {
					find := GetServerWebsocket(msg.id)
					if find != nil {
						list = append(list, find)
					}
				} else if msg.userId != "" {
					list = GetUserServerWebsocketList(msg.userId)
				} else {
					list = GetServerWebsocketList()
				}
				go func() {

					defer func() {
						if e := recover(); e != nil {
							util.Logger.Error("WSWriteText error", zap.Any("error", e))
						}
					}()

					bs, _ := json.Marshal(msg)
					for _, one := range list {
						one.WSWriteText(bs)
					}
				}()
			}

		}
	}()

}

func ServerWebsocketOutEvent(event string, data interface{}) {
	websocketMessageLock.Lock()
	defer websocketMessageLock.Unlock()

	msg := &WebsocketMessage{
		IsEvent: true,
		Event:   event,
		Data:    data,
	}
	websocketMessageOut <- msg
}

func AddServerWebsocket(id string, userId string, ws *websocket.Conn) {
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	serverWebsocket := &ServerWebsocket{
		id:     id,
		userId: userId,
		ws:     ws,
	}

	idServerWebsocket[id] = serverWebsocket
	_, ok := userServerWebsocket[userId]
	if !ok {
		userServerWebsocket[userId] = []*ServerWebsocket{}
	}
	userServerWebsocket[userId] = append(userServerWebsocket[userId], serverWebsocket)
}

func GetServerWebsocket(id string) (serverWebsocket *ServerWebsocket) {
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	serverWebsocket = idServerWebsocket[id]
	return
}

func GetServerWebsocketList() (list []*ServerWebsocket) {
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	for _, one := range idServerWebsocket {
		list = append(list, one)
	}
	return
}

func GetUserServerWebsocketList(userId string) (list []*ServerWebsocket) {
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	list = userServerWebsocket[userId]
	return
}

func CloseServerWebsocket(serverWebsocket *ServerWebsocket) {
	if serverWebsocket == nil {
		return
	}
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	_ = serverWebsocket.ws.Close()

	delete(idServerWebsocket, serverWebsocket.id)

	list, ok := userServerWebsocket[serverWebsocket.userId]
	if ok {
		var newList []*ServerWebsocket
		for _, one := range list {
			if one.id != serverWebsocket.id {
				newList = append(newList, one)
			}
		}
		if len(newList) == 0 {
			delete(userServerWebsocket, serverWebsocket.userId)
		} else {
			userServerWebsocket[serverWebsocket.userId] = newList
		}
	}
}

type WebsocketMessage struct {
	IsMessage      bool        `json:"isMessage,omitempty"`
	IsEvent        bool        `json:"isEvent,omitempty"`
	ErrorMessage   string      `json:"errorMessage,omitempty"`
	WarnMessage    string      `json:"warnMessage,omitempty"`
	InfoMessage    string      `json:"infoMessage,omitempty"`
	SuccessMessage string      `json:"successMessage,omitempty"`
	Event          string      `json:"event,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	id             string
	userId         string
}

type ServerWebsocket struct {
	id          string
	userId      string
	ws          *websocket.Conn
	wsWriteLock sync.Mutex
}

func (this_ *ServerWebsocket) WSWriteText(bs []byte) {
	this_.WSWriteByType(websocket.TextMessage, bs)
	return
}

func (this_ *ServerWebsocket) WSWriteBinary(bs []byte) {
	this_.WSWriteByType(websocket.BinaryMessage, bs)
	return
}

func (this_ *ServerWebsocket) WSWriteByType(messageType int, bs []byte) {

	this_.wsWriteLock.Lock()
	defer this_.wsWriteLock.Unlock()

	//fmt.Println("write message:", string(msg.Data))
	err := this_.ws.WriteMessage(messageType, bs)

	if err != nil {
		CloseServerWebsocket(this_)
		return
	}

}
