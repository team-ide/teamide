package userService

import (
	"server/base"
	"server/component"
	"server/model/modelSql"
)

var (
	selectUserPageSql = &modelSql.Select{
		Table: "TM_USER",
		Columns: []*modelSql.SelectColumn{
			{Name: "userId"},
			{Name: "name"},
			{Name: "email"},
		},
		Wheres: []*modelSql.Where{
			{Name: "account"},
			{Piece: true,
				Wheres: []*modelSql.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
		UnionSelects: []*modelSql.Select{
			{
				Table: "TM_USER",
				Columns: []*modelSql.SelectColumn{
					{Name: "userId"},
					{Name: "name"},
					{Name: "email"},
				},
				Wheres: []*modelSql.Where{
					{Name: "name"},
					{Name: "email"},
				},
			},
		},
	}
)

func UserQueryPage(data map[string]interface{}) (page *base.PageBean, err error) {
	var sqlParam base.SqlParam
	sqlParam, err = selectUserPageSql.GetSqlParam(data)
	if err != nil {
		return
	}

	page, err = component.DB.QueryPage(sqlParam, base.NewUserEntityInterface)
	if err != nil {
		return
	}
	return
}
