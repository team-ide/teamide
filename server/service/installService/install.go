package installService

import (
	"server/base"
	"server/component"
)

func Install(info *base.InstallInfo) {
	if info == nil || info.Stages == nil {
		return
	}
	for _, stage := range info.Stages {
		InstallStage(info, stage)
	}
}

func InstallStage(info *base.InstallInfo, stage *base.InstallStageInfo) {
	if info == nil || stage == nil {
		return
	}
	var err error
	var res int64
	res, err = component.DB.Count(base.SqlParam{
		Sql:    "SELECT count(1) FROM  " + base.TABLE_INSTALL + " WHERE module=? AND stage=? ",
		Params: []interface{}{info.Module, stage.Stage},
	})
	if err != nil {
		panic(err)
	}
	if res > 0 {
		return
	}
	detailInfo := make(map[string]interface{})
	if stage.SqlParam.Sql != "" {
		detailInfo["sqlParam"] = stage.SqlParam
		_, err = component.DB.Exec(stage.SqlParam)
		if err != nil {
			panic(err)
		}
	}
	detail := base.ToJSON(detailInfo)
	// 加密detail
	detail = base.AesEncryptCBC(detail)
	sqlParam := base.SqlParam{
		Sql:    "INSERT INTO  " + base.TABLE_INSTALL + " (module, stage, detail, createTime) VALUES (?, ?, ?, ?) ",
		Params: []interface{}{info.Module, stage.Stage, detail, base.Now()},
	}

	_, err = component.DB.Exec(sqlParam)
	if err != nil {
		panic(err)
	}
}
