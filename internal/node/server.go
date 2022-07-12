package node

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	ServerHost string
	ServerPort int
}

func (this_ *Server) Start() (err error) {
	address := fmt.Sprintf("%s:%d", this_.ServerHost, this_.ServerPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err.Error())
			continue
		}
		go this_.ProcessConn(conn)
	}
}

func (this_ *Server) ProcessConn(conn net.Conn) {
	var err error
	defer func() {
		err = conn.Close()
		fmt.Println("server conn close err:", err.Error())
	}()
	for {
		buf := make([]byte, 1025)
		var n int
		n, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("server conn read err:", err.Error())
			continue
		}
		if n <= 0 {
			continue
		}
		bs := buf[0:n]
		this_.OnRead(bs)
	}
}

func (this_ *Server) OnRead(bs []byte) {
	if len(bs) == 0 {
		return
	}

	fmt.Println("server on read:", string(bs))
}
