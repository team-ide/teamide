package main

import (
	"os"
	"sync"
	"teamide/server"
	"teamide/window"
)

var (
	waitGroupForStop sync.WaitGroup
	serverTitle           = "Team · IDE"
	serverUrl             = ""
	IsDev            bool = false
)

func main() {
	var err error

	for _, v := range os.Args {
		if v == "--isDev" {
			IsDev = true
			continue
		}
	}

	waitGroupForStop.Add(1)

	serverUrl, err = server.Start()
	if err != nil {
		panic(err)
	}
	if IsDev {
		serverUrl = "http://127.0.0.1:21081/"
	}

	if !IsDev {
		err = window.Start(serverTitle, serverUrl, func() {
			waitGroupForStop.Done()
		})
		if err != nil {
			panic(err)
		}
	}

	waitGroupForStop.Wait()
}
