package node

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var (
	OK = []byte{0, 0, 0, 0}
)

type Message struct {
	Method string
	Body   []byte
}

type Processor struct {
	conn   net.Conn
	isStop bool
}

func (this_ *Processor) Start() {
	var err error
	this_.isStop = false
	defer func() {
		this_.isStop = true
		if x := recover(); x != nil {
			fmt.Println("conn read err:", x)
			return
		}
		err = this_.conn.Close()
		fmt.Println("conn close err:", err.Error())
	}()
	for {
		buf := make([]byte, 1024*8)
		var n int
		n, err = this_.conn.Read(buf)
		if err != nil {
			fmt.Println("conn read err:", err.Error())
			return
		}
		if n <= 4 {
			return
		}
		methodLengthBytes := buf[0:4]
		methodLength := int(binary.LittleEndian.Uint32(methodLengthBytes))
		methodBytes := buf[4:methodLength]
		methodName := string(methodBytes)

		bodyLengthBytes := buf[4+methodLength : 4]
		bodyLength := int(binary.LittleEndian.Uint32(bodyLengthBytes))
		bodyBytes := buf[4+methodLength+4 : bodyLength]

		msg := &Message{
			Method: methodName,
			Body:   bodyBytes,
		}
		this_.onMethod(msg)
	}
}

func (this_ *Processor) onRead(bs []byte) {

	fmt.Println("processor on read:", string(bs))
}

func (this_ *Processor) onMethod(msg *Message) {
	if msg == nil {
		return
	}

	fmt.Println("processor on method:", msg.Method)
	fmt.Println("processor on body  :", string(msg.Body))
}

func (this_ *Processor) write(bs []byte) (err error) {
	if this_.needStop() {
		err = errors.New("processor is stopped")
		return
	}
	fmt.Println("processor write:", string(bs))
	_, err = this_.conn.Write(bs)
	return
}

func (this_ *Processor) Send(msg *Message) (err error) {
	if msg == nil {
		return
	}
	var bs []byte
	methodBytes := []byte(msg.Method)
	methodLengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(methodLengthBytes, uint32(len(methodBytes)))
	bs = append(bs, methodLengthBytes...)
	bs = append(bs, methodBytes...)

	bodyLengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(bodyLengthBytes, uint32(len(msg.Body)))
	bs = append(bs, bodyLengthBytes...)
	bs = append(bs, msg.Body...)

	err = this_.write(bs)
	return
}

func (this_ *Processor) needStop() bool {

	return this_.isStop || this_.conn == nil
}

func ByteContains(x, y []byte) (n []byte, contain bool) {
	index := bytes.Index(x, y)
	if index == -1 {
		return
	}
	lastIndex := index + len(y)
	n = append(x[:index], x[lastIndex:]...)
	return n, true
}
