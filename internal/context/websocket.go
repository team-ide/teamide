package context

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strings"
	"sync"
	"teamide/pkg/util"
)

var (
	idServerWebsocket   = make(map[string]*ServerWebsocket)
	userServerWebsocket = make(map[string][]*ServerWebsocket)
	serverWebsocketLock sync.Mutex

	websocketMessageLock sync.Mutex
)

func ServerWebsocketOutEvent(event string, data interface{}) {
	websocketMessageLock.Lock()
	defer websocketMessageLock.Unlock()

	msg := &WebsocketMessage{
		IsEvent: true,
		Event:   event,
		Data:    data,
	}

	//util.Logger.Info("ServerWebsocketOutEvent start", zap.Any("event", event), zap.Any("data", data))

	go func(m *WebsocketMessage) {
		var list = GetServerWebsocketList()
		bs, err := json.Marshal(m)
		if err != nil {
			util.Logger.Error("ServerWebsocketOutEvent json marshal error", zap.Error(err))
			return
		}
		for _, one := range list {
			one.WSWriteText(bs)
		}
	}(msg)
}

func AddServerWebsocket(id string, userId string, ws *websocket.Conn) {
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	serverWebsocket := &ServerWebsocket{
		id:     id,
		userId: userId,
		ws:     ws,
	}
	go serverWebsocket.start()

	idServerWebsocket[id] = serverWebsocket

	if userId != "" {
		list := userServerWebsocket[userId]
		list = append(list, serverWebsocket)
		userServerWebsocket[userId] = list
	}
}

func ChangeServerWebsocketUserId(serverWebsocket *ServerWebsocket, userId string) {
	if serverWebsocket.userId == userId {
		return
	}
	serverWebsocketLock.Lock()
	defer serverWebsocketLock.Unlock()

	if serverWebsocket.userId != "" {
		list := userServerWebsocket[serverWebsocket.userId]
		var newList []*ServerWebsocket
		for _, one := range list {
			if one != serverWebsocket {
				newList = append(newList, one)
			}
		}
		userServerWebsocket[serverWebsocket.userId] = newList
	}

	serverWebsocket.userId = userId
	if userId != "" {
		list := userServerWebsocket[userId]
		list = append(list, serverWebsocket)
		userServerWebsocket[userId] = list
	}

	return
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

	if serverWebsocket.userId != "" {
		list := userServerWebsocket[serverWebsocket.userId]
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

	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("WSWriteByType error", zap.Any("error", e))
		}
	}()
	//fmt.Println("write message:", string(msg.Data))
	err := this_.ws.WriteMessage(messageType, bs)

	if err != nil {
		CloseServerWebsocket(this_)
		return
	}

}

func (this_ *ServerWebsocket) start() {

	defer func() {
		_ = this_.ws.Close()
	}()
	for {
		//读取ws中的数据
		_, bytes, err := this_.ws.ReadMessage()
		if err != nil {
			break
		}
		text := string(bytes)
		if strings.HasSuffix(text, "change userId:") {
			userId := strings.TrimSuffix(text, "change userId:")
			userId = strings.TrimSpace(userId)
			util.Logger.Info("web socket change userId", zap.Any("oldUserId", this_.userId), zap.Any("userId", userId))
			ChangeServerWebsocketUserId(this_, userId)
		}
	}

}
