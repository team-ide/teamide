package main

import (
	"go.uber.org/zap"
	"os"
	"strings"
	"sync"
	"teamide/internal"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/pkg/util"
	"teamide/pkg/window"
)

var (
	waitGroupForStop sync.WaitGroup
	serverTitle      = "Team · IDE"
	serverUrl        = ""

	// buildFlags go build -ldflags '-X main.buildFlags=--isStandAlone' .
	buildFlags = ""
)

func init() {
	if strings.Contains(buildFlags, "--isStandAlone") {
		base.IsStandAlone = true
	}
	if strings.Contains(buildFlags, "--isHtmlDev") {
		base.IsHtmlDev = true
	}
}

func main() {
	var err error
	var serverContext *context.ServerContext

	defer func() {
		if err := recover(); err != nil {
			if serverContext != nil {
				serverContext.Logger.Error("启动失败", zap.Any("error", err))
			}
			waitGroupForStop.Done()
		}
	}()

	for _, v := range os.Args {
		if v == "--isStandAlone" {
			base.IsStandAlone = true
			continue
		}
		if v == "--isHtmlDev" {
			base.IsHtmlDev = true
			continue
		}
	}

	waitGroupForStop.Add(1)

	serverConf, err := GetServerConf()
	if err != nil {
		panic(err)
	}

	serverContext, err = context.NewServerContext(serverConf)
	if err != nil {
		panic(err)
	}

	serverUrl, err = internal.Start(serverContext)
	if err != nil {
		panic(err)
	}
	if base.IsHtmlDev {
		serverUrl = "http://127.0.0.1:21081/"
	}

	if base.IsStandAlone {
		err = window.Start(serverUrl, func() {
			waitGroupForStop.Done()
		})
		if err != nil {
			panic(err)
		}
	}

	waitGroupForStop.Wait()
}

func GetServerConf() (serverConf context.ServerConf, err error) {
	serverConf = context.ServerConf{
		Server:     base.RootDir + "conf/config.yaml",
		PublicKey:  base.RootDir + "conf/publicKey.pem",
		PrivateKey: base.RootDir + "conf/privateKey.pem",
	}
	exists, err := util.PathExists(serverConf.Server)
	if err != nil {
		return
	}
	if !exists {
		serverConf.Server = ""
	}
	exists, err = util.PathExists(serverConf.PublicKey)
	if err != nil {
		return
	}
	if !exists {
		serverConf.PublicKey = ""
	}
	exists, err = util.PathExists(serverConf.PrivateKey)
	if err != nil {
		return
	}
	if !exists {
		serverConf.PrivateKey = ""
	}
	return
}
