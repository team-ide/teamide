package golang

import (
	"github.com/team-ide/go-tool/util"
	"sort"
	"strings"
)

var (
	mainCode = `
func main() {
	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			println("Release version : " + common.ReleaseVersion)
			println("   Release time : " + common.ReleaseTime)
			println("     Git commit : " + common.GitCommit)
			println("          Go os : " + runtime.GOOS)
			println("        Go arch : " + runtime.GOARCH)
			println("     Go version : " + runtime.Version())
			return
		}
	}
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("奔溃异常:" + fmt.Sprint(e))
		}
		if err != nil {
			fmt.Println("启动失败:", err)
			logger.Logger.Error("启动失败", zap.Error(err))
		}

	}()

	go common.OnSignal()

	var conf = flag.String("config", "conf/application.yml", "配置文件地址")
	flag.Parse()
	
	err = start.InitApp(*conf)
	if err != nil {
		return
	}
	
	err = start.RunServer()
	if err != nil {
		return
	}

	common.Wait()
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
		logger.Logger.Error("配置 ` + componentType + ` 为空，请检查配置")
		return
	}
`
	*code += `
	if err = ` + pack + `.Init(config.GetConfig().` + configName + `); err != nil {
		logger.Logger.Error("初始化 ` + componentType + ` 失败", zap.Error(err))
		return
	}
`
	addImport(imports, importName)
}

func addImport(imports *[]string, importName string) {
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

	builder.AppendTabLine("package main")
	builder.NewLine()

	var imports []string

	imports = append(imports, this_.golang.GetLoggerImport())
	imports = append(imports, this_.golang.GetCommonImport())
	imports = append(imports, this_.golang.GetStartImport())
	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"flag"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"runtime"
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
		as := strings.Split(im, " ")
		if len(as) == 2 {
			builder.AppendTabLine(as[1] + " \"" + as[0] + "\"")
		} else {
			builder.AppendTabLine("\"" + im + "\"")
		}
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	code := mainCode

	builder.AppendCode(code)
	err = this_.GenStart()
	if err != nil {
		return
	}
	return
}
