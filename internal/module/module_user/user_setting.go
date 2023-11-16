package module_user

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewUserSettingService 根据库配置创建 UserSettingService
func NewUserSettingService(ServerContext *context.ServerContext) (res *UserSettingService) {

	idService := module_id.NewIDService(ServerContext)

	res = &UserSettingService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// UserSettingService 用户授权服务
type UserSettingService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Save 保持
func (this_ *UserSettingService) Save(userId int64, setting map[string]string) (userSetting map[string]string, err error) {

	userSetting, err = this_.Query(userId)
	if err != nil {
		return
	}
	for name, value := range setting {
		if name == "" {
			continue
		}
		if value == "" {
			sql := `DELETE FROM ` + TableUserSetting + ` WHERE userId=? AND name=? `
			_, err = this_.DatabaseWorker.Exec(sql, []interface{}{userId, name})
			if err != nil {
				return
			}
			delete(userSetting, name)
			continue
		}
		_, find := userSetting[name]
		if find {
			sql := `UPDATE ` + TableUserSetting + ` SET value=?,updateTime=? WHERE userId=? AND name=? `
			_, err = this_.DatabaseWorker.Exec(sql, []interface{}{value, time.Now(), userId, name})
			if err != nil {
				return
			}
		} else {
			sql := `INSERT INTO ` + TableUserSetting + `(userId, name, value, createTime) VALUES (?, ?, ?, ?) `

			_, err = this_.DatabaseWorker.Exec(sql, []interface{}{userId, name, value, time.Now()})
			if err != nil {
				return
			}
		}
		userSetting[name] = value
	}

	return
}

// Query 新增
func (this_ *UserSettingService) Query(userId int64) (setting map[string]string, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableUserSetting + ` WHERE userId=? `
	values = append(values, userId)

	var list []*UserSettingModel
	err = this_.DatabaseWorker.Query(sql, values, &list)
	if err != nil {
		return
	}

	setting = map[string]string{}
	for _, one := range list {
		setting[one.Name] = one.Value
	}
	//fmt.Println("list:", list)
	//fmt.Println("setting:", setting)
	return
}
