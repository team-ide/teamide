package common

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
