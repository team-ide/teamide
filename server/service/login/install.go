package loginService

import "server/base"

func (this_ *LoginService) GetInstall() (info *base.InstallInfo) {

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
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, loginId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
