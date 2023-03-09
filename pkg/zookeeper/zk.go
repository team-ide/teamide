package zookeeper

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Address  string `json:"address"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func CreateZKService(config Config) (service *ZKService, err error) {
	service = &ZKService{
		address:  config.Address,
		username: config.Username,
		password: config.Password,
	}
	err = service.init()
	return
}

// ZKService 注册处理器在线信息等
type ZKService struct {
	address     string
	username    string
	password    string
	zkConn      *zk.Conn        //zk连接
	zkEvent     <-chan zk.Event // zk事件通知管道
	lastUseTime int64
}

func (this_ *ZKService) GetZkEvent() <-chan zk.Event {
	return this_.zkEvent
}
func (this_ *ZKService) init() (err error) {
	this_.zkConn, this_.zkEvent, err = zk.Connect(this_.GetServers(), time.Second*10, func(c *zk.Conn) {
		c.SetLogger(defaultLogger{})
	})
	if err != nil {
		util.Logger.Error("zk.Connect error", zap.Any("servers", this_.GetServers()), zap.Error(err))
		if this_.zkConn != nil {
			this_.zkConn.Close()
		}
		return
	}
	if this_.username != "" || this_.password != "" {
		err = this_.zkConn.AddAuth(this_.username, []byte(this_.password))
		if err != nil {
			util.Logger.Error("zk.Connect AddAuth error", zap.Any("servers", this_.GetServers()), zap.Error(err))
			this_.zkConn.Close()
			return
		}
	}
	return
}

type defaultLogger struct{}

func (defaultLogger) Printf(format string, args ...interface{}) {
	util.Logger.Info(format, zap.Any("args", args))
}

func (this_ *ZKService) GetServers() []string {
	var servers []string
	if strings.Contains(this_.address, ",") {
		servers = strings.Split(this_.address, ",")
	} else if strings.Contains(this_.address, ";") {
		servers = strings.Split(this_.address, ";")
	} else {
		servers = []string{this_.address}
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
func (this_ *ZKService) SetLastUseTime() {
	this_.lastUseTime = util.GetNowTime()
}

func (this_ *ZKService) Stop() {
	this_.GetConn().Close()
}

// Create 创建节点
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

type Info struct {
	Server    string   `json:"server"`
	SessionID int64    `json:"sessionID"`
	State     zk.State `json:"state"`
}

func (this_ *ZKService) Info() (info *Info, err error) {
	info = &Info{}
	info.SessionID = this_.GetConn().SessionID()
	info.Server = this_.GetConn().Server()
	info.State = this_.GetConn().State()
	return
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

// CreateIfNotExists 一层层检查，如果不存在则创建父节点
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

// Exists 判断节点是否存在
func (this_ *ZKService) Exists(path string) (isExist bool, err error) {
	isExist, _, err = this_.GetConn().Exists(path)
	return
}

type StatInfo struct {
	Czxid          int64 `json:"czxid,omitempty"`
	Mzxid          int64 `json:"mzxid,omitempty"`
	Ctime          int64 `json:"ctime,omitempty"`
	Mtime          int64 `json:"mtime,omitempty"`
	Version        int32 `json:"version,omitempty"`
	Cversion       int32 `json:"cversion,omitempty"`
	Aversion       int32 `json:"aversion,omitempty"`
	EphemeralOwner int64 `json:"ephemeralOwner,omitempty"`
	DataLength     int32 `json:"dataLength,omitempty"`
	NumChildren    int32 `json:"numChildren,omitempty"`
	Pzxid          int64 `json:"pzxid,omitempty"`
}

type NodeInfo struct {
	Path string    `json:"path"`
	Data string    `json:"data"`
	Stat *StatInfo `json:"stat"`
}

// Get 判断节点是否存在
func (this_ *ZKService) Get(path string) (info *NodeInfo, err error) {

	data, stat, err := this_.GetConn().Get(path)
	if err != nil {
		return
	}
	info = &NodeInfo{}
	info.Data = string(data)

	if stat != nil {
		info.Stat = &StatInfo{
			Czxid:          stat.Czxid,
			Mzxid:          stat.Mzxid,
			Ctime:          stat.Ctime,
			Mtime:          stat.Mtime,
			Version:        stat.Version,
			Cversion:       stat.Cversion,
			Aversion:       stat.Aversion,
			EphemeralOwner: stat.EphemeralOwner,
			DataLength:     stat.DataLength,
			NumChildren:    stat.NumChildren,
			Pzxid:          stat.Pzxid,
		}
	}
	return
}

// Stat 判断节点是否存在
func (this_ *ZKService) Stat(path string) (info *StatInfo, err error) {
	_, stat, err := this_.GetConn().Exists(path)
	if err != nil {
		return
	}
	if stat != nil {
		info = &StatInfo{
			Czxid:          stat.Czxid,
			Mzxid:          stat.Mzxid,
			Ctime:          stat.Ctime,
			Mtime:          stat.Mtime,
			Version:        stat.Version,
			Cversion:       stat.Cversion,
			Aversion:       stat.Aversion,
			EphemeralOwner: stat.EphemeralOwner,
			DataLength:     stat.DataLength,
			NumChildren:    stat.NumChildren,
			Pzxid:          stat.Pzxid,
		}
	}
	return
}

// GetChildren 判断节点是否存在
func (this_ *ZKService) GetChildren(path string) (children []string, err error) {
	children, _, err = this_.GetConn().Children(path)
	sort.Strings(children)
	return
}

// Delete 判断节点是否存在
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
