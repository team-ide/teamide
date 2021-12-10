package loginService

import (
	"server/base"
	"server/userService"
)

func LoginByAccount(account string, password string) (user *base.UserEntity, err error) {
	user, err = userService.UserGetByAccount(account)
	if err != nil {
		return
	}
	return
}
