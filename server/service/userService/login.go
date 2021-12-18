package userService

import (
	"server/base"
)

func LoginByAccount(account string, password string) (res *base.UserEntity, err error) {
	var user *base.UserEntity
	user, err = UserGetByAccount(account)
	if err != nil {
		return
	}
	if user == nil {
		return
	}
	var check bool
	check, err = UserPasswordCheck(user.UserId, password)
	if err != nil {
		return
	}
	if !check {
		return
	}
	res = user
	return
}
