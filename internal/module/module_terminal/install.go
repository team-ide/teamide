package module_terminal

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建登录表
		{
			Version: "1.0",
			Module:  ModuleTerminalLog,
			Stage:   `创建表[` + TableTerminalLog + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableTerminalLog + ` (
	terminalLogId bigint(20) NOT NULL COMMENT '日志ID',
	loginId bigint(20) DEFAULT NULL COMMENT '登录ID',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	userName varchar(50) DEFAULT NULL COMMENT '用户名称',
	userAccount varchar(50) DEFAULT NULL COMMENT '用户账号',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	userAgent text DEFAULT NULL COMMENT 'User-Agent',
	place varchar(20) DEFAULT NULL COMMENT '位置',
	placeId varchar(20) DEFAULT NULL COMMENT '位置ID',
	command text DEFAULT NULL COMMENT '命令',
	createTime datetime NOT NULL COMMENT '创建时间',
	PRIMARY KEY (terminalLogId),
	KEY index_loginId (loginId),
	KEY index_userId (userId),
	KEY index_userName (userName),
	KEY index_userAccount (userAccount),
	KEY index_ip (ip),
	KEY index_place (place),
	KEY index_placeId (placeId),
	KEY index_useTime (createTime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableTerminalLogComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableTerminalLog + ` (
	terminalLogId bigint(20) NOT NULL,
	loginId bigint(20) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	userName varchar(50) DEFAULT NULL,
	userAccount varchar(50) DEFAULT NULL,
	ip varchar(50) DEFAULT NULL,
	userAgent text DEFAULT NULL,
	place varchar(20) DEFAULT NULL,
	placeId varchar(20) DEFAULT NULL,
	command text DEFAULT NULL,
	createTime datetime NOT NULL,
	PRIMARY KEY (terminalLogId)
);
`,
					`CREATE INDEX ` + TableTerminalLog + `_index_loginId on ` + TableTerminalLog + ` (loginId);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_userId on ` + TableTerminalLog + ` (userId);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_userName on ` + TableTerminalLog + ` (userName);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_userAccount on ` + TableTerminalLog + ` (userAccount);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_ip on ` + TableTerminalLog + ` (ip);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_place on ` + TableTerminalLog + ` (place);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_placeId on ` + TableTerminalLog + ` (placeId);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_createTime on ` + TableTerminalLog + ` (createTime);`,
				},
			},
		},
	}
}
