package golang

import (
	"github.com/team-ide/go-tool/util"
	"sort"
	"strings"
)

var (
	startCode = `
var (
	appInit bool
)

func InitApp(conf string) (err error) {
	if appInit {
		return
	}
	appInit = true

	initBase()

	common.CallEvent(common.EventAppStart)

	common.CallEvent(common.EventAppConfigBefore)
	err = config.InitConfig(conf)
	if err != nil {
		logger.Logger.Error("初始化配置文件失败", zap.Error(err))
		return
	}
	common.CallEvent(common.EventAppConfigAfter)

	logger.Init(config.GetConfig().Log)

	common.CallEvent(common.EventAppComponentBefore)
	err = initComponent()
	if err != nil {
		logger.Logger.Error("初始化组件失败", zap.Error(err))
		return
	}
	common.CallEvent(common.EventAppComponentAfter)

	common.CallEvent(common.EventAppIFaceBefore)
	err = initIFace()
	if err != nil {
		logger.Logger.Error("初始化接口失败", zap.Error(err))
		return
	}
	common.CallEvent(common.EventAppIFaceAfter)

	common.CallEvent(common.EventAppReady)

	return
}

func initBase() {
	var err error

	common.RootDir, err = os.Getwd()
	if err != nil {
		logger.Logger.Error("os get wd error", zap.Error(err))
		panic(err)
	}

	common.RootDir, err = filepath.Abs(common.RootDir)
	if err != nil {
		logger.Logger.Error("filepath abs error", zap.Error(err))
		panic(err)
	}
	common.RootDir = filepath.ToSlash(common.RootDir)
	if !strings.HasSuffix(common.RootDir, "/") {
		common.RootDir += "/"
	}
	current, err := user.Current()
	if err != nil {
		logger.Logger.Error("user current error", zap.Error(err))
		panic(err)
	}

	common.UserHomeDir = current.HomeDir
	if common.UserHomeDir != "" {
		common.UserHomeDir, err = filepath.Abs(common.UserHomeDir)
		if err != nil {
			logger.Logger.Error("filepath abs error", zap.Error(err))
			panic(err)
		}
		common.UserHomeDir = filepath.ToSlash(common.UserHomeDir)
		if !strings.HasSuffix(common.UserHomeDir, "/") {
			common.UserHomeDir += "/"
		}
	}
}

func initComponent()(err error){
	
{componentContent}
	return
}

func initIFace()(err error){
	
{iFaceContent}
	return
}

func RunServer()(err error) {

	return
}
`
)

func (this_ *Generator) GenStart() (err error) {
	dir := this_.golang.GetStartDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "start.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	var imports []string
	imports = append(imports, this_.golang.GetLoggerImport())
	imports = append(imports, this_.golang.GetConfigImport())
	imports = append(imports, this_.golang.GetCommonImport())

	componentContent := ""
	for _, one := range this_.GetConfigDbList() {
		this_.appendMainInit(&componentContent, &imports, "db", one.Name)
	}
	for _, one := range this_.GetConfigRedisList() {
		this_.appendMainInit(&componentContent, &imports, "redis", one.Name)
	}
	for _, one := range this_.GetConfigZkList() {
		this_.appendMainInit(&componentContent, &imports, "zk", one.Name)
	}
	for _, one := range this_.GetConfigKafkaList() {
		this_.appendMainInit(&componentContent, &imports, "kafka", one.Name)
	}
	for _, one := range this_.GetConfigEsList() {
		this_.appendMainInit(&componentContent, &imports, "es", one.Name)
	}
	for _, one := range this_.GetConfigMongodbList() {
		this_.appendMainInit(&componentContent, &imports, "mongodb", one.Name)
	}

	iFaceContent := ""
	for _, one := range this_.iFaceClassList {
		iFaceContent += "\t"

		var asName = one.spacePack + "_" + one.implPack
		asNames := strings.Split(asName, "_")
		asName = ""
		for i, n := range asNames {
			if i > 0 {
				asName += util.FirstToUpper(n)
			} else {
				asName += n
			}
		}

		iFaceContent += "// 初始化 接口 I" + one.GetClassName() + " 实现" + "\n"
		iFaceContent += "\t"
		iFaceContent += one.spacePack + "." + one.GetClassBeanName() + " = " + asName + ".New()" + "\n"

		iFaceContent += "\n"
		addImport(&imports, one.implImport+" "+asName)
		addImport(&imports, one.spaceImport)
	}

	pack := this_.golang.GetStartPack()

	builder.AppendTabLine("package " + pack)
	builder.NewLine()

	builder.AppendTabLine("import(")
	builder.Tab()

	ss := strings.Split(`
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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

	code := strings.ReplaceAll(startCode, "{componentContent}", componentContent)
	code = strings.ReplaceAll(code, "{iFaceContent}", iFaceContent)

	builder.AppendCode(code)
	return
}
