package coders

import (
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

func NewApplicationCoder(app *modelers.Application, coderFactory common.ICoderFactory) (applicationCoder *ApplicationCoder) {
	applicationCoder = &ApplicationCoder{
		app:          app,
		coderFactory: coderFactory,
		appCoder:     coderFactory.GetApp(),
	}
	return
}

type ApplicationCoder struct {
	app          *modelers.Application
	coderFactory common.ICoderFactory
	appCoder     common.IAppCoder
}
