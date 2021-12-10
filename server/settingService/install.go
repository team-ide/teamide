package settingService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "setting"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SETTING",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SETTING (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	settingId bigint(20) NOT NULL COMMENT '设置ID',
	PRIMARY KEY (serverId, settingId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设置';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
