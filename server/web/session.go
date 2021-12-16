package web

import (
	"server/base"
	"server/userService"

	"github.com/gin-gonic/gin"
)

const (
	JWT_key = "JWT"
)

func getJWT(c *gin.Context) *base.JWTBean {
	jwt := c.GetHeader(JWT_key)
	if jwt == "" {
		return nil
	}
	jwt = base.AesDecryptCBC(jwt)
	if jwt == "" {
		return nil
	}
	res := &base.JWTBean{}
	base.JSON.Unmarshal([]byte(jwt), res)
	return res
}

func getJWTStr(jwt *base.JWTBean) string {
	if jwt == nil {
		return ""
	}
	jwtStr := base.ToJSON(jwt)
	jwtStr = base.AesEncryptCBC(jwtStr)
	if jwtStr == "" {
		return ""
	}
	return jwtStr
}

type SessionResponse struct {
	User   *base.UserEntity `json:"user"`
	Powers []string         `json:"powers"`
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

	response.Powers = getPowersByUserId(userId)

	res = response
	return
}
