package userService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "user"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_USER (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	name varchar(50) NOT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	account varchar(20) NOT NULL COMMENT '账号',
	email varchar(50) DEFAULT NULL COMMENT '邮箱',
	activedState int(1) NOT NULL DEFAULT 2 COMMENT '激活状态:1-激活、2-未激活',
	lockedState int(1) NOT NULL DEFAULT 2 COMMENT '锁定状态:1-锁定、2-未锁定',
	enabledState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-启用、2-停用',
	deletedState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId),
	KEY index_name (name),
	KEY index_account (account),
	KEY index_email (email),
	KEY index_activedState (activedState),
	KEY index_lockedState (lockedState),
	KEY index_enabledState (enabledState),
	KEY index_deletedState (deletedState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_METADATA",
		Sql: &base.InstallSql{
			MySql: `
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
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_AUTH",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_USER_AUTH (
	authId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	authType int(2) NOT NULL COMMENT '授权类型',
	openId varchar(100) NOT NULL COMMENT 'OpenID',
	name varchar(100) DEFAULT NULL COMMENT '名称',
	avatar varchar(200) DEFAULT NULL COMMENT '头像',
	activedState int(1) NOT NULL DEFAULT 2 COMMENT '激活状态:1-激活、2-未激活',
	lockedState int(1) NOT NULL DEFAULT 2 COMMENT '锁定状态:1-锁定、2-未锁定',
	enabledState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-启用、2-停用',
	deletedState int(1) NOT NULL DEFAULT 1 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (authId),
	KEY index_userId (userId),
	KEY index_authType_openId (authType, openId),
	KEY index_activedState (activedState),
	KEY index_lockedState (lockedState),
	KEY index_enabledState (enabledState),
	KEY index_deletedState (deletedState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户第三方授权';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_PASSWORD",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_USER_PASSWORD (
	userId bigint(20) NOT NULL COMMENT '用户ID',
	salt varchar(20) NOT NULL COMMENT '盐',
	password varchar(100) NOT NULL COMMENT '密码',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户密码';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_AUTH",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_USER_AUTH (
	authId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (authId),
	KEY index_userId (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户授权';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_LOCK",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_USER_LOCK (
	lockId bigint(20) NOT NULL COMMENT '锁定ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (lockId),
	KEY index_userId (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户锁定';
				`,
		},
	})

	info.Stages = stages

	return
}
