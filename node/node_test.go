package node

import (
	"sync"
	"testing"
)

var ()

func init() {
}

func TestNode(t *testing.T) {

	port := 11001

	root := &Info{
		Name:  "root",
		Code:  "root",
		Ip:    "127.0.0.1",
		Port:  port,
		Token: "root_token",
	}
	rootWorker := testStartNode(root)
	port++
	node1 := &Info{
		Name:       "node1",
		Code:       "node1",
		Ip:         "127.0.0.1",
		Port:       port,
		Token:      "node1_token",
		ParentCode: "root",
	}
	_ = testStartNode(node1)
	port++
	node2 := &Info{
		Name:       "node2",
		Code:       "node2",
		Ip:         "127.0.0.1",
		Port:       port,
		Token:      "node2_token",
		ParentCode: "node1",
	}
	_ = testStartNode(node2)
	port++
	node3 := &Info{
		Name:       "node3",
		Code:       "node3",
		Ip:         "127.0.0.1",
		Port:       port,
		Token:      "node3_token",
		ParentCode: "node1",
	}
	_ = testStartNode(node3)
	port++
	node4 := &Info{
		Name:       "node4",
		Code:       "node4",
		Ip:         "127.0.0.1",
		Port:       port,
		Token:      "node4_token",
		ParentCode: "node2",
	}
	_ = testStartNode(node4)

	rootWorker.AddNode(node1)
	rootWorker.AddNode(node2)
	rootWorker.AddNode(node3)
	rootWorker.AddNode(node4)

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
