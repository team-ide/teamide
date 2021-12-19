package certificateService

import "server/base"

func (this_ *CertificateService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "certificate"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_CERTIFICATE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_CERTIFICATE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	certificateId bigint(20) NOT NULL COMMENT '组ID',
	PRIMARY KEY (serverId, certificateId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
