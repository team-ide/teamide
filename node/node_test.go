package node

import (
	"fmt"
	"sync"
	"teamide/pkg/util"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	var err error
	port := 11000

	root := &Info{
		Id:      "root",
		Name:    "root",
		Address: fmt.Sprintf("127.0.0.1:%d", port),
		Token:   "root_token",
	}
	rootWorker := testStartNode(root)
	err = rootWorker.AddNode(root)
	if err != nil {
		panic(err)
	}
	for n := 1; n <= 3; n++ {
		port++

		node := &Info{
			Id:      fmt.Sprintf("node-%d", n),
			Name:    fmt.Sprintf("node-%d", n),
			Address: fmt.Sprintf("127.0.0.1:%d", port),
			Token:   fmt.Sprintf("node-%d-token", n),
		}
		if util.ContainsInt([]int{1, 2}, n) >= 0 {
			node.ParentId = root.Id
		} else if util.ContainsInt([]int{3, 4}, n) >= 0 {
			node.ParentId = fmt.Sprintf("node-%d", 2)
		} else if util.ContainsInt([]int{5, 6}, n) >= 0 {
			node.ParentId = fmt.Sprintf("node-%d", 3)
		} else if util.ContainsInt([]int{7, 8}, n) >= 0 {
			node.ParentId = fmt.Sprintf("node-%d", 1)
		} else if util.ContainsInt([]int{9, 10}, n) >= 0 {
			node.ParentId = fmt.Sprintf("node-%d", 8)
		}
		_ = testStartNode(node)

		err = rootWorker.AddNode(node)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 5)

	err = rootWorker.AddNetProxy(&NetProxy{
		Id: util.UUID(),
		Inner: &NetConfig{
			Address: ":8088",
			NodeId:  fmt.Sprintf("node-%d", 1),
			//NodeId: fmt.Sprintf("root"),
		},
		Outer: &NetConfig{
			NodeId:  fmt.Sprintf("node-%d", 3),
			Address: "teamide.com:22",
		},
	})
	if err != nil {
		panic(err)
	}

	//rootWorker.RemoveNode(&Info{
	//	Id: fmt.Sprintf("node-%d", 5),
	//})

	var waitGroupForStop sync.WaitGroup

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}

func testStartNode(node *Info) (worker *Worker) {
	println("启动节点 " + node.GetNodeStr())
	worker = &Worker{
		Node: node,
	}
	err := worker.Start()
	if err != nil {
		panic(err)
	}
	return
}
