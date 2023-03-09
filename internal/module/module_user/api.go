package module_user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"teamide/internal/context"
	"teamide/pkg/base"
)

type Api struct {
	*context.ServerContext
	UserService         *UserService
	UserPasswordService *UserPasswordService
}

func NewApi(UserService *UserService) *Api {
	return &Api{
		ServerContext:       UserService.ServerContext,
		UserService:         UserService,
		UserPasswordService: NewUserPasswordService(UserService.ServerContext),
	}
}

var (
	// 用户 权限

	// Power 用户基本 权限
	Power               = base.AppendPower(&base.PowerAction{Action: "user", Text: "用户", ShouldLogin: true, StandAlone: true})
	getPower            = base.AppendPower(&base.PowerAction{Action: "get", Text: "登录用户信息", Parent: Power, ShouldLogin: true, StandAlone: true})
	updatePower         = base.AppendPower(&base.PowerAction{Action: "update", Text: "登录用户信息修改", Parent: Power, ShouldLogin: true, StandAlone: true})
	updatePasswordPower = base.AppendPower(&base.PowerAction{Action: "updatePassword", Text: "登录用户密码修改", Parent: Power, ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: getPower, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: updatePower, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Power: updatePasswordPower, Do: this_.updatePassword})

	return
}

type GetRequest struct {
}

type GetResponse struct {
	User *UserModel `json:"user"`
}

func (this_ *Api) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &GetRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &GetResponse{}

	response.User, err = this_.UserService.Get(requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Account string `json:"account,omitempty"`
	Email   string `json:"email,omitempty"`
}

type UpdateResponse struct {
	User *UserModel `json:"user"`
}

func (this_ *Api) update(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdateRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdateResponse{}

	_, err = this_.UserService.Update(&UserModel{
		UserId:  requestBean.JWT.UserId,
		Name:    request.Name,
		Avatar:  request.Avatar,
		Account: request.Account,
		Email:   request.Email,
	})
	if err != nil {
		return
	}

	response.User, err = this_.UserService.Get(requestBean.JWT.UserId)
	if err != nil {
		return
	}

	res = response
	return
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword,omitempty"`
	Password    string `json:"password,omitempty"`
}

type UpdatePasswordResponse struct {
}

func (this_ *Api) updatePassword(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &UpdatePasswordRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &UpdatePasswordResponse{}

	if request.OldPassword == "" {
		err = base.NewValidateError("原密码不能为空!")
		return
	}
	if request.Password == "" {
		err = base.NewValidateError("登录密码不能为空!")
		return
	}

	var oldPassword string
	oldPassword, err = util.AesDecryptCBCByKey(request.OldPassword, this_.HttpAesKey)
	if err != nil {
		return
	}
	if oldPassword == "" {
		err = base.NewValidateError("原密码错误!")
		return
	}

	var password string
	password, err = util.AesDecryptCBCByKey(request.Password, this_.HttpAesKey)
	if err != nil {
		return
	}
	if password == "" {
		err = base.NewValidateError("密码错误!")
		return
	}

	checked, err := this_.UserPasswordService.CheckPassword(requestBean.JWT.UserId, oldPassword)
	if err != nil {
		return
	}

	if !checked {
		err = errors.New(fmt.Sprintf("原密码错误!"))
		return
	}

	_, err = this_.UserPasswordService.UpdatePassword(requestBean.JWT.UserId, password)
	if err != nil {
		return
	}

	res = response
	return
}
