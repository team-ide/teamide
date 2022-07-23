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
	serverId varchar(50) NOT NULL COMMENT '节点服务ID',
	name varchar(50) NOT NULL COMMENT '名称',
	comment varchar(200) DEFAULT NULL COMMENT '说明',
	bindAddress varchar(50) DEFAULT NULL COMMENT '说明',
	bindToken varchar(50) DEFAULT NULL COMMENT '说明',
	connAddress varchar(100) DEFAULT NULL COMMENT '连接节点地址',
	connToken varchar(50) DEFAULT NULL COMMENT '连接节点Token',
	connServerIds varchar(2000) DEFAULT NULL COMMENT '连接节点服务ID',
	historyConnServerIds varchar(2000) DEFAULT NULL COMMENT '历史连接节点服务ID',
	option varchar(2000) DEFAULT NULL COMMENT '配置',
	isLocal int(1) DEFAULT NULL COMMENT '是本地',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	enabled int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-启用、2-停用',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	deleteUserId bigint(20) DEFAULT NULL COMMENT '删除用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (nodeId),
	KEY index_userId (userId),
	KEY index_name (name),
	KEY index_enabled (enabled),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableNodeComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableNode + ` (
	nodeId bigint(20) NOT NULL,
	serverId varchar(50) NOT NULL,
	name varchar(50) NOT NULL,
	comment varchar(200) DEFAULT NULL,
	bindAddress varchar(50) DEFAULT NULL,
	bindToken varchar(50) DEFAULT NULL,
	connAddress varchar(100) DEFAULT NULL,
	connToken varchar(50) DEFAULT NULL,
	connServerIds varchar(2000) DEFAULT NULL,
	historyConnServerIds varchar(2000) DEFAULT NULL,
	option varchar(2000) DEFAULT NULL,
	isLocal int(1) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	enabled int(1) NOT NULL DEFAULT 1,
	deleted int(1) NOT NULL DEFAULT 2,
	deleteUserId bigint(20) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (nodeId)
);
`,
					`CREATE INDEX ` + TableNode + `_index_userId on ` + TableNode + ` (userId);`,
					`CREATE INDEX ` + TableNode + `_index_name on ` + TableNode + ` (name);`,
					`CREATE INDEX ` + TableNode + `_index_enabled on ` + TableNode + ` (enabled);`,
					`CREATE INDEX ` + TableNode + `_index_deleted on ` + TableNode + ` (deleted);`,
				},
			},
		},

		// 创建 节点网络代理 表
		{
			Version: "1.1.0",
			Module:  ModuleNode,
			Stage:   `创建表[` + TableNodeNetProxy + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableNodeNetProxy + ` (
	netProxyId bigint(20) NOT NULL COMMENT '节点网络代理ID',
	name varchar(50) NOT NULL COMMENT '名称',
	comment varchar(200) DEFAULT NULL COMMENT '说明',
	code varchar(50) NOT NULL COMMENT '标识码',
	innerServerId varchar(50) NOT NULL COMMENT '输入服务ID',
	innerType varchar(20) NOT NULL COMMENT '输入类型',
	innerAddress varchar(200) NULL NULL COMMENT '输入地址',
	outerServerId varchar(50) NOT NULL COMMENT '输出服务ID',
	outerType varchar(20) NOT NULL COMMENT '输出类型',
	outerAddress varchar(200) NULL NULL COMMENT '输出地址',
	lineServerIds varchar(2000) DEFAULT NULL COMMENT '输入输出节点线',
	option varchar(2000) DEFAULT NULL COMMENT '配置',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	enabled int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-启用、2-停用',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	deleteUserId bigint(20) DEFAULT NULL COMMENT '删除用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (netProxyId),
	KEY index_userId (userId),
	KEY index_name (name),
	KEY index_enabled (enabled),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableNodeNetProxyComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableNodeNetProxy + ` (
	netProxyId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	comment varchar(200) DEFAULT NULL,
	code varchar(50) NOT NULL,
	innerServerId varchar(50) NOT NULL,
	innerType varchar(20) NOT NULL,
	innerAddress varchar(200) NULL NULL,
	outerServerId varchar(50) NOT NULL,
	outerType varchar(20) NOT NULL,
	outerAddress varchar(200) NULL NULL,
	lineServerIds varchar(2000) DEFAULT NULL,
	option varchar(2000) DEFAULT NULL,
	userId bigint(20) DEFAULT NULL,
	enabled int(1) NOT NULL DEFAULT 1,
	deleted int(1) NOT NULL DEFAULT 2,
	deleteUserId bigint(20) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (netProxyId)
);
`,
					`CREATE INDEX ` + TableNodeNetProxy + `_index_userId on ` + TableNodeNetProxy + ` (userId);`,
					`CREATE INDEX ` + TableNodeNetProxy + `_index_name on ` + TableNodeNetProxy + ` (name);`,
					`CREATE INDEX ` + TableNodeNetProxy + `_index_enabled on ` + TableNode + ` (enabled);`,
					`CREATE INDEX ` + TableNodeNetProxy + `_index_deleted on ` + TableNodeNetProxy + ` (deleted);`,
				},
			},
		},
	}

}
