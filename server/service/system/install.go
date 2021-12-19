package systemService

import "server/base"

func (this_ *SystemService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "system"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SYSTEM",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SYSTEM (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	systemId bigint(20) NOT NULL COMMENT '系统ID',
	PRIMARY KEY (serverId, systemId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
