package service

import (
	"github.com/gin-gonic/gin"
	base2 "teamide/internal/server/base"
	"teamide/internal/server/component"
	"teamide/internal/server/factory"
)

type RegisterRequest struct {
	Name     string `json:"name,omitempty"`
	Account  string `json:"account,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func apiRegister(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
	registerRequest := &RegisterRequest{}
	base2.RequestJSON(registerRequest, c)

	pwd, err := base2.AesDecryptCBCByKey(registerRequest.Password, component.HTTP_AES_KEY)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base2.NewValidateError("密码不能为空!")
		return
	}
	user := &base2.UserEntity{
		Name:    registerRequest.Name,
		Account: registerRequest.Account,
		Email:   registerRequest.Email,
	}
	password := &base2.UserPasswordEntity{
		Password: pwd,
	}
	userBean := &base2.UserTotalBean{
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
