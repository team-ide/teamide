package spaceService

import "server/base"

func (this_ *SpaceService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "space"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SPACE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	PRIMARY KEY (serverId, spaceId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE_USER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SPACE_USER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	PRIMARY KEY (serverId, spaceId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间用户';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SPACE_POWER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SPACE_POWER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	spaceId bigint(20) NOT NULL COMMENT '空间ID',
	powerId bigint(20) NOT NULL COMMENT '用户ID',
	PRIMARY KEY (serverId, spaceId, powerId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='空间权限';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
