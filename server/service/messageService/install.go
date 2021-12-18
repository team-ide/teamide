package messageService

import "server/base"

func GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "message"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_MESSAGE",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_MESSAGE (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	messageId bigint(20) NOT NULL COMMENT '消息ID',
	PRIMARY KEY (serverId, messageId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
