package logService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "log"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_LOG",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_LOG (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	logId bigint(20) NOT NULL COMMENT '日志ID',
	PRIMARY KEY (serverId, logId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
