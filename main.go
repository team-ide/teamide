package main

import (
	"fmt"
	"go.uber.org/zap"
	"io"
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

	// buildFlags go build -ldflags '-X main.buildFlags=--isServer' .
	buildFlags  = ""
	version     = ""
	isServer    = false
	isHtmlDev   = false
	isServerDev = false
	rootDir     string
	userHomeDir string
	isElectron  = false
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
	if strings.Contains(buildFlags, "--isServer") {
		isServer = true
	}
	if strings.Contains(buildFlags, "--isDev") || strings.Contains(buildFlags, "--isHtmlDev") {
		isHtmlDev = true
	}
	if strings.Contains(buildFlags, "--isDev") || strings.Contains(buildFlags, "--isServerDev") {
		isServerDev = true
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
			fmt.Println("启动失败:", err)
			if serverContext != nil {
				serverContext.Logger.Error("启动失败", zap.Any("error", err))
			}
			waitGroupForStop.Done()
		}
	}()

	for _, v := range os.Args {
		if v == "--isServer" {
			isServer = true
		}
		if v == "--isDev" || v == "--isHtmlDev" {
			isHtmlDev = true
		}
		if v == "--isDev" || v == "--isServerDev" {
			isServerDev = true
		}
		if v == "--isElectron" {
			isElectron = true
		}

	}

	waitGroupForStop.Add(1)

	serverConf := &context.ServerConf{
		Version:     version,
		IsServer:    isServer,
		IsHtmlDev:   isHtmlDev,
		IsServerDev: isServerDev,
		RootDir:     rootDir,
		UserHomeDir: userHomeDir,
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
		serverUrl = "http://localhost:21081/"
	}

	// 如果是  Electron 打开该程序，则监听控制台
	if isElectron {
		os.Stdout.Write([]byte("TeamIDE:event:serverUrl:" + serverUrl))

		go func() {
			for {
				var bs = make([]byte, 1024)
				_, err := os.Stdin.Read(bs)
				if err != nil {
					if err == io.EOF {
						err = nil
						break
					}
					panic(err)
				}
				util.Logger.Info("On Electron：", zap.Any("msg", string(bs)))
			}
		}()
	} else {
		if !serverContext.IsServer {
			err = window.Start(serverUrl, func() {
				waitGroupForStop.Done()
			})
			if err != nil {
				panic(err)
			}
		}
	}

	waitGroupForStop.Wait()
}

func formatServerConf(serverConf *context.ServerConf) (err error) {
	if serverConf.IsServer {

		serverConf.Server = serverConf.RootDir + "conf/sqlite.yaml"
		serverConf.PublicKey = serverConf.RootDir + "conf/publicKey.pem"
		serverConf.PrivateKey = serverConf.RootDir + "conf/privateKey.pem"

		var exists bool
		exists, err = util.PathExists(serverConf.Server)
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
