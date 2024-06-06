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

	compileErrors := this_.Compile(false)
	if len(compileErrors) > 0 {
		err = errors.New(compileErrors[0].Method.GetKey() + " error:" + compileErrors[0].Err.Error())
		return
	}

	err = generator.GenBase()
	if err != nil {
		return
	}

	err = generator.GenCommon()
	if err != nil {
		return
	}

	for _, one := range this_.GetConfigDbList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentDb(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.GetConfigRedisList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentRedis(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.GetConfigZkList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentZk(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.GetConfigEsList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentEs(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.GetConfigKafkaList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentKafka(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.GetConfigMongodbList() {
		name := one.Name
		if name == "default" {
			name = ""
		}
		err = generator.GenComponentMongodb(name, one)
		if err != nil {
			return
		}
	}

	for _, one := range this_.SpaceList {
		err = generator.GenSpace(one)
		if err != nil {
			return
		}
	}

	err = generator.GenMain()
	if err != nil {
		return
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
	dir := path[0:strings.LastIndex(path, "/")]
	err = this_.Mkdir(dir)
	if err != nil {
		return
	}
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
	GenCommon() (err error)
	GenComponentDb(name string, model *modelers.ConfigDbModel) (err error)
	GenComponentRedis(name string, model *modelers.ConfigRedisModel) (err error)
	GenComponentZk(name string, model *modelers.ConfigZkModel) (err error)
	GenComponentKafka(name string, model *modelers.ConfigKafkaModel) (err error)
	GenComponentEs(name string, model *modelers.ConfigEsModel) (err error)
	GenComponentMongodb(name string, model *modelers.ConfigMongodbModel) (err error)
	GenSpace(space *maker.CompilerSpace) (err error)
	GenMain() (err error)
}
