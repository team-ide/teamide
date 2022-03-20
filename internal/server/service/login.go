package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os/user"
	base2 "teamide/internal/server/base"
	component2 "teamide/internal/server/component"
	"teamide/internal/server/factory"
)

var (
	SystemUser_Uid      string // 用户的 ID
	SystemUser_Gid      string // 用户所属组的 ID，如果属于多个组，那么此 ID 为主组的 ID
	SystemUser_Username string // 用户名
	SystemUser_Name     string // 属组名称，如果属于多个组，那么此名称为主组的名称
	SystemUser_HomeDir  string // 用户的宿主目录
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	SystemUser_Username = u.Username
	SystemUser_Name = u.Name
	SystemUser_HomeDir = u.HomeDir
	SystemUser_Gid = u.Gid
	SystemUser_Uid = u.Uid
}

type LoginRequest struct {
	Account  string `json:"account,omitempty"`
	Password string `json:"password,omitempty"`
}

func apiLogin(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
	loginRequest := &LoginRequest{}
	base2.RequestJSON(loginRequest, c)
	if loginRequest.Account == "" {
		err = base2.NewValidateError("登录账号不能为空!")
		return
	}
	if loginRequest.Password == "" {
		err = base2.NewValidateError("登录密码不能为空!")
		return
	}
	pwd, err := base2.AesDecryptCBCByKey(loginRequest.Password, component2.HTTP_AES_KEY)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base2.NewValidateError("用户名或密码错误!")
		return
	}
	var user *base2.UserEntity
	user, err = factory.UserService.LoginByAccount(loginRequest.Account, pwd)
	if err != nil {
		return
	}
	if user == nil {
		err = base2.NewValidateError("用户名或密码错误!")
		return
	}

	res = getJWTStr(user)
	return
}

func apiLogout(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

const (
	JWT_key = "JWT"
)

func getJWT(c *gin.Context) *base2.JWTBean {
	jwt := c.GetHeader(JWT_key)
	if jwt == "" {
		return nil
	}
	jwt = component2.RSADecrypt(jwt)
	if jwt == "" {
		return nil
	}
	res := &base2.JWTBean{}
	json.Unmarshal([]byte(jwt), res)
	if res.Time == 0 {
		return nil
	}
	// 超过两小时
	if res.Time < (base2.GetNowTime() - 1000*60*60*2) {
		return nil
	}
	return res
}

func getJWTStr(user *base2.UserEntity) string {
	if user == nil {
		return ""
	}
	jwt := &base2.JWTBean{
		Sign:   base2.GenerateUUID(),
		UserId: user.UserId,
		Name:   user.Name,
		Time:   base2.GetNowTime(),
	}
	jwtStr := base2.ToJSON(jwt)
	jwtStr = component2.RSAEncrypt(jwtStr)
	if jwtStr == "" {
		return ""
	}
	return jwtStr
}

type SessionResponse struct {
	User   *base2.UserEntity `json:"user,omitempty"`
	Powers []string          `json:"powers,omitempty"`
	JWT    string           `json:"JWT,omitempty"`
}

func apiSession(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
	response := &SessionResponse{}

	if base2.IS_STAND_ALONE {
		response.User = &base2.UserEntity{
			UserId: 1,
			Name:   SystemUser_Username,
		}
		response.Powers = []string{}
		ps := base2.GetPowers()
		for _, power := range ps {
			if power.AllowNative {
				response.Powers = append(response.Powers, power.Action)
			}
		}
	} else {
		var userId int64 = 0
		if request.JWT != nil {
			userId = request.JWT.UserId
		}
		if userId > 0 {
			var user *base2.UserEntity
			user, err = factory.UserService.Get(request.JWT.UserId)
			if err != nil {
				return
			}
			response.User = user
		}
		response.Powers = getPowersByUserId(userId)
	}
	response.JWT = getJWTStr(response.User)

	json := base2.ToJSON(response)
	res, err = base2.AesEncryptCBCByKey(json, component2.HTTP_AES_KEY)
	if err != nil {
		return
	}
	return
}
