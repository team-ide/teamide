package service

import (
	"teamide/internal/model"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

// NewUserPasswordService 根据库配置创建UserPasswordService
func NewUserPasswordService(dbWorker db.DatabaseWorker) (res *UserPasswordService) {

	idService := NewIDService(dbWorker)

	res = &UserPasswordService{
		dbWorker:  dbWorker,
		idService: idService,
	}
	return
}

// UserPasswordService 用户密码服务
type UserPasswordService struct {
	dbWorker  db.DatabaseWorker
	idService *IDService
}

// CheckPassword 检测密码是否一致
func (this_ *UserPasswordService) CheckPassword(userId int64, password string) (res bool, err error) {
	sql := `SELECT salt,password FROM ` + model.TableUserPassword + ` WHERE userId=? `
	list, err := this_.dbWorker.Query(sql, []interface{}{userId}, util.GetStructFieldTypes(model.UserPasswordModel{}))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	pwd := util.EncodePassword(list[0]["salt"].(string), password)

	res = pwd == list[0]["password"].(string)

	return
}

// Insert 新增
func (this_ *UserPasswordService) Insert(userId int64, password string) (rowsAffected int64, err error) {

	salt := util.GenerateUUID()[2:12]
	pwd := util.EncodePassword(salt, password)

	sql := `INSERT INTO ` + model.TableUserPassword + `(userId, salt, password, createTime) VALUES (?, ?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{userId, salt, pwd, time.Now()})
	if err != nil {
		return
	}

	return
}

// UpdatePassword 修改密码
func (this_ *UserPasswordService) UpdatePassword(userId int64, password string) (rowsAffected int64, err error) {

	salt := util.GenerateUUID()[2:12]
	pwd := util.EncodePassword(salt, password)

	sql := `UPDATE ` + model.TableUserPassword + ` SET salt=?,password=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{salt, pwd, time.Now(), userId})
	if err != nil {
		return
	}
	return
}
