package userService

import "teamide/server/base"

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "user"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
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
	PRIMARY KEY (serverId, userId),
	KEY index_serverId_name (serverId, name),
	KEY index_serverId_account (serverId, account),
	KEY index_serverId_email (serverId, email),
	KEY index_serverId_activedState (serverId, activedState),
	KEY index_serverId_lockedState (serverId, lockedState),
	KEY index_serverId_enabledState (serverId, enabledState),
	KEY index_serverId_deletedState (serverId, deletedState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_METADATA",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER_METADATA (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	metadataId bigint(20) NOT NULL COMMENT '元数据ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	metadataStruct int(10) NOT NULL COMMENT '元数据结构',
	metadataField int(10) NOT NULL COMMENT '元数据字段',
	metadataValue text NOT NULL COMMENT '元数据值',
	parentId bigint(20) DEFAULT NULL COMMENT '父ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, metadataId),
	KEY index_serverId_userId (serverId, userId),
	KEY index_serverId_metadataStruct_metadataField (serverId, metadataStruct, metadataField),
	KEY index_serverId_userId_metadataStruct_metadataField (serverId, userId, metadataStruct, metadataField),
	KEY index_serverId_parentId (serverId, parentId),
	KEY index_serverId_userId_parentId (serverId, userId, parentId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户元数据';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_AUTH",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER_AUTH (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
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
	PRIMARY KEY (serverId, authId),
	KEY index_serverId_userId (serverId, userId),
	KEY index_serverId_authType_openId (serverId, authType, openId),
	KEY index_serverId_activedState (serverId, activedState),
	KEY index_serverId_lockedState (serverId, lockedState),
	KEY index_serverId_enabledState (serverId, enabledState),
	KEY index_serverId_deletedState (serverId, deletedState)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户第三方授权';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_PASSWORD",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER_PASSWORD (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	salt varchar(20) NOT NULL COMMENT '盐',
	password varchar(100) NOT NULL COMMENT '密码',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户密码';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_AUTH",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER_AUTH (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	authId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, authId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户授权';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_USER_LOCK",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_USER_LOCK (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	lockId bigint(20) NOT NULL COMMENT '授权ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, lockId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户授权';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
