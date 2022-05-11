package module_user

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/pkg/util"
	"time"
)

// NewUserPasswordService 根据库配置创建UserPasswordService
func NewUserPasswordService(ServerContext *context.ServerContext) (res *UserPasswordService) {

	idService := module_id.NewIDService(ServerContext)

	res = &UserPasswordService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// UserPasswordService 用户密码服务
type UserPasswordService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// CheckPassword 检测密码是否一致
func (this_ *UserPasswordService) CheckPassword(userId int64, password string) (res bool, err error) {
	var list []*UserPasswordModel

	sql := `SELECT salt,password FROM ` + TableUserPassword + ` WHERE userId=? `
	err = this_.DatabaseWorker.Query(sql, []interface{}{userId}, &list)
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	pwd := util.EncodePassword(list[0].Salt, password)

	res = pwd == list[0].Password

	return
}

// Insert 新增
func (this_ *UserPasswordService) Insert(userId int64, password string) (rowsAffected int64, err error) {

	salt := util.GenerateUUID()[2:12]
	pwd := util.EncodePassword(salt, password)

	sql := `INSERT INTO ` + TableUserPassword + `(userId, salt, password, createTime) VALUES (?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{userId, salt, pwd, time.Now()})
	if err != nil {
		return
	}

	return
}

// UpdatePassword 修改密码
func (this_ *UserPasswordService) UpdatePassword(userId int64, password string) (rowsAffected int64, err error) {

	salt := util.GenerateUUID()[2:12]
	pwd := util.EncodePassword(salt, password)

	sql := `UPDATE ` + TableUserPassword + ` SET salt=?,password=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{salt, pwd, time.Now(), userId})
	if err != nil {
		return
	}
	return
}
