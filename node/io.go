package node

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
)

var (
	LengthError     = errors.New("读取流长度错误")
	ConnClosedError = errors.New("连接已关闭")
)

type Message struct {
	Token          string      `json:"token,omitempty"`
	Id             string      `json:"id,omitempty"`
	FromNodeId     string      `json:"fromNodeId,omitempty"`
	Method         int         `json:"method,omitempty"`
	Error          string      `json:"error,omitempty"`
	Ok             bool        `json:"ok,omitempty"`
	Node           *Info       `json:"node,omitempty"`
	NodeList       []*Info     `json:"nodeList,omitempty"`
	NetProxyId     string      `json:"netProxyId,omitempty"`
	ConnId         string      `json:"connId,omitempty"`
	IsReverse      bool        `json:"isReverse,omitempty"`
	LineNodeIdList []string    `json:"lineNodeIdList,omitempty"`
	NetProxy       *NetProxy   `json:"netProxy,omitempty"`
	NetProxyList   []*NetProxy `json:"netProxyList,omitempty"`
	Bytes          []byte      `json:"bytes,omitempty"`
	listener       *MessageListener
}

func (this_ *Message) ReturnError(error string) (err error) {
	err = this_.Return(&Message{
		Error: error,
	})
	if err != nil {
		return
	}

	return
}

func (this_ *Message) Return(msg *Message) (err error) {
	msg.Id = this_.Id
	msg.Method = this_.Method
	msg.Token = this_.Token
	err = this_.listener.Send(msg)
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
	id        string
}

func (this_ *MessageListener) stop() {
	this_.isStop = true
	_ = this_.conn.Close()
}

func (this_ *MessageListener) listen(onClose func()) {
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
			msg, err = ReadMessage(this_.conn)
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

func (this_ *MessageListener) Send(msg *Message) (err error) {
	if msg == nil {
		return
	}
	if this_.isClose {
		err = ConnClosedError
		return
	}
	this_.writeMu.Lock()
	defer this_.writeMu.Unlock()
	err = WriteMessage(this_.conn, msg)
	return
}

func ReadMessage(reader io.Reader) (message *Message, err error) {
	var bytes []byte

	bytes, err = ReadBytes(reader)
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

func WriteMessage(writer io.Writer, message *Message) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(message)
	if err != nil {
		return
	}

	err = WriteBytes(writer, bytes)
	return
}

func ReadData(reader io.Reader) (data map[string]interface{}, err error) {
	var bytes []byte

	bytes, err = ReadBytes(reader)
	if err != nil {
		return
	}
	data = map[string]interface{}{}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return
	}
	return
}

func WriteData(writer io.Writer, data map[string]interface{}) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(data)
	if err != nil {
		return
	}

	err = WriteBytes(writer, bytes)
	return
}

func ReadBytes(reader io.Reader) (bytes []byte, err error) {
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
	return
}

func WriteBytes(writer io.Writer, bytes []byte) (err error) {
	var n int
	var length = len(bytes)

	lengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

	n, err = writer.Write(lengthBytes)
	if err != nil {
		return
	}
	if n < 4 {
		err = LengthError
		return
	}
	if length == 0 {
		return
	}
	n, err = writer.Write(bytes)
	if err != nil {
		return
	}
	if n < length {
		err = LengthError
		return
	}
	return
}
