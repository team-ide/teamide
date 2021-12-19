package jobService

import "server/base"

func (this_ *JobService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "job"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_JOB",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_JOB (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	jobId bigint(20) NOT NULL COMMENT 'JobID',
	PRIMARY KEY (serverId, jobId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Job';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
