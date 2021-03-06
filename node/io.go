package node

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"teamide/pkg/util"
)

var (
	LengthError     = errors.New("读取流长度错误")
	ConnClosedError = errors.New("连接已关闭")
)

type Message struct {
	Id                            string          `json:"id,omitempty"`
	Method                        int             `json:"method,omitempty"`
	Error                         string          `json:"error,omitempty"`
	Ok                            bool            `json:"ok,omitempty"`
	NotifyChildren                bool            `json:"notifyChildren,omitempty"`
	NotifyParent                  bool            `json:"notifyParent,omitempty"`
	NotifyAll                     bool            `json:"notifyAll,omitempty"`
	NodeId                        string          `json:"nodeId,omitempty"`
	NodeIdList                    []string        `json:"nodeIdList,omitempty"`
	RemoveNodeIdList              []string        `json:"removeNodeIdList,omitempty"`
	NotifiedNodeIdList            []string        `json:"notifiedNodeIdList,omitempty"`
	RemoveConnNodeIdList          []string        `json:"removeConnNodeIdList,omitempty"`
	Node                          *Info           `json:"node,omitempty"`
	NodeList                      []*Info         `json:"nodeList,omitempty"`
	NetProxyId                    string          `json:"netProxyId,omitempty"`
	NetProxyIdList                []string        `json:"netProxyIdList,omitempty"`
	RemoveNetProxyIdList          []string        `json:"removeNetProxyIdList,omitempty"`
	ConnId                        string          `json:"connId,omitempty"`
	IsReverse                     bool            `json:"isReverse,omitempty"`
	LineNodeIdList                []string        `json:"lineNodeIdList,omitempty"`
	NetProxy                      *NetProxy       `json:"netProxy,omitempty"`
	NetProxyList                  []*NetProxy     `json:"netProxyList,omitempty"`
	Bytes                         []byte          `json:"bytes,omitempty"`
	ClientIndex                   int             `json:"clientIndex,omitempty"`
	NodeStatusChangeList          []*StatusChange `json:"nodeStatusChangeList,omitempty"`
	NetProxyInnerStatusChangeList []*StatusChange `json:"netProxyInnerStatusChangeList,omitempty"`
	NetProxyOuterStatusChangeList []*StatusChange `json:"netProxyOuterStatusChangeList,omitempty"`
	Version                       string          `json:"version,omitempty"`
	MonitorData                   *MonitorData    `json:"monitorData,omitempty"`
	listener                      *MessageListener
}

type StatusChange struct {
	Id          string `json:"id,omitempty"`
	Status      int8   `json:"status,omitempty"`
	StatusError string `json:"statusError,omitempty"`
}

func (this_ *Message) ReturnError(error string, MonitorData *MonitorData) (err error) {
	err = this_.Return(&Message{
		Error: error,
	}, MonitorData)
	if err != nil {
		return
	}

	return
}

func (this_ *Message) Return(msg *Message, MonitorData *MonitorData) (err error) {
	if this_.listener == nil {
		err = errors.New("消息监听器丢失")
		return
	}
	msg.Id = this_.Id
	err = this_.listener.Send(msg, MonitorData)
	if err != nil {
		return
	}
	return
}

type MessageListener struct {
	conn      net.Conn
	onMessage func(msg *Message)
	isClose   bool
	isStop    bool
	writeMu   sync.Mutex
}

func (this_ *MessageListener) stop() {
	this_.isStop = true
	_ = this_.conn.Close()
}

func (this_ *MessageListener) listen(onClose func(), MonitorData *MonitorData) {
	var err error
	this_.isClose = false
	go func() {
		defer func() {
			this_.isClose = true
			if x := recover(); x != nil {
				Logger.Error("message listen error", zap.Error(err))
			}
			_ = this_.conn.Close()
			onClose()
		}()

		for {
			if this_.isStop {
				return
			}
			var msg *Message
			msg, err = ReadMessage(this_.conn, MonitorData)
			if err != nil {
				if this_.isStop {
					return
				}
				if err == io.EOF {
					return
				}
				//Logger.Error("message read error", zap.Error(err))
				return
			}
			msg.listener = this_
			go this_.onMessage(msg)
		}
	}()
}

func (this_ *MessageListener) Send(msg *Message, MonitorData *MonitorData) (err error) {
	if msg == nil {
		return
	}
	if this_.isClose {
		err = ConnClosedError
		return
	}
	this_.writeMu.Lock()
	defer this_.writeMu.Unlock()
	err = WriteMessage(this_.conn, msg, MonitorData)
	return
}

func ReadMessage(reader io.Reader, MonitorData *MonitorData) (message *Message, err error) {
	var bytes []byte

	bytes, err = ReadBytes(reader, MonitorData)
	if err != nil {
		return
	}
	message = &Message{}

	err = json.Unmarshal(bytes, &message)
	if err != nil {
		return
	}
	return
}

func WriteMessage(writer io.Writer, message *Message, MonitorData *MonitorData) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(message)
	if err != nil {
		return
	}

	err = WriteBytes(writer, bytes, MonitorData)
	return
}

func ReadBytes(reader io.Reader, MonitorData *MonitorData) (bytes []byte, err error) {

	start := util.Now().UnixNano()

	var buf []byte
	var n int

	buf = make([]byte, 4)
	n, err = reader.Read(buf)
	if err != nil {
		return
	}
	if n < 4 {
		err = LengthError
		return
	}

	length := int(binary.LittleEndian.Uint32(buf))
	if length < 0 {
		err = LengthError
		return
	}

	if length > 0 {
		bytes = make([]byte, length)
		n, err = reader.Read(bytes)
		if err != nil {
			return
		}
		if n < length {
			err = LengthError
			return
		}
	}
	end := util.Now().UnixNano()
	MonitorData.monitorWrite(int64(length+4), end-start)
	return
}

func WriteBytes(writer io.Writer, bytes []byte, MonitorData *MonitorData) (err error) {
	var n int
	var length = len(bytes)

	start := util.Now().UnixNano()

	writeBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(writeBytes, uint32(length))

	writeBytes = append(writeBytes, bytes...)
	length += 4
	n, err = writer.Write(writeBytes)
	if err != nil {
		return
	}
	if n < length {
		err = LengthError
		return
	}
	end := util.Now().UnixNano()
	MonitorData.monitorWrite(int64(length), end-start)
	return
}
