package main

import (
	"os"
	"sync"
	"teamide/server"
	"teamide/server/base"
	"teamide/window"
)

var (
	waitGroupForStop sync.WaitGroup
	serverTitle      = "Team · IDE"
	serverUrl        = ""
	isHtmlDev        = false
)

func main() {
	var err error

	for _, v := range os.Args {
		if v == "--isStandAlone" {
			base.IS_STAND_ALONE = true
			continue
		}
		if v == "--isHtmlDev" {
			isHtmlDev = true
			continue
		}
	}

	waitGroupForStop.Add(1)

	serverUrl, err = server.Start()
	if err != nil {
		panic(err)
	}
	if isHtmlDev {
		serverUrl = "http://127.0.0.1:21081/"
	}

	if base.IS_STAND_ALONE {
		err = window.Start(serverTitle, serverUrl, func() {
			waitGroupForStop.Done()
		})
		if err != nil {
			panic(err)
		}
	}

	waitGroupForStop.Wait()
}
