package powerService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "power"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_ROLE",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_POWER_ROLE (
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (roleId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_USER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_POWER_USER (
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (roleId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_DATA",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_POWER_DATA (
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	dataId bigint(20) NOT NULL COMMENT '数据ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (roleId, dataId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据权限';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_ACTION",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_POWER_ACTION (
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	actionId bigint(20) NOT NULL COMMENT '功能ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (roleId, actionId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='功能权限';
				`,
		},
	})

	info.Stages = stages

	return
}
