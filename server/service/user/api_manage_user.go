package userService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindManageUserApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/page"}, Power: base.PowerManageUserPage, Do: apiManageUserPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/list"}, Power: base.PowerManageUserPage, Do: apiManageUserList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/one"}, Power: base.PowerManageUserPage, Do: apiManageUserOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/insert"}, Power: base.PowerManageUserInsert, Do: apiManageUserInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/update"}, Power: base.PowerManageUserUpdate, Do: apiManageUserUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/delete"}, Power: base.PowerManageUserDelete, Do: apiManageUserDelete})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/active"}, Power: base.PowerManageUserActive, Do: apiManageUserActive})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/lock"}, Power: base.PowerManageUserLock, Do: apiManageUserLock})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/unlock"}, Power: base.PowerManageUserUnlock, Do: apiManageUserUnlock})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/disable"}, Power: base.PowerManageUserDisable, Do: apiManageUserDisable})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/user/enable"}, Power: base.PowerManageUserEnable, Do: apiManageUserEnable})

}

var (
	sqlManageUserPage = &sqlModel.Select{
		Table: TABLE_USER,
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

	sqlManageUserList = &sqlModel.Select{
		Table: TABLE_USER,
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

	sqlManageUserOne = &sqlModel.Select{
		Table: TABLE_USER,
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

	sqlManageUserInsert = &sqlModel.Insert{
		Table: TABLE_USER,
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

	sqlManageUserUpdate = &sqlModel.Update{
		Table: TABLE_USER,
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

	sqlManageUserDelete = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "deleteState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserActive = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "activedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserLock = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserUnlock = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "lockedState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserDisable = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "2"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageUserEnable = &sqlModel.Update{
		Table: TABLE_USER,
		Columns: []*sqlModel.UpdateColumn{
			{Name: "enabledState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}
)

func apiManageUserPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserPage.GetSqlParam(data)
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

func apiManageUserList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserList.GetSqlParam(data)
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

func apiManageUserOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserOne.GetSqlParam(data)
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

func apiManageUserInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserInsert.GetSqlParam(data)
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

func apiManageUserUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserUpdate.GetSqlParam(data)
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

func apiManageUserDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserDelete.GetSqlParam(data)
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

func apiManageUserActive(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserActive.GetSqlParam(data)
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

func apiManageUserLock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserLock.GetSqlParam(data)
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

func apiManageUserUnlock(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserUnlock.GetSqlParam(data)
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

func apiManageUserDisable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserDisable.GetSqlParam(data)
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

func apiManageUserEnable(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageUserEnable.GetSqlParam(data)
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
