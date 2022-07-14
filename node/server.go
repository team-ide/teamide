package node

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
	*Worker
}

func (this_ *Server) Start() (err error) {
	address := fmt.Sprintf("%s:%d", this_.Ip, this_.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	go func() {
		for {
			var conn net.Conn
			conn, err = listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err.Error())
				continue
			}
			processor := &Processor{
				conn:     conn,
				Node:     this_.Node,
				doMethod: this_.Worker.doMethod,
			}
			processor.listen()
		}
	}()
	return
}
