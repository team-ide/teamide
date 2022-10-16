package common

import "teamide/pkg/maker/modelers"

type IAppCoder interface {
}

type IConstantCoder interface {
	Gen(code *Code, model *modelers.ConstantModel) (err error)
}

type IStructCoder interface {
	Gen(code *Code, model *modelers.StructModel) (err error)
}

type IServiceCoder interface {
	Gen(code *Code, model *modelers.ServiceModel) (err error)
}

type IDaoCoder interface {
	Gen(code *Code, model *modelers.DaoModel) (err error)
}

type IErrorCoder interface {
	Gen(code *Code, model *modelers.ErrorModel) (err error)
}

type IFuncCoder interface {
	Gen(code *Code, model *modelers.FuncModel) (err error)
}
