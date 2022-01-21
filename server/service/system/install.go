package systemService

import "teamide/server/base"

func (this_ *SystemService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "system"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SYSTEM_LOG",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SYSTEM_LOG (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	logId bigint(20) NOT NULL COMMENT '日志ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, logId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统日志';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_SYSTEM_SETTING",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_SYSTEM_SETTING (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	settingId bigint(20) NOT NULL COMMENT '设置ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, settingId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统设置';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
