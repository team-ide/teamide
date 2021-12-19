package powerService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindManagePowerRoleApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/page"}, Power: base.PowerManagePowerRolePage, Do: apiManagePowerRolePage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/list"}, Power: base.PowerManagePowerRolePage, Do: apiManagePowerRoleList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/one"}, Power: base.PowerManagePowerRolePage, Do: apiManagePowerRoleOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/insert"}, Power: base.PowerManagePowerRoleInsert, Do: apiManagePowerRoleInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/update"}, Power: base.PowerManagePowerRoleUpdate, Do: apiManagePowerRoleUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/power/role/delete"}, Power: base.PowerManagePowerRoleDelete, Do: apiManagePowerRoleDelete})

}

var (
	sqlManagePowerRolePage = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_ROLE,
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

	sqlManagePowerRoleList = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_ROLE,
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

	sqlManagePowerRoleOne = newSelectSql(&sqlModel.Select{
		Table: TABLE_POWER_ROLE,
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

	sqlManagePowerRoleInsert = newInsertSql(&sqlModel.Insert{
		Table: TABLE_POWER_ROLE,
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

	sqlManagePowerRoleUpdate = newUpdateSql(&sqlModel.Update{
		Table: TABLE_POWER_ROLE,
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

	sqlManagePowerRoleDelete = newDeleteSql(&sqlModel.Delete{
		Table: TABLE_POWER_ROLE,
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	})
)

func apiManagePowerRolePage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRolePage.GetSqlParam(role)
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

func apiManagePowerRoleList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRoleList.GetSqlParam(role)
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

func apiManagePowerRoleOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRoleOne.GetSqlParam(role)
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

func apiManagePowerRoleInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRoleInsert.GetSqlParam(role)
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

func apiManagePowerRoleUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRoleUpdate.GetSqlParam(role)
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

func apiManagePowerRoleDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	role := make(map[string]interface{})
	err = c.BindJSON(&role)
	if err != nil {
		return
	}
	sqlParam, err := sqlManagePowerRoleDelete.GetSqlParam(role)
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
