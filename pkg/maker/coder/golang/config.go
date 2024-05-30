package golang

import (
	"github.com/team-ide/go-tool/util"
	"sort"
	"strings"
)

var (
	configCode = `
var (
	config     *Config
	configPath string

	RootDir     = ""
	UserHomeDir = ""
	Version     = "0.0.1"
	onStops     []func()
)

func AddOnStop(onStop func()) {
	onStops = append(onStops, onStop)
}

func OnStop() {
	for _, onStop := range onStops {
		callOnStop(onStop)
	}
}

func callOnStop(onStop func()) {
	defer func() {
		if e := recover(); e != nil {
			err := errors.New(fmt.Sprint(e))
			fmt.Println("callOnStop:", err)
			util.Logger.Error("callOnStop", zap.Error(err))
		}
	}()
	onStop()
}

{config}

type Log struct {
	Console    bool   ` + "`" + `json:"console,omitempty" yaml:"console,omitempty"` + "`" + `
	Filename   string ` + "`" + `json:"filename,omitempty" yaml:"filename,omitempty"` + "`" + `
	MaxSize    int    ` + "`" + `json:"maxSize,omitempty" yaml:"maxSize,omitempty"` + "`" + `
	MaxAge     int    ` + "`" + `json:"maxAge,omitempty" yaml:"maxAge,omitempty"` + "`" + `
	MaxBackups int    ` + "`" + `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty"` + "`" + `
	Level      string ` + "`" + `json:"level,omitempty" yaml:"level,omitempty"` + "`" + `
}

func GetConfig() *Config {
	return config
}

func InitConfig(conf string) (err error) {

	config = &Config{}
	configPath = conf
	var exists bool
	exists, err = util.PathExists(configPath)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("配置文件[", configPath, "]不存在"))
		util.Logger.Error("配置文件不存在", zap.Error(err))
		return
	}
	var f *os.File
	f, err = os.Open(configPath)
	if err != nil {
		return
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		util.Logger.Error("配置文件["+configPath+"]读取异常", zap.Error(err))
		return
	}
	configMap := map[string]interface{}{}
	err = yaml.Unmarshal(bs, &configMap)
	if err != nil {
		util.Logger.Error("配置文件["+configPath+"] yaml转map异常", zap.Error(err))
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
	re, _ = regexp.Compile(` + "`" + `[$]+{(.+?)}` + "`" + `)
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

`
)

func (this_ *Generator) appendConfig(code *string, imports *[]string, modelType string, modelName string, importName string) {
	*code += "\t"
	var name = modelType
	if modelName != "default" {
		name += "_" + modelName
	}
	fieldName := util.FirstToUpper(name)
	*code += fieldName + " *" + importName + ".Config  " + "`json:\"" + name + ",omitempty\" yaml:\"" + name + ",omitempty\"`"
	var find bool
	var imp = "github.com/team-ide/go-tool/" + importName
	for _, s := range *imports {
		if s == imp {
			find = true
		}
	}
	if !find {
		*imports = append(*imports, "github.com/team-ide/go-tool/"+importName)
	}
	*code += "\n"
}
func (this_ *Generator) GenConfig() (err error) {
	dir := this_.golang.GetConfigDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "config.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string

	configPack := this_.golang.GetConfigPack()

	configStruct := "type Config struct{" + "\n"
	for _, one := range this_.GetConfigDbList() {
		this_.appendConfig(&configStruct, &imports, "db", one.Name, "db")
	}
	for _, one := range this_.GetConfigRedisList() {
		this_.appendConfig(&configStruct, &imports, "redis", one.Name, "redis")
	}
	for _, one := range this_.GetConfigZkList() {
		this_.appendConfig(&configStruct, &imports, "zk", one.Name, "zookeeper")
	}
	for _, one := range this_.GetConfigKafkaList() {
		this_.appendConfig(&configStruct, &imports, "kafka", one.Name, "kafka")
	}
	for _, one := range this_.GetConfigEsList() {
		this_.appendConfig(&configStruct, &imports, "es", one.Name, "elasticsearch")
	}
	for _, one := range this_.GetConfigMongodbList() {
		this_.appendConfig(&configStruct, &imports, "mongodb", one.Name, "mongodb")
	}
	configStruct += "\t" + "Log   *Log  `" + `json:"log,omitempty" yaml:"log,omitempty"` + "`" + "\n"
	configStruct += "}"

	builder.AppendTabLine("package " + configPack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
`, "\n")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		imports = append(imports, s)
	}

	sort.Strings(imports)
	for _, im := range imports {
		builder.AppendTabLine("\"" + im + "\"")
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := strings.ReplaceAll(configCode, "{pack}", this_.golang.GetConfigPack())
	code = strings.ReplaceAll(code, "{config}", configStruct)

	builder.AppendCode(code)
	return
}
