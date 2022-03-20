package logService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "log"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_LOG",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_LOG (
	logId bigint(20) NOT NULL COMMENT '日志ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (logId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志';
				`,
		},
	})

	info.Stages = stages

	return
}
