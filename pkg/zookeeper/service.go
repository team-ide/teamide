package zookeeper

import "github.com/go-zookeeper/zk"

type Service interface {
	GetConn() *zk.Conn
	CreateIfNotExists(path string, data []byte) (err error)
}
