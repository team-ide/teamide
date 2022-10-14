package common

import "teamide/pkg/maker/model"

type GeneratorFactory interface {
	GetApp() *AppGenerator
	GetConstant() []*ConstantGenerator
	GetStruct() []*StructGenerator
	GetService() []*ServiceGenerator
	GetDao() []*DaoGenerator
	GetError() []*ErrorGenerator
	GetFunc() []*FuncGenerator
	GetStep() *StepGenerator
}

type AppGenerator interface {
}

type ConstantGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.ConstantModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
}

type StructGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.StructModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
}

type ServiceGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.ServiceModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
	*StepGenerator
}

type DaoGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.DaoModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
}

type ErrorGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.ErrorModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
}

type FuncGenerator struct {
	App          AppGenerator
	Code         *Code
	Model        *model.FuncModel
	Gen          func() (err error)
	GetNamespace func() (namespace string)
}

type Code struct {
	Dir       string
	Name      string
	Namespace string
	Header    string
	Body      string
	Footer    string
	Tab       int
}
