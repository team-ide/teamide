package javascript

import (
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

func NewFactory(app *modelers.Application) (factory common.ICoderFactory) {
	appCoder := &AppCoder{app: app}
	factory = &Factory{
		appCoder: appCoder,
	}
	return
}

type Factory struct {
	appCoder *AppCoder
}

func (this_ Factory) GetApp() common.IAppCoder {

	return this_.appCoder
}

func (this_ Factory) GetConstant() common.IConstantCoder {

	return nil
}

func (this_ Factory) GetStruct() common.IStructCoder {

	return nil
}

func (this_ Factory) GetService() common.IServiceCoder {

	return nil
}

func (this_ Factory) GetDao() common.IDaoCoder {

	return nil
}

func (this_ Factory) GetError() common.IErrorCoder {

	return nil
}

func (this_ Factory) GetFunc() common.IFuncCoder {

	return nil
}

func (this_ Factory) GetStep() common.IStepCoder {

	return nil
}
