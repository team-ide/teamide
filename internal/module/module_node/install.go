package module_node

import "teamide/internal/install"

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建 节点 表
		{
			Version: "1.1.0",
			Module:  ModuleNode,
			Stage:   `创建表[` + TableNode + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableNode + ` (
	nodeId bigint(20) NOT NULL COMMENT '节点ID',
	nodeServerId varchar(50) NOT NULL COMMENT '节点服务ID',
	name varchar(50) NOT NULL COMMENT '名称',
	comment varchar(200) DEFAULT NULL COMMENT '说明',
	option varchar(2000) NOT NULL COMMENT '配置',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	deleteUserId bigint(20) DEFAULT NULL COMMENT '删除用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (nodeId),
	KEY index_nodeServerId (nodeServerId),
	KEY index_userId (userId),
	KEY index_name (name),
	KEY index_deleteUserId (deleteUserId),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableNodeComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableNode + ` (
	nodeId bigint(20) NOT NULL,
	nodeServerId varchar(50) NOT NULL,
	name varchar(50) NOT NULL,
	comment varchar(200) DEFAULT NULL,
	option varchar(2000) NOT NULL,
	userId bigint(20) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	deleteUserId bigint(20) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (nodeId)
);
`,
					`CREATE INDEX ` + TableNode + `_index_nodeServerId on ` + TableNode + ` (nodeServerId);`,
					`CREATE INDEX ` + TableNode + `_index_userId on ` + TableNode + ` (userId);`,
					`CREATE INDEX ` + TableNode + `_index_name on ` + TableNode + ` (name);`,
					`CREATE INDEX ` + TableNode + `_index_deleted on ` + TableNode + ` (deleted);`,
				},
			},
		},
	}

}
