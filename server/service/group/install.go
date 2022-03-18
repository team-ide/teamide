package groupService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "group"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_GROUP",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_GROUP (
	groupId bigint(20) NOT NULL COMMENT '组ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (groupId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_GROUP_USER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_GROUP_USER (
	groupId bigint(20) NOT NULL COMMENT '组ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (groupId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组用户';
				`,
		},
	})

	info.Stages = stages

	return
}
