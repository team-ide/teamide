package web

import (
	"server/base"
	"server/service/userService"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func apiRegister(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	registerRequest := &RegisterRequest{}
	base.RequestJSON(registerRequest, c)

	pwd := base.AesDecryptCBCByKey(registerRequest.Password, base.HTTP_AES_KEY)
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

	err = userService.UserTotalInsert(userBean)
	if err != nil {
		return
	}
	return
}
