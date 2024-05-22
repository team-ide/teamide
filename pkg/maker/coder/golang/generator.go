package golang

import (
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
	return
}

func (this_ *Generator) GenConstant(model *modelers.ConstantModel) (err error) {
	constantDir := this_.golang.GetConstantDir(this_.Dir)
	if err = this_.Mkdir(constantDir); err != nil {
		return
	}
	path := constantDir + model.Name + ".go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	if err = builder.AppendLine("package constant"); err != nil {
		return
	}
	if err = builder.NewLine(); err != nil {
		return
	}
	return
}
