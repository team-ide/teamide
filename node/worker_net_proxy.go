package node

import (
	"errors"
	"go.uber.org/zap"
)

func (this_ *Worker) AddNetProxy(netProxy *NetProxy) (err error) {
	this_.netProxyLock.Lock()
	defer this_.netProxyLock.Unlock()

	if netProxy == nil {
		err = errors.New("网络代理配置不能为空")
		return
	}
	if netProxy.Inner == nil {
		err = errors.New("网络代理输入配置不能为空")
		return
	}
	if netProxy.Outer == nil {
		err = errors.New("网络代理输出配置不能为空")
		return
	}

	var lineNodeIdList = this_.GetNodeLineByFromTo(netProxy.Inner.NodeId, netProxy.Outer.NodeId)
	if len(lineNodeIdList) == 0 {
		err = errors.New("无法正确解析输入输出节点关系")
		return
	}
	Logger.Info("网络代理输入输出节点线", zap.Any("lineNodeIdList", lineNodeIdList))

	return
}
