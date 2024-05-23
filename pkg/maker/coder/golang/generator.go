package golang

import (
	"github.com/team-ide/go-tool/util"
	"strings"
	"teamide/pkg/maker/coder"
	"teamide/pkg/maker/modelers"
)

func FullGenerator(coder *coder.Coder) (err error) {
	res := &Generator{
		Coder: coder,
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

	if err = builder.AppendLine("module " + this_.golang.GetModuleName()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}
	if err = builder.AppendLine("go " + this_.golang.GetGoVersion()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}
	return
}

func (this_ *Generator) GenBase() (err error) {
	if err = this_.GenMod(); err != nil {
		return
	}
	if err = this_.GenCommon(); err != nil {
		return
	}
	return
}

var (
	commonCode = `package {pack}

import "fmt"

type Error struct {
	code string
	msg  string
}

func (this_ *Error) Error() string {
	return fmt.Sprintf("code:%s , msg:%s", this_.code, this_.msg)
}

// NewError 构造异常对象，code为错误码，msg为错误信息
func NewError(code string, msg string) *Error {
	err := &Error{
		code: code,
		msg:  msg,
	}
	return err
}

`
)

func (this_ *Generator) GenCommon() (err error) {
	dir := this_.golang.GetCommonDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "common.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	commonCode = strings.ReplaceAll(commonCode, "{pack}", this_.golang.GetCommonPack())

	if err = builder.AppendCode(commonCode); err != nil {
		return
	}
	return
}

func (this_ *Generator) GenConstant(model *modelers.ConstantModel) (err error) {
	dir := this_.golang.GetConstantDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + model.Name + ".go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	if err = builder.AppendLine("package " + this_.golang.GetStructPack()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	if err = builder.AppendLine("var("); err != nil {
		return
	}

	builder.Tab()
	for _, one := range model.Options {
		var str string
		str, err = this_.GetTypeStr(one.Type)
		if err != nil {
			return
		}
		name := util.FirstToUpper(one.Name)
		if err = builder.AppendLine("// " + name + " " + one.Comment + ""); err != nil {
			return
		}
		if str == "string" {
			if err = builder.AppendLine("" + name + " = \"" + one.Value + "\""); err != nil {
				return
			}
		} else {
			if str == "int" {
				str = ""
			}
			if err = builder.AppendLine("" + name + " " + str + " = " + one.Value); err != nil {
				return
			}
		}
		if err = builder.NewLine(); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine(")"); err != nil {
		return
	}
	return
}

func (this_ *Generator) GenError(model *modelers.ErrorModel) (err error) {
	dir := this_.golang.GetErrorDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + model.Name + ".go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	if err = builder.AppendLine("package " + this_.golang.GetErrorPack()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	var imports []string

	commonPack := this_.golang.GetCommonPack()
	imports = append(imports, this_.golang.GetCommonImport())

	if err = builder.AppendLine("import("); err != nil {
		return
	}
	builder.Tab()
	for _, im := range imports {
		if err = builder.AppendLine("\"" + im + "\""); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine(")"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	if err = builder.AppendLine("var("); err != nil {
		return
	}

	builder.Tab()
	for _, one := range model.Options {
		name := util.FirstToUpper(one.Name)
		if err = builder.AppendLine("// " + name + " " + one.Comment + ""); err != nil {
			return
		}
		if err = builder.AppendLine("" + name + " = " + commonPack + ".NewError(\"" + one.Code + "\", \"" + one.Msg + "\")"); err != nil {
			return
		}
		if err = builder.NewLine(); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine(")"); err != nil {
		return
	}
	return
}

func (this_ *Generator) GenStruct(model *modelers.StructModel) (err error) {
	dir := this_.golang.GetStructDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + model.Name + ".go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	if err = builder.AppendLine("package " + this_.golang.GetStructPack()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	var imports []string

	//commonPack := this_.golang.GetCommonPack()
	//imports = append(imports, this_.golang.GetCommonImport())

	if err = builder.AppendLine("import("); err != nil {
		return
	}
	builder.Tab()
	for _, im := range imports {
		if err = builder.AppendLine("\"" + im + "\""); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine(")"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	structName := util.FirstToUpper(model.Name)
	if err = builder.AppendLine("// " + structName + " " + model.Comment + ""); err != nil {
		return
	}
	if err = builder.AppendLine("type " + structName + " struct {"); err != nil {
		return
	}

	builder.Tab()
	for _, one := range model.Fields {
		name := util.FirstToUpper(one.Name)
		if err = builder.AppendLine("// " + name + " " + one.Comment + ""); err != nil {
			return
		}
		var str string
		str, err = this_.GetTypeStr(one.Type)
		if err != nil {
			return
		}
		name_ := util.FirstToLower(one.Name)
		if err = builder.AppendLine("" + name + " " + str + "`json:\"" + name_ + "\"`"); err != nil {
			return
		}
		if err = builder.NewLine(); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine("}"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	if err = builder.AppendLine("// New" + structName + " 新建 " + structName + "对象"); err != nil {
		return
	}
	if err = builder.AppendLine("func New" + structName + "() *" + structName + " { "); err != nil {
		return
	}
	builder.Tab()
	if err = builder.AppendLine("st := &" + structName + "{ "); err != nil {
		return
	}
	builder.Tab()
	builder.Indent()
	if err = builder.AppendLine("}"); err != nil {
		return
	}
	if err = builder.AppendLine("return st"); err != nil {
		return
	}
	builder.Indent()
	if err = builder.AppendLine("}"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	if err = builder.AppendLine("// Copy 复制 " + structName + "对象"); err != nil {
		return
	}
	if err = builder.AppendLine("func (this_ *" + structName + ") Copy() *" + structName + " { "); err != nil {
		return
	}
	builder.Tab()
	if err = builder.AppendLine("st := &" + structName + "{ "); err != nil {
		return
	}
	builder.Tab()
	for _, one := range model.Fields {
		name := util.FirstToUpper(one.Name)
		if err = builder.AppendLine("" + name + " : this_." + name + ","); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine("}"); err != nil {
		return
	}
	if err = builder.AppendLine("return st"); err != nil {
		return
	}
	builder.Indent()
	if err = builder.AppendLine("}"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	return
}

func (this_ *Generator) GenFunc(funcPath string, models []*modelers.FuncModel) (err error) {
	dir := this_.golang.GetFuncDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	fileName := funcPath
	if fileName == "" {
		fileName = "tool"
	}
	path := dir + fileName + ".go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	if err = builder.AppendLine("package " + this_.golang.GetFuncPack()); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	var imports []string

	//commonPack := this_.golang.GetCommonPack()
	//imports = append(imports, this_.golang.GetCommonImport())

	if err = builder.AppendLine("import("); err != nil {
		return
	}
	builder.Tab()
	for _, im := range imports {
		if err = builder.AppendLine("\"" + im + "\""); err != nil {
			return
		}
	}
	builder.Indent()
	if err = builder.AppendLine(")"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}

	for _, model := range models {

		funcName := util.FirstToUpper(model.Name)
		if err = builder.AppendLine("// " + funcName + " " + model.Comment + ""); err != nil {
			return
		}
		if err = builder.AppendLine("func " + funcName + "() {"); err != nil {
			return
		}

		builder.Tab()

		builder.Indent()
		if err = builder.AppendLine("}"); err != nil {
			return
		}
		if err = builder.NewLine(); err != nil {
			return
		}

	}
	return
}
