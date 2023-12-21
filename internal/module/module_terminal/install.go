package module_terminal

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建 终端日志 表 开始 已废弃
		{
			Version: "1.0",
			Module:  ModuleTerminalLog,
			Stage:   `创建表[` + TableTerminalLog + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableTerminalLog + ` (
	terminalLogId bigint(20) NOT NULL COMMENT '日志ID',
	loginId bigint(20) DEFAULT NULL COMMENT '登录ID',
	workerId varchar(50) DEFAULT NULL COMMENT '工作ID',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	userName varchar(50) DEFAULT NULL COMMENT '用户名称',
	userAccount varchar(50) DEFAULT NULL COMMENT '用户账号',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	place varchar(20) DEFAULT NULL COMMENT '位置',
	placeId varchar(20) DEFAULT NULL COMMENT '位置ID',
	command text DEFAULT NULL COMMENT '命令',
	userAgent text DEFAULT NULL COMMENT 'User-Agent',
	createTime datetime NOT NULL COMMENT '创建时间',
	PRIMARY KEY (terminalLogId),
	KEY index_loginId (loginId),
	KEY index_workerId (workerId),
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
	workerId varchar(50) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	userName varchar(50) DEFAULT NULL,
	userAccount varchar(50) DEFAULT NULL,
	ip varchar(50) DEFAULT NULL,
	place varchar(20) DEFAULT NULL,
	placeId varchar(20) DEFAULT NULL,
	command text DEFAULT NULL,
	userAgent text DEFAULT NULL,
	createTime datetime NOT NULL,
	PRIMARY KEY (terminalLogId)
);
`,
					`CREATE INDEX ` + TableTerminalLog + `_index_loginId on ` + TableTerminalLog + ` (loginId);`,
					`CREATE INDEX ` + TableTerminalLog + `_index_workerId on ` + TableTerminalLog + ` (workerId);`,
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
		// 创建 终端日志 表 结束 已废弃

		// 创建 终端命令 表 开始
		{
			Version: "1.0",
			Module:  ModuleTerminalCommand,
			Stage:   `创建表[` + TableTerminalCommand + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableTerminalCommand + ` (
	terminalCommandId bigint(20) NOT NULL COMMENT '日志ID',
	loginId bigint(20) DEFAULT NULL COMMENT '登录ID',
	workerId varchar(50) DEFAULT NULL COMMENT '工作ID',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	userName varchar(50) DEFAULT NULL COMMENT '用户名称',
	userAccount varchar(50) DEFAULT NULL COMMENT '用户账号',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	place varchar(20) DEFAULT NULL COMMENT '位置',
	placeId varchar(20) DEFAULT NULL COMMENT '位置ID',
	command text DEFAULT NULL COMMENT '命令',
	userAgent text DEFAULT NULL COMMENT 'User-Agent',
	createTime datetime NOT NULL COMMENT '创建时间',
	PRIMARY KEY (terminalCommandId),
	KEY index_loginId (loginId),
	KEY index_workerId (workerId),
	KEY index_userId (userId),
	KEY index_userName (userName),
	KEY index_userAccount (userAccount),
	KEY index_ip (ip),
	KEY index_place (place),
	KEY index_placeId (placeId),
	KEY index_useTime (createTime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableTerminalCommandComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableTerminalCommand + ` (
	terminalCommandId bigint(20) NOT NULL,
	loginId bigint(20) DEFAULT NULL,
	workerId varchar(50) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	userName varchar(50) DEFAULT NULL,
	userAccount varchar(50) DEFAULT NULL,
	ip varchar(50) DEFAULT NULL,
	place varchar(20) DEFAULT NULL,
	placeId varchar(20) DEFAULT NULL,
	command text DEFAULT NULL,
	userAgent text DEFAULT NULL,
	createTime datetime NOT NULL,
	PRIMARY KEY (terminalCommandId)
);
`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_loginId on ` + TableTerminalCommand + ` (loginId);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_workerId on ` + TableTerminalCommand + ` (workerId);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_userId on ` + TableTerminalCommand + ` (userId);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_userName on ` + TableTerminalCommand + ` (userName);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_userAccount on ` + TableTerminalCommand + ` (userAccount);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_ip on ` + TableTerminalCommand + ` (ip);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_place on ` + TableTerminalCommand + ` (place);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_placeId on ` + TableTerminalCommand + ` (placeId);`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_createTime on ` + TableTerminalCommand + ` (createTime);`,
				},
			},
		},

		// 创建 终端命令 表 结束

		/** 终端命令 添加 类型、注释 开始**/
		{
			Version: "1.0.4",
			Module:  ModuleTerminalCommand,
			Stage:   `终端命令[` + ModuleTerminalCommand + `]添加类型[commandType]、说明[comment]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableTerminalCommand + ` ADD COLUMN comment varchar(500);`,
					`ALTER TABLE ` + TableTerminalCommand + ` ADD COLUMN commandType int(10) DEFAULT '1';`,
					`ALTER TABLE ` + TableTerminalCommand + ` ADD INDEX ` + TableTerminalCommand + `_index_commandType(commandType);`,
					`UPDATE ` + TableTerminalCommand + ` SET commandType=1 WHERE commandType=0 OR commandType IS NULL;`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableTerminalCommand + ` ADD comment varchar(500);`,
					`ALTER TABLE ` + TableTerminalCommand + ` ADD commandType int(10) DEFAULT '1';`,
					`CREATE INDEX ` + TableTerminalCommand + `_index_commandType on ` + TableTerminalCommand + ` (commandType);`,
					`UPDATE ` + TableTerminalCommand + ` SET commandType=1 WHERE commandType=0 OR commandType IS NULL;`,
				},
			},
		},
		/** 终端命令 添加 类型、注释 结束**/
	}
}
