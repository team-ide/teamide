package jobService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "job"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_JOB",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_JOB (
	jobId bigint(20) NOT NULL COMMENT 'JobID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (jobId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Job';
				`,
		},
	})

	info.Stages = stages

	return
}
