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
		Id:          "root",
		BindAddress: fmt.Sprintf("127.0.0.1:%d", port),
		BindToken:   "root_token",
	}
	rootServer := testStartServer(root.Id, root.BindAddress, root.BindToken, "", "")
	if err != nil {
		panic(err)
	}

	var nodeList []*Info
	nodeList = append(nodeList, root)
	for n := 1; n <= 3; n++ {
		port++

		node := &Info{
			Id:          fmt.Sprintf("node-%d", n),
			ConnAddress: fmt.Sprintf("127.0.0.1:%d", port),
			BindToken:   fmt.Sprintf("node-%d-token", n),
		}
		if util.ContainsInt([]int{1, 2}, n) >= 0 {
			node.addConnNodeId(root.Id)
		} else if util.ContainsInt([]int{3, 4}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 2))
		} else if util.ContainsInt([]int{5, 6}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 3))
		} else if util.ContainsInt([]int{7, 8}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 1))
		} else if util.ContainsInt([]int{9, 10}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 8))
		}
		nodeList = append(nodeList, node)
		_ = testStartServer(node.Id, node.BindAddress, node.BindToken, "", "")

	}

	time.Sleep(time.Second * 5)
	println("开始添加节点")
	for _, one := range nodeList {
		err = rootServer.AddNodeList([]*Info{one})
		if err != nil {
			panic(err)
		}
	}
	println("节点添加完成")

	time.Sleep(time.Second * 5)

	println("开始添加代理")
	err = rootServer.AddNetProxyList([]*NetProxy{
		{
			Id: "1",
			Inner: &NetConfig{
				Address: ":8088",
				NodeId:  fmt.Sprintf("node-%d", 1),
				//NodeId: fmt.Sprintf("root"),
			},
			Outer: &NetConfig{
				NodeId:  fmt.Sprintf("node-%d", 3),
				Address: "teamide.com:22",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	println("代理添加完成")

	//rootWorker.RemoveNode(&Info{
	//	Id: fmt.Sprintf("node-%d", 5),
	//})

	var waitGroupForStop sync.WaitGroup

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}

func testStartServer(id, address, token, connAddress, connToken string) (server *Server) {
	server = &Server{
		Id:          id,
		BindAddress: address,
		BindToken:   token,
		ConnAddress: connAddress,
		ConnToken:   connToken,
	}
	err := server.Start()
	if err != nil {
		panic(err)
	}
	return
}

func TestNodeLine(t *testing.T) {
	nodeList := getTestNodeList()
	server := &Server{}
	err := server.Start()
	if err != nil {
		panic(err)
	}
	err = server.AddNodeList(nodeList)
	if err != nil {
		panic(err)
	}

	line := server.GetNodeLineByFromTo("node-2", "node-9")
	println("节点线:", fmt.Sprint(line))

}

func getTestNodeList() (nodeList []*Info) {
	port := 11000
	root := &Info{
		Id:          "root",
		BindAddress: fmt.Sprintf("127.0.0.1:%d", port),
		BindToken:   "root_token",
	}
	nodeList = append(nodeList, root)

	for n := 1; n <= 10; n++ {
		port++

		node := &Info{
			Id:          fmt.Sprintf("node-%d", n),
			ConnAddress: fmt.Sprintf("127.0.0.1:%d", port),
			BindToken:   fmt.Sprintf("node-%d-token", n),
		}
		if util.ContainsInt([]int{1, 2}, n) >= 0 {
			node.addConnNodeId(root.Id)
		} else if util.ContainsInt([]int{3, 4}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 2))
		} else if util.ContainsInt([]int{5, 6}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 3))
		} else if util.ContainsInt([]int{7, 8}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 1))
		} else if util.ContainsInt([]int{9, 10}, n) >= 0 {
			node.addConnNodeId(fmt.Sprintf("node-%d", 8))
		}
		nodeList = append(nodeList, node)

	}

	return
}
