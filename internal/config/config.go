package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"strings"
	base2 "teamide/internal/server/base"
	"teamide/pkg/db"
	"teamide/pkg/util"
)

type ServerConfig struct {
	Server         *server            `json:"server,omitempty" yaml:"server,omitempty"`
	Mysql          *mysql             `json:"mysql,omitempty" yaml:"mysql,omitempty"`
	Log            *log               `json:"log,omitempty" yaml:"log,omitempty"`
	DatabaseConfig *db.DatabaseConfig `json:"-" yaml:"-"`
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

var (
	Config *ServerConfig = &ServerConfig{}
)

func init() {

	filePath := base2.BaseDir + "conf/config.yaml"
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
		formatMap(configMap)

		bs, err = json.Marshal(configMap)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, Config)
	} else {
		base2.IsStandAlone = true
	}
	formatConfig()
	fmt.Println(util.ToJSON(Config))
}

//formatConfig 格式化配置，填充默认值
func formatConfig() {

	if Config.Log == nil {
		Config.Log = &log{
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
			Level:      "info",
		}
	}
	if base2.IsStandAlone {
		Config.Server = &server{}
		if base2.UserHomeDir != "" {
			Config.Server.Data = base2.UserHomeDir + "TeamIDE/data"
			Config.Log.Filename = base2.UserHomeDir + "TeamIDE/log/server.log"
		} else {
			Config.Server.Data = base2.BaseDir + "data"
			Config.Log.Filename = base2.BaseDir + "log/server.log"
		}
		Config.DatabaseConfig = &db.DatabaseConfig{
			Type: "sqlite",
		}
		if base2.UserHomeDir != "" {
			Config.DatabaseConfig.Database = base2.UserHomeDir + "TeamIDE/data/database"
		} else {
			Config.DatabaseConfig.Database = base2.BaseDir + "data/database"
		}
	} else {
		if Config.Server == nil || Config.Server.Host == "" || Config.Server.Port == 0 {
			panic("请检查Server配置是否正确")
		}
		if Config.Mysql == nil || Config.Mysql.Host == "" || Config.Mysql.Port == 0 {
			panic("请检查MySql配置是否正确")
		}

		Config.DatabaseConfig = &db.DatabaseConfig{
			Type:     "mysql",
			Host:     Config.Mysql.Host,
			Port:     Config.Mysql.Port,
			Database: Config.Mysql.Database,
			Username: Config.Mysql.Username,
			Password: Config.Mysql.Password,
		}

		if Config.Server.Data == "" {
			Config.Server.Data = base2.BaseDir + "data"
		} else {
			Config.Server.Data = base2.BaseDir + strings.TrimPrefix(Config.Server.Data, "./")
		}
		if Config.Log.Filename == "" {
			Config.Log.Filename = base2.BaseDir + "log/server.log"
		} else {
			Config.Log.Filename = base2.BaseDir + strings.TrimPrefix(Config.Log.Filename, "./")
		}
	}

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
