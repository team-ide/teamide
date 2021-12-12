package powerService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "power"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_POWER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	powerId bigint(20) NOT NULL COMMENT '权限ID',
	PRIMARY KEY (serverId, powerId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
