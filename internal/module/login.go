package module

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_user"
	"teamide/pkg/util"
)

type LoginRequest struct {
	Account  string `json:"account,omitempty"`
	Password string `json:"password,omitempty"`
}

func (this_ *Api) apiLogin(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	loginRequest := &LoginRequest{}
	if !base.RequestJSON(loginRequest, c) {
		return
	}

	var loginUser *module_user.UserModel
	var loginId int64
	if this_.IsServer {

		if loginRequest.Account == "" {
			err = base.NewValidateError("登录账号不能为空!")
			return
		}
		if loginRequest.Password == "" {
			err = base.NewValidateError("登录密码不能为空!")
			return
		}
		var pwd string
		pwd, err = util.AesDecryptCBCByKey(loginRequest.Password, this_.HttpAesKey)
		if err != nil {
			return
		}
		if pwd == "" {
			err = base.NewValidateError("用户名或密码错误!")
			return
		}
		userAgentStr := c.Request.UserAgent()

		var source string
		if userAgentStr != "" {
			userAgent := user_agent.New(userAgentStr)
			source += userAgent.OS()
			s1, s2 := userAgent.Browser()
			source += ";" + s1 + "@" + s2
			s1, s2 = userAgent.Engine()
			source += ";" + s1 + "@" + s2
		}
		login := &module_login.LoginModel{
			Account:    loginRequest.Account,
			Password:   pwd,
			Ip:         c.ClientIP(),
			SourceType: module_login.SourceTypeWeb,
			Source:     source,
			UserAgent:  userAgentStr,
		}

		loginUser, err = this_.loginService.Login(login)
		if err != nil {
			return
		}
		loginId = login.LoginId
	} else {
		loginUser, err = this_.userService.Get(base.StandAloneUserId)
		if err != nil {
			return
		}
		if loginUser == nil {
			err = base.NewValidateError("单机版用户信息不存在!")
			return
		}
	}

	if loginUser == nil {
		err = base.NewValidateError("用户名或密码错误!")
		return
	}

	res, err = this_.getJWTStr(loginId, loginUser)
	if err != nil {
		return
	}

	listener := context.GetListener(request.ClientTabKey)
	if listener != nil && listener.UserId != loginUser.UserId {
		context.ChangeListenerUserId(listener, loginUser.UserId)
	}
	return
}

func (this_ *Api) apiLogout(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	listener := context.GetListener(request.ClientTabKey)
	if listener != nil && listener.UserId != 0 {
		context.ChangeListenerUserId(listener, 0)
	}

	if request.JWT == nil {
		return
	}
	if request.JWT.LoginId != 0 {
		_, _ = this_.loginService.Logout(request.JWT.LoginId)
	}

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
	jwt, err := util.AesDecryptCBCByKey(jwt, this_.JWTAesKey)
	if err != nil {
		this_.Logger.Error("jwt decrypt error", zap.Error(err))
		return nil
	}
	if jwt == "" {
		return nil
	}
	res := &base.JWTBean{}
	_ = json.Unmarshal([]byte(jwt), res)
	if res.Time == 0 {
		return nil
	}
	if this_.IsServer {
		// 超过两小时
		if res.Time < (util.GetNowTime() - 1000*60*60*2) {

			return nil
		}
	}
	return res
}

func (this_ *Api) getJWTStr(loginId int64, user *module_user.UserModel) (jwtStr string, err error) {
	if user == nil {
		err = errors.New("用户信息不存在")
		return
	}
	jwt := &base.JWTBean{
		Sign:    util.UUID(),
		UserId:  user.UserId,
		Name:    user.Name,
		Account: user.Account,
		Time:    util.GetNowTime(),
		LoginId: loginId,
	}
	jwtJSONBytes, err := json.Marshal(jwt)
	if err != nil {
		return
	}
	jwtJSON := string(jwtJSONBytes)
	jwtStr, err = util.AesEncryptCBCByKey(jwtJSON, this_.JWTAesKey)
	if err != nil {
		return
	}
	if jwtStr == "" {
		err = errors.New("jwt加密失败")
		return
	}
	return
}

type SessionResponse struct {
	User   *module_user.UserModel `json:"user,omitempty"`
	Powers []string               `json:"powers,omitempty"`
	JWT    string                 `json:"JWT,omitempty"`
}

func (this_ *Api) apiSession(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	response := &SessionResponse{}

	var userId int64 = 0
	var loginId int64 = 0

	if this_.IsServer {
		if request.JWT != nil {
			userId = request.JWT.UserId
			loginId = request.JWT.LoginId
		}
	} else {
		userId = base.StandAloneUserId
	}
	if userId > 0 {
		var find *module_user.UserModel
		find, err = this_.userService.Get(userId)
		if err != nil {
			return
		}
		response.User = find
	}
	response.Powers = this_.getPowersByUserId(userId)

	if response.User != nil {
		response.JWT, err = this_.getJWTStr(loginId, response.User)
		if err != nil {
			return
		}
	}

	jsonString := util.ToJSON(response)
	res, err = util.AesEncryptCBCByKey(jsonString, this_.HttpAesKey)
	if err != nil {
		return
	}
	return
}
