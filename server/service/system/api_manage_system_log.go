package systemService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindManageSystemLogApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/page"}, Power: base.PowerManageSystemLogPage, Do: apiManageSystemLogPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/list"}, Power: base.PowerManageSystemLogPage, Do: apiManageSystemLogList})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/one"}, Power: base.PowerManageSystemLogPage, Do: apiManageSystemLogOne})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/insert"}, Power: base.PowerManageSystemLogInsert, Do: apiManageSystemLogInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/update"}, Power: base.PowerManageSystemLogUpdate, Do: apiManageSystemLogUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/log/delete"}, Power: base.PowerManageSystemLogDelete, Do: apiManageSystemLogDelete})

}

var (
	sqlManageSystemLogPage = &sqlModel.Select{
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

	sqlManageSystemLogList = &sqlModel.Select{
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

	sqlManageSystemLogOne = &sqlModel.Select{
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

	sqlManageSystemLogInsert = &sqlModel.Insert{
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

	sqlManageSystemLogUpdate = &sqlModel.Update{
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

	sqlManageSystemLogDelete = &sqlModel.Update{
		Table: "TM_USER",
		Columns: []*sqlModel.UpdateColumn{
			{Name: "deleteState", ValueScript: "1"},
		},
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	}

)

func apiManageSystemLogPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogPage.GetSqlParam(data)
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

func apiManageSystemLogList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogList.GetSqlParam(data)
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

func apiManageSystemLogOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogOne.GetSqlParam(data)
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

func apiManageSystemLogInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogInsert.GetSqlParam(data)
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

func apiManageSystemLogUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogUpdate.GetSqlParam(data)
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

func apiManageSystemLogDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemLogDelete.GetSqlParam(data)
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