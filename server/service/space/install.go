package spaceService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "space"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_SPACE (
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (spaceId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE_USER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_SPACE_USER (
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (spaceId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间用户';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE_POWER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_SPACE_POWER (
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	powerId bigint(20) NOT NULL COMMENT '权限ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (spaceId, powerId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间权限';
				`,
		},
	})

	info.Stages = stages

	return
}
