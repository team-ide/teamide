package module_user

import "teamide/internal/install"

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建用户表
		{
			Version: "1.0",
			Module:  ModuleUser,
			Stage:   `创建表[` + TableUser + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableUser + ` (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	name varchar(50) NOT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	account varchar(20) NOT NULL COMMENT '账号',
	email varchar(50) DEFAULT NULL COMMENT '邮箱',
	activated int(1) NOT NULL DEFAULT 2 COMMENT '激活状态:1-激活、2-未激活',
	locked int(1) NOT NULL DEFAULT 2 COMMENT '锁定状态:1-锁定、2-未锁定',
	enabled int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-启用、2-停用',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (userId),
	KEY index_name (name),
	KEY index_account (account),
	KEY index_email (email),
	KEY index_activated (activated),
	KEY index_locked (locked),
	KEY index_enabled (enabled),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableUser + ` (
	userId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	avatar varchar(200) DEFAULT NULL,
	account varchar(20) NOT NULL,
	email varchar(50) DEFAULT NULL,
	activated int(1) NOT NULL DEFAULT 2,
	locked int(1) NOT NULL DEFAULT 2,
	enabled int(1) NOT NULL DEFAULT 1,
	deleted int(1) NOT NULL DEFAULT 2,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (userId)
);
`,
					`CREATE INDEX ` + TableUser + `_index_name on ` + TableUser + ` (name);`,
					`CREATE INDEX ` + TableUser + `_index_account on ` + TableUser + ` (account);`,
					`CREATE INDEX ` + TableUser + `_index_email on ` + TableUser + ` (email);`,
					`CREATE INDEX ` + TableUser + `_index_activated on ` + TableUser + ` (activated);`,
					`CREATE INDEX ` + TableUser + `_index_locked on ` + TableUser + ` (locked);`,
					`CREATE INDEX ` + TableUser + `_index_enabled on ` + TableUser + ` (enabled);`,
					`CREATE INDEX ` + TableUser + `_index_deleted on ` + TableUser + ` (deleted);`,
				},
			},
		},

		// 创建用户授权表
		{
			Version: "1.0",
			Module:  ModuleUser,
			Stage:   `创建表[` + TableUserAuth + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableUserAuth + ` (
	authId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	authType int(2) NOT NULL COMMENT '授权类型',
	openId varchar(100) NOT NULL COMMENT 'OpenID',
	name varchar(100) DEFAULT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	homepage varchar(200) DEFAULT NULL COMMENT '主页',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (authId),
	KEY index_userId (userId),
	KEY index_authType (authType),
	KEY index_openId (openId),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户授权';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableUserAuth + ` (
	authId bigint(20) NOT NULL,
	userId bigint(20) NOT NULL,
	authType int(2) NOT NULL,
	openId varchar(100) NOT NULL,
	name varchar(100) DEFAULT NULL,
	avatar varchar(200) DEFAULT NULL,
	homepage varchar(200) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (authId)
);
`,
					`CREATE INDEX ` + TableUserAuth + `_index_userId on ` + TableUserAuth + ` (userId);`,
					`CREATE INDEX ` + TableUserAuth + `_index_authType on ` + TableUserAuth + ` (authType);`,
					`CREATE INDEX ` + TableUserAuth + `_index_openId on ` + TableUserAuth + ` (openId);`,
					`CREATE INDEX ` + TableUserAuth + `_deleted on ` + TableUserAuth + ` (deleted);`,
				},
			},
		},

		// 创建用户密码表
		{
			Version: "1.0",
			Module:  ModuleUser,
			Stage:   `创建表[` + TableUserPassword + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableUserPassword + ` (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	salt varchar(20) NOT NULL COMMENT '盐',
	password varchar(100) NOT NULL COMMENT '密码',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户密码';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableUserPassword + ` (
	userId bigint(20) NOT NULL,
	salt varchar(20) NOT NULL,
	password varchar(100) NOT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (userId)
);
`,
				},
			},
		},

		// 创建用户设置表
		{
			Version: "1.0",
			Module:  ModuleUser,
			Stage:   `创建表[` + TableUserSetting + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableUserSetting + ` (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	name varchar(100) NOT NULL COMMENT '名称',
	value varchar(1000) DEFAULT NULL COMMENT '值',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId, name),
	KEY index_userId (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户设置';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableUserSetting + ` (
	userId bigint(20) NOT NULL,
	name varchar(100) NOT NULL,
	value varchar(1000) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (userId, name)
);
`,
					`CREATE INDEX ` + TableUserSetting + `_index_userId on ` + TableUserSetting + ` (userId);`,
				},
			},
		},
	}
}
