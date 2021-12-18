package userService

import (
	"server/base"
	"server/component"
	"server/model/sqlModel"

	"github.com/gin-gonic/gin"
)

func bindManageSystemSettingApi(appendApi func(apis ...*base.ApiWorker)) {

	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/setting/page"}, Power: base.PowerManageSystemSettingPage, Do: apiManageSystemSettingPage})
	appendApi(&base.ApiWorker{Apis: []string{"/manage/system/setting/update"}, Power: base.PowerManageSystemSettingUpdate, Do: apiManageSystemSettingUpdate})

}

var (
	sqlManageSystemSettingUpdate = &sqlModel.Update{
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

func apiManageSystemSettingPage(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	return
}

func apiManageSystemSettingUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	data := make(map[string]interface{})
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	sqlParam, err := sqlManageSystemSettingUpdate.GetSqlParam(data)
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
