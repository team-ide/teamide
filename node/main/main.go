package main

import (
	"sync"
	"teamide/node"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {

	var err error

	nodeInfo := &node.Info{}
	worker := &node.Worker{
		Node: nodeInfo,
	}
	println("启动节点 " + nodeInfo.GetNodeStr() + " 开始")
	err = worker.Start()
	if err != nil {
		println("启动节点 " + nodeInfo.GetNodeStr() + " 异常")
		panic(err)
	}
	println("启动节点 " + nodeInfo.GetNodeStr() + " 成功")

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
