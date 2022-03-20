package module

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os/user"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_user"
	base2 "teamide/internal/server/base"
	component2 "teamide/internal/server/component"
	"teamide/pkg/util"
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

func (this_ *Api) apiLogin(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
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
	pwd, err := util.AesDecryptCBCByKey(loginRequest.Password, component2.HttpAesKey)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base2.NewValidateError("用户名或密码错误!")
		return
	}

	login := &module_login.LoginModel{
		Account:  loginRequest.Account,
		Password: pwd,
	}

	var user *module_user.UserModel
	user, err = this_.loginService.Login(login)
	if err != nil {
		return
	}
	if user == nil {
		err = base2.NewValidateError("用户名或密码错误!")
		return
	}

	res = this_.getJWTStr(user)
	return
}

func (this_ *Api) apiLogout(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

const (
	JWT_key = "JWT"
)

func (this_ *Api) getJWT(c *gin.Context) *base2.JWTBean {
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
	if res.Time < (util.GetNowTime() - 1000*60*60*2) {
		return nil
	}
	return res
}

func (this_ *Api) getJWTStr(user *module_user.UserModel) string {
	if user == nil {
		return ""
	}
	jwt := &base2.JWTBean{
		Sign:   util.GenerateUUID(),
		UserId: user.UserId,
		Name:   user.Name,
		Time:   util.GetNowTime(),
	}
	jwtStr := util.ToJSON(jwt)
	jwtStr = component2.RSAEncrypt(jwtStr)
	if jwtStr == "" {
		return ""
	}
	return jwtStr
}

type SessionResponse struct {
	User   *module_user.UserModel `json:"user,omitempty"`
	Powers []string               `json:"powers,omitempty"`
	JWT    string                 `json:"JWT,omitempty"`
}

func (this_ *Api) apiSession(request *base2.RequestBean, c *gin.Context) (res interface{}, err error) {
	response := &SessionResponse{}

	if base2.IsStandAlone {
		response.User = &module_user.UserModel{
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
			var user *module_user.UserModel
			user, err = this_.userService.Get(request.JWT.UserId)
			if err != nil {
				return
			}
			response.User = user
		}
		response.Powers = this_.getPowersByUserId(userId)
	}
	response.JWT = this_.getJWTStr(response.User)

	json := util.ToJSON(response)
	res, err = util.AesEncryptCBCByKey(json, component2.HttpAesKey)
	if err != nil {
		return
	}
	return
}
