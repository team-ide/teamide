package service

import (
	"teamide/server/base"
	"teamide/server/component"
	"teamide/server/factory"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name,omitempty"`
	Account  string `json:"account,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func apiRegister(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	registerRequest := &RegisterRequest{}
	base.RequestJSON(registerRequest, c)

	pwd, err := base.AesDecryptCBCByKey(registerRequest.Password, component.HTTP_AES_KEY)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base.NewValidateError("密码不能为空!")
		return
	}
	user := &base.UserEntity{
		Name:    registerRequest.Name,
		Account: registerRequest.Account,
		Email:   registerRequest.Email,
	}
	password := &base.UserPasswordEntity{
		Password: pwd,
	}
	userBean := &base.UserTotalBean{
		User:       user,
		Password:   password,
		Persona:    nil,
		Enterprise: nil,
	}

	err = factory.UserService.TotalInsert(userBean)
	if err != nil {
		return
	}
	return
}
