package groupService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

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
	PRIMARY KEY (serverId, groupId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
