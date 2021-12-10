package spaceService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

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

	info.Stages = stages

	return
}
