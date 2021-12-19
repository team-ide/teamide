package powerService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindManagePowerUserApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/page"}, Power: base.PowerManagePowerUserPage, Do: apiManagePowerUserPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/list"}, Power: base.PowerManagePowerUserPage, Do: apiManagePowerUserList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/one"}, Power: base.PowerManagePowerUserPage, Do: apiManagePowerUserOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/insert"}, Power: base.PowerManagePowerUserInsert, Do: apiManagePowerUserInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/update"}, Power: base.PowerManagePowerUserUpdate, Do: apiManagePowerUserUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/user/delete"}, Power: base.PowerManagePowerUserDelete, Do: apiManagePowerUserDelete})

}

var (
	sqlManagePowerUserPage = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_USER,
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
	})

	sqlManagePowerUserList = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_USER,
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
	})

	sqlManagePowerUserOne = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_USER,
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
	})

	sqlManagePowerUserInsert = newInsertSql(&sqlModel.Insert{
		Table: TABLE_POWER_USER,
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
	})

	sqlManagePowerUserUpdate = newUpdateSql(&sqlModel.Update{
		Table: TABLE_POWER_USER,
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
	})

	sqlManagePowerUserDelete = newDeleteSql(&sqlModel.Delete{
		Table: TABLE_POWER_USER,
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	})
)

func apiManagePowerUserPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserPage.GetSqlParam(user)
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

func apiManagePowerUserList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserList.GetSqlParam(user)
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

func apiManagePowerUserOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserOne.GetSqlParam(user)
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

func apiManagePowerUserInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserInsert.GetSqlParam(user)
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

func apiManagePowerUserUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserUpdate.GetSqlParam(user)
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

func apiManagePowerUserDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	user := make(map[string]interface{})
	err = c.BindJSON(&user)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerUserDelete.GetSqlParam(user)
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
