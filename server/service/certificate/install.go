package certificateService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "certificate"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_CERTIFICATE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_CERTIFICATE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	certificateId bigint(20) NOT NULL COMMENT '凭证ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, certificateId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='凭证';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
