package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
)

type ServerConfig struct {
	Server          *server `json:"server,omitempty" yaml:"server,omitempty"`
	Mysql           *mysql  `json:"mysql,omitempty" yaml:"mysql,omitempty"`
	Log             *log    `json:"log,omitempty" yaml:"log,omitempty"`
	Github          *Github `json:"github,omitempty" yaml:"github,omitempty"`
	LogDataSaveDays int     `json:"logDataSaveDays,omitempty" yaml:"logDataSaveDays,omitempty"`
}

type server struct {
	Host       string     `json:"host,omitempty" yaml:"host,omitempty"`
	Port       int        `json:"port,omitempty" yaml:"port,omitempty"`
	Context    string     `json:"context,omitempty" yaml:"context,omitempty"`
	Data       string     `json:"data,omitempty" yaml:"data,omitempty"`
	TLS        *ServerTLS `json:"tls,omitempty" yaml:"tls,omitempty"`
	CertKey    string     `json:"certKey,omitempty" yaml:"certKey,omitempty"`
	BackupsDir string     `json:"-" yaml:"-"`
	TempDir    string     `json:"-" yaml:"-"`
}

type ServerTLS struct {
	Open bool   `json:"open,omitempty" yaml:"open,omitempty"`
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}
type mysql struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Database string `json:"database,omitempty" yaml:"database,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}
type Github struct {
	ClientId     string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
}

type log struct {
	Filename   string `json:"filename,omitempty" yaml:"filename,omitempty"`
	MaxSize    int    `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`
	MaxAge     int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`
	MaxBackups int    `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty"`
	Level      string `json:"level,omitempty" yaml:"level,omitempty"`
}

func CreateServerConfig(configPath string) (config *ServerConfig, err error) {

	config = &ServerConfig{
		LogDataSaveDays: 15,
	}
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
			return
		}
		var bs []byte
		bs, err = io.ReadAll(f)
		if err != nil {
			util.Logger.Error("yaml read error", zap.Any("path", configPath), zap.Error(err))
			return
		}
		configMap := map[string]interface{}{}
		err = yaml.Unmarshal(bs, &configMap)
		if err != nil {
			util.Logger.Error("yaml to map error", zap.Any("path", configPath), zap.Error(err))
			return
		}
		formatMap(configMap)

		if configMap["logDataSaveDays"] == nil {
			configMap["logDataSaveDays"] = 15
		}

		bs, err = json.Marshal(configMap)
		if err != nil {
			util.Logger.Error("config map to bytes error", zap.Any("configMap", configMap), zap.Error(err))
			return
		}
		err = yaml.Unmarshal(bs, config)
		if err != nil {
			util.Logger.Error("config bytes to config error", zap.Any("config", string(bs)), zap.Error(err))
			return
		}
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
		case []interface{}:
			for i, one := range v {
				switch oneV := one.(type) {
				case map[string]interface{}:
					formatMap(oneV)
					v[i] = oneV
				default:
					res := formatValue(oneV)
					v[i] = res
				}
			}
		default:
			res := formatValue(value)
			mapValue[key] = res
		}
	}

}
func formatValue(value interface{}) (v interface{}) {
	if value == nil {
		return
	}
	stringValue, stringValueOk := value.(string)
	if !stringValueOk {
		v = value
		return
	}
	res := ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexList := re.FindAllIndex([]byte(stringValue), -1)
	var lastIndex int = 0
	for _, indexes := range indexList {
		res += stringValue[lastIndex:indexes[0]]

		lastIndex = indexes[1]

		key := stringValue[indexes[0]+2 : indexes[1]-1]
		envValue := GetFromSystem(key)
		if value == "" {
			return
		}
		res += envValue
	}
	res += stringValue[lastIndex:]

	v = res
	return
}

func GetFromSystem(key string) string {
	return os.Getenv(key)
}

/*
PathExists

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
