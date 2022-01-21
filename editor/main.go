package main

import (
	"io"
	"os"
	"sync"
	"teamide/application/base"
)

var (
	waitGroupForStop sync.WaitGroup
)

const (
	K = "l@FQsnTKEdtty1@w"
)

var (
	dataDir = "data"
	appsDir = dataDir + "/apps"
	isDev   = false
)

func onWindowClose() {
	waitGroupForStop.Done()
}
func init() {
	var err error
	dirPath := dataDir
	var exists bool
	exists, err = base.PathExists(dirPath)
	if err != nil {
		panic(err)
	}
	if !exists {
		os.MkdirAll(dirPath, 0777)
	}

	dirPath = appsDir
	exists, err = base.PathExists(dirPath)
	if err != nil {
		panic(err)
	}
	if !exists {
		os.MkdirAll(dirPath, 0777)
	}
	isDev = true
	checkFiles := []string{"main.go", "window.go", "package.go"}
	for _, file := range checkFiles {
		dirPath = "main.go"
		exists, err = base.PathExists(file)
		if err != nil {
			panic(err)
		}
		if !exists {
			isDev = false
			break
		}
	}
}

// 取消控制台 go build -ldflags "-H windowsgui"
func main() {
	waitGroupForStop.Add(1)
	var err error

	err = startWebServer()
	if err != nil {
		panic(err)
	}

	// if !isDev {
	err = openWindow(htmlServerUrl)
	if err != nil {
		panic(err)
	}
	// }

	waitGroupForStop.Wait()
}

func WriteFile(filename string, bs []byte) (err error) {
	var f *os.File
	var exists bool
	exists, err = base.PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	}
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(bs)

	if err != nil {
		return
	}
	return
}

func ReadFile(filename string) (bs []byte, err error) {
	var f *os.File
	var exists bool
	exists, err = base.PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		return
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		return
	}
	defer f.Close()
	bs, err = io.ReadAll(f)
	if err != nil {
		return
	}
	return
}
