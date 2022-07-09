package module_power

import "teamide/internal/install"

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建权限角色表
		{
			Version: "1.0",
			Module:  ModulePower,
			Stage:   `创建表[` + TablePowerRole + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TablePowerRole + ` (
	powerRoleId bigint(20) NOT NULL COMMENT '权限角色ID',
	name varchar(50) NOT NULL COMMENT '名称',
	expirationTime datetime DEFAULT NULL COMMENT '过期时间',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (powerRoleId),
	KEY index_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TablePowerRoleComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TablePowerRole + ` (
	powerRoleId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	expirationTime datetime DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (powerRoleId)
);
`,
					`CREATE INDEX ` + TablePowerRole + `_index_name on ` + TablePowerRole + ` (name);`,
				},
			},
		},

		// 创建权限路由表
		{
			Version: "1.0",
			Module:  ModulePower,
			Stage:   `创建表[` + TablePowerRoute + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TablePowerRoute + ` (
	powerRouteId bigint(20) NOT NULL COMMENT '权限路由ID',
	powerRoleId bigint(20) NOT NULL COMMENT '权限角色',
	name varchar(50) NOT NULL COMMENT '名称',
	route varchar(50) NOT NULL COMMENT '路由',
	expirationTime datetime DEFAULT NULL COMMENT '过期时间',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (powerRouteId),
	KEY index_powerRoleId (powerRoleId),
	KEY index_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TablePowerRouteComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TablePowerRoute + ` (
	powerRouteId bigint(20) NOT NULL,
	powerRoleId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	route varchar(50) NOT NULL,
	expirationTime datetime DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (powerRoleId)
);
`,
					`CREATE INDEX ` + TablePowerRoute + `_index_powerRoleId on ` + TablePowerRoute + ` (powerRoleId);`,
					`CREATE INDEX ` + TablePowerRoute + `_index_name on ` + TablePowerRoute + ` (name);`,
				},
			},
		},

		// 创建权限用户表
		{
			Version: "1.0",
			Module:  ModulePower,
			Stage:   `创建表[` + TablePowerUser + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TablePowerUser + ` (
	powerUserId bigint(20) NOT NULL COMMENT '权限用户ID',
	userId bigint(20) NOT NULL COMMENT '用户',
	expirationTime datetime DEFAULT NULL COMMENT '过期时间',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (powerUserId),
	KEY index_userId (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TablePowerUserComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TablePowerUser + ` (
	powerUserId bigint(20) NOT NULL,
	userId bigint(20) NOT NULL,
	expirationTime datetime DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (powerUserId)
);
`,
					`CREATE INDEX ` + TablePowerUser + `_index_userId on ` + TablePowerUser + ` (userId);`,
				},
			},
		},
	}
}
