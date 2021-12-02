package zookeeper

import (
	"config"
	"fmt"
)

var (
	ZKService ZookeeperService
)

func init() {
	var service interface{}
	var err error
	address := config.Config.Zookeeper.Address
	fmt.Println("zookeeper service init address:", address)
	service, err = CreateZookeeperService(address)
	if err != nil {
		panic(err)
	}
	ZKService = *service.(*ZookeeperService)

	_, err = ZKService.Exists("/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Zookeeper连接成功")
}
func Init() {
}
