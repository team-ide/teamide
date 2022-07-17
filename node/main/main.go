package main

import (
	"flag"
	"fmt"
	"sync"
	"teamide/node"
	"teamide/pkg/util"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {

	var err error

	var id string
	var name string
	var network string
	var address string
	var token string
	flag.StringVar(&id, "id", "", "节点ID，需要唯一")
	flag.StringVar(&name, "name", "", "节点名称")
	flag.StringVar(&network, "network", "tcp", "节点启动网络类型")
	flag.StringVar(&address, "address", "", "节点启动监听地址")
	flag.StringVar(&token, "token", "", "节点Token，用于验证")

	//解析
	flag.Parse()

	if id == "" {
		println("请配置有效参数，可以使用-help查看")
		panic("请设置 -id")
	}
	if name == "" {
		println("请配置有效参数，可以使用-help查看")
		panic("请设置 -name")
	}
	if address == "" {
		println("请配置有效参数，可以使用-help查看")
		panic("请设置 -address")
	}
	if token == "" {
		println("请配置有效参数，可以使用-help查看")
		panic("请设置 -token")
	}

	worker := &node.Worker{
		Node: &node.Info{
			Id:      id,
			Name:    name,
			Network: network,
			Address: address,
			Token:   token,
		},
	}
	println("启动节点 [" + name + "][" + network + "][" + address + "] 开始")
	err = worker.Start()
	if err != nil {
		println("启动节点 [" + name + "][" + network + "][" + address + "] 异常")
		panic(err)
	}
	println("启动节点 [" + name + "][" + network + "][" + address + "] 成功")

	_ = worker.AddNode(&node.Info{
		Id:       "node1",
		Name:     "node1",
		Address:  "127.0.0.1:11001",
		Token:    "xxx",
		ParentId: "root",
	})
	_ = worker.AddNode(&node.Info{
		Id:       "node2",
		Name:     "node2",
		Address:  "127.0.0.1:11002",
		Token:    "xxxx",
		ParentId: "root",
	})

	_ = worker.AddNetProxy(&node.NetProxy{
		Id: util.UUID(),
		Inner: &node.NetConfig{
			NodeId: fmt.Sprintf("node%d", 1),
			//NodeId: fmt.Sprintf("root"),
			Address: ":8088",
		},
		Outer: &node.NetConfig{
			NodeId:  fmt.Sprintf("node%d", 2),
			Address: "teamide.com:22",
		},
	})

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
