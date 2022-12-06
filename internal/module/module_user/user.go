package module_user

import (
	"errors"
	"fmt"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_lock"
	"time"
)

// NewUserService 根据库配置创建UserService
func NewUserService(ServerContext *context.ServerContext) (res *UserService) {

	idService := module_id.NewIDService(ServerContext)

	res = &UserService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// UserService 用户服务
type UserService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Get 查询单个
func (this_ *UserService) Get(userId int64) (res *UserModel, err error) {
	var list []*UserModel
	sql := `SELECT * FROM ` + TableUser + ` WHERE userId=? `
	err = this_.DatabaseWorker.Query(sql, []interface{}{userId}, &list)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}
	return
}

// Query 查询
func (this_ *UserService) Query(user *UserModel) (res []*UserModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableUser + ` WHERE deleted=2 `
	if user.Activated != 0 {
		sql += " AND activated = ?"
		values = append(values, user.Activated)
	}
	if user.Locked != 0 {
		sql += " AND locked = ?"
		values = append(values, user.Locked)
	}
	if user.Enabled != 0 {
		sql += " AND enabled = ?"
		values = append(values, user.Enabled)
	}
	if user.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", user.Name, "%"))
	}
	if user.Account != "" {
		sql += " AND account like ?"
		values = append(values, fmt.Sprint("%", user.Account, "%"))
	}
	if user.Email != "" {
		sql += " AND email like ?"
		values = append(values, fmt.Sprint("%", user.Email, "%"))
	}

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}

	return
}

// QueryByAccount 查询
func (this_ *UserService) QueryByAccount(account string) (res []*UserModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableUser + ` WHERE deleted=2 `
	sql += " AND account = ?"
	values = append(values, account)
	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}

	return
}

// CheckExist 查询
func (this_ *UserService) CheckExist(userId int64, account string, email string) (res bool, err error) {

	var values []interface{}
	sql := `SELECT COUNT(1) FROM ` + TableUser + ` WHERE deleted=2 AND (1=2`

	if account != "" {
		sql += " OR account = ?"
		values = append(values, account)
	}
	if email != "" {
		sql += " OR email = ?"
		values = append(values, email)
	}
	sql += ")"
	if userId != 0 {
		sql += " AND userId != ?"
		values = append(values, userId)
	}

	count, err := this_.DatabaseWorker.Count(sql, values)
	if err != nil {
		return
	}

	res = count > 0

	return
}

// GetByAccount 根据账号查询
func (this_ *UserService) GetByAccount(account string) (res *UserModel, err error) {

	var list []*UserModel

	sql := `SELECT * FROM ` + TableUser + ` WHERE deleted=2 AND (account = ? OR email = ?)`
	err = this_.DatabaseWorker.Query(sql, []interface{}{account, account}, &list)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}
	return
}

// Insert 新增
func (this_ *UserService) Insert(user *UserModel) (rowsAffected int64, err error) {

	if user.UserId == 0 {
		user.UserId, err = this_.idService.GetNextID(module_id.IDTypeUser)
		if err != nil {
			return
		}
	}
	if user.Activated == 0 {
		user.Activated = 2
	}
	if user.CreateTime.IsZero() {
		user.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableUser + `(userId, name, avatar, account, email, activated, createTime) VALUES (?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{user.UserId, user.Name, user.Avatar, user.Account, user.Email, user.Activated, user.CreateTime})
	if err != nil {
		return
	}

	return
}

// Update 更新
func (this_ *UserService) Update(user *UserModel) (rowsAffected int64, err error) {
	if user.Account != "" || user.Email != "" {
		checkExist := func() error {

			exist, err := this_.CheckExist(user.UserId, user.Account, user.Email)
			if err != nil {
				return err
			}
			if exist {
				err = errors.New(fmt.Sprintf("用户账号[%s],[%s]已存在!", user.Account, user.Email))
				return err
			}
			return nil
		}

		err = checkExist()
		if err != nil {
			return
		}

		accountLock := module_lock.GetLock("user:account:" + user.Account)
		accountLock.Lock()
		defer accountLock.Unlock()

		emailLock := module_lock.GetLock("user:email:" + user.Email)
		emailLock.Lock()
		defer emailLock.Unlock()

		err = checkExist()
		if err != nil {
			return
		}
	}

	var values []interface{}

	sql := `UPDATE ` + TableUser + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if user.Name != "" {
		sql += "name=?,"
		values = append(values, user.Name)
	}
	sql += "avatar=?,"
	values = append(values, user.Avatar)
	if user.Account != "" {
		sql += "account=?,"
		values = append(values, user.Account)
	}
	if user.Email != "" {
		sql += "email=?,"
		values = append(values, user.Email)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE userId=? "
	values = append(values, user.UserId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		return
	}

	return
}

// Active 激活
func (this_ *UserService) Active(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET activated=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// UnActive 不激活
func (this_ *UserService) UnActive(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET activated=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{2, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// Lock 锁定
func (this_ *UserService) Lock(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET locked=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// Unlock 解锁
func (this_ *UserService) Unlock(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET locked=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{2, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// Enable 启用
func (this_ *UserService) Enable(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET enabled=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// Disable 禁用
func (this_ *UserService) Disable(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET enabled=?,updateTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{2, time.Now(), userId})
	if err != nil {
		return
	}

	return
}

// Delete 更新
func (this_ *UserService) Delete(userId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableUser + ` SET deleted=?,deleteTime=? WHERE userId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), userId})
	if err != nil {
		return
	}

	return
}
