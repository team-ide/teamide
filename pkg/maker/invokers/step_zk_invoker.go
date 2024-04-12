package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"github.com/team-ide/go-tool/zookeeper"
	"go.uber.org/zap"
	"teamide/pkg/maker/modelers"
)

func (this_ *Invoker) invokeZkStep(step *modelers.StepZkModel, invokeData *InvokeData) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			util.Logger.Error("invoke zk step error", zap.Any("error", err))
		}
		util.Logger.Debug("invoke zk step end", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))
	}()
	util.Logger.Debug("invoke zk step start", zap.Any("args", invokeData.args), zap.Any("vars", invokeData.vars))

	var path string
	path, err = invokeData.InvokeByStringRule(step.Path)
	if err != nil {
		util.Logger.Error("invoke zk step get path error", zap.Any("path", step.Path), zap.Any("error", err))
		return
	}

	zkService, err := this_.GetZkServiceByName(step.Config)
	if err != nil {
		util.Logger.Error("invoke zk step get zk service error", zap.Any("config", step.Config), zap.Any("error", err))
		return
	}

	switch step.GetType() {
	case modelers.ZkWatchChildren:

		if path == "" {
			err = errors.New("zk childrenW path is empty")
			util.Logger.Error("invoke zk step error", zap.Any("error", err))
			return
		}

		if step.CreatePathIfNotExists {
			err = zkService.CreateIfNotExists(path, "")
			if err != nil {
				util.Logger.Error("invoke zk step createIfNotExists error", zap.Any("path", path), zap.Any("error", err))
				return
			}
		}
		util.Logger.Debug("invoke zk childrenW start", zap.Any("path", path))

		err = zkService.WatchChildren(path, func(data *zookeeper.WatchChildrenListenData) (finish bool) {
			util.Logger.Debug("on zk WatchChildren", zap.Any("path", path), zap.Any("data", data))
			return
		})
		if err != nil {
			util.Logger.Error("invoke zk step WatchChildren error", zap.Any("path", path), zap.Any("error", err))
			return
		}

		break
	default:
		err = errors.New("invoke zk [" + step.Zk + "] can not be support")
		util.Logger.Error("invoke zk error", zap.Any("error", err))
		return
	}

	return
}
