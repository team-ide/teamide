package main

import (
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"teamide/internal"
	"teamide/internal/context"
	"teamide/pkg/util"
	"teamide/pkg/window"
)

var (
	waitGroupForStop sync.WaitGroup
	serverTitle      = "Team · IDE"
	serverUrl        = ""

	// buildFlags go build -ldflags '-X main.buildFlags=--isStandAlone' .
	buildFlags   = ""
	isStandAlone = false
	isHtmlDev    = false
	rootDir      string
	userHomeDir  string
)

func getUserHome() string {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir
	}
	return ""
}
func init() {
	var err error
	if strings.Contains(buildFlags, "--isStandAlone") {
		isStandAlone = true
	}
	if strings.Contains(buildFlags, "--isHtmlDev") {
		isHtmlDev = true
	}
	rootDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		panic(err)
	}
	rootDir = filepath.ToSlash(rootDir)
	if !strings.HasSuffix(rootDir, "/") {
		rootDir += "/"
	}

	userHome := getUserHome()
	if userHome != "" {
		userHome, err = filepath.Abs(userHome)
		if err != nil {
			panic(err)
		}
		userHomeDir = filepath.ToSlash(userHome)
		if !strings.HasSuffix(userHomeDir, "/") {
			userHomeDir += "/"
		}

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
			isStandAlone = true
			continue
		}
		if v == "--isHtmlDev" {
			isHtmlDev = true
			continue
		}
	}

	waitGroupForStop.Add(1)

	serverConf := &context.ServerConf{
		IsStandAlone: isStandAlone,
		IsHtmlDev:    isHtmlDev,
		RootDir:      rootDir,
		UserHomeDir:  userHomeDir,
	}
	err = formatServerConf(serverConf)
	if err != nil {
		panic(err)
	}

	serverContext, err = context.NewServerContext(*serverConf)
	if err != nil {
		panic(err)
	}
	serverUrl, err = internal.Start(serverContext)
	if err != nil {
		panic(err)
	}
	if serverContext.IsHtmlDev {
		serverUrl = "http://127.0.0.1:21081/"
	}

	if serverContext.IsStandAlone {
		err = window.Start(serverUrl, func() {
			waitGroupForStop.Done()
		})
		if err != nil {
			panic(err)
		}
	}

	waitGroupForStop.Wait()
}

func formatServerConf(serverConf *context.ServerConf) (err error) {
	if !serverConf.IsStandAlone {

		serverConf.Server = serverConf.RootDir + "conf/config.yaml"
		serverConf.PublicKey = serverConf.RootDir + "conf/publicKey.pem"
		serverConf.PrivateKey = serverConf.RootDir + "conf/privateKey.pem"

		var exists bool
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
	}
	return
}
