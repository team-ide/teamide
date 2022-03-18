package userService

import (
	"errors"
	"fmt"
	"teamide/server/base"
	"teamide/server/component"
	"teamide/server/factory"
)

func check(user *base.UserEntity) (err error) {
	if user.Name == "" {
		err = base.NewValidateError("用户名称不能为空!")
		return
	}
	if len(user.Name) > 50 {
		err = base.NewValidateError("用户:", user.Name, "，名称长度不能大于50个字节!")
		return
	}
	if user.Account == "" {
		err = base.NewValidateError("用户:", user.Name, "，账号不能为空!")
		return
	}
	if len(user.Account) < 4 {
		err = base.NewValidateError("用户:", user.Name, "，账号:", user.Account, "，长度不能小于4个字节!")
		return
	}
	if len(user.Account) > 20 {
		err = base.NewValidateError("用户:", user.Name, "，账号:", user.Account, "，长度不能大于20个字节!")
		return
	}
	if !base.MatchAccount(user.Account) {
		err = base.NewValidateError("用户:", user.Name, "，账号:", user.Account, "，格式不正确!")
		return
	}
	if base.MatchNumber(user.Account) {
		err = base.NewValidateError("用户:", user.Name, "，账号:", user.Account, "，不能为数字!")
		return
	}
	if user.Email != "" {
		if len(user.Email) > 50 {
			err = base.NewValidateError("用户:", user.Name, "，邮箱:", user.Email, "，长度不能大于50个字节!")
			return
		}
		if !base.MatchEmail(user.Email) {
			err = base.NewValidateError("用户:", user.Name, "，邮箱:", user.Email, "，格式不正确!")
			return
		}
	}

	return
}

// 用户全量信息新增
func (this_ *Service) TotalInsert(userTotal *base.UserTotalBean) (err error) {

	user := userTotal.User
	err = insert(user)
	if err != nil {
		return
	}
	password := userTotal.Password

	password.UserId = user.UserId

	err = setPassword(password)
	if err != nil {
		return
	}

	err = setMetadata(userTotal)
	if err != nil {
		return
	}

	return
}

func setMetadata(userTotal *base.UserTotalBean) (err error) {

	user := userTotal.User
	metadata := base.BeanToMap(userTotal)

	// lock := base.GetLock(fmt.Sprint("user:metadata:", user.UserId, ":lock"))
	// lock.Lock()

	// defer lock.Unlock()

	lockKey := component.GetUserMetadataLockRedisKey(user.UserId)
	var unlock func() (err error)
	unlock, err = component.Redis.Lock(lockKey, 10, 1000)
	if err != nil {
		return
	}
	defer unlock()

	err = setMetadataByMap(user.UserId, metadata)
	unlock()
	if err != nil {
		return
	}

	return
}

func setMetadataByMap(userId int64, metadata map[string]interface{}) (err error) {
	if userId == 0 {
		err = errors.New("userId is null")
		return
	}
	inserts := []base.UserMetadataEntity{}
	for _, mStruct := range base.U_M {
		data := metadata[mStruct.Name].(map[string]interface{})
		var inserts_ []base.UserMetadataEntity
		inserts_, err = getInsertMetadataEntity(userId, data, mStruct)
		if err != nil {
			return
		}
		inserts = append(inserts, inserts_...)
	}
	size := len(inserts)
	if size > 0 {
		var ids []int64
		ids, err = factory.IdService.GetIDs(component.ID_TYPE_USER_METADATA, int64(size))
		if err != nil {
			return
		}
		datas := []interface{}{}
		for index, one := range inserts {
			one.MetadataId = ids[index]
			datas = append(datas, one)
		}
		err = component.DB.BatchInsertBean(TABLE_USER_METADATA, datas)

		if err != nil {
			return
		}
	}

	return
}

func getInsertMetadataEntity(userId int64, metadata map[string]interface{}, mStruct *base.MStruct) (inserts []base.UserMetadataEntity, err error) {
	if len(metadata) == 0 {
		return
	}

	inserts = []base.UserMetadataEntity{}

	for _, mSField := range mStruct.Fields {

		metadataValue := metadata[mSField.Name]
		if metadataValue == nil {
			continue
		}
		if mSField.Struct == nil {
			value := base.GetStringValue(metadataValue)
			if value == "" {
				continue
			}
			insert := base.UserMetadataEntity{
				UserId:         userId,
				MetadataStruct: mStruct.StructCode,
				MetadataField:  mSField.StructFieldCode,
				MetadataValue:  value,
				CreateTime:     base.Now(),
			}
			inserts = append(inserts, insert)
		} else {

			if mSField.FieldType == base.F_T_LIST_STRUCT {
				list := metadata[mSField.Name].([]map[string]interface{})
				for index, one := range list {
					insert := base.UserMetadataEntity{
						UserId:         userId,
						MetadataStruct: mStruct.StructCode,
						MetadataField:  mSField.StructFieldCode,
						MetadataValue:  fmt.Sprint("list[", index, "]"),
						CreateTime:     base.Now(),
					}
					inserts = append(inserts, insert)

					var inserts_ []base.UserMetadataEntity
					inserts_, err = getInsertMetadataEntity(userId, one, mSField.Struct)
					if err != nil {
						return
					}
					inserts = append(inserts, inserts_...)
				}
			} else {
				one := metadata[mSField.Name].(map[string]interface{})
				insert := base.UserMetadataEntity{
					UserId:         userId,
					MetadataStruct: mStruct.StructCode,
					MetadataField:  mSField.StructFieldCode,
					MetadataValue:  "",
					CreateTime:     base.Now(),
				}
				inserts = append(inserts, insert)

				var inserts_ []base.UserMetadataEntity
				inserts_, err = getInsertMetadataEntity(userId, one, mSField.Struct)
				if err != nil {
					return
				}
				inserts = append(inserts, inserts_...)
			}
		}

	}
	return
}

// 用户全量信息批量新增
func (this_ *Service) TotalBatchInsert(userTotals []*base.UserTotalBean) (successUserTotals []*base.UserTotalBean, errUserTotals []*base.UserTotalBean, errs []error, err error) {
	for _, one := range userTotals {
		e := this_.TotalInsert(one)
		if e != nil {
			errUserTotals = append(errUserTotals, one)
			errs = append(errs, e)
		} else {
			successUserTotals = append(successUserTotals, one)
		}
	}

	return
}

// 用户信息新增
func insert(user *base.UserEntity) (err error) {

	err = check(user)
	if err != nil {
		return
	}
	// lock := base.GetLock("user:insert:lock")
	// lock.Lock()

	// defer lock.Unlock()
	var accountUnlock func() (err error)
	accountUnlock, err = component.Redis.Lock(component.GetUserInsertLockRedisKey(user.Account), 10, 1000)
	if err != nil {
		return
	}
	defer accountUnlock()

	if user.Email != "" {
		var emailUnlock func() (err error)
		emailUnlock, err = component.Redis.Lock(component.GetUserInsertLockRedisKey(user.Email), 10, 1000)
		if err != nil {
			return
		}
		defer emailUnlock()
	}

	var exist bool
	exist, err = existByAccount(user.Account, user.Email)
	if err != nil {
		return
	}
	if exist {
		err = base.NewValidateError("用户:", user.Name, "，账号或邮箱信息已存在!")
		return
	}
	var userId int64
	userId, err = factory.IdService.GetID(component.ID_TYPE_USER)
	if err != nil {
		return
	}
	user.UserId = userId
	user.EnabledState = 1
	user.ActivedState = 2
	user.LockedState = 2
	user.CreateTime = base.Now()

	err = component.DB.InsertBean(TABLE_USER, *user)

	if err != nil {
		return
	}

	return
}

// 用户设置密码
func setPassword(password *base.UserPasswordEntity) (err error) {
	password.CreateTime = base.Now()
	password.UpdateTime = base.Now()
	salt := base.GenerateUUID()[0:6]
	pwd := base.EncodePassword(salt, password.Password)
	password.Password = ""
	password.Salt = ""

	sql := "INSERT INTO " + TABLE_USER_PASSWORD + " (userId, salt, password, createTime) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE salt=?, password=?, updateTime=?"
	params := []interface{}{password.UserId, salt, pwd, password.CreateTime, salt, pwd, password.UpdateTime}

	sqlParam := base.NewSqlParam(sql, params)

	_, err = component.DB.Exec(sqlParam)

	if err != nil {
		return
	}

	return
}

// 查询用户密码
func passwordCheck(userId int64, password string) (check bool, err error) {
	sql := "SELECT * FROM " + TABLE_USER_PASSWORD + " WHERE userId=? "
	params := []interface{}{userId}

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewUserPasswordEntityInterface)

	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}
	userPassword := res[0].(*base.UserPasswordEntity)

	pwd := base.EncodePassword(userPassword.Salt, password)
	if pwd != userPassword.Password {
		return
	}
	check = true
	return
}

func (this_ *Service) Query(user base.UserEntity) (users []*base.UserEntity, err error) {
	sql := "SELECT * FROM " + TABLE_USER + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := base.NewSqlParam(sql, params)

	UserAppendWhere(user, &sqlParam)

	var res []interface{}
	_, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	users = []*base.UserEntity{}
	for _, one := range res {
		user := one.(*base.UserEntity)
		users = append(users, user)
	}
	return
}

func (this_ *Service) Count(user base.UserEntity) (count int64, err error) {
	sql := "SELECT COUNT(*) FROM " + TABLE_USER + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := base.NewSqlParam(sql, params)

	UserAppendWhere(user, &sqlParam)

	count, err = component.DB.Count(sqlParam)
	if err != nil {
		return
	}
	return
}
func UserAppendWhere(user base.UserEntity, sqlParam *base.SqlParam) {

	if user.EnabledState != 0 {
		sqlParam.Sql += " AND enabledState=? "
		sqlParam.Params = append(sqlParam.Params, user.EnabledState)
	}
	if user.ActivedState != 0 {
		sqlParam.Sql += " AND activedState=? "
		sqlParam.Params = append(sqlParam.Params, user.ActivedState)
	}
	if user.LockedState != 0 {
		sqlParam.Sql += " AND lockedState=? "
		sqlParam.Params = append(sqlParam.Params, user.LockedState)
	}
	if user.Name != "" {
		sqlParam.Sql += " AND name LIKE ? "
		sqlParam.Params = append(sqlParam.Params, "%"+user.Name+"%")
	}
	if user.Account != "" {
		sqlParam.Sql += " AND account LIKE ? "
		sqlParam.Params = append(sqlParam.Params, "%"+user.Account+"%")
	}
	if user.Email != "" {
		sqlParam.Sql += " AND email LIKE ? "
		sqlParam.Params = append(sqlParam.Params, "%"+user.Email+"%")
	}
}

//用户搜索，只搜索有效用户
func (this_ *Service) Userh(name string) (users []*base.UserEntity, err error) {
	sql := "SELECT userId,name,avatar FROM " + TABLE_USER + " WHERE enabledState=1 AND activedState=1 AND lockedState=2 AND (name LIKE ? OR account LIKE ? OR email LIKE ?)"
	params := []interface{}{"" + name + "%", "" + name + "%", "" + name + "%"}

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	users = []*base.UserEntity{}
	for _, one := range res {
		user := one.(*base.UserEntity)
		users = append(users, user)
	}
	return
}

//查询单个用户
func (this_ *Service) Get(userId int64) (user *base.UserEntity, err error) {
	sql := "SELECT * FROM " + TABLE_USER + " WHERE userId=? "
	params := []interface{}{userId}

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		user = res[0].(*base.UserEntity)
	}
	return
}

// 根据登录名称 或 邮箱 或 手机 查询单个用户
func getByAccount(account string) (user *base.UserEntity, err error) {
	sql := "SELECT * FROM " + TABLE_USER + " WHERE enabledState=1 AND (account=? OR email=?)"
	params := []interface{}{account, account}

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		user = res[0].(*base.UserEntity)
	}
	return
}

// 根据 登录名称 邮箱 手机 查询UserId
func getUserIdByAccount(account string, email string) (userId int64, err error) {
	sql := "SELECT userId FROM " + TABLE_USER + " WHERE enabledState=1 AND (account=? "
	params := []interface{}{account}

	if email != "" {
		sql += "OR email=? "
		params = append(params, email)
	}
	sql += ")"

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		userId = res[0].(*base.UserEntity).UserId
	}
	return
}

// 根据 登录名称 邮箱 手机 统计
func existByAccount(account string, email string) (exist bool, err error) {
	sql := "SELECT COUNT(userId) FROM " + TABLE_USER + " WHERE enabledState=1 AND (account=? "
	params := []interface{}{account}

	if email != "" {
		sql += "OR email=? "
		params = append(params, email)
	}
	sql += ")"

	sqlParam := base.NewSqlParam(sql, params)

	var res int64
	res, err = component.DB.Count(sqlParam)

	if err != nil {
		return
	}
	exist = res > 0
	return
}
