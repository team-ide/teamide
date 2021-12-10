package enterpriseService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "enterprise"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ENTERPRISE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	PRIMARY KEY (serverId, enterpriseId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
