package coder

import "teamide/pkg/maker/model"

type Coder interface {
}

type Code struct {
}

func GetServiceCode(coder Coder, app *model.Application, service *model.ServiceModel) (code *Code, err error) {

	return
}
