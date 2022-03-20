package service

import (
	"errors"
	"fmt"
	"teamide/internal/model"
	"teamide/pkg/db"
	"time"
)

// NewLoginService 根据库配置创建LoginService
func NewLoginService(dbWorker db.DatabaseWorker) (res *LoginService) {

	idService := NewIDService(dbWorker)

	userService := NewUserService(dbWorker)

	userPasswordService := NewUserPasswordService(dbWorker)

	res = &LoginService{
		dbWorker:            dbWorker,
		idService:           idService,
		userService:         userService,
		userPasswordService: userPasswordService,
	}
	return
}

// LoginService 注册服务
type LoginService struct {
	dbWorker            db.DatabaseWorker
	idService           *IDService
	userService         *UserService
	userPasswordService *UserPasswordService
}

// Login 注册
func (this_ *LoginService) Login(login *model.LoginModel) (rowsAffected int64, err error) {

	accountLock := GetLock("user:login:" + login.Account)
	accountLock.Lock()
	defer accountLock.Unlock()

	user, err := this_.userService.GetByAccount(login.Account)
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
func (this_ *LoginService) insert(login *model.LoginModel) (rowsAffected int64, err error) {

	if login.LoginId == 0 {
		login.LoginId, err = this_.idService.GetNextID(model.IDTypeLogin)
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

	sql := `INSERT INTO ` + model.TableLogin + `(loginId, account, ip, sourceType, source, userId, loginTime, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{login.LoginId, login.Account, login.Ip, login.SourceType, login.Source, login.UserId, login.LoginTime, login.CreateTime})
	if err != nil {
		return
	}

	return
}

// Logout 登出
func (this_ *LoginService) Logout(loginId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + model.TableLogin + ` SET logoutTime=?,updateTime=? WHERE loginId=? `
	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{time.Now(), time.Now(), loginId})
	if err != nil {
		return
	}

	return
}
