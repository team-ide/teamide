package coders

import (
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

func NewApplicationCoder(app *modelers.Application, coderFactory common.ICoderFactory) (res *applicationCoder) {
	res = &applicationCoder{
		app:          app,
		coderFactory: coderFactory,
		appCoder:     coderFactory.GetApp(),
	}
	return
}

type applicationCoder struct {
	app          *modelers.Application
	coderFactory common.ICoderFactory
	appCoder     common.IAppCoder
}

func (this_ *applicationCoder) GenConstant(model *modelers.ConstantModel) (code *common.Code, err error) {
	code = &common.Code{}
	err = this_.coderFactory.GetConstant().Gen(code, model)
	return
}

func (this_ *applicationCoder) GenError(model *modelers.ErrorModel) (code *common.Code, err error) {
	code = &common.Code{}

	return
}

func (this_ *applicationCoder) GenStruct(model *modelers.StructModel) (code *common.Code, err error) {
	code = &common.Code{}

	return
}

func (this_ *applicationCoder) GenService(model *modelers.ServiceModel) (code *common.Code, err error) {
	code = &common.Code{}

	return
}

func (this_ *applicationCoder) GenDao(model *modelers.DaoModel) (code *common.Code, err error) {
	code = &common.Code{}

	return
}

func (this_ *applicationCoder) GenFunc(model *modelers.FuncModel) (code *common.Code, err error) {
	code = &common.Code{}

	return
}
