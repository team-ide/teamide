package node

import (
	"sync"
	"testing"
)

var ()

func init() {
}

func TestNode(t *testing.T) {

	root := &Info{
		Name:  "root",
		Code:  "root",
		Ip:    "127.0.0.1",
		Port:  11001,
		Token: "root_token",
	}
	rootWorker := testStartNode(root)

	node1 := &Info{
		Name:       "node1",
		Code:       "node1",
		Ip:         "127.0.0.1",
		Port:       11002,
		Token:      "node1_token",
		ParentCode: "root",
	}
	_ = testStartNode(node1)

	node2 := &Info{
		Name:       "node2",
		Code:       "node2",
		Ip:         "127.0.0.1",
		Port:       11003,
		Token:      "node2_token",
		ParentCode: "node1",
	}
	_ = testStartNode(node2)

	node3 := &Info{
		Name:       "node3",
		Code:       "node3",
		Ip:         "127.0.0.1",
		Port:       11004,
		Token:      "node3_token",
		ParentCode: "node1",
	}
	_ = testStartNode(node3)

	rootWorker.AddNode(node1)
	rootWorker.AddNode(node2)
	rootWorker.AddNode(node3)

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
