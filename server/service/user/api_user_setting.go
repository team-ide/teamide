package userService

import (
	"server/base"
	"server/component"
	"server/model/sqlModel"

	"github.com/gin-gonic/gin"
)

func bindUserSettingApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/user/setting/page"}, Power: base.PowerUserSettingPage, Do: apiUserSettingPage})
	appendApi(&base.ApiWorker{Apis: []string{"/user/setting/update"}, Power: base.PowerUserSettingUpdate, Do: apiUserSettingUpdate})

}

var (
	sqlUserSettingUpdate = &sqlModel.Update{
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

func apiUserSettingPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

func apiUserSettingUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlUserSettingUpdate.GetSqlParam(data)
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
