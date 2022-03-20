package service

import (
	"errors"
	"fmt"
	"teamide/internal/model"
	"teamide/pkg/db"
	"time"
)

// NewRegisterService 根据库配置创建RegisterService
func NewRegisterService(dbWorker db.DatabaseWorker) (res *RegisterService) {

	idService := NewIDService(dbWorker)

	userService := NewUserService(dbWorker)

	userPasswordService := NewUserPasswordService(dbWorker)

	res = &RegisterService{
		dbWorker:            dbWorker,
		idService:           idService,
		userService:         userService,
		userPasswordService: userPasswordService,
	}
	return
}

// RegisterService 注册服务
type RegisterService struct {
	dbWorker            db.DatabaseWorker
	idService           *IDService
	userService         *UserService
	userPasswordService *UserPasswordService
}

// Register 注册
func (this_ *RegisterService) Register(register *model.RegisterModel) (rowsAffected int64, err error) {

	checkExist := func() error {

		exist, err := this_.userService.CheckExist(register.Account, register.Email)
		if err != nil {
			return err
		}
		if exist {
			err = errors.New(fmt.Sprintf("用户账号[%s],[%s]已存在!", register.Account, register.Email))
			return err
		}
		return nil
	}

	err = checkExist()
	if err != nil {
		return
	}

	accountLock := GetLock("user:account:" + register.Account)
	accountLock.Lock()
	defer accountLock.Unlock()

	emailLock := GetLock("user:email:" + register.Email)
	emailLock.Lock()
	defer emailLock.Unlock()

	err = checkExist()
	if err != nil {
		return
	}

	_, err = this_.insert(register)
	if err != nil {
		return
	}
	user := &model.UserModel{
		Name:    register.Name,
		Account: register.Account,
		Email:   register.Email,
	}
	_, err = this_.userService.Insert(user)
	if err != nil {
		return
	}

	_, err = this_.userPasswordService.Insert(user.UserId, register.Password)
	if err != nil {
		return
	}

	register.UserId = user.UserId

	_, err = this_.bindUser(register.RegisterId, register.UserId)
	if err != nil {
		return
	}

	return
}

// insert 新增
func (this_ *RegisterService) insert(register *model.RegisterModel) (rowsAffected int64, err error) {

	if register.RegisterId == 0 {
		register.RegisterId, err = this_.idService.GetNextID(model.IDTypeRegister)
		if err != nil {
			return
		}
	}
	if register.CreateTime.IsZero() {
		register.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + model.TableRegister + `(registerId, name, account, email, ip, sourceType, source, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{register.RegisterId, register.Name, register.Account, register.Email, register.Ip, register.SourceType, register.Source, register.CreateTime})
	if err != nil {
		return
	}

	return
}

// bindUser 绑定用户
func (this_ *RegisterService) bindUser(registerId int64, userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + model.TableRegister + ` SET userId=?,updateTime=? WHERE registerId=? `
	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{userId, time.Now(), registerId})
	if err != nil {
		return
	}

	return
}
