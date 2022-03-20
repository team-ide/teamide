package certificateService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "certificate"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_CERTIFICATE",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_CERTIFICATE (
	certificateId bigint(20) NOT NULL COMMENT '凭证ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (certificateId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='凭证';
				`,
		},
	})

	info.Stages = stages

	return
}
