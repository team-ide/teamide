package installService

import (
	base2 "teamide/internal/server/base"
	component2 "teamide/internal/server/component"
)

func install(info *base2.InstallInfo) {
	if info == nil || info.Stages == nil {
		return
	}
	for _, stage := range info.Stages {
		installStage(info, stage)
	}
}

func installStage(info *base2.InstallInfo, stage *base2.InstallStageInfo) {
	if info == nil || stage == nil {
		return
	}

	if base2.IS_STAND_ALONE {
		return
	}
	var err error
	var res int64
	res, err = component2.DB.Count(base2.SqlParam{
		Sql:    "SELECT count(1) FROM  " + TABLE_INSTALL + " WHERE module=? AND stage=? ",
		Params: []interface{}{info.Module, stage.Stage},
	})
	if err != nil {
		panic(err)
	}
	if res > 0 {
		return
	}
	detailInfo := make(map[string]interface{})
	if stage.Sql != nil {
		detailInfo["sql"] = stage.Sql
		_, err = component2.DB.Exec(base2.SqlParam{Sql: stage.Sql.MySql})
		if err != nil {
			panic(err)
		}
	}
	detail := base2.ToJSON(detailInfo)
	// 加密detail
	detail = component2.RSAEncrypt(detail)
	sqlParam := base2.SqlParam{
		Sql:    "INSERT INTO  " + TABLE_INSTALL + " (module, stage, detail, createTime) VALUES (?, ?, ?, ?) ",
		Params: []interface{}{info.Module, stage.Stage, detail, base2.Now()},
	}

	_, err = component2.DB.Exec(sqlParam)
	if err != nil {
		panic(err)
	}
}

func (this_ *Service) GetInstall() (info *base2.InstallInfo) {

	info = &base2.InstallInfo{}

	info.Module = "install"
	stages := []*base2.InstallStageInfo{}

	info.Stages = stages

	return
}
