package main

import (
	"sync"
	"teamide/server"
	"teamide/server/base"
	"teamide/window"
)

var (
	waitGroupForStop sync.WaitGroup
	ServerTitle      = "Team · IDE"
	ServerUrl        = ""
)

func main() {
	var err error

	base.IsLocalStartup = true

	waitGroupForStop.Add(1)

	ServerUrl, err = server.Start()
	if err != nil {
		panic(err)
	}
	err = window.Start(ServerTitle, ServerUrl, func() {
		waitGroupForStop.Done()
	})
	if err != nil {
		panic(err)
	}

	waitGroupForStop.Wait()
}
