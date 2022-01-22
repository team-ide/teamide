package component

import (
	"errors"
	"strings"
	"teamide/server/config"
	"time"

	"github.com/go-zookeeper/zk"
)

var (
	Zookeeper ZookeeperService
)

func init() {

	if config.Config.IsNative {
		return
	}

	var service interface{}
	var err error
	address := config.Config.Zookeeper.Address
	Logger.Info(LogStr("Zookeeper初始化:address:", address))
	service, err = CreateZookeeperService(address)
	if err != nil {
		panic(err)
	}
	Zookeeper = *service.(*ZookeeperService)

	_, err = Zookeeper.Exists("/")
	if err != nil {
		panic(err)
	}
	Logger.Info(LogStr("Zookeeper连接成功!"))
}

//注册处理器在线信息等
type ZookeeperService struct {
	conn      *zk.Conn        //zk连接
	connEvent <-chan zk.Event // zk事件通知管道
}

func CreateZookeeperService(address string) (*ZookeeperService, error) {
	service := &ZookeeperService{}
	err := service.init(address)
	return service, err
}

func (service *ZookeeperService) init(address string) error {
	var err error
	service.conn, service.connEvent, err = zk.Connect(strings.Split(address, ","), time.Second*3, zk.WithLogInfo(false))
	return err
}

//创建节点
func (service *ZookeeperService) Create(path string, data []byte, mode int32) (err error) {
	isExist, err := service.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("node:" + path + " already exists")
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = service.CreateIfNotExists(parentPath, []byte{})
		if err != nil {
			return err
		}
	}
	if _, err = service.conn.Create(path, data, mode, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (service *ZookeeperService) SetData(path string, data []byte) (err error) {
	isExist, state, err := service.conn.Exists(path)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("node:" + path + " not exists")
	}
	if _, err = service.conn.Set(path, data, state.Version); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//一层层检查，如果不存在则创建父节点
func (service *ZookeeperService) CreateIfNotExists(path string, data []byte) (err error) {
	isExist, err := service.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = service.CreateIfNotExists(parentPath, data)
		if err != nil {
			return err
		}
	}
	if _, err = service.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//判断节点是否存在
func (service *ZookeeperService) Exists(path string) (isExist bool, err error) {
	isExist, _, err = service.conn.Exists(path)
	return
}

//判断节点是否存在
func (service *ZookeeperService) Get(path string) (data []byte, err error) {
	data, _, err = service.conn.Get(path)
	return
}

//判断节点是否存在
func (service *ZookeeperService) GetChildren(path string) (children []string, err error) {
	children, _, err = service.conn.Children(path)
	return
}

//判断节点是否存在
func (service *ZookeeperService) Delete(path string) (err error) {
	var isExist bool
	var stat *zk.Stat
	isExist, stat, err = service.conn.Exists(path)
	if !isExist {
		return
	}
	var children []string
	children, _, err = service.conn.Children(path)
	if err != nil {
		return
	}
	if len(children) > 0 {
		for _, one := range children {
			err = service.Delete(path + "/" + one)
			if err != nil {
				return
			}
		}
	}
	err = service.conn.Delete(path, stat.Version)
	return
}
