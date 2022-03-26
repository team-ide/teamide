package module_register

import (
	"errors"
	"fmt"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_lock"
	"teamide/internal/module/module_user"
	"time"
)

// NewRegisterService 根据库配置创建RegisterService
func NewRegisterService(ServerContext *context.ServerContext) (res *RegisterService) {

	idService := module_id.NewIDService(ServerContext)

	userService := module_user.NewUserService(ServerContext)

	userPasswordService := module_user.NewUserPasswordService(ServerContext)

	res = &RegisterService{
		ServerContext:       ServerContext,
		idService:           idService,
		userService:         userService,
		userPasswordService: userPasswordService,
	}
	return
}

// RegisterService 注册服务
type RegisterService struct {
	*context.ServerContext
	idService           *module_id.IDService
	userService         *module_user.UserService
	userPasswordService *module_user.UserPasswordService
}

// Register 注册
func (this_ *RegisterService) Register(register *RegisterModel) (rowsAffected int64, err error) {

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

	accountLock := module_lock.GetLock("user:account:" + register.Account)
	accountLock.Lock()
	defer accountLock.Unlock()

	emailLock := module_lock.GetLock("user:email:" + register.Email)
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
	user := &module_user.UserModel{
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
func (this_ *RegisterService) insert(register *RegisterModel) (rowsAffected int64, err error) {

	if register.RegisterId == 0 {
		register.RegisterId, err = this_.idService.GetNextID(module_id.IDTypeRegister)
		if err != nil {
			return
		}
	}
	if register.CreateTime.IsZero() {
		register.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableRegister + `(registerId, name, account, email, ip, sourceType, source, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{register.RegisterId, register.Name, register.Account, register.Email, register.Ip, register.SourceType, register.Source, register.CreateTime})
	if err != nil {
		return
	}

	return
}

// bindUser 绑定用户
func (this_ *RegisterService) bindUser(registerId int64, userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableRegister + ` SET userId=?,updateTime=? WHERE registerId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{userId, time.Now(), registerId})
	if err != nil {
		return
	}

	return
}
