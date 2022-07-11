package main

import (
	"sync"
	"teamide/pkg/util"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {

	var err error
	server := &Server{
		ServerHost: "0.0.0.0",
		ServerPort: 0,
	}
	serverUrl, err := server.Start()
	if err != nil {
		panic(err)
	}
	util.Logger.Info("Node Server Url:" + serverUrl)

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
