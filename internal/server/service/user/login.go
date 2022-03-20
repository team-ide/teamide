package userService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) LoginByAccount(account string, password string) (res *base.UserEntity, err error) {
	var user *base.UserEntity
	user, err = getByAccount(account)
	if err != nil {
		return
	}
	if user == nil {
		return
	}
	var check bool
	check, err = passwordCheck(user.UserId, password)
	if err != nil {
		return
	}
	if !check {
		return
	}
	res = user
	return
}
