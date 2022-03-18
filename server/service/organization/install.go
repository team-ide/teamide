package organizationService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "organization"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ORGANIZATION",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ORGANIZATION (
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	organizationId bigint(20) NOT NULL COMMENT '企业组织ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (enterpriseId, organizationId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业组织';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ORGANIZATION_USER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ORGANIZATION_USER (
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	organizationId bigint(20) NOT NULL COMMENT '企业组织ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (enterpriseId, organizationId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业用户';
				`,
		},
	})

	info.Stages = stages

	return
}
