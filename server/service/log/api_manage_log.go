package logService

import (
	"server/base"
	"server/component"
	"server/model/modelSql"

	"github.com/gin-gonic/gin"
)

func bindManageLogApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/page"}, Power: base.PowerManageLogPage, Do: apiManageLogPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/list"}, Power: base.PowerManageLogPage, Do: apiManageLogList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/one"}, Power: base.PowerManageLogPage, Do: apiManageLogOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/insert"}, Power: base.PowerManageLogInsert, Do: apiManageLogInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/update"}, Power: base.PowerManageLogUpdate, Do: apiManageLogUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/log/delete"}, Power: base.PowerManageLogDelete, Do: apiManageLogDelete})

}

var (
	sqlManageLogPage = &modelSql.Select{
		Table: "TM_USER",
		Columns: []*modelSql.SelectColumn{
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
		Wheres: []*modelSql.Where{
			{Name: "deletedState", ValueScript: "1"},
			{
				Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "name", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "email", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "account", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "activedState"},
					{Name: "lockedState"},
					{Name: "enabledState"},
				}},
		},
	}

	sqlManageLogList = &modelSql.Select{
		Table: "TM_USER",
		Columns: []*modelSql.SelectColumn{
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
		Wheres: []*modelSql.Where{
			{Name: "deletedState", ValueScript: "1"},
			{
				Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "name", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "email", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "account", Operator: modelSql.LIKE_BEFORE.Value},
					{Name: "activedState"},
					{Name: "lockedState"},
					{Name: "enabledState"},
				}},
		},
	}

	sqlManageLogOne = &modelSql.Select{
		Table: "TM_USER",
		Columns: []*modelSql.SelectColumn{
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
		Wheres: []*modelSql.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageLogInsert = &modelSql.Insert{
		Table: "TM_USER",
		Columns: []*modelSql.InsertColumn{
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

	sqlManageLogUpdate = &modelSql.Update{
		Table: "TM_USER",
		Columns: []*modelSql.UpdateColumn{
			{Name: "name"},
			{Name: "avatar"},
			{Name: "email"},
			{Name: "account"},
			{Name: "updateTime", ValueScript: "now()"},
		},
		Wheres: []*modelSql.Where{
			{Name: "userId", Required: true},
		},
	}

	sqlManageLogDelete = &modelSql.Update{
		Table: "TM_USER",
		Columns: []*modelSql.UpdateColumn{
			{Name: "deleteState", ValueScript: "1"},
		},
		Wheres: []*modelSql.Where{
			{Name: "userId", Required: true},
		},
	}

)

func apiManageLogPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogPage.GetSqlParam(data)
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

func apiManageLogList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogList.GetSqlParam(data)
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

func apiManageLogOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogOne.GetSqlParam(data)
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

func apiManageLogInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogInsert.GetSqlParam(data)
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

func apiManageLogUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogUpdate.GetSqlParam(data)
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

func apiManageLogDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageLogDelete.GetSqlParam(data)
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