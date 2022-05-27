package zookeeper

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"go.uber.org/zap"
	"sort"
	"strings"
	"teamide/pkg/util"
	"time"
)

type Config struct {
	Address string `json:"address"`
}

func CreateZKService(config Config) (*ZKService, error) {
	service := &ZKService{
		Address: config.Address,
	}
	err := service.init()
	return service, err
}

//ZKService 注册处理器在线信息等
type ZKService struct {
	Address     string
	zkConn      *zk.Conn        //zk连接
	zkConnEvent <-chan zk.Event // zk事件通知管道
	lastUseTime int64
}

func (this_ *ZKService) init() error {
	var err error
	this_.zkConn, this_.zkConnEvent, err = zk.Connect(this_.GetServers(), time.Second*60, func(c *zk.Conn) {
		c.SetLogger(defaultLogger{})
	})
	return err
}

type defaultLogger struct{}

func (defaultLogger) Printf(format string, args ...interface{}) {
	util.Logger.Info(format, zap.Any("args", args))
}

func (this_ *ZKService) GetServers() []string {
	var servers []string
	if strings.Contains(this_.Address, ",") {
		servers = strings.Split(this_.Address, ",")
	} else if strings.Contains(this_.Address, ";") {
		servers = strings.Split(this_.Address, ";")
	} else {
		servers = []string{this_.Address}
	}
	return servers
}
func (this_ *ZKService) GetConn() *zk.Conn {
	defer func() {
		this_.lastUseTime = util.GetNowTime()
	}()
	return this_.zkConn
}

func (this_ *ZKService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *ZKService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *ZKService) Stop() {
	this_.GetConn().Close()
}

//Create 创建节点
func (this_ *ZKService) Create(path string, data []byte, mode int32) (err error) {
	isExist, err := this_.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("node:" + path + " already exists")
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = this_.CreateIfNotExists(parentPath, []byte{})
		if err != nil {
			return err
		}
	}
	if _, err = this_.GetConn().Create(path, data, mode, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (this_ *ZKService) SetData(path string, data []byte) (err error) {
	isExist, state, err := this_.GetConn().Exists(path)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("node:" + path + " not exists")
	}
	if _, err = this_.GetConn().Set(path, data, state.Version); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//CreateIfNotExists 一层层检查，如果不存在则创建父节点
func (this_ *ZKService) CreateIfNotExists(path string, data []byte) (err error) {
	isExist, err := this_.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = this_.CreateIfNotExists(parentPath, []byte{})
		if err != nil {
			return err
		}
	}
	if _, err = this_.GetConn().Create(path, data, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//Exists 判断节点是否存在
func (this_ *ZKService) Exists(path string) (isExist bool, err error) {
	isExist, _, err = this_.GetConn().Exists(path)
	return
}

//Get 判断节点是否存在
func (this_ *ZKService) Get(path string) (data []byte, err error) {
	data, _, err = this_.GetConn().Get(path)
	return
}

//GetChildren 判断节点是否存在
func (this_ *ZKService) GetChildren(path string) (children []string, err error) {
	children, _, err = this_.GetConn().Children(path)
	sort.Strings(children)
	return
}

//Delete 判断节点是否存在
func (this_ *ZKService) Delete(path string) (err error) {
	var isExist bool
	var stat *zk.Stat
	isExist, stat, err = this_.GetConn().Exists(path)
	if !isExist {
		return
	}
	var children []string
	children, _, err = this_.GetConn().Children(path)
	if err != nil {
		return
	}
	if len(children) > 0 {
		for _, one := range children {
			err = this_.Delete(path + "/" + one)
			if err != nil {
				return
			}
		}
	}
	err = this_.GetConn().Delete(path, stat.Version)
	return
}
