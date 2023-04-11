package module

import (
	"go.uber.org/zap"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_user"
	"teamide/pkg/base"
)

func (this_ *Api) initStandAlone() (err error) {
	err = this_.initStandAloneUser()
	if err != nil {
		return
	}
	return
}

// 初始化 单机用户
func (this_ *Api) initStandAloneUser() (err error) {

	// 单机版用户 ID

	var standAloneUserId = this_.Setting.StandAloneUserId

	var user *module_user.UserModel
	this_.Logger.Info("检测单机用户")

	account := "admin"
	if standAloneUserId == 0 {

		user, err = this_.userService.QueryByAccount(account)
		if err != nil {
			return
		}
		if user == nil {
			standAloneUserId, err = this_.idService.GetNextID(module_id.IDTypeUser)
			if err != nil {
				return
			}
			this_.Logger.Info("生成单机用户ID", zap.Any("standAloneUserId", standAloneUserId))
		} else {
			standAloneUserId = user.UserId
			this_.Logger.Info("已有用户，使用该用户作为单机用户", zap.Any("standAloneUserId", standAloneUserId))
		}
		this_.Logger.Info("保存单机用户ID", zap.Any("standAloneUserId", standAloneUserId))
		err = this_.settingService.Save(map[string]interface{}{
			"standAloneUserId": standAloneUserId,
		})
		if err != nil {
			return
		}
	} else {
		this_.Logger.Info("已有单机用户ID，查询用户", zap.Any("standAloneUserId", standAloneUserId))
		user, err = this_.userService.Get(standAloneUserId)
		if err != nil {
			return
		}
	}

	if user == nil {
		this_.Logger.Info("单机用户不存在，插入用户", zap.Any("standAloneUserId", standAloneUserId))
		register := &module_register.RegisterModel{
			UserId:     standAloneUserId,
			Name:       base.SystemUserName,
			Account:    account,
			Email:      "admin@teamide.com",
			Password:   "admin123",
			SourceType: 1,
			Ip:         "127.0.0.1",
		}
		_, err = this_.registerService.Register(register)
		if err != nil {
			return
		}
	}
	return
}
