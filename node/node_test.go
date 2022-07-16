package node

import (
	"fmt"
	"sync"
	"teamide/pkg/util"
	"testing"
	"time"
)

func TestNode(t *testing.T) {

	port := 11000

	root := &Info{
		Id:      "root",
		Name:    "root",
		Address: fmt.Sprintf("127.0.0.1:%d", port),
		Token:   "root_token",
	}
	rootWorker := testStartNode(root)
	rootWorker.AddNode(root)
	for n := 1; n <= 10; n++ {
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

		rootWorker.AddNode(node)
	}

	time.Sleep(time.Second * 5)

	err := rootWorker.AddNetProxy(&NetProxy{
		Inner: &NetConfig{
			NodeId: fmt.Sprintf("node-%d", 8),
			//NodeId: fmt.Sprintf("root"),
		},
		Outer: &NetConfig{
			NodeId: fmt.Sprintf("node-%d", 4),
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
