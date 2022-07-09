package module_id

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建 ID 表
		{
			Version: "1.0",
			Module:  ModuleID,
			Stage:   `创建表[` + TableID + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableID + ` (
	idType int(10) NOT NULL COMMENT '类型',
	value bigint(20) NOT NULL COMMENT '值',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (idType)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableIDComment + `';
`,
				},
				Sqlite: []string{`
CREATE TABLE ` + TableID + ` (
	idType int(10) NOT NULL,
	value bigint(20) NOT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (idType)
);
`,
				},
			},
		},
	}

}
