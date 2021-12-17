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

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE_POSITION",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ENTERPRISE_POSITION (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	positionId bigint(20) NOT NULL COMMENT '企业职位ID',
	PRIMARY KEY (serverId, enterpriseId, positionId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业职位';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE_LEVEL",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ENTERPRISE_LEVEL (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	levelId bigint(20) NOT NULL COMMENT '企业级别ID',
	PRIMARY KEY (serverId, enterpriseId, levelId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业级别';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
