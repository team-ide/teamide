package node

import (
	"fmt"
	"io"
	"net"
)

type Client struct {
	Host string
	Port int
}

func (this_ *Client) Start() (err error) {
	address := fmt.Sprintf("%s:%d", this_.Host, this_.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("client dial err:", err.Error())
		return
	}
	this_.ProcessConn(conn)
	return
}

func (this_ *Client) ProcessConn(conn net.Conn) {
	var err error
	defer func() {
		err = conn.Close()
		fmt.Println("client conn close err:", err.Error())
	}()
	for {
		buf := make([]byte, 1025)
		var n int
		n, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("client conn read err:", err.Error())
			continue
		}
		if n <= 0 {
			continue
		}
		bs := buf[0:n]
		this_.OnRead(bs)
	}
}

func (this_ *Client) OnRead(bs []byte) {
	if len(bs) == 0 {
		return
	}

	fmt.Println("client on read:", string(bs))
}
