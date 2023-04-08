package module_setting

import (
	"teamide/internal/install"
)

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建设置表
		{
			Version: "1.4",
			Module:  ModuleSetting,
			Stage:   `创建表[` + TableSetting + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableSetting + ` (
	name varchar(200) NOT NULL COMMENT '名称',
	value varchar(1000) NOT NULL COMMENT '值',
	createTime bigint(20) NOT NULL COMMENT '创建时间',
	updateTime bigint(20) DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableSettingComment + `';
`,
				},
				Sqlite: []string{`
CREATE TABLE ` + TableSetting + ` (
	name varchar(200) NOT NULL,
	value varchar(200) NOT NULL,
	createTime bigint(20) NOT NULL,
	updateTime bigint(20) DEFAULT NULL,
	PRIMARY KEY (name)
);
`,
				},
			},
		},
	}
}
