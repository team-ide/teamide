package install

import "db"

func getUser() (info *InstallInfo) {

	info = &InstallInfo{}

	info.Module = "user"
	stages := []*InstallStageInfo{}

	stages = append(stages, &InstallStageInfo{
		Stage: "CREATE TABLE TM_USER",
		SqlParam: db.SqlParam{
			Sql: `
CREATE TABLE TM_USER (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	name varchar(50) NOT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	account varchar(20) NOT NULL COMMENT '账号',
	email varchar(50) DEFAULT NULL COMMENT '邮箱',
	phone varchar(20) DEFAULT NULL COMMENT '手机',
	activedState int(1) NOT NULL DEFAULT 2 COMMENT '激活状态：1-激活、2-未激活',
	lockedState int(1) NOT NULL DEFAULT 2 COMMENT '锁定状态：1-锁定、2-未锁定',
	enabledState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态：1-启用、2-停用',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId),
	KEY index_name (name),
	KEY index_account (account),
	KEY index_email (email),
	KEY index_phone (phone),
	KEY index_activedState (activedState),
	KEY index_lockedState (lockedState),
	KEY index_enabledState (enabledState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_METADATA",
		SqlParam: db.SqlParam{
			Sql: `
CREATE TABLE TM_USER_METADATA (
	metadataId bigint(20) NOT NULL COMMENT '元数据ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	metadataStruct int(10) NOT NULL COMMENT '元数据结构',
	metadataField int(10) NOT NULL COMMENT '元数据字段',
	metadataValue text NOT NULL COMMENT '元数据值',
	parentId bigint(20) DEFAULT NULL COMMENT '父ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (metadataId),
	KEY index_userId (userId),
	KEY index_metadataStruct_metadataField (metadataStruct, metadataField),
	KEY index_userId_metadataStruct_metadataField (userId, metadataStruct, metadataField),
	KEY index_parentId (parentId),
	KEY index_userId_parentId (userId, parentId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户元数据';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_AUTH",
		SqlParam: db.SqlParam{
			Sql: `
CREATE TABLE TM_USER_AUTH (
	authId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	authType int(2) NOT NULL COMMENT '授权类型',
	openId varchar(100) NOT NULL COMMENT 'OpenID',
	name varchar(100) DEFAULT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	activedState int(1) NOT NULL DEFAULT 2 COMMENT '激活状态：1-激活、2-未激活',
	lockedState int(1) NOT NULL DEFAULT 2 COMMENT '锁定状态：1-锁定、2-未锁定',
	enabledState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态：1-启用、2-停用',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (authId),
	KEY index_userId (userId),
	KEY index_authType_openId (authType, openId),
	KEY index_activedState (activedState),
	KEY index_lockedState (lockedState),
	KEY index_enabledState (enabledState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户第三方授权';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_PASSWORD",
		SqlParam: db.SqlParam{
			Sql: `
CREATE TABLE TM_USER_PASSWORD (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	salt varchar(20) NOT NULL COMMENT '盐',
	password varchar(100) NOT NULL COMMENT '密码',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户密码';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
