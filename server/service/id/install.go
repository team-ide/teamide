package idService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "id"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_ID",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_ID (
	type int(2) NOT NULL COMMENT 'ID类型',
	id bigint(20) NOT NULL COMMENT 'ID',
	PRIMARY KEY (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ID';
				`,
		},
	})

	info.Stages = stages

	return
}
