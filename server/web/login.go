package web

import (
	"server/base"
	"server/userService"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func apiLogin(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	loginRequest := &LoginRequest{}
	base.RequestJSON(loginRequest, c)
	if loginRequest.Account == "" {
		err = base.NewValidateError("登录账号不能为空!")
		return
	}
	if loginRequest.Password == "" {
		err = base.NewValidateError("登录密码不能为空!")
		return
	}
	pwd := base.AesDecryptCBCByKey(loginRequest.Password, base.HTTP_AES_KEY)
	if pwd == "" {
		err = base.NewValidateError("用户名或密码错误!")
		return
	}
	var user *base.UserEntity
	user, err = userService.LoginByAccount(loginRequest.Account, pwd)
	if err != nil {
		return
	}
	if user == nil {
		err = base.NewValidateError("用户名或密码错误!")
		return
	}
	jwt := &base.JWTBean{
		UserId:   user.UserId,
		Name:     user.Name,
		ServerId: user.ServerId,
		Time:     base.GetNowTime(),
	}
	res = getJWTStr(jwt)
	return
}

func apiLogout(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}
