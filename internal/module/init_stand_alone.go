package module

import (
	"teamide/internal/base"
	"teamide/internal/module/module_register"
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

	user, err := this_.userService.Get(base.StandAloneUserId)
	if err != nil {
		return
	}
	if user == nil {
		register := &module_register.RegisterModel{
			UserId:     base.StandAloneUserId,
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
