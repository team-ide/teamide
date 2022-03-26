package module_login

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建登录表
		&install.StageModel{
			Version: "1.0",
			Module:  ModuleLogin,
			Stage:   `创建表[` + TableLogin + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableLogin + ` (
	loginId bigint(20) NOT NULL COMMENT '登录ID',
	account varchar(20) NOT NULL COMMENT '账号',
	ip varchar(50) DEFAULT NULL COMMENT 'IP',
	sourceType int(10) NOT NULL COMMENT '来源类型',
	source varchar(100) DEFAULT NULL COMMENT '来源',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	loginTime datetime NOT NULL COMMENT '登录时间',
	logoutTime datetime DEFAULT NULL COMMENT '登出时间',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (loginId),
	KEY index_userId (userId),
	KEY index_account (account),
	KEY index_sourceType (sourceType),
	KEY index_source (source),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableLogin + ` (
	loginId bigint(20) NOT NULL,
	account varchar(20) NOT NULL,
	ip varchar(50) DEFAULT NULL,
	sourceType int(10) NOT NULL,
	source varchar(100) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	loginTime datetime NOT NULL,
	logoutTime datetime DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (loginId)
);
`,
					`CREATE INDEX ` + TableLogin + `_index_userId on ` + TableLogin + ` (userId);`,
					`CREATE INDEX ` + TableLogin + `_index_account on ` + TableLogin + ` (account);`,
					`CREATE INDEX ` + TableLogin + `_index_sourceType on ` + TableLogin + ` (sourceType);`,
					`CREATE INDEX ` + TableLogin + `_index_source on ` + TableLogin + ` (source);`,
					`CREATE INDEX ` + TableLogin + `_index_deleted on ` + TableLogin + ` (deleted);`,
				},
			},
		},
	}
}
