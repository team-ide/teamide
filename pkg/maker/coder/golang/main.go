package golang

import (
	"github.com/team-ide/go-tool/util"
	"sort"
	"strings"
)

var (
	mainCode = `
func Init() {
	var err error

	config.RootDir, err = os.Getwd()
	if err != nil {
		util.Logger.Error("os get wd error", zap.Error(err))
		panic(err)
	}

	config.RootDir, err = filepath.Abs(config.RootDir)
	if err != nil {
		util.Logger.Error("filepath abs error", zap.Error(err))
		panic(err)
	}
	config.RootDir = filepath.ToSlash(config.RootDir)
	if !strings.HasSuffix(config.RootDir, "/") {
		config.RootDir += "/"
	}
	current, err := user.Current()
	if err != nil {
		util.Logger.Error("user current error", zap.Error(err))
		panic(err)
	}

	config.UserHomeDir = current.HomeDir
	if config.UserHomeDir != "" {
		config.UserHomeDir, err = filepath.Abs(config.UserHomeDir)
		if err != nil {
			util.Logger.Error("filepath abs error", zap.Error(err))
			panic(err)
		}
		config.UserHomeDir = filepath.ToSlash(config.UserHomeDir)
		if !strings.HasSuffix(config.UserHomeDir, "/") {
			config.UserHomeDir += "/"
		}

	}
}

func main() {
	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			println("app version:" + config.Version)
			println("Go os:" + runtime.GOOS)
			println("Go arch:" + runtime.GOARCH)
			println("Go compiler:" + runtime.Compiler)
			println("Go version:" + runtime.Version())
			return
		}
	}
	var err error
	var waitGroupForStop sync.WaitGroup

	defer func() {
		if e := recover(); e != nil {
			err = errors.New("奔溃异常:" + fmt.Sprint(e))
		}
		if err!=nil {
			fmt.Println("启动失败:", err)
			util.Logger.Error("启动失败", zap.Error(err))
		}
		waitGroupForStop.Done()

	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c)
		for s := range c {
			switch s {
			case os.Kill: // kill -9 pid，下面的无效
				fmt.Println("强制退出", s)
				common.OnStop()
				os.Exit(0)
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT: // ctrl + c
				fmt.Println("退出", s)
				common.OnStop()
				os.Exit(0)
			}
		}
	}()

	waitGroupForStop.Add(1)

	var conf = flag.String("config", "conf/application.yml", "配置文件地址")

	flag.Parse()

	err = config.InitConfig(*conf)
	if err != nil {
		util.Logger.Error("初始化配置文件失败", zap.Error(err))
		return
	}

	logger.Init(config.GetConfig().Log)

{component_init}

}
`
)

func (this_ *Generator) appendMainInit(code *string, imports *[]string, componentType string, modelName string) {
	*code += "\t"

	importName := this_.golang.GetComponentImport(componentType, modelName)
	pack := this_.golang.GetComponentPack(componentType, modelName)
	var configName = componentType
	if modelName != "default" {
		configName += "_" + modelName
	}
	configName = util.FirstToUpper(configName)
	*code += `
	if config.GetConfig().` + configName + ` == nil {
		util.Logger.Error("配置 ` + componentType + ` 为空，请检查配置")
		return
	}
`
	*code += `
	if err = ` + pack + `.Init(config.GetConfig().` + configName + `); err != nil {
		util.Logger.Error("初始化 ` + componentType + ` 失败", zap.Error(err))
		return
	}
`
	var find bool
	for _, s := range *imports {
		if s == importName {
			find = true
		}
	}
	if !find {
		*imports = append(*imports, importName)
	}
}
func (this_ *Generator) GenMain() (err error) {
	dir := this_.Dir
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "main.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string
	component_init := ""
	for _, one := range this_.GetConfigDbList() {
		this_.appendMainInit(&component_init, &imports, "db", one.Name)
	}
	for _, one := range this_.GetConfigRedisList() {
		this_.appendMainInit(&component_init, &imports, "redis", one.Name)
	}
	for _, one := range this_.GetConfigZkList() {
		this_.appendMainInit(&component_init, &imports, "zk", one.Name)
	}
	for _, one := range this_.GetConfigKafkaList() {
		this_.appendMainInit(&component_init, &imports, "kafka", one.Name)
	}
	for _, one := range this_.GetConfigEsList() {
		this_.appendMainInit(&component_init, &imports, "es", one.Name)
	}
	for _, one := range this_.GetConfigMongodbList() {
		this_.appendMainInit(&component_init, &imports, "mongodb", one.Name)
	}

	builder.AppendTabLine("package main")
	builder.NewLine()

	imports = append(imports, this_.golang.GetConfigImport())
	imports = append(imports, this_.golang.GetLoggerImport())
	imports = append(imports, this_.golang.GetCommonImport())
	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"errors"
	"flag"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
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

	code := strings.ReplaceAll(mainCode, "{component_init}", component_init)

	builder.AppendCode(code)
	return
}
