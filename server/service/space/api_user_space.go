package spaceService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindUserSpaceApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/user/space/page"}, Power: base.PowerUserSpacePage, Do: apiUserSpacePage})
	appendApi(&base.ApiWorker{Apis: []string{"/user/space/list"}, Power: base.PowerUserSpacePage, Do: apiUserSpaceList})
	appendApi(&base.ApiWorker{Apis: []string{"/user/space/one"}, Power: base.PowerUserSpacePage, Do: apiUserSpaceOne})
	appendApi(&base.ApiWorker{Apis: []string{"/user/space/insert"}, Power: base.PowerUserSpaceInsert, Do: apiUserSpaceInsert})
	appendApi(&base.ApiWorker{Apis: []string{"/user/space/update"}, Power: base.PowerUserSpaceUpdate, Do: apiUserSpaceUpdate})
	appendApi(&base.ApiWorker{Apis: []string{"/user/space/delete"}, Power: base.PowerUserSpaceDelete, Do: apiUserSpaceDelete})

}

var (
	sqlUserSpacePage = newSelectSql(&sqlModel.Select{
		Table: TABLE_SPACE,
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

	sqlUserSpaceList = newSelectSql(&sqlModel.Select{
		Table: TABLE_SPACE,
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

	sqlUserSpaceOne = newSelectSql(&sqlModel.Select{
		Table: TABLE_SPACE,
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

	sqlUserSpaceInsert = newInsertSql(&sqlModel.Insert{
		Table: TABLE_SPACE,
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

	sqlUserSpaceUpdate = newUpdateSql(&sqlModel.Update{
		Table: TABLE_SPACE,
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

	sqlUserSpaceDelete = newDeleteSql(&sqlModel.Delete{
		Table: TABLE_SPACE,
		Wheres: []*sqlModel.Where{
			{Name: "userId", Required: true},
		},
	})
)

func apiUserSpacePage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpacePage.GetSqlParam(data)
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

func apiUserSpaceList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpaceList.GetSqlParam(data)
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

func apiUserSpaceOne(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpaceOne.GetSqlParam(data)
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

func apiUserSpaceInsert(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpaceInsert.GetSqlParam(data)
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

func apiUserSpaceUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpaceUpdate.GetSqlParam(data)
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

func apiUserSpaceDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSpaceDelete.GetSqlParam(data)
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
