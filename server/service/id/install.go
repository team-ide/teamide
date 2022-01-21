package idService

import "teamide/server/base"

func (this_ *IdService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "id"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ID",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_ID (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	type int(2) NOT NULL COMMENT 'ID类型',
	id bigint(20) NOT NULL COMMENT 'ID',
	PRIMARY KEY (serverId, type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ID';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
