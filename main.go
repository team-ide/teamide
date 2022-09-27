package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"teamide/internal"
	"teamide/internal/context"
	"teamide/pkg/util"
	"teamide/pkg/window"
)

var (
	waitGroupForStop sync.WaitGroup
	serverUrl        = ""

	// buildFlags go build -ldflags '-X main.buildFlags=--isServer' .
	buildFlags  = ""
	isServer    = false
	isHtmlDev   = false
	isServerDev = false
	rootDir     = ""
	userHomeDir = ""
	isElectron  = false
)

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

	rootDir, err = os.Getwd()
	if err != nil {
		util.Logger.Error("os get wd error", zap.Error(err))
		panic(err)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		util.Logger.Error("filepath abs error", zap.Error(err))
		panic(err)
	}
	rootDir = filepath.ToSlash(rootDir)
	if !strings.HasSuffix(rootDir, "/") {
		rootDir += "/"
	}
	current, err := user.Current()
	if err != nil {
		util.Logger.Error("user current error", zap.Error(err))
		panic(err)
	}

	userHomeDir = current.HomeDir
	if userHomeDir != "" {
		userHomeDir, err = filepath.Abs(userHomeDir)
		if err != nil {
			util.Logger.Error("filepath abs error", zap.Error(err))
			panic(err)
		}
		userHomeDir = filepath.ToSlash(userHomeDir)
		if !strings.HasSuffix(userHomeDir, "/") {
			userHomeDir += "/"
		}

	}
}

func main() {
	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			println(util.GetVersion())
			return
		}
	}
	var err error
	var serverContext *context.ServerContext

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("启动失败:", e)
			if serverContext != nil {
				serverContext.Logger.Error("启动失败", zap.Any("error", e))
			}
			waitGroupForStop.Done()
		}
	}()

	// 开启 cpu 采集分析：
	if isServerDev {
		var f *os.File
		//-cpu-profile=cpu.pprof -mem-profile=mem.pprof
		f, err = os.Create("cpu.pprof")
		if err != nil {
			util.Logger.Error("could not create CPU profile", zap.Error(err))
			return
		}
		defer func() { _ = f.Close() }()

		util.Logger.Info("CPU profile create success")
		if err = pprof.StartCPUProfile(f); err != nil {
			util.Logger.Error("could not start CPU profile", zap.Error(err))
		}
		defer pprof.StopCPUProfile()
	}

	defer func() {

		if isServerDev {
			var f *os.File
			f, err = os.Create("mem.pprof")
			if err != nil {
				util.Logger.Error("could not create memory profile", zap.Error(err))
				return
			}
			defer func() { _ = f.Close() }()
			util.Logger.Info("memory profile create success")

			runtime.GC()

			if err = pprof.WriteHeapProfile(f); err != nil {
				util.Logger.Error("could not write memory profile", zap.Error(err))
			}
		}
	}()

	waitGroupForStop.Add(1)

	serverConf := &context.ServerConf{
		Version:     util.GetVersion(),
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
		_, _ = os.Stdout.Write([]byte("TeamIDE:event:serverUrl:" + serverUrl))
		go func() {
			var buf = make([]byte, 1024)
			err = util.Read(os.Stdin, buf, func(n int) (err error) {
				util.Logger.Info("On Electron：", zap.Any("msg", string(buf[:n])))
				return
			})
			if err != nil {
				err = errors.New("electron window closed")
			}
			waitGroupForStop.Done()
			panic(err)
		}()
	} else {
		if !serverContext.IsServer {
			err = window.Start(serverUrl, func() {
				util.Logger.Info("TeamIDE stopped")
				waitGroupForStop.Done()
			})
			if err != nil {
				panic(err)
			}
		}
	}

	waitGroupForStop.Wait()

	if isServerDev {
		pprof.StopCPUProfile()
	}
}

func formatServerConf(serverConf *context.ServerConf) (err error) {
	if serverConf.IsServer {

		serverConf.Server = serverConf.RootDir + "conf/config.yaml"
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
