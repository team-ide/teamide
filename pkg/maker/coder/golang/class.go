package golang

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sort"
	"strings"
	"teamide/pkg/maker"
	"teamide/pkg/maker/coder"
)

type SpaceBuilder struct {
	*Generator
	*maker.CompilerSpace
	spaceDir    string
	spacePath   string
	spacePack   string
	spaceImport string
}

type PackBuilder struct {
	*SpaceBuilder
	*maker.CompilerPack
	packDir  string
	packPath string
	packPack string
}

type ClassBuilder struct {
	*coder.Builder
	*maker.CompilerClass
	*PackBuilder
	filePath      string
	className     string
	classBeanName string

	classFileName string
	iFilePath     string
	iBuilder      *coder.Builder

	implDir    string
	implPath   string
	implPack   string
	implImport string
}

func (this_ *Generator) getSpaceBuilder(space *maker.CompilerSpace) (builder *SpaceBuilder) {

	builder = this_.spaceCache[space.GetKey()]
	if builder != nil {
		return
	}
	builder = &SpaceBuilder{
		CompilerSpace: space,
		Generator:     this_,
	}
	switch space.Space {
	case "constant":
		builder.spacePath = this_.golang.GetConstantPath()
		builder.spacePack = this_.golang.GetConstantPack()
		builder.spaceDir = this_.golang.GetConstantDir(this_.Dir)
		break
	case "error":
		builder.spacePath = this_.golang.GetErrorPath()
		builder.spacePack = this_.golang.GetErrorPack()
		builder.spaceDir = this_.golang.GetErrorDir(this_.Dir)
		break
	case "struct":
		builder.spacePath = this_.golang.GetStructPath()
		builder.spacePack = this_.golang.GetStructPack()
		builder.spaceDir = this_.golang.GetStructDir(this_.Dir)
		break
	case "func":
		builder.spacePath = this_.golang.GetFuncIFacePath()
		builder.spacePack = this_.golang.GetFuncIFacePack()
		builder.spaceDir = this_.golang.GetFuncIFaceDir(this_.Dir)
		builder.spaceImport = this_.golang.GetFuncIFaceImport()

		break
	case "dao":
		builder.spacePath = this_.golang.GetDaoIFacePath()
		builder.spacePack = this_.golang.GetDaoIFacePack()
		builder.spaceDir = this_.golang.GetDaoIFaceDir(this_.Dir)
		builder.spaceImport = this_.golang.GetDaoIFaceImport()
		break
	case "service":
		builder.spacePath = this_.golang.GetServiceIFacePath()
		builder.spacePack = this_.golang.GetServiceIFacePack()
		builder.spaceDir = this_.golang.GetServiceIFaceDir(this_.Dir)
		builder.spaceImport = this_.golang.GetServiceIFaceImport()
		break
	default:
		panic("space [" + space.Space + "] 不支持")
	}
	this_.spaceCache[space.GetKey()] = builder
	return
}

func (this_ *Generator) getPackBuilder(pack *maker.CompilerPack) (builder *PackBuilder) {
	builder = this_.packCache[pack.GetKey()]
	if builder != nil {
		return
	}
	spaceBuilder := this_.getSpaceBuilder(pack.CompilerSpace)

	builder = &PackBuilder{
		SpaceBuilder: spaceBuilder,
		CompilerPack: pack,
	}
	builder.packDir = builder.spaceDir
	builder.packPath = builder.spacePath
	builder.packPack = builder.spacePack
	this_.packCache[pack.GetKey()] = builder
	return
}

func (this_ *Generator) getClassBuilder(class *maker.CompilerClass) (builder *ClassBuilder) {
	builder = this_.classCache[class.GetKey()]
	if builder != nil {
		return
	}
	packBuilder := this_.getPackBuilder(class.CompilerPack)

	builder = &ClassBuilder{
		PackBuilder:   packBuilder,
		CompilerClass: class,
	}
	builder.filePath = builder.packDir

	if class.Pack != "" {
		builder.classFileName += strings.ReplaceAll(class.Pack, ".", "_")
		if class.Class != nil {
			builder.classFileName += "_" + strings.Join(class.Class, "_")
		}
	} else {
		if class.Class != nil {
			builder.classFileName += strings.Join(class.Class, "_")
		} else {
			builder.classFileName += builder.packPack
		}
	}
	builder.filePath = builder.packDir + builder.classFileName + ".go"

	if class.Space == "func" {
		builder.implDir = this_.golang.GetFuncImplDir(this_.Dir, builder.classFileName)
		builder.implPath = this_.golang.GetFuncImplPath(builder.classFileName)
		builder.implPack = this_.golang.GetFuncImplPack(builder.classFileName)
		builder.implImport = this_.golang.GetFuncImplImport(builder.classFileName)
	} else if class.Space == "dao" {
		builder.implDir = this_.golang.GetDaoImplDir(this_.Dir, builder.classFileName)
		builder.implPath = this_.golang.GetDaoImplPath(builder.classFileName)
		builder.implPack = this_.golang.GetDaoImplPack(builder.classFileName)
		builder.implImport = this_.golang.GetDaoImplImport(builder.classFileName)
	} else if class.Space == "service" {
		builder.implDir = this_.golang.GetServiceImplDir(this_.Dir, builder.classFileName)
		builder.implPath = this_.golang.GetServiceImplPath(builder.classFileName)
		builder.implPack = this_.golang.GetServiceImplPack(builder.classFileName)
		builder.implImport = this_.golang.GetServiceImplImport(builder.classFileName)
	}

	if builder.implDir != "" {
		builder.filePath = builder.implDir + "impl.go"
		builder.iFilePath = builder.packDir + builder.classFileName + ".go"
	}

	this_.classCache[class.GetKey()] = builder
	return
}

func (this_ *Generator) GenSpace(space *maker.CompilerSpace) (err error) {
	for _, one := range space.PackList {
		err = this_.GenPack(one)
	}
	return
}

func (this_ *Generator) GenPack(pack *maker.CompilerPack) (err error) {
	for _, one := range pack.ClassList {
		err = this_.GenClass(one)
	}

	return
}

func (this_ *ClassBuilder) GetClassName() (res string) {
	res = this_.className
	if res != "" {
		return
	}
	for _, name := range this_.Class {
		res += util.FirstToUpper(name)
	}
	res += util.FirstToUpper(this_.spacePack)
	res = util.FirstToUpper(res)
	this_.className = res
	return
}

func (this_ *ClassBuilder) GetImplClassName() (res string) {
	switch this_.CompilerClass.Space {
	case "dao":
		res = "Dao"
		return
	case "service":
		res = "Service"
		return
	case "func":
		res = "Tool"
		return
	}
	res = this_.GetClassName()
	return
}

func (this_ *ClassBuilder) GetClassBeanName() (res string) {
	res = this_.classBeanName
	if res != "" {
		return
	}
	res = this_.GetClassName()
	res += ""
	this_.classBeanName = res
	return
}

func (this_ *Generator) GetImportAsNameFromValueType(v *maker.ValueType) (impl string, asName string) {
	if v == nil {
		return
	}
	if v == maker.ValueTypeContext {
		return this_.GetImportAsName("context")
	} else if v.Struct != nil {
		return this_.GetImportAsName("struct")
	}
	return
}
func (this_ *Generator) GetImportAsName(name string) (impl string, asName string) {
	switch name {
	case "logger":
		impl = this_.golang.GetLoggerImport()
		asName = this_.golang.GetLoggerPack()
		break
	case "common":
		impl = this_.golang.GetCommonImport()
		asName = this_.golang.GetCommonPack()
		break
	case "constant":
		impl = this_.golang.GetConstantImport()
		asName = this_.golang.GetConstantPack()
		break
	case "error":
		impl = this_.golang.GetErrorImport()
		asName = this_.golang.GetErrorPack()
		break
	case "struct":
		impl = this_.golang.GetStructImport()
		asName = this_.golang.GetStructPack()
		break
	case "func":
		impl = this_.golang.GetFuncIFaceImport()
		asName = this_.golang.GetFuncIFacePack()
		break
	case "dao":
		impl = this_.golang.GetDaoIFaceImport()
		asName = this_.golang.GetDaoIFacePack()
		break
	case "service":
		impl = this_.golang.GetServiceIFaceImport()
		asName = this_.golang.GetServiceIFacePack()
		break
	case "util":
		impl = "github.com/team-ide/go-tool/util"
		asName = "util"
		break
	case "context":
		impl = "context"
		asName = "context"
		break
	case "fmt":
		impl = "fmt"
		asName = "fmt"
		break
	default:
		var componentType string
		var componentName = "default"
		if strings.HasPrefix(name, "component_") {
			ss := strings.Split(strings.TrimPrefix(name, "component_"), "_")
			componentType = ss[0]
			componentName = "default"
			if len(ss) > 1 {
				componentName = ss[1]
			}
		} else {
			var componentTypes = []string{"db", "redis", "zk", "kafka", "es", "mongodb"}
			for _, s := range componentTypes {
				if name == s {
					componentType = s
				} else if strings.HasPrefix(name, s+"_") {
					componentType = s
					componentName = strings.TrimPrefix(name, s+"_")
				}
			}
		}
		if componentType != "" {
			impl = this_.golang.GetComponentImport(componentType, componentName)
			asName = this_.golang.GetComponentPack(componentType, componentName)
		}
	}
	return
}
func (this_ *Generator) GenClass(class *maker.CompilerClass) (err error) {
	builder := this_.getClassBuilder(class)
	util.Logger.Debug("gen "+class.GetKey(), zap.Any("path", builder.filePath))
	builder.Builder, err = this_.NewBuilder(builder.filePath)
	if err != nil {
		return
	}
	defer builder.Close()

	if builder.implPack != "" {
		builder.AppendTabLine("package " + builder.implPack)
	} else {
		builder.AppendTabLine("package " + builder.packPack)
	}
	builder.NewLine()

	if class.Constant != nil {
		err = this_.GenConstant(builder)
	} else if class.Error != nil {
		err = this_.GenError(builder)
	} else if class.Struct != nil {
		err = this_.GenStruct(builder)
	} else {
		if builder.iFilePath != "" {
			err = builder.GenIFace()
			if err != nil {
				return
			}
		}
		err = builder.GenImpl()
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	return
}

func (this_ *ClassBuilder) GenImpl() (err error) {
	var imports []string

	for _, impl := range this_.ImportList {
		if impl.Import != "" {
			implPath, asName := this_.GetImportAsName(impl.Import)
			if implPath != "" {
				imports = append(imports, implPath)
				impl.AsName = asName
			}
		}
	}

	this_.AppendTabLine("import(")
	this_.Tab()
	sort.Strings(imports)
	for _, im := range imports {
		this_.AppendTabLine("\"" + im + "\"")
	}
	this_.Indent()
	this_.AppendTabLine(")")
	this_.NewLine()

	//this_.AppendTabLine("// ", this_.GetClassBeanName(), " ", this_.GetClassName(), "对象实例")
	//this_.AppendTabLine("var ", this_.GetClassBeanName(), " = New", this_.GetClassName(), "()")
	//this_.NewLine()

	this_.AppendTabLine("// New 新建", this_.GetImplClassName(), "对象实例")
	this_.AppendTabLine("func New() (res *", this_.GetImplClassName(), ") {")
	this_.Tab()
	this_.AppendTabLine("res = &" + this_.GetImplClassName() + "{}")
	this_.AppendTabLine("return")
	this_.Indent()
	this_.AppendTabLine("}")
	this_.NewLine()

	this_.AppendTabLine("// ", this_.GetImplClassName(), " 接口 I", this_.GetClassName(), " 实现")
	this_.AppendTabLine("type ", this_.GetImplClassName(), " struct {")
	this_.NewLine()
	this_.AppendTabLine("}")
	this_.NewLine()

	for _, method := range this_.MethodList {

		methodBuilder := &MethodBuilder{
			ClassBuilder:   this_,
			CompilerMethod: method,
		}
		err = methodBuilder.Gen()
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ClassBuilder) GenIFace() (err error) {
	this_.iFaceClassList = append(this_.iFaceClassList, this_)
	util.Logger.Debug("gen "+this_.GetKey(), zap.Any("path", this_.iFilePath))
	this_.iBuilder, err = this_.NewBuilder(this_.iFilePath)
	if err != nil {
		return
	}
	defer this_.iBuilder.Close()

	this_.iBuilder.AppendTabLine("package " + this_.packPack)
	this_.iBuilder.NewLine()

	var imports []string
	var importCache = map[string]string{}

	for _, method := range this_.MethodList {
		var valueTypes []*maker.ValueType
		for _, arg := range method.ParamList {
			valueTypes = append(valueTypes, arg.GetValueType())
		}
		valueTypes = append(valueTypes, method.Result.GetValueType())
		for _, v := range valueTypes {
			implPath, _ := this_.Generator.GetImportAsNameFromValueType(v)
			if implPath != "" && importCache[implPath] == "" {
				importCache[implPath] = implPath
				imports = append(imports, implPath)
			}
		}
	}

	this_.iBuilder.AppendTabLine("import(")
	this_.iBuilder.Tab()
	sort.Strings(imports)
	for _, im := range imports {
		this_.iBuilder.AppendTabLine("\"" + im + "\"")
	}
	this_.iBuilder.Indent()
	this_.iBuilder.AppendTabLine(")")
	this_.iBuilder.NewLine()

	this_.iBuilder.AppendTabLine("// ", this_.GetClassBeanName(), " I", this_.GetClassName(), " 接口实现")
	this_.iBuilder.AppendTabLine("var ", this_.GetClassBeanName(), " I"+this_.GetClassName())
	this_.iBuilder.NewLine()

	this_.iBuilder.AppendTabLine("type I", this_.GetClassName(), " interface {")
	this_.iBuilder.NewLine()

	this_.iBuilder.Tab()
	for _, method := range this_.MethodList {

		methodName := util.FirstToUpper(method.Method)
		this_.iBuilder.AppendTab()
		this_.iBuilder.AppendCode("// " + methodName + " ")
		this_.iBuilder.AppendComment(method.Comment)
		this_.iBuilder.NewLine()
		var str string
		str += "" + methodName
		str += "("
		for i, param := range method.ParamList {

			var typeS string
			typeS, err = this_.GetTypeStr(param.CompilerValueType.GetValueType())
			if err != nil {
				return
			}
			if i > 0 {
				str += ", "
			}

			str += param.Name + " " + typeS
		}
		str += ")"
		str += " ("
		if method.Result.GetValueType() != nil {
			var typeS string
			typeS, err = this_.GetTypeStr(method.Result.GetValueType())
			if err != nil {
				return
			}
			str += "res " + typeS
			str += ", "
		}
		str += "err error"
		str += ")"

		this_.iBuilder.AppendTabLine(str)
	}

	this_.iBuilder.Indent()

	this_.iBuilder.AppendTabLine("}")
	this_.iBuilder.NewLine()

	return
}
