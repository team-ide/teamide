package module

import (
	"github.com/gin-gonic/gin"
	"teamide/internal/base"
	"teamide/internal/module/module_register"
	"teamide/pkg/util"
)

type RegisterRequest struct {
	Name     string `json:"name,omitempty"`
	Account  string `json:"account,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (this_ *Api) apiRegister(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	registerRequest := &RegisterRequest{}

	if !base.RequestJSON(registerRequest, c) {
		return
	}

	pwd, err := util.AesDecryptCBCByKey(registerRequest.Password, this_.HttpAesKey)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base.NewValidateError("密码不能为空!")
		return
	}
	register := &module_register.RegisterModel{
		Name:     registerRequest.Name,
		Account:  registerRequest.Account,
		Email:    registerRequest.Email,
		Password: pwd,
	}

	_, err = this_.registerService.Register(register)
	if err != nil {
		return
	}
	return
}
