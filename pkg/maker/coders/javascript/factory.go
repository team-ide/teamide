package javascript

import (
	"teamide/pkg/maker"
	"teamide/pkg/maker/coders/common"
)

func NewFactory(app *maker.Application) (res common.ICoderFactory) {
	appCoder_ := &appCoder{Application: app, codeType: common.CodeTypeJs, LanguageJavascriptModel: app.GetLanguageJavascript()}
	res = &factory{
		appCoder:      appCoder_,
		constantCoder: &constantCoder{appCoder: appCoder_},
		errorCoder:    &errorCoder{appCoder: appCoder_},
		structCoder:   &structCoder{appCoder: appCoder_},
		serviceCoder:  &serviceCoder{appCoder: appCoder_},
		daoCoder:      &daoCoder{appCoder: appCoder_},
		funcCoder:     &funcCoder{appCoder: appCoder_},
		stepCoder:     &stepCoder{appCoder: appCoder_},
	}
	return
}

type factory struct {
	appCoder      *appCoder
	constantCoder *constantCoder
	errorCoder    *errorCoder
	structCoder   *structCoder
	serviceCoder  *serviceCoder
	daoCoder      *daoCoder
	funcCoder     *funcCoder
	stepCoder     *stepCoder
}

func (this_ *factory) GetApp() common.IAppCoder {

	return this_.appCoder
}

func (this_ *factory) GetConstant() common.IConstantCoder {

	return this_.constantCoder
}

func (this_ *factory) GetStruct() common.IStructCoder {

	return this_.structCoder
}

func (this_ *factory) GetService() common.IServiceCoder {

	return this_.serviceCoder
}

func (this_ *factory) GetDao() common.IDaoCoder {

	return this_.daoCoder
}

func (this_ *factory) GetError() common.IErrorCoder {

	return this_.errorCoder
}

func (this_ *factory) GetFunc() common.IFuncCoder {

	return this_.funcCoder
}

func (this_ *factory) GetStep() common.IStepCoder {

	return this_.stepCoder
}
