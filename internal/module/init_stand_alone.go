package module

import (
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_user"
	"teamide/pkg/base"
)

func (this_ *Api) initStandAlone() (err error) {
	err = this_.initStandAloneUser()
	if err != nil {
		return
	}
	return
}

// 初始化 单机用户
func (this_ *Api) initStandAloneUser() (err error) {

	// 单机版用户 ID

	var standAloneUserId = this_.Setting.StandAloneUserId

	var user *module_user.UserModel
	if standAloneUserId == 0 {
		user, err = this_.userService.Get(1)
		if err != nil {
			return
		}
		if user == nil {
			standAloneUserId, err = this_.idService.GetNextID(module_id.IDTypeUser)
			if err != nil {
				return
			}
		} else {
			standAloneUserId = user.UserId
		}
		err = this_.settingService.Save(map[string]interface{}{
			"standAloneUserId": standAloneUserId,
		})
		if err != nil {
			return
		}
	} else {
		user, err = this_.userService.Get(standAloneUserId)
		if err != nil {
			return
		}
	}

	if user == nil {
		register := &module_register.RegisterModel{
			UserId:     standAloneUserId,
			Name:       base.SystemUserName,
			Account:    "admin",
			Email:      "admin@teamide.com",
			Password:   "admin123",
			SourceType: 1,
			Ip:         "127.0.0.1",
		}
		_, err = this_.registerService.Register(register)
		if err != nil {
			return
		}
	}
	return
}
