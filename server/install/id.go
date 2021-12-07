package install

import "db"

func getID() (info *InstallInfo) {

	info = &InstallInfo{}

	info.Module = "id"
	stages := []*InstallStageInfo{}

	stages = append(stages, &InstallStageInfo{
		Stage: "CREATE TABLE TM_ID",
		SqlParam: db.SqlParam{
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
