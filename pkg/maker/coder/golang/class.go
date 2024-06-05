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
	spaceDir  string
	spacePath string
	spacePack string
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
		builder.spacePath = this_.golang.GetFuncPath()
		builder.spacePack = this_.golang.GetFuncPack()
		builder.spaceDir = this_.golang.GetFuncDir(this_.Dir)
		break
	case "dao":
		builder.spacePath = this_.golang.GetDaoPath()
		builder.spacePack = this_.golang.GetDaoPack()
		builder.spaceDir = this_.golang.GetDaoDir(this_.Dir)
		break
	case "service":
		builder.spacePath = this_.golang.GetServicePath()
		builder.spacePack = this_.golang.GetServicePack()
		builder.spaceDir = this_.golang.GetServiceDir(this_.Dir)
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
		builder.filePath += strings.ReplaceAll(class.Pack, ".", "_")
		if class.Class != "" {
			builder.filePath += "_" + class.Class
		}
	} else {
		if class.Class != "" {
			builder.filePath += class.Class
		} else {
			builder.filePath += builder.packPack
		}
	}
	builder.filePath += ".go"

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
	res = this_.Class + util.FirstToUpper(this_.spacePack)
	res = util.FirstToUpper(res)
	this_.className = res
	return
}

func (this_ *ClassBuilder) GetClassBeanName() (res string) {
	res = this_.classBeanName
	if res != "" {
		return
	}
	res = this_.GetClassName()
	res += "Obj"
	this_.classBeanName = res
	return
}

func (this_ *Generator) GetImportAsName(name string) (impl string, asName string) {
	switch name {
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
		impl = this_.golang.GetFuncImport()
		asName = this_.golang.GetFuncPack()
		break
	case "dao":
		impl = this_.golang.GetDaoImport()
		asName = this_.golang.GetDaoPack()
		break
	case "service":
		impl = this_.golang.GetServiceImport()
		asName = this_.golang.GetServicePack()
		break
	case "util":
		impl = "github.com/team-ide/go-tool/util"
		asName = "util"
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

	builder.AppendTabLine("package " + builder.packPack)
	builder.NewLine()

	if class.Constant != nil {
		err = this_.GenConstant(builder)
	} else if class.Error != nil {
		err = this_.GenError(builder)
	} else if class.Struct != nil {
		err = this_.GenStruct(builder)
	} else {

		var imports []string
		imports = append(imports, "context")

		for _, impl := range class.ImportList {
			if impl.Import != "" {
				implPath, asName := this_.GetImportAsName(impl.Import)
				if implPath != "" {
					imports = append(imports, implPath)
					impl.AsName = asName
				}
			}
		}

		builder.AppendTabLine("import(")
		builder.Tab()
		sort.Strings(imports)
		for _, im := range imports {
			builder.AppendTabLine("\"" + im + "\"")
		}
		builder.Indent()
		builder.AppendTabLine(")")
		builder.NewLine()

		builder.AppendTabLine("// ", builder.GetClassBeanName(), " ", builder.GetClassName(), "对象实例")
		builder.AppendTabLine("var ", builder.GetClassBeanName(), " = New", builder.GetClassName(), "()")
		builder.NewLine()

		builder.AppendTabLine("// New", builder.GetClassName(), " 新建", builder.GetClassName(), "对象实例")
		builder.AppendTabLine("func New", builder.GetClassName(), "() (res ", builder.GetClassName(), ") {")
		builder.Tab()
		builder.NewLine()
		builder.AppendTabLine("return")
		builder.Indent()
		builder.AppendTabLine("}")
		builder.NewLine()

		builder.AppendTabLine("type ", builder.GetClassName(), " struct {")
		builder.NewLine()
		builder.AppendTabLine("}")
		builder.NewLine()

		for _, method := range class.MethodList {

			methodBuilder := &MethodBuilder{
				ClassBuilder:   builder,
				CompilerMethod: method,
			}

			err = methodBuilder.Gen()
			if err != nil {
				return
			}
		}
	}
	if err != nil {
		return
	}
	return
}
