package wbsService

import "teamide/server/base"

func (this_ *WbsService) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "wbs"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_PROJECT",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_WBS_PROJECT (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	projectId bigint(20) NOT NULL COMMENT '项目ID',
	name varchar(50) NOT NULL COMMENT '名称',
	projectNo varchar(50) NOT NULL COMMENT '项目编号',
	ownerId bigint(20) NOT NULL COMMENT '所有者ID',
	ownerName varchar(50) NOT NULL COMMENT '所有者名称',
	managerId bigint(20) NOT NULL COMMENT '管理者ID',
	managerName varchar(50) NOT NULL COMMENT '管理者名称',
	description varchar(1000) DEFAULT NULL COMMENT '描述',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, projectId),
	KEY index_serverId_name (serverId, name),
	KEY index_serverId_projectNo (serverId, projectNo),
	KEY index_serverId_ownerId (serverId, ownerId),
	KEY index_serverId_ownerName (serverId, ownerName),
	KEY index_serverId_managerId (serverId, managerId),
	KEY index_serverId_managerName (serverId, managerName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='项目';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_VERSION",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_WBS_VERSION (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	versionId bigint(20) NOT NULL COMMENT '版本ID',
	name varchar(50) NOT NULL COMMENT '名称',
	projectId bigint(20) NOT NULL COMMENT '项目ID',
	ownerId bigint(20) NOT NULL COMMENT '所有者ID',
	ownerName varchar(50) NOT NULL COMMENT '所有者名称',
	managerId bigint(20) NOT NULL COMMENT '管理者ID',
	managerName varchar(50) NOT NULL COMMENT '管理者名称',
	description varchar(1000) DEFAULT NULL COMMENT '描述',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, versionId),
	KEY index_serverId_name (serverId, name),
	KEY index_serverId_projectId (serverId, projectId),
	KEY index_serverId_ownerId (serverId, ownerId),
	KEY index_serverId_ownerName (serverId, ownerName),
	KEY index_serverId_managerId (serverId, managerId),
	KEY index_serverId_managerName (serverId, managerName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本';
				`,
			Params: []interface{}{},
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_TASK",
		SqlParam: base.SqlParam{
			Sql: `
CREATE TABLE TM_WBS_TASK (
	serverId bigint(20) NOT NULL COMMENT '服务ID',
	taskId bigint(20) NOT NULL COMMENT '任务ID',
	name varchar(50) NOT NULL COMMENT '名称',
	projectId bigint(20) NOT NULL COMMENT '项目ID',
	versionId bigint(20) NOT NULL COMMENT '版本ID',
	parentId bigint(20) NOT NULL COMMENT '父任务ID',
	priority int(10) NOT NULL COMMENT '优先级',
	description varchar(1000) DEFAULT NULL COMMENT '描述',
	ownerId bigint(20) NOT NULL COMMENT '所有者ID',
	ownerName varchar(50) NOT NULL COMMENT '所有者名称',
	managerId bigint(20) NOT NULL COMMENT '管理者ID',
	managerName varchar(50) NOT NULL COMMENT '管理者名称',
	executorId bigint(20) NOT NULL COMMENT '执行者ID',
	executorName varchar(50) NOT NULL COMMENT '执行者名称',
	planStartTime datetime NOT NULL COMMENT '计划开始时间',
	planFinishTime datetime NOT NULL COMMENT '计划完成时间',
	actualStartTime datetime NOT NULL COMMENT '实际开始时间',
	actualFinishTime datetime NOT NULL COMMENT '实际完成时间',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (serverId, taskId),
	KEY index_serverId_name (serverId, name),
	KEY index_serverId_projectId (serverId, projectId),
	KEY index_serverId_versionId (serverId, versionId),
	KEY index_serverId_projectId_versionId (serverId, projectId, versionId),
	KEY index_serverId_parentId (serverId, parentId),
	KEY index_serverId_ownerId (serverId, ownerId),
	KEY index_serverId_ownerName (serverId, ownerName),
	KEY index_serverId_managerId (serverId, managerId),
	KEY index_serverId_managerName (serverId, managerName),
	KEY index_serverId_executorId (serverId, executorId),
	KEY index_serverId_executorName (serverId, executorName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务';
				`,
			Params: []interface{}{},
		},
	})

	info.Stages = stages

	return
}
