package main

import (
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
	"teamide/pkg/util"
)

func main() {
	//if err := test(); err != nil {
	//	log.Fatal(err)
	//}

	var err error
	obj := &ptyMaster{}
	err = obj.Start("bash", nil, nil, 0, 0)
	if err != nil {
		util.Logger.Error("start error", zap.Error(err))
		return
	}

	go func() {
		_, err = io.Copy(obj, os.Stdin)
		if err != nil {
			util.Logger.Error("Stdin Copy error", zap.Error(err))
			return
		}
	}()

	go func() {
		_, err = io.Copy(os.Stdout, obj)
		if err != nil {
			util.Logger.Error("Stdout Copy error", zap.Error(err))
			return
		}
	}()
	var waitGroupForStop sync.WaitGroup
	waitGroupForStop.Add(1)
	waitGroupForStop.Wait()
}
