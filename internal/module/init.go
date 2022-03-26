package module

import "teamide/internal/module/module_register"

func (this_ *Api) InitStandAlone() (err error) {
	err = this_.InitStandAloneUser()
	if err != nil {
		return
	}
	return
}

func (this_ *Api) InitStandAloneUser() (err error) {

	standAloneUserId := this_.getStandAloneUserId()

	user, err := this_.userService.Get(standAloneUserId)
	if err != nil {
		return
	}
	if user == nil {
		register := &module_register.RegisterModel{
			UserId:     standAloneUserId,
			Name:       SystemUserName,
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
