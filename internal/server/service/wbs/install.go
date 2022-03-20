package wbsService

import (
	"teamide/internal/server/base"
)

func (this_ *Service) GetInstall() (info *base.InstallInfo) {

	info = &base.InstallInfo{}

	info.Module = "wbs"
	stages := []*base.InstallStageInfo{}

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_PROJECT",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_WBS_PROJECT (
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
	PRIMARY KEY (projectId),
	KEY index_name (name),
	KEY index_projectNo (projectNo),
	KEY index_ownerId (ownerId),
	KEY index_ownerName (ownerName),
	KEY index_managerId (managerId),
	KEY index_managerName (managerName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='项目';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_VERSION",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_WBS_VERSION (
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
	PRIMARY KEY (versionId),
	KEY index_name (name),
	KEY index_projectId (projectId),
	KEY index_ownerId (ownerId),
	KEY index_ownerName (ownerName),
	KEY index_managerId (managerId),
	KEY index_managerName (managerName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本';
				`,
		},
	})

	stages = append(stages, &base.InstallStageInfo{
		Stage: "CREATE TABLE TM_WBS_TASK",
		Sql: &base.InstallSql{
			MySql: `
CREATE TABLE TM_WBS_TASK (
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
	PRIMARY KEY (taskId),
	KEY index_name (name),
	KEY index_projectId (projectId),
	KEY index_versionId (versionId),
	KEY index_projectId_versionId (projectId, versionId),
	KEY index_projectId_name (projectId, name),
	KEY index_projectId_versionId_name (projectId, versionId, name),
	KEY index_parentId (parentId),
	KEY index_ownerId (ownerId),
	KEY index_ownerName (ownerName),
	KEY index_managerId (managerId),
	KEY index_managerName (managerName),
	KEY index_executorId (executorId),
	KEY index_executorName (executorName)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务';
				`,
		},
	})

	info.Stages = stages

	return
}
