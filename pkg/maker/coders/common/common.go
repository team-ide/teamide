package common

import "teamide/pkg/maker/modelers"

type ICoderFactory interface {
	GetApp() IAppCoder
	GetConstant() IConstantCoder
	GetStruct() IStructCoder
	GetService() IServiceCoder
	GetDao() IDaoCoder
	GetError() IErrorCoder
	GetFunc() IFuncCoder
	GetStep() IStepCoder
}

type IAppCoder interface {
}

type IConstantCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.ConstantModel) (err error)
}

type IStructCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.StructModel) (err error)
}

type IServiceCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.ServiceModel) (err error)
}

type IDaoCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.DaoModel) (err error)
}

type IErrorCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.ErrorModel) (err error)
}

type IFuncCoder interface {
	Gen(appCoder IAppCoder, code *Code, model *modelers.FuncModel) (err error)
}

type Code struct {
	Dir       string
	Name      string
	Namespace string
	Header    string
	Body      string
	Footer    string
	Tab       int
	Other     map[string]*Code
}

func (this_ *Code) GetOther(key string) (other *Code) {
	if this_.Other == nil {
		this_.Other = make(map[string]*Code)
	}
	other = this_.Other[key]
	return
}

func (this_ *Code) SetOther(key string, other *Code) {
	if this_.Other == nil {
		this_.Other = make(map[string]*Code)
	}
	this_.Other[key] = other
	return
}
