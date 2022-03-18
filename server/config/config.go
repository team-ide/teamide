package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"teamide/server/base"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
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
		configMap := map[string]interface{}{}
		err = yaml.Unmarshal(bs, &configMap)
		if err != nil {
			panic(err)
		}
		foramtMap(configMap)

		bs, err = json.Marshal(configMap)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, Config)
	} else {
		base.IS_STAND_ALONE = true
	}
	formatConfig()
	fmt.Println(base.ToJSON(Config))
}

//格式化配置，填充默认值
func formatConfig() {

	if Config.Server == nil {
		base.IS_STAND_ALONE = true
	}
	if Config.Server != nil {
		if Config.Server.Data == "" {
			if base.IS_STAND_ALONE && base.UserHomeDir != "" {
				Config.Server.Data = base.UserHomeDir + "TeamIDE/data"
			} else {
				Config.Server.Data = base.BaseDir + "data"
			}
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

		if base.IS_STAND_ALONE && base.UserHomeDir != "" {
			Config.Log.Filename = base.UserHomeDir + "TeamIDE/log/server.log"
		} else {
			Config.Log.Filename = base.BaseDir + "log/server.log"
		}
	} else {
		Config.Log.Filename = base.BaseDir + strings.TrimPrefix(Config.Log.Filename, "./")
	}

	if !base.IS_STAND_ALONE {
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

func foramtMap(mapValue map[string]interface{}) {
	if mapValue == nil {
		return
	}
	for key, value := range mapValue {
		switch v := value.(type) {
		case map[string]interface{}:
			foramtMap(v)
		default:
			res := foramtValue(value)
			mapValue[key] = res
		}
	}

}
func foramtValue(value interface{}) (res string) {
	if value == nil {
		return
	}
	stringValue, stringValueOk := value.(string)
	if stringValueOk {
		res = stringValue
		return
	}
	res = ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexsList := re.FindAllIndex([]byte(stringValue), -1)
	var lastIndex int = 0
	for _, indexs := range indexsList {
		res += stringValue[lastIndex:indexs[0]]

		lastIndex = indexs[1]

		key := stringValue[indexs[0]+1 : indexs[1]-1]
		value := GetFromSystem(key)
		if value == "" {
			return
		}
		res += value
	}
	res += stringValue[lastIndex:]

	return
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
