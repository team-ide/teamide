package userService

import (
	"server/base"
	"server/component"
	sqlModel "server/model/sql"

	"github.com/gin-gonic/gin"
)

func bindUserPasswordApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/user/password/page"}, Power: base.PowerUserPasswordPage, Do: apiUserPasswordPage})
	appendApi(&base.ApiWorker{Apis: []string{"/user/password/update"}, Power: base.PowerUserPasswordUpdate, Do: apiUserPasswordUpdate})

}

var (
	sqlUserPasswordUpdate = &sqlModel.Update{
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
)

func apiUserPasswordPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

func apiUserPasswordUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserPasswordUpdate.GetSqlParam(data)
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
