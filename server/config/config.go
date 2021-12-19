package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Server    *server   `json:"server" toml:"server,omitempty"`
	Redis     redis     `json:"redis" toml:"redis,omitempty"`
	Zookeeper zookeeper `json:"zookeeper" toml:"zookeeper,omitempty"`
	Kafka     kafka     `json:"kafka" toml:"kafka,omitempty"`
	Mysql     mysql     `json:"mysql" toml:"mysql,omitempty"`
	Log       log       `json:"log" toml:"log,omitempty"`
}

type server struct {
	Host    string `toml:"server,omitempty"`
	Port    int    `json:"port,omitempty"`
	Context string `json:"context,omitempty"`
	Data    string `json:"data,omitempty"`
}

type redis struct {
	Address string `json:"address,omitempty"`
	Auth    string `json:"auth,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
}

type zookeeper struct {
	Address   string `json:"address,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type kafka struct {
	Address string `json:"address,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
}

type mysql struct {
	Host     string `json:"host,omitempty"`
	Port     int32  `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type log struct {
	Filename   string `json:"filename,omitempty"`
	MaxSize    int    `json:"maxSize,omitempty"`
	MaxAge     int    `json:"maxAge,omitempty"`
	MaxBackups int    `json:"maxBackups,omitempty"`
	Level      string `json:"level,omitempty"`
}

var (
	Config *TomlConfig
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path + "/./conf/config.toml"
	exists, err := PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("配置文件[" + filePath + "]不存在")
	}
	if _, err := toml.DecodeFile(filePath, &Config); err != nil {
		panic(err)
	}
	Config = formatConfig(Config)
}

//格式化配置，填充默认值
func formatConfig(config *TomlConfig) *TomlConfig {
	if config == nil {
		config = &TomlConfig{}
	}
	if config.Server == nil {
		config.Server = &server{}
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 19000
	}
	if config.Server.Context == "" {
		config.Server.Context = "/teamide"
	}
	if config.Server.Data == "" {
		config.Server.Data = "./data"
	}
	return config
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
