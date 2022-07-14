package node

import (
	"sync"
	"testing"
)

var ()

func init() {
}

func TestNode(t *testing.T) {

	node1 := &Info{
		Name:  "node1",
		Ip:    "127.0.0.1",
		Port:  11001,
		Token: "node1_Token",
	}
	node1Worker := testStartNode(node1)

	node2 := &Info{
		Name:  "node2",
		Ip:    "127.0.0.1",
		Port:  11002,
		Token: "node2_Token",
	}
	node2Worker := testStartNode(node2)

	node3 := &Info{
		Name:  "node3",
		Ip:    "127.0.0.1",
		Port:  11003,
		Token: "node3_Token",
	}
	node3Worker := testStartNode(node3)

	node1Worker.AddChildrenNode(node2)

	node2Worker.AddChildrenNode(node3)

	node3Worker.AddChildrenNode(node1)

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
