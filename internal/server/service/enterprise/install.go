package enterpriseService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "enterprise"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ENTERPRISE (
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (enterpriseId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE_POSITION",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ENTERPRISE_POSITION (
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	positionId bigint(20) NOT NULL COMMENT '企业职位ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (enterpriseId, positionId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业职位';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ENTERPRISE_LEVEL",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ENTERPRISE_LEVEL (
	enterpriseId bigint(20) NOT NULL COMMENT '企业ID',
	levelId bigint(20) NOT NULL COMMENT '企业级别ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (enterpriseId, levelId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='企业级别';
				`,
		},
	})

	info.Stages = stages

	return
}
