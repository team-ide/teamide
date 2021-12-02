package service

import (
	"base"
	"db"
	"errors"
	"fmt"
	"redis"
)

func UserCheck(user *base.UserEntity) (err error) {
	if user.Name == "" {
		err = base.NewErrorValdate("用户名称不能为空！")
		return
	}
	if len(user.Name) > 50 {
		err = base.NewErrorValdate("用户：", user.Name, "，名称长度不能大于50个字节！")
		return
	}
	if user.Account == "" {
		err = base.NewErrorValdate("用户：", user.Name, "，账号不能为空！")
		return
	}
	if len(user.Account) < 4 {
		err = base.NewErrorValdate("用户：", user.Name, "，账号：", user.Account, "，长度不能小于4个字节！")
		return
	}
	if len(user.Account) > 20 {
		err = base.NewErrorValdate("用户：", user.Name, "，账号：", user.Account, "，长度不能大于20个字节！")
		return
	}
	if !base.MatchAccount(user.Account) {
		err = base.NewErrorValdate("用户：", user.Name, "，账号：", user.Account, "，格式不正确！")
		return
	}
	if base.MatchNumber(user.Account) {
		err = base.NewErrorValdate("用户：", user.Name, "，账号：", user.Account, "，不能为数字！")
		return
	}
	if user.Email != "" {
		if len(user.Email) > 50 {
			err = base.NewErrorValdate("用户：", user.Name, "，邮箱：", user.Email, "，长度不能大于50个字节！")
			return
		}
		if !base.MatchEmail(user.Email) {
			err = base.NewErrorValdate("用户：", user.Name, "，邮箱：", user.Email, "，格式不正确！")
			return
		}
	}
	if user.Phone != "" {
		if len(user.Phone) > 20 {
			err = base.NewErrorValdate("用户：", user.Name, "，手机：", user.Phone, "，长度不能大于20个字节！")
			return
		}
		if !base.MatchPhone(user.Phone) {
			err = base.NewErrorValdate("用户：", user.Name, "，手机：", user.Phone, "，格式不正确！")
			return
		}
	}

	return
}

// 用户全量信息新增
func UserTotalInsert(userTotal *base.UserTotalBean) (err error) {

	user := userTotal.User
	err = UserInsert(user)
	if err != nil {
		return
	}
	password := userTotal.Password

	password.UserId = user.UserId

	err = UserSetPassword(password)
	if err != nil {
		return
	}

	err = UserSetMetadata(userTotal)
	if err != nil {
		return
	}

	return
}

func UserSetMetadata(userTotal *base.UserTotalBean) (err error) {

	user := userTotal.User
	metadata := base.BeanToMap(userTotal)

	// lock := base.GetLock(fmt.Sprint("user:metadata:", user.UserId, ":lock"))
	// lock.Lock()

	// defer lock.Unlock()

	lockKey := base.GetUserMetadataLockRedisKey(user.UserId)
	var unlock func() (err error)
	unlock, err = redis.RedisService.Lock(lockKey, 10, 1000)
	if err != nil {
		return
	}
	defer unlock()

	err = UserSetMetadataByMap(user.UserId, metadata)
	unlock()
	if err != nil {
		return
	}

	return
}

func UserSetMetadataByMap(userId int64, metadata map[string]interface{}) (err error) {
	if userId == 0 {
		err = errors.New("userId is null.")
		return
	}
	inserts := []base.UserMetadataEntity{}
	for _, mStruct := range base.U_M {
		data := metadata[mStruct.Name].(map[string]interface{})
		var inserts_ []base.UserMetadataEntity
		inserts_, err = UserGetInsertMetadataEntity(userId, data, mStruct)
		if err != nil {
			return
		}
		for _, one := range inserts_ {
			inserts = append(inserts, one)
		}
	}
	size := len(inserts)
	if size > 0 {
		var ids []int64
		ids, err = GetIDs(base.ID_TYPE_USER_METADATA, int64(size))
		if err != nil {
			return
		}
		datas := []interface{}{}
		for index, one := range inserts {
			one.MetadataId = ids[index]
			datas = append(datas, one)
		}
		err = BatchInsert(db.TABLE_USER_METADATA, datas)

		if err != nil {
			return
		}
	}

	return
}

func UserGetInsertMetadataEntity(userId int64, metadata map[string]interface{}, mStruct *base.MStruct) (inserts []base.UserMetadataEntity, err error) {
	if metadata == nil || len(metadata) == 0 {
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
					inserts_, err = UserGetInsertMetadataEntity(userId, one, mSField.Struct)
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
					MetadataValue:  fmt.Sprint("map"),
					CreateTime:     base.Now(),
				}
				inserts = append(inserts, insert)

				var inserts_ []base.UserMetadataEntity
				inserts_, err = UserGetInsertMetadataEntity(userId, one, mSField.Struct)
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
func UserTotalBatchInsert(userTotals []*base.UserTotalBean) (successUserTotals []*base.UserTotalBean, errUserTotals []*base.UserTotalBean, errs []error, err error) {
	for _, one := range userTotals {
		e := UserTotalInsert(one)
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
func UserInsert(user *base.UserEntity) (err error) {

	err = UserCheck(user)
	if err != nil {
		return
	}
	// lock := base.GetLock("user:insert:lock")
	// lock.Lock()

	// defer lock.Unlock()
	var accountUnlock func() (err error)
	accountUnlock, err = redis.RedisService.Lock(base.GetUserInsertLockRedisKey(user.Account), 10, 1000)
	if err != nil {
		return
	}
	defer accountUnlock()

	if user.Email != "" {
		var emailUnlock func() (err error)
		emailUnlock, err = redis.RedisService.Lock(base.GetUserInsertLockRedisKey(user.Email), 10, 1000)
		if err != nil {
			return
		}
		defer emailUnlock()
	}

	if user.Phone != "" {
		var phoneUnlock func() (err error)
		phoneUnlock, err = redis.RedisService.Lock(base.GetUserInsertLockRedisKey(user.Phone), 10, 1000)
		if err != nil {
			return
		}
		defer phoneUnlock()
	}

	var exist bool
	exist, err = UserExistByAccount(user.Account, user.Email, user.Phone)
	if err != nil {
		return
	}
	if exist {
		err = base.NewErrorValdate("用户：", user.Name, "，账号信息已存在！")
		return
	}
	var userId int64
	userId, err = GetID(base.ID_TYPE_USER)
	if err != nil {
		return
	}
	user.UserId = userId
	user.EnabledState = 1
	user.ActivedState = 2
	user.LockedState = 2
	user.CreateTime = base.Now()

	err = Insert(db.TABLE_USER, *user)

	if err != nil {
		return
	}

	return
}

// 用户设置密码
func UserSetPassword(password *base.UserPasswordEntity) (err error) {
	password.CreateTime = base.Now()
	password.UpdateTime = base.Now()
	salt := base.GenerateUUID()[0:6]
	pwd := base.EncodePassword(salt, password.Password)
	password.Password = ""
	password.Salt = ""

	sql := "INSERT INTO " + db.TABLE_USER_PASSWORD + " (userId, salt, password, createTime) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE salt=?, password=?, updateTime=?"
	params := []interface{}{password.UserId, salt, pwd, password.CreateTime, salt, pwd, password.UpdateTime}

	sqlParam := db.NewSqlParam(sql, params)

	_, err = db.DBService.Exec(sqlParam)

	if err != nil {
		return
	}

	return
}
func UserQuery(user base.UserEntity) (users []*base.UserEntity, err error) {
	sql := "SELECT * FROM " + db.TABLE_USER + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := db.NewSqlParam(sql, params)

	UserAppendWhere(user, &sqlParam)

	var res []interface{}
	_, err = db.DBService.Query(sqlParam, base.NewUserEntityInterface)

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

func UserCount(user base.UserEntity) (count int64, err error) {
	sql := "SELECT COUNT(*) FROM " + db.TABLE_USER + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := db.NewSqlParam(sql, params)

	UserAppendWhere(user, &sqlParam)

	count, err = db.DBService.Count(sqlParam)
	if err != nil {
		return
	}
	return
}
func UserAppendWhere(user base.UserEntity, sqlParam *db.SqlParam) {

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
func UserSearch(name string) (users []*base.UserEntity, err error) {
	sql := "SELECT userId,name,avatar FROM " + db.TABLE_USER + " WHERE enabledState=1 AND activedState=1 AND lockedState=2 AND (name LIKE ? OR account LIKE ? OR email LIKE ?)"
	params := []interface{}{"" + name + "%", "" + name + "%", "" + name + "%"}

	sqlParam := db.NewSqlParam(sql, params)

	var res []interface{}
	res, err = db.DBService.Query(sqlParam, base.NewUserEntityInterface)
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
func UserGet(userId int64) (user *base.UserEntity, err error) {
	sql := "SELECT * FROM " + db.TABLE_USER + " WHERE userId=? "
	params := []interface{}{userId}

	sqlParam := db.NewSqlParam(sql, params)

	var res []interface{}
	res, err = db.DBService.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		user = res[0].(*base.UserEntity)
	}
	return
}

// 根据登录名称 或 邮箱 或 手机 查询单个用户
func UserGetByAccount(account string) (user *base.UserEntity, err error) {
	sql := "SELECT * FROM " + db.TABLE_USER + " WHERE enabledState=1 AND (account=? OR email=? OR phone=?)"
	params := []interface{}{account, account, account}

	sqlParam := db.NewSqlParam(sql, params)

	var res []interface{}
	res, err = db.DBService.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		user = res[0].(*base.UserEntity)
	}
	return
}

// 根据 登录名称 邮箱 手机 查询UserId
func UserGetUserIdByAccount(account string, email string, phone string) (userId int64, err error) {
	sql := "SELECT userId FROM " + db.TABLE_USER + " WHERE enabledState=1 AND (account=? "
	params := []interface{}{account}

	if email != "" {
		sql += "OR email=? "
		params = append(params, email)
	}
	if phone != "" {
		sql += "OR phone=? "
		params = append(params, phone)
	}
	sql += ")"

	sqlParam := db.NewSqlParam(sql, params)

	var res []interface{}
	res, err = db.DBService.Query(sqlParam, base.NewUserEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		userId = res[0].(*base.UserEntity).UserId
	}
	return
}

// 根据 登录名称 邮箱 手机 统计
func UserExistByAccount(account string, email string, phone string) (exist bool, err error) {
	sql := "SELECT COUNT(userId) FROM " + db.TABLE_USER + " WHERE enabledState=1 AND (account=? "
	params := []interface{}{account}

	if email != "" {
		sql += "OR email=? "
		params = append(params, email)
	}
	if phone != "" {
		sql += "OR phone=? "
		params = append(params, phone)
	}
	sql += ")"

	sqlParam := db.NewSqlParam(sql, params)

	var res int64
	res, err = db.DBService.Count(sqlParam)

	if err != nil {
		return
	}
	exist = res > 0
	return
}
