package module

import (
	"encoding/json"
	"errors"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"teamide/internal/module/module_power"
	"teamide/internal/module/module_register"
	"teamide/pkg/base"
)

func (this_ *Api) initServer() (err error) {
	err = this_.initServerPower()
	if err != nil {
		return
	}
	err = this_.initServerUser()
	if err != nil {
		return
	}
	return
}

// 初始化 服务 权限
func (this_ *Api) initServerPower() (err error) {

	list, err := this_.powerRoleService.QueryByRoleType(base.SuperRoleType)
	if err != nil {
		return
	}
	if len(list) > 0 {
		return
	}
	this_.Logger.Info("not found super roles,now to create")
	powerRole := &module_power.PowerRoleModel{
		Name:     base.SuperRoleName,
		RoleType: base.SuperRoleType,
	}
	_, err = this_.powerRoleService.Insert(powerRole)
	if err != nil {
		this_.Logger.Error("super role create error", zap.Error(err))
		return
	}
	this_.Logger.Info("super role create success")
	return
}

// 初始化 服务 用户
func (this_ *Api) initServerUser() (err error) {
	list, err := this_.powerRoleService.QueryPowerUsersByRoleType(base.SuperRoleType)
	if err != nil {
		return
	}
	if len(list) > 0 {
		return
	}
	this_.Logger.Info("not found super users,now to create")

	var powerRoles []*module_power.PowerRoleModel
	powerRoles, err = this_.powerRoleService.QueryByRoleType(base.SuperRoleType)
	if err != nil {
		return
	}
	if len(powerRoles) == 0 {
		err = errors.New("super roles not found")
		return
	}
	users, err := this_.userService.QueryByAccount(base.SuperUserAccount)
	if err != nil {
		return
	}
	if len(users) == 0 {
		password := util.GetUUID()[0:10]
		register := &module_register.RegisterModel{
			Name:       base.SuperUserAccount,
			Account:    base.SuperUserAccount,
			Email:      "admin@teamide.com",
			Password:   password,
			SourceType: 1,
			Ip:         "127.0.0.1",
		}
		this_.Logger.Info("not found super user account,now to create")
		userInfoFile := this_.ServerConfig.Server.Data + "init-user-info.json"
		var infoFile *os.File
		infoFile, err = os.Create(userInfoFile)
		if err != nil {
			return
		}
		defer func() { _ = infoFile.Close() }()
		_, err = this_.registerService.Register(register)
		if err != nil {
			return
		}
		bs, _ := json.MarshalIndent(register, "", "  ")
		_, err = infoFile.WriteString(string(bs))
		if err != nil {
			return
		}
		this_.Logger.Info("super user account create success,user password saved to:" + userInfoFile)
		users, err = this_.userService.QueryByAccount(base.SuperUserAccount)
		if len(users) == 0 {
			err = errors.New("super user account not found")
		}
	}
	powerUser := &module_power.PowerUserModel{
		PowerRoleId: powerRoles[0].PowerRoleId,
		UserId:      users[0].UserId,
	}
	_, err = this_.powerUserService.Insert(powerUser)
	if err != nil {
		return
	}
	return
}
