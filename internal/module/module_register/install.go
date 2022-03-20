package module_register

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建注册表
		&install.StageModel{
			Version: "1.0",
			Module:  ModuleRegister,
			Stage:   `创建表[` + TableRegister + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableRegister + ` (
	registerId bigint(20) NOT NULL COMMENT '注册ID',
	name varchar(50) NOT NULL COMMENT '名称',
	account varchar(20) NOT NULL COMMENT '账号',
	email varchar(50) DEFAULT NULL COMMENT '邮箱',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	sourceType int(10) NOT NULL COMMENT '来源类型',
	source varchar(100) DEFAULT NULL COMMENT '来源',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (registerId),
	KEY index_userId (userId),
	KEY index_name (name),
	KEY index_account (account),
	KEY index_email (email),
	KEY index_sourceType (sourceType),
	KEY index_source (source),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='注册';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableRegister + ` (
	registerId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	account varchar(20) NOT NULL,
	email varchar(50) DEFAULT NULL,
	ip varchar(50) DEFAULT NULL,
	sourceType int(10) NOT NULL,
	source varchar(100) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (registerId)
);
`,
					`CREATE INDEX ` + TableRegister + `_index_userId on ` + TableRegister + ` (userId);`,
					`CREATE INDEX ` + TableRegister + `_index_name on ` + TableRegister + ` (name);`,
					`CREATE INDEX ` + TableRegister + `_index_account on ` + TableRegister + ` (account);`,
					`CREATE INDEX ` + TableRegister + `_index_email on ` + TableRegister + ` (email);`,
					`CREATE INDEX ` + TableRegister + `_index_sourceType on ` + TableRegister + ` (sourceType);`,
					`CREATE INDEX ` + TableRegister + `_index_source on ` + TableRegister + ` (source);`,
					`CREATE INDEX ` + TableRegister + `_index_deleted on ` + TableRegister + ` (deleted);`,
				},
			},
		},
	}

}
