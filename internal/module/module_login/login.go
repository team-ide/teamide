package module_login

import (
	"errors"
	"fmt"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_lock"
	"teamide/internal/module/module_user"
	"teamide/pkg/db"
	"time"
)

// NewLoginService 根据库配置创建LoginService
func NewLoginService(dbWorker db.DatabaseWorker) (res *LoginService) {

	idService := module_id.NewIDService(dbWorker)

	userService := module_user.NewUserService(dbWorker)

	userPasswordService := module_user.NewUserPasswordService(dbWorker)

	userAuthService := module_user.NewUserAuthService(dbWorker)

	res = &LoginService{
		dbWorker:            dbWorker,
		idService:           idService,
		userService:         userService,
		userPasswordService: userPasswordService,
		userAuthService:     userAuthService,
	}
	return
}

// LoginService 注册服务
type LoginService struct {
	dbWorker            db.DatabaseWorker
	idService           *module_id.IDService
	userService         *module_user.UserService
	userPasswordService *module_user.UserPasswordService
	userAuthService     *module_user.UserAuthService
}

// Login 注册
func (this_ *LoginService) Login(login *LoginModel) (user *module_user.UserModel, err error) {

	accountLock := module_lock.GetLock("user:login:" + login.Account)
	accountLock.Lock()
	defer accountLock.Unlock()

	user, err = this_.userService.GetByAccount(login.Account)
	if err != nil {
		return
	}

	if user == nil {
		err = errors.New(fmt.Sprintf("用户名或密码错误!"))
		return
	}

	checked, err := this_.userPasswordService.CheckPassword(user.UserId, login.Password)
	if err != nil {
		return
	}

	if !checked {
		err = errors.New(fmt.Sprintf("用户名或密码错误!"))
		return
	}
	login.UserId = user.UserId

	_, err = this_.insert(login)
	if err != nil {
		return
	}

	return
}

// insert 新增
func (this_ *LoginService) insert(login *LoginModel) (rowsAffected int64, err error) {

	if login.LoginId == 0 {
		login.LoginId, err = this_.idService.GetNextID(module_id.IDTypeLogin)
		if err != nil {
			return
		}
	}
	if login.LoginTime.IsZero() {
		login.LoginTime = time.Now()
	}
	if login.CreateTime.IsZero() {
		login.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableLogin + `(loginId, account, ip, sourceType, source, userId, loginTime, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{login.LoginId, login.Account, login.Ip, login.SourceType, login.Source, login.UserId, login.LoginTime, login.CreateTime})
	if err != nil {
		return
	}

	return
}

// Logout 登出
func (this_ *LoginService) Logout(loginId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableLogin + ` SET logoutTime=?,updateTime=? WHERE loginId=? `
	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{time.Now(), time.Now(), loginId})
	if err != nil {
		return
	}

	return
}
