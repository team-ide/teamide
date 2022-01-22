package config

import (
	"io"
	"os"
	"strings"
	"teamide/server/base"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	IsNative  bool       `json:"-" yaml:"-"` // 是否是本机运行
	Server    *server    `json:"server,omitempty" yaml:"server,omitempty"`
	Redis     *redis     `json:"redis,omitempty" yaml:"redis,omitempty"`
	Zookeeper *zookeeper `json:"zookeeper,omitempty" yaml:"zookeeper,omitempty"`
	Mysql     *mysql     `json:"mysql,omitempty" yaml:"mysql,omitempty"`
	Log       *log       `json:"log,omitempty" yaml:"log,omitempty"`
}

type server struct {
	Host    string `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int    `json:"port,omitempty" yaml:"port,omitempty"`
	Context string `json:"context,omitempty" yaml:"context,omitempty"`
	Data    string `json:"data,omitempty" yaml:"data,omitempty"`
}

type redis struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	Auth    string `json:"auth,omitempty" yaml:"auth,omitempty"`
	Prefix  string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

type zookeeper struct {
	Address   string `json:"address,omitempty" yaml:"address,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type mysql struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int32  `json:"port,omitempty" yaml:"port,omitempty"`
	Database string `json:"database,omitempty" yaml:"database,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type log struct {
	Filename   string `json:"filename,omitempty" yaml:"filename,omitempty"`
	MaxSize    int    `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`
	MaxAge     int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`
	MaxBackups int    `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty"`
	Level      string `json:"level,omitempty" yaml:"level,omitempty"`
}

var (
	Config *ServerConfig = &ServerConfig{}
)

func init() {

	filePath := base.BaseDir + "conf/config.yaml"
	exists, err := PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if exists {
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		bs, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(bs, Config)
		if err != nil {
			panic(err)
		}
	}
	formatConfig()
}

//格式化配置，填充默认值
func formatConfig() {

	if Config.Mysql == nil {
		Config.IsNative = true
	}
	if Config.Server != nil {
		if Config.Server.Data == "" {
			Config.Server.Data = base.BaseDir + "data"
		} else {
			Config.Server.Data = base.BaseDir + strings.TrimPrefix(Config.Server.Data, "./")
		}
	}
	if Config.Log == nil {
		Config.Log = &log{
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
			Level:      "info",
		}
	}

	if Config.Log.Filename == "" {
		Config.Log.Filename = base.BaseDir + "log/server.log"
	} else {
		Config.Log.Filename = base.BaseDir + strings.TrimPrefix(Config.Log.Filename, "./")
	}

	if !Config.IsNative {
		if Config.Mysql == nil || Config.Mysql.Host == "" || Config.Mysql.Port == 0 {
			panic("请检查MySql配置是否正确")
		}
		if Config.Redis == nil || Config.Redis.Address == "" {
			panic("请检查Redis配置是否正确")
		}
		if Config.Zookeeper == nil || Config.Zookeeper.Address == "" {
			panic("请检查Zookeeper配置是否正确")
		}
	}
}

func GetFromSystem(key string) string {
	return os.Getenv(key)
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
