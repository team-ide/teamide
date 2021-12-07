package install

import (
	"base"
	"db"
	"server"
)

type InstallInfo struct {
	Module string              `json:"module"`
	Stages []*InstallStageInfo `json:"stages"`
}

type InstallStageInfo struct {
	Stage    string      `json:"stage"`
	SqlParam db.SqlParam `json:"sqlParam"`
}

func Install(info *InstallInfo) {
	if info == nil || info.Stages == nil {
		return
	}
	for _, stage := range info.Stages {
		InstallStage(info, stage)
	}
}

func InstallStage(info *InstallInfo, stage *InstallStageInfo) {
	if info == nil || stage == nil {
		return
	}
	var err error
	var res int64
	res, err = db.DBService.Count(db.SqlParam{
		Sql:    "SELECT count(1) FROM  " + db.TABLE_INSTALL + " WHERE module=? AND stage=? ",
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
		_, err = db.DBService.Exec(stage.SqlParam)
		if err != nil {
			panic(err)
		}
	}
	detail := base.ToJSON(detailInfo)
	// 加密detail
	detail = server.Encrypt(detail)
	sqlParam := db.SqlParam{
		Sql:    "INSERT INTO  " + db.TABLE_INSTALL + " (module, stage, detail, createTime) VALUES (?, ?, ?, ?) ",
		Params: []interface{}{info.Module, stage.Stage, detail, base.Now()},
	}

	_, err = db.DBService.Exec(sqlParam)
	if err != nil {
		panic(err)
	}
}
