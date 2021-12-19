package powerService

import "server/base"

func (this_ *PowerService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "power"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_ROLE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_POWER_ROLE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, roleId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_USER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_POWER_USER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, roleId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_DATA",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_POWER_DATA (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	dataId bigint(20) NOT NULL COMMENT '数据ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, roleId, dataId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据权限';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_POWER_ACTION",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_POWER_ACTION (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	roleId bigint(20) NOT NULL COMMENT '角色ID',
	actionId bigint(20) NOT NULL COMMENT '功能ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, roleId, actionId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='功能权限';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
