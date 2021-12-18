package userService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindManageUserAuthApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/page"}, Power: base.PowerManageUserAuthPage, Do: apiManageUserAuthPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/list"}, Power: base.PowerManageUserAuthPage, Do: apiManageUserAuthList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/one"}, Power: base.PowerManageUserAuthPage, Do: apiManageUserAuthOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/insert"}, Power: base.PowerManageUserAuthInsert, Do: apiManageUserAuthInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/update"}, Power: base.PowerManageUserAuthUpdate, Do: apiManageUserAuthUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/delete"}, Power: base.PowerManageUserAuthDelete, Do: apiManageUserAuthDelete})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/active"}, Power: base.PowerManageUserAuthActive, Do: apiManageUserAuthActive})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/lock"}, Power: base.PowerManageUserAuthLock, Do: apiManageUserAuthLock})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/unlock"}, Power: base.PowerManageUserAuthUnlock, Do: apiManageUserAuthUnlock})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/disable"}, Power: base.PowerManageUserAuthDisable, Do: apiManageUserAuthDisable})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/auth/enable"}, Power: base.PowerManageUserAuthEnable, Do: apiManageUserAuthEnable})

}

var (
	sqlManageUserAuthPage = &sqlModel.Select{
		Table: "TM_USER",
		Columns: []*sqlModel.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "avatar"},
			{Name: "email"},
			{Name: "account"},
			{Name: "activedState"},
			{Name: "lockedState"},
			{Name: "enabledState"},
			{Name: "deletedState"},
			{Name: "createTime"},
			{Name: "updateTime"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "deletedState", ValueScript: "1"},
			{
				Piece: true,
				Wheres: []*sqlModel.Where{
					{Name: "name", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "email", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "account", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "activedState"},
					{Name: "lockedState"},
					{Name: "enabledState"},
				}},
		},
	}

	sqlManageUserAuthList = &sqlModel.Select{
		Table: "TM_USER",
		Columns: []*sqlModel.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "avatar"},
			{Name: "email"},
			{Name: "account"},
			{Name: "activedState"},
			{Name: "lockedState"},
			{Name: "enabledState"},
			{Name: "deletedState"},
			{Name: "createTime"},
			{Name: "updateTime"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "deletedState", ValueScript: "1"},
			{
				Piece: true,
				Wheres: []*sqlModel.Where{
					{Name: "name", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "email", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "account", Operator: sqlModel.LIKE_BEFORE.Value},
					{Name: "activedState"},
					{Name: "lockedState"},
					{Name: "enabledState"},
				}},
		},
	}

	sqlManageUserAuthOne = &sqlModel.Select{
		Table: "TM_USER",
		Columns: []*sqlModel.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "avatar"},
			{Name: "email"},
			{Name: "account"},
			{Name: "activedState"},
			{Name: "lockedState"},
			{Name: "enabledState"},
			{Name: "deletedState"},
			{Name: "createTime"},
			{Name: "updateTime"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthInsert = &sqlModel.Insert{
		Table: "TM_USER",
		Columns: []*sqlModel.InsertColumn{
			{Name: "userId", Required: true},
			{Name: "name", Required: true},
			{Name: "avatar"},
			{Name: "email", Required: true},
			{Name: "account", Required: true},
			{Name: "activedState"},
			{Name: "lockedState"},
			{Name: "enabledState"},
			{Name: "createTime", ValueScript: "now()", Required: true},
		},
	}

	sqlManageUserAuthUpdate = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "name"},
			{Name: "avatar"},
			{Name: "email"},
			{Name: "account"},
			{Name: "updateTime", ValueScript: "now()"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthDelete = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "deleteState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthActive = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "activedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthLock = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthUnlock = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthDisable = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserAuthEnable = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}
)

func apiManageUserAuthPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthPage.GetSqlParam(data)
	if err != nil {
		return
	}
	page, err := component.DB.QueryPage(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	res = page
	return
}

func apiManageUserAuthList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthList.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Query(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthOne.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.QueryOne(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthInsert.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthUpdate.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthDelete.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthActive(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthActive.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthLock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthLock.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthUnlock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthUnlock.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthDisable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthDisable.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}

func apiManageUserAuthEnable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserAuthEnable.GetSqlParam(data)
	if err != nil {
		return
	}
	one, err := component.DB.Exec(sqlParam)
	if err != nil {
		return
	}
	res = one
	return
}
