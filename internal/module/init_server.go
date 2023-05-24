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
	err = this_.initPower(base.SuperRoleType, base.SuperRoleName)
	if err != nil {
		return
	}
	_, err = this_.initPowerUser(base.SuperRoleType, base.SuperRoleName, "admin", "admin", "admin@teamide.com", true)
	if err != nil {
		return
	}

	err = this_.initPower(base.AnonymousRoleType, base.AnonymousRoleName)
	if err != nil {
		return
	}

	return
}

// 初始化 服务 权限
func (this_ *Api) initPower(roleType int, roleName string) (err error) {

	list, err := this_.powerRoleService.QueryByRoleType(roleType)
	if err != nil {
		return
	}
	if len(list) > 0 {
		return
	}
	this_.Logger.Info("not found role [" + roleName + "],now to create")
	powerRole := &module_power.PowerRoleModel{
		Name:     roleName,
		RoleType: roleType,
	}
	_, err = this_.powerRoleService.Insert(powerRole)
	if err != nil {
		this_.Logger.Error("role ["+roleName+"] create error", zap.Error(err))
		return
	}
	this_.Logger.Info("role [" + roleName + "] create success")
	return
}

// 初始化 服务 用户
func (this_ *Api) initPowerUser(roleType int, roleName string, name string, account string, email string, saveToFile bool) (userId int64, err error) {
	list, err := this_.powerRoleService.QueryPowerUsersByRoleType(roleType)
	if err != nil {
		return
	}
	if len(list) > 0 {
		userId = list[0].UserId
		return
	}
	this_.Logger.Info("not found role [" + roleName + "] users,now to create")

	var powerRoles []*module_power.PowerRoleModel
	powerRoles, err = this_.powerRoleService.QueryByRoleType(roleType)
	if err != nil {
		return
	}
	if len(powerRoles) == 0 {
		err = errors.New("role [" + roleName + "] not found")
		return
	}
	user, err := this_.userService.QueryByAccount(account)
	if err != nil {
		return
	}
	if user == nil {
		password := util.GetUUID()[0:10]
		register := &module_register.RegisterModel{
			Name:       name,
			Account:    account,
			Email:      email,
			Password:   password,
			SourceType: 1,
			Ip:         "127.0.0.1",
		}
		this_.Logger.Info("not found rule [" + roleName + "] user account,now to create")
		_, err = this_.registerService.Register(register)
		if err != nil {
			return
		}

		user, err = this_.userService.QueryByAccount(account)
		if user == nil {
			err = errors.New("rule [" + roleName + "] user account not found")
			return
		}

		if saveToFile {
			bs, _ := json.MarshalIndent(register, "", "  ")

			userInfoFile := this_.ServerConfig.Server.Data + "init-user-info.json"
			var infoFile *os.File
			infoFile, err = os.Create(userInfoFile)
			if err != nil {
				return
			}
			defer func() { _ = infoFile.Close() }()

			_, err = infoFile.WriteString(string(bs))
			if err != nil {
				return
			}
			this_.Logger.Info("rule [" + roleName + "] user account create success,user password saved to:" + userInfoFile)
		}
	}
	userId = user.UserId
	powerUser := &module_power.PowerUserModel{
		PowerRoleId: powerRoles[0].PowerRoleId,
		UserId:      userId,
	}
	_, err = this_.powerUserService.Insert(powerUser)
	if err != nil {
		return
	}
	return
}
