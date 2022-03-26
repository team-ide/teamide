package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"teamide/pkg/util"
)

type ServerConfig struct {
	Server *server `json:"server,omitempty" yaml:"server,omitempty"`
	Mysql  *mysql  `json:"mysql,omitempty" yaml:"mysql,omitempty"`
	Log    *log    `json:"log,omitempty" yaml:"log,omitempty"`
}

type server struct {
	Host    string `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int    `json:"port,omitempty" yaml:"port,omitempty"`
	Context string `json:"context,omitempty" yaml:"context,omitempty"`
	Data    string `json:"data,omitempty" yaml:"data,omitempty"`
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

func CreateServerConfig(configPath string) (config *ServerConfig, err error) {

	//filePath := base.RootDir + "conf/config.yaml"
	config = &ServerConfig{}
	if configPath != "" {
		var exists bool
		exists, err = util.PathExists(configPath)
		if err != nil {
			panic(err)
		}
		if !exists {
			err = errors.New(fmt.Sprint("服务配置文件[", configPath, "]不存在"))
			return
		}
		var f *os.File
		f, err = os.Open(configPath)
		if err != nil {
			panic(err)
		}
		var bs []byte
		bs, err = io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		configMap := map[string]interface{}{}
		err = yaml.Unmarshal(bs, &configMap)
		if err != nil {
			panic(err)
		}
		formatMap(configMap)

		bs, err = json.Marshal(configMap)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, config)
	}

	if config.Server == nil {
		config.Server = &server{}
	}
	if config.Log == nil {
		config.Log = &log{
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
			Level:      "info",
		}
	}
	return
}

func formatMap(mapValue map[string]interface{}) {
	if mapValue == nil {
		return
	}
	for key, value := range mapValue {
		switch v := value.(type) {
		case map[string]interface{}:
			formatMap(v)
		default:
			res := formatValue(value)
			mapValue[key] = res
		}
	}

}
func formatValue(value interface{}) (res string) {
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
	indexList := re.FindAllIndex([]byte(stringValue), -1)
	var lastIndex int = 0
	for _, indexes := range indexList {
		res += stringValue[lastIndex:indexes[0]]

		lastIndex = indexes[1]

		key := stringValue[indexes[0]+1 : indexes[1]-1]
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

/*PathExists
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
