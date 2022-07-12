package main

import (
	"sync"
	"teamide/internal/node"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {

	var err error
	server := &node.Server{
		ServerHost: "0.0.0.0",
		ServerPort: 0,
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
