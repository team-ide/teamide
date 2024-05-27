package golang

import (
	"github.com/team-ide/go-tool/util"
	"teamide/pkg/maker/coder"
	"teamide/pkg/maker/modelers"
)

func FullGenerator(coder *coder.Coder) (err error) {
	res := &Generator{
		Coder: coder,

		spaceCache: make(map[string]*SpaceBuilder),
		packCache:  make(map[string]*PackBuilder),
		classCache: make(map[string]*ClassBuilder),
	}

	err = res.init()
	if err != nil {
		return
	}
	coder.SetGenerator(res)
	return
}

type Generator struct {
	*coder.Coder
	module string
	golang *modelers.LanguageGolangModel

	spaceCache map[string]*SpaceBuilder
	packCache  map[string]*PackBuilder
	classCache map[string]*ClassBuilder
}

func (this_ *Generator) init() (err error) {
	this_.golang = this_.GetLanguageGolang()
	if this_.Dir == "" {
		this_.Dir = this_.golang.Dir
	}
	return
}

func (this_ *Generator) GenMod() (err error) {
	path := this_.Dir + "go.mod"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	builder.AppendTabLine("module " + this_.golang.GetModuleName())
	builder.NewLine()
	builder.AppendTabLine("go " + this_.golang.GetGoVersion())
	builder.NewLine()

	builder.AppendTabLine("require (")
	builder.Tab()
	builder.AppendTabLine("github.com/team-ide/go-tool v1.2.14")
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	return
}

func (this_ *Generator) GenBase() (err error) {
	if err = this_.GenMod(); err != nil {
		return
	}
	return
}

func (this_ *Generator) GenConstant(builder *ClassBuilder) (err error) {

	builder.AppendTabLine("var(")

	builder.Tab()
	for _, one := range builder.FieldList {
		var str string
		str, err = this_.GetTypeStr(one.CompilerValueType.GetValueType())
		if err != nil {
			return
		}
		name := util.FirstToUpper(one.Name)
		builder.AppendTabLine("// " + name + " " + one.ConstantOption.Comment + "")
		if str == "string" {
			builder.AppendTabLine("" + name + " = \"" + one.ConstantOption.Value + "\"")
		} else {
			if str == "int" {
				str = ""
			}
			builder.AppendTabLine("" + name + " " + str + " = " + one.ConstantOption.Value)
		}
		builder.NewLine()
	}
	builder.Indent()
	builder.AppendTabLine(")")
	return
}

func (this_ *Generator) GenError(builder *ClassBuilder) (err error) {

	var imports []string

	commonPack := this_.golang.GetCommonPack()
	imports = append(imports, this_.golang.GetCommonImport())

	builder.AppendTabLine("import(")
	builder.Tab()
	for _, im := range imports {
		builder.AppendTabLine("\"" + im + "\"")
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	builder.AppendTabLine("var(")

	builder.Tab()
	for _, one := range builder.FieldList {
		name := util.FirstToUpper(one.Name)
		builder.AppendTabLine("// " + name + " " + one.ErrorOption.Comment + "")
		builder.AppendTabLine("" + name + " = " + commonPack + ".NewError(\"" + one.ErrorOption.Code + "\", \"" + one.ErrorOption.Msg + "\")")
		builder.NewLine()
	}
	builder.Indent()
	builder.AppendTabLine(")")
	return
}

func (this_ *Generator) GenStruct(builder *ClassBuilder) (err error) {

	var imports []string

	//commonPack := this_.golang.GetCommonPack()
	//imports = append(imports, this_.golang.GetCommonImport())

	builder.AppendTabLine("import(")
	builder.Tab()
	for _, im := range imports {
		builder.AppendTabLine("\"" + im + "\"")
	}
	builder.Indent()
	builder.AppendTabLine(")")
	builder.NewLine()

	structName := util.FirstToUpper(builder.Struct.Name)
	builder.AppendTabLine("// " + structName + " " + builder.Struct.Comment + "")
	builder.AppendTabLine("type " + structName + " struct {")

	builder.Tab()
	for _, one := range builder.FieldList {
		name := util.FirstToUpper(one.Name)
		builder.AppendTabLine("// " + name + " " + one.StructField.Comment + "")
		var str string
		str, err = this_.GetTypeStr(one.CompilerValueType.GetValueType())
		if err != nil {
			return
		}
		name_ := util.FirstToLower(one.Name)
		builder.AppendTabLine("" + name + " " + str + "`json:\"" + name_ + "\"`")
		builder.NewLine()
	}
	builder.Indent()
	builder.AppendTabLine("}")
	builder.NewLine()

	builder.AppendTabLine("// New" + structName + " 新建 " + structName + "对象")
	builder.AppendTabLine("func New" + structName + "() *" + structName + " { ")
	builder.Tab()
	builder.AppendTabLine("st := &" + structName + "{ ")
	builder.Tab()
	builder.Indent()
	builder.AppendTabLine("}")
	builder.AppendTabLine("return st")
	builder.Indent()
	builder.AppendTabLine("}")
	builder.NewLine()

	builder.AppendTabLine("// Copy 复制 " + structName + "对象")
	builder.AppendTabLine("func (this_ *" + structName + ") Copy() *" + structName + " { ")
	builder.Tab()
	builder.AppendTabLine("st := &" + structName + "{ ")
	builder.Tab()
	for _, one := range builder.FieldList {
		name := util.FirstToUpper(one.Name)
		builder.AppendTabLine("" + name + " : this_." + name + ",")
	}
	builder.Indent()
	builder.AppendTabLine("}")
	builder.AppendTabLine("return st")
	builder.Indent()
	builder.AppendTabLine("}")
	builder.NewLine()

	return
}
