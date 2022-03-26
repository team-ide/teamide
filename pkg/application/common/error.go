package common

import (
	"errors"
	"fmt"
	"teamide/pkg/application/base"
	"teamide/pkg/application/model"
)

func GetErrorModel(app IApplication, errorName string, errorCode string, errorMsg string) (res *model.ErrorModel, err error) {
	res = &model.ErrorModel{
		Name: errorName,
		Code: errorCode,
		Msg:  errorMsg,
	}
	if res.Name != "" {
		errorModel_ := app.GetContext().GetError(res.Name)
		if errorModel_ == nil {
			err = errors.New(fmt.Sprint("error [", res.Name, "] not defind"))
			return
		}
		if res.Code == "" {
			res.Code = errorModel_.Code
		}
		if res.Msg == "" {
			res.Msg = errorModel_.Msg
		}
	}
	return
}

func GoErrorByErrorModel(errorModel *model.ErrorModel) (err error) {
	err = base.NewError(errorModel.Code, errorModel.Msg)
	return
}
