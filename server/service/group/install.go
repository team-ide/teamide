package groupService

import "teamide/server/base"

func (this_ *GroupService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "group"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_GROUP",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_GROUP (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	groupId bigint(20) NOT NULL COMMENT '组ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, groupId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_GROUP_USER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_GROUP_USER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	groupId bigint(20) NOT NULL COMMENT '组ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, groupId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组用户';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
