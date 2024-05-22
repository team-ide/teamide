package coder

import (
	"errors"
	"github.com/team-ide/go-tool/util"
	"os"
	"strings"
	"teamide/pkg/maker"
	"teamide/pkg/maker/modelers"
)

func NewCoder(compiler *maker.Compiler, options *Options) (res *Coder, err error) {
	if options == nil {
		options = &Options{}
	}
	res = &Coder{
		Compiler: compiler,
		Options:  options,
	}

	err = res.init()

	return
}

type Options struct {
	Dir string `json:"dir"`
}

type Coder struct {
	*maker.Compiler
	*Options
	generator IGenerator
}

func (this_ *Coder) init() (err error) {
	return
}

func (this_ *Coder) Gen() (err error) {
	var generator = this_.generator
	if generator == nil {
		err = errors.New("generator为空，请配置生成器")
		return
	}
	if this_.Dir == "" {
		err = errors.New("请配置生成目录")
		return
	}
	this_.Dir = util.FormatPath(this_.Dir)
	if !strings.HasSuffix(this_.Dir, "/") {
		this_.Dir += "/"
	}
	err = this_.Mkdir(this_.Dir)
	if err != nil {
		return
	}
	err = generator.GenBase()
	if err != nil {
		return
	}
	for _, one := range this_.GetConstantList() {
		err = generator.GenConstant(one)
		if err != nil {
			return
		}
	}
	return
}
func (this_ *Coder) SetGenerator(generator IGenerator) {
	this_.generator = generator
	return
}
func (this_ *Coder) Mkdir(path string) (err error) {
	ex, err := util.PathExists(path)
	if err != nil {
		return
	}
	if !ex {
		err = os.MkdirAll(path, os.ModePerm)
	}
	return
}
func (this_ *Coder) CreateAndOpen(path string) (f *os.File, err error) {
	ex, err := util.PathExists(path)
	if err != nil {
		return
	}
	if !ex {
		f, err = os.Create(path)
	} else {
		f, err = os.Create(path)
	}
	return
}

func (this_ *Coder) AppendLine(content *string, line string, tab int) {
	for i := 0; i < tab; i++ {
		*content += "    "
	}
	*content += line
	*content += "\n"
}

type IGenerator interface {
	GenBase() (err error)
	GenConstant(model *modelers.ConstantModel) (err error)
}
