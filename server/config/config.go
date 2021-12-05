package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Server    *server   `toml:"server"`
	Redis     redis     `toml:"redis"`
	Zookeeper zookeeper `toml:"zookeeper"`
	Kafka     kafka     `toml:"kafka"`
	Mysql     mysql     `toml:"mysql"`
}

type server struct {
	Host    string
	Port    int
	Context string
	Data    string
}

type redis struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Prefix  string `json:"prefix"`
}

type zookeeper struct {
	Address   string `json:"address"`
	Namespace string `json:"namespace"`
}

type kafka struct {
	Address string `json:"address"`
	Prefix  string `json:"prefix"`
}

type mysql struct {
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
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
