package userService

import (
	"server/base"
	"server/component"
	"server/model/sqlModel"

	"github.com/gin-gonic/gin"
)

func bindUserAuthApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/page"}, Power: base.PowerUserAuthPage, Do: apiUserAuthPage})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/list"}, Power: base.PowerUserAuthPage, Do: apiUserAuthList})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/one"}, Power: base.PowerUserAuthPage, Do: apiUserAuthOne})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/insert"}, Power: base.PowerUserAuthInsert, Do: apiUserAuthInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/update"}, Power: base.PowerUserAuthUpdate, Do: apiUserAuthUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/delete"}, Power: base.PowerUserAuthDelete, Do: apiUserAuthDelete})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/active"}, Power: base.PowerUserAuthActive, Do: apiUserAuthActive})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/lock"}, Power: base.PowerUserAuthLock, Do: apiUserAuthLock})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/unlock"}, Power: base.PowerUserAuthUnlock, Do: apiUserAuthUnlock})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/disable"}, Power: base.PowerUserAuthDisable, Do: apiUserAuthDisable})
	appendApi(&base.ApiWorker{Apis: []string{"/user/auth/enable"}, Power: base.PowerUserAuthEnable, Do: apiUserAuthEnable})

}

var (
	sqlUserAuthPage = &sqlModel.Select{
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

	sqlUserAuthList = &sqlModel.Select{
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

	sqlUserAuthOne = &sqlModel.Select{
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

	sqlUserAuthInsert = &sqlModel.Insert{
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

	sqlUserAuthUpdate = &sqlModel.Update{
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

	sqlUserAuthDelete = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "deleteState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlUserAuthActive = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "activedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlUserAuthLock = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlUserAuthUnlock = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlUserAuthDisable = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlUserAuthEnable = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}
)

func apiUserAuthPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthPage.GetSqlParam(data)
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

func apiUserAuthList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthList.GetSqlParam(data)
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

func apiUserAuthOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthOne.GetSqlParam(data)
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

func apiUserAuthInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthInsert.GetSqlParam(data)
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

func apiUserAuthUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthUpdate.GetSqlParam(data)
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

func apiUserAuthDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthDelete.GetSqlParam(data)
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

func apiUserAuthActive(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthActive.GetSqlParam(data)
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

func apiUserAuthLock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthLock.GetSqlParam(data)
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

func apiUserAuthUnlock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthUnlock.GetSqlParam(data)
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

func apiUserAuthDisable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthDisable.GetSqlParam(data)
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

func apiUserAuthEnable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserAuthEnable.GetSqlParam(data)
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
