package organizationService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "organization"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ORGANIZATION",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ORGANIZATION (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	organizationId bigint(20) NOT NULL COMMENT '企业组织ID',
	PRIMARY KEY (serverId, enterpriseId, organizationId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业组织';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ORGANIZATION_USER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ORGANIZATION_USER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	organizationId bigint(20) NOT NULL COMMENT '企业组织ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	PRIMARY KEY (serverId, enterpriseId, organizationId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业用户';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
