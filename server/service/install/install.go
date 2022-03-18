package installService

import (
	"teamide/server/base"
	"teamide/server/component"
)

func install(info *base.InstallInfo) {
	if info == nil || info.Stages == nil {
		return
	}
	for _, stage := range info.Stages {
		installStage(info, stage)
	}
}

func installStage(info *base.InstallInfo, stage *base.InstallStageInfo) {
	if info == nil || stage == nil {
		return
	}

	if base.IS_STAND_ALONE {
		return
	}
	var err error
	var res int64
	res, err = component.DB.Count(base.SqlParam{
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
		_, err = component.DB.Exec(base.SqlParam{Sql: stage.Sql.MySql})
		if err != nil {
			panic(err)
		}
	}
	detail := base.ToJSON(detailInfo)
	// 加密detail
	detail = component.RSAEncrypt(detail)
	sqlParam := base.SqlParam{
		Sql:    "INSERT INTO  " + TABLE_INSTALL + " (module, stage, detail, createTime) VALUES (?, ?, ?, ?) ",
		Params: []interface{}{info.Module, stage.Stage, detail, base.Now()},
	}

	_, err = component.DB.Exec(sqlParam)
	if err != nil {
		panic(err)
	}
}

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "install"
	stages := []*base.InstallStageInfo{}

	info.Stages = stages

	return
}
