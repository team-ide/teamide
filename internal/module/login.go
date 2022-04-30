package module

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os/user"
	"strings"
	"teamide/internal/base"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_user"
	"teamide/pkg/util"
)

var (
	SystemUserUid      string // 用户的 ID
	SystemUserGid      string // 用户所属组的 ID，如果属于多个组，那么此 ID 为主组的 ID
	SystemUserUsername string // 用户名
	SystemUserName     string // 属组名称，如果属于多个组，那么此名称为主组的名称
	SystemUserHomeDir  string // 用户的宿主目录
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	SystemUserUsername = u.Username
	SystemUserName = u.Name
	SystemUserHomeDir = u.HomeDir
	SystemUserGid = u.Gid
	SystemUserUid = u.Uid
}

type LoginRequest struct {
	Account  string `json:"account,omitempty"`
	Password string `json:"password,omitempty"`
}

func (this_ *Api) apiLogin(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
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
	pwd, err := util.AesDecryptCBCByKey(loginRequest.Password, this_.HttpAesKey)
	if err != nil {
		return
	}
	if pwd == "" {
		err = base.NewValidateError("用户名或密码错误!")
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
		err = base.NewValidateError("用户名或密码错误!")
		return
	}

	res = this_.getJWTStr(user)
	return
}

func (this_ *Api) apiLogout(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

const (
	JwtKey = "JWT"
)

func (this_ *Api) getJWT(c *gin.Context) *base.JWTBean {
	jwt := c.GetHeader(JwtKey)
	if jwt == "" {
		if strings.EqualFold(c.Request.Method, "get") {
			jwt = c.Query("jwt")
		}
	}
	if jwt == "" {
		return nil
	}
	jwt = this_.Decryption.Decrypt(jwt)
	if jwt == "" {
		return nil
	}
	res := &base.JWTBean{}
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
	jwt := &base.JWTBean{
		Sign:   util.GenerateUUID(),
		UserId: user.UserId,
		Name:   user.Name,
		Time:   util.GetNowTime(),
	}
	jwtStr := util.ToJSON(jwt)
	jwtStr = this_.Decryption.Encrypt(jwtStr)
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

func (this_ *Api) getStandAloneUserId() (userId int64) {
	userId = 1
	return
}

func (this_ *Api) apiSession(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	response := &SessionResponse{}

	var userId int64
	if request.JWT != nil {
		userId = request.JWT.UserId
	} else {
		if !this_.IsServer {
			userId = this_.getStandAloneUserId()
		}
	}
	if userId > 0 {
		var user *module_user.UserModel
		user, err = this_.userService.Get(userId)
		if err != nil {
			return
		}
		response.User = user
	}
	response.Powers = this_.getPowersByUserId(userId)
	response.JWT = this_.getJWTStr(response.User)

	json := util.ToJSON(response)
	res, err = util.AesEncryptCBCByKey(json, this_.HttpAesKey)
	if err != nil {
		return
	}
	return
}
