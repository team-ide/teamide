package loginService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "log"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_LOGIN",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_LOGIN (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	loginId bigint(20) NOT NULL COMMENT '登录ID',
	PRIMARY KEY (serverId, loginId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
