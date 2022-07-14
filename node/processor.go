package node

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"time"
)

type MethodMessage struct {
	ID     string
	Method string
	Error  string
	Body   []byte
	*Processor
}

func (this_ *MethodMessage) ReturnError(error string) (err error) {
	msg := &MethodMessage{
		ID:        this_.ID,
		Method:    this_.Method,
		Processor: this_.Processor,
		Error:     error,
	}

	err = this_.Send(msg)
	if err != nil {
		return
	}

	return
}

func (this_ *MethodMessage) Return(body []byte) (err error) {
	msg := &MethodMessage{
		ID:        this_.ID,
		Method:    this_.Method,
		Processor: this_.Processor,
		Body:      body,
	}

	err = this_.Send(msg)
	if err != nil {
		return
	}

	return
}

type Processor struct {
	conn          net.Conn
	isStop        bool
	Node          *Info
	doMethod      func(msg *MethodMessage) (body []byte, err error)
	CallbackCache map[string]func(msg *MethodMessage)
}

func (this_ *Processor) readField() (bytes []byte, err error) {
	if this_.needStop() {
		return
	}
	var buf []byte
	var n int

	buf = make([]byte, 4)
	n, err = this_.conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			this_.isStop = true
		}
		if this_.needStop() {
			return
		}
		fmt.Println("conn read field length err:", err.Error())
		return
	}
	if n < 4 {
		err = errors.New("conn read field length bytes less than 4")
		fmt.Println(err)
		return
	}

	length := int(binary.LittleEndian.Uint32(buf))
	if n < 0 {
		err = errors.New("conn read field length less than 0")
		fmt.Println(err)
		return
	}
	bytes = make([]byte, length)
	if len(bytes) == 0 {
		return
	}
	n, err = this_.conn.Read(bytes)
	if err != nil {
		if err == io.EOF {
			this_.isStop = true
		}
		if this_.needStop() {
			return
		}
		fmt.Println("conn read field err:", err.Error())
		return
	}
	if n < length {
		err = errors.New(fmt.Sprint("conn read field bytes less than ", length))
		fmt.Println(err)
		return
	}

	return
}

func (this_ *Processor) writeField(bytes []byte) (err error) {
	if this_.needStop() {
		return
	}
	var n int
	var length = len(bytes)

	lengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

	n, err = this_.conn.Write(lengthBytes)
	if err != nil {
		if this_.needStop() {
			return
		}
		fmt.Println("conn write field length err:", err.Error())
		return
	}
	if n < 4 {
		err = errors.New(fmt.Sprint("conn write field length bytes less than ", 4))
		fmt.Println(err)
		return
	}
	if length == 0 {
		return
	}
	n, err = this_.conn.Write(bytes)
	if err != nil {
		if this_.needStop() {
			return
		}
		fmt.Println("conn write field err:", err.Error())
		return
	}
	if n < length {
		err = errors.New(fmt.Sprint("conn write field bytes less than ", length))
		fmt.Println(err)
		return
	}

	return
}

func (this_ *Processor) stopListen() {
	this_.isStop = true
	err := this_.conn.Close()
	if err != nil {
		fmt.Println("conn close err:", err.Error())
	}
}
func (this_ *Processor) listen() {
	this_.CallbackCache = make(map[string]func(msg *MethodMessage))
	var err error
	this_.isStop = false
	go func() {
		defer func() {
			this_.isStop = true
			if x := recover(); x != nil {
				fmt.Println("conn read err:", x)
				return
			}
			_ = this_.conn.Close()
		}()

		for {
			if this_.needStop() {
				return
			}

			var buf []byte
			msg := &MethodMessage{
				Processor: this_,
			}
			var token []byte
			token, err = this_.readField()
			if err != nil {
				return
			}

			if !this_.Node.checkToken(token) {
				return
			}

			buf, err = this_.readField()
			if err != nil {
				return
			}
			msg.ID = string(buf)

			buf, err = this_.readField()
			if err != nil {
				return
			}
			msg.Method = string(buf)

			buf, err = this_.readField()
			if err != nil {
				return
			}
			msg.Error = string(buf)

			msg.Body, err = this_.readField()
			if err != nil {
				return
			}

			go this_.onMethod(msg)
		}
	}()
}

func (this_ *Processor) onMethod(msg *MethodMessage) {
	if msg == nil {
		return
	}
	callback, ok := this_.CallbackCache[msg.ID]
	if ok {
		callback(msg)
	} else {

		body, err := this_.doMethod(msg)
		if err != nil {
			fmt.Println("doMethod err:", err.Error())
			err = msg.ReturnError(err.Error())
			if err != nil {
				fmt.Println("msg return err:", err.Error())
				return
			}
			return
		}
		err = msg.Return(body)
		if err != nil {
			fmt.Println("msg return err:", err.Error())
			return
		}
	}

}

func (this_ *Processor) Send(msg *MethodMessage) (err error) {
	if msg == nil {
		return
	}

	err = this_.writeField([]byte(this_.Node.Token))
	if err != nil {
		fmt.Println("writ field token err:", err.Error())
		return
	}

	err = this_.writeField([]byte(msg.ID))
	if err != nil {
		return
	}

	err = this_.writeField([]byte(msg.Method))
	if err != nil {
		return
	}

	err = this_.writeField([]byte(msg.Error))
	if err != nil {
		return
	}

	err = this_.writeField(msg.Body)
	if err != nil {
		return
	}

	return
}

func (this_ *Processor) Call(method string, body []byte) (result []byte, err error) {
	msg := &MethodMessage{
		ID:     uuid.NewString(),
		Method: method,
		Body:   body,
	}
	defer func() {
		delete(this_.CallbackCache, msg.ID)
	}()
	waitResult := make(chan *MethodMessage, 1)
	this_.CallbackCache[msg.ID] = func(msg *MethodMessage) {
		waitResult <- msg
	}
	err = this_.Send(msg)
	if err != nil {
		return
	}
	var isEnd bool
	go func() {
		time.Sleep(time.Second * 5)
		if isEnd {
			return
		}
		waitResult <- &MethodMessage{}
	}()
	res := <-waitResult
	if res.Error != "" {
		err = errors.New(res.Error)
		return
	}
	result = res.Body
	return
}

func (this_ *Processor) needStop() bool {

	return this_.isStop || this_.conn == nil
}
