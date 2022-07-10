package module_ssh_tunnel

import "teamide/internal/install"

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建 隧道 表
		{
			Version: "1.1.0",
			Module:  ModuleTunnel,
			Stage:   `创建表[` + TableTunnel + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableTunnel + ` (
	tunnelId bigint(20) NOT NULL COMMENT '工具箱ID',
	tunnelType varchar(10) NOT NULL COMMENT '工具箱类型',
	name varchar(50) NOT NULL COMMENT '名称',
	option varchar(2000) NOT NULL COMMENT '配置',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	deleteUserId bigint(20) DEFAULT NULL COMMENT '删除用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (tunnelId),
	KEY index_userId (userId),
	KEY index_tunnelType (tunnelType),
	KEY index_name (name),
	KEY index_deleteUserId (deleteUserId),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableTunnelComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableTunnel + ` (
	tunnelId bigint(20) NOT NULL,
	tunnelType varchar(10) NOT NULL,
	name varchar(50) NOT NULL,
	option varchar(2000) NOT NULL,
	userId bigint(20) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	deleteUserId bigint(20) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (tunnelId)
);
`,
					`CREATE INDEX ` + TableTunnel + `_index_userId on ` + TableTunnel + ` (userId);`,
					`CREATE INDEX ` + TableTunnel + `_index_tunnelType on ` + TableTunnel + ` (tunnelType);`,
					`CREATE INDEX ` + TableTunnel + `_index_name on ` + TableTunnel + ` (name);`,
					`CREATE INDEX ` + TableTunnel + `_index_deleted on ` + TableTunnel + ` (deleted);`,
				},
			},
		},
	}

}
