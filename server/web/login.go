package web

import (
	"encoding/json"
	"server/base"
	"server/component"
	userService "server/service/user"

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
	pwd := component.AesDecryptCBCByKey(loginRequest.Password, component.HTTP_AES_KEY)
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

	res = getJWTStr(user)
	return
}

func apiLogout(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

const (
	JWT_key = "JWT"
)

func getJWT(c *gin.Context) *base.JWTBean {
	jwt := c.GetHeader(JWT_key)
	if jwt == "" {
		return nil
	}
	jwt = component.AesDecryptCBC(jwt)
	if jwt == "" {
		return nil
	}
	res := &base.JWTBean{}
	json.Unmarshal([]byte(jwt), res)
	if res.Time == 0 {
		return nil
	}
	// 超过两小时
	if res.Time < (base.GetNowTime() - 1000*60*60*2) {
		return nil
	}
	return res
}

func getJWTStr(user *base.UserEntity) string {
	if user == nil {
		return ""
	}
	jwt := &base.JWTBean{
		Sign:     base.GenerateUUID(),
		UserId:   user.UserId,
		Name:     user.Name,
		ServerId: user.ServerId,
		Time:     base.GetNowTime(),
	}
	jwtStr := base.ToJSON(jwt)
	jwtStr = component.AesEncryptCBC(jwtStr)
	if jwtStr == "" {
		return ""
	}
	return jwtStr
}

type SessionResponse struct {
	User   *base.UserEntity `json:"user"`
	Powers []string         `json:"powers"`
	JWT    string           `json:"JWT"`
}

func apiSession(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	response := &SessionResponse{}
	var userId int64 = 0
	if request.JWT != nil {
		userId = request.JWT.UserId
	}
	if userId > 0 {
		var user *base.UserEntity
		user, err = userService.UserGet(request.JWT.UserId)
		if err != nil {
			return
		}
		response.User = user
	}
	response.JWT = getJWTStr(response.User)
	response.Powers = getPowersByUserId(userId)

	json := base.ToJSON(response)
	res = component.AesEncryptCBCByKey(json, component.HTTP_AES_KEY)
	return
}
