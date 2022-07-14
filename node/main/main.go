package main

import (
	"sync"
	"teamide/node"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {

	var err error
	server := &node.Server{
		Ip:   "0.0.0.0",
		Port: 0,
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
