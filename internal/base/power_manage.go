package base

var (
	// 系统管理 权限

	// 管理用户 权限
	PowerManageUserPage    = addPower(&PowerAction{Action: "manage_user_page", Text: "管理用户页面", ShouldLogin: true, StandAlone: false})
	PowerManageUserInsert  = addPower(&PowerAction{Action: "manage_user_insert", Text: "管理用户新增", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserUpdate  = addPower(&PowerAction{Action: "manage_user_update", Text: "管理用户修改", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserDelete  = addPower(&PowerAction{Action: "manage_user_delete", Text: "管理用户删除", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserActive  = addPower(&PowerAction{Action: "manage_user_active", Text: "管理用户激活", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserLock    = addPower(&PowerAction{Action: "manage_user_lock", Text: "管理用户锁定", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserUnlock  = addPower(&PowerAction{Action: "manage_user_unlock", Text: "管理用户解锁", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserEnable  = addPower(&PowerAction{Action: "manage_user_enable", Text: "管理用户启用", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserDisable = addPower(&PowerAction{Action: "manage_user_disable", Text: "管理用户禁用", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})

	// 管理用户密码 权限
	PowerManageUserPasswordReset  = addPower(&PowerAction{Action: "manage_user_password_reset", Text: "管理用户密码重置", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserPasswordUpdate = addPower(&PowerAction{Action: "manage_user_password_update", Text: "管理用户密码修改", Parent: PowerManageUserPage, ShouldLogin: true, StandAlone: false})

	// 管理锁定用户 权限
	PowerManageUserLockPage   = addPower(&PowerAction{Action: "manage_user_lock_page", Text: "管理锁定用户页面", ShouldLogin: true, StandAlone: false})
	PowerManageUserLockInsert = addPower(&PowerAction{Action: "manage_user_lock_insert", Text: "管理锁定用户新增", Parent: PowerManageUserLockPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserLockUpdate = addPower(&PowerAction{Action: "manage_user_lock_update", Text: "管理锁定用户修改", Parent: PowerManageUserLockPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserLockDelete = addPower(&PowerAction{Action: "manage_user_lock_delete", Text: "管理锁定用户删除", Parent: PowerManageUserLockPage, ShouldLogin: true, StandAlone: false})

	// 管理用户授权 权限
	PowerManageUserAuthPage    = addPower(&PowerAction{Action: "manage_user_auth_page", Text: "管理用户授权页面", ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthInsert  = addPower(&PowerAction{Action: "manage_user_auth_insert", Text: "管理用户授权新增", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthUpdate  = addPower(&PowerAction{Action: "manage_user_auth_update", Text: "管理用户授权修改", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthDelete  = addPower(&PowerAction{Action: "manage_user_auth_delete", Text: "管理用户授权删除", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthActive  = addPower(&PowerAction{Action: "manage_user_auth_active", Text: "管理用户授权激活", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthLock    = addPower(&PowerAction{Action: "manage_user_auth_lock", Text: "管理用户授权锁定", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthUnlock  = addPower(&PowerAction{Action: "manage_user_auth_unlock", Text: "管理用户授权解锁", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthEnable  = addPower(&PowerAction{Action: "manage_user_auth_enable", Text: "管理用户授权启用", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerManageUserAuthDisable = addPower(&PowerAction{Action: "manage_user_auth_disable", Text: "管理用户授权禁用", Parent: PowerManageUserAuthPage, ShouldLogin: true, StandAlone: false})

	// 管理凭证 权限
	PowerManageCertificatePage   = addPower(&PowerAction{Action: "manage_certificate_page", Text: "管理凭证页面", ShouldLogin: true, StandAlone: false})
	PowerManageCertificateInsert = addPower(&PowerAction{Action: "manage_certificate_insert", Text: "管理凭证新增", Parent: PowerManageCertificatePage, ShouldLogin: true, StandAlone: false})
	PowerManageCertificateUpdate = addPower(&PowerAction{Action: "manage_certificate_update", Text: "管理凭证修改", Parent: PowerManageCertificatePage, ShouldLogin: true, StandAlone: false})
	PowerManageCertificateDelete = addPower(&PowerAction{Action: "manage_certificate_delete", Text: "管理凭证删除", Parent: PowerManageCertificatePage, ShouldLogin: true, StandAlone: false})

	// 管理角色权限 权限
	PowerManagePowerRolePage   = addPower(&PowerAction{Action: "manage_power_role_page", Text: "管理权限角色页面", ShouldLogin: true, StandAlone: false})
	PowerManagePowerRoleInsert = addPower(&PowerAction{Action: "manage_power_role_insert", Text: "管理权限角色新增", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerRoleUpdate = addPower(&PowerAction{Action: "manage_power_role_update", Text: "管理权限角色修改", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerRoleDelete = addPower(&PowerAction{Action: "manage_power_role_delete", Text: "管理权限角色删除", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})

	// 管理用户权限 权限
	PowerManagePowerUserPage   = addPower(&PowerAction{Action: "manage_power_user_page", Text: "管理用户权限页面", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerUserInsert = addPower(&PowerAction{Action: "manage_power_user_insert", Text: "管理用户权限新增", Parent: PowerManagePowerUserPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerUserUpdate = addPower(&PowerAction{Action: "manage_power_user_update", Text: "管理用户权限修改", Parent: PowerManagePowerUserPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerUserDelete = addPower(&PowerAction{Action: "manage_power_user_delete", Text: "管理用户权限删除", Parent: PowerManagePowerUserPage, ShouldLogin: true, StandAlone: false})

	// 管理功能权限 权限
	PowerManagePowerActionPage   = addPower(&PowerAction{Action: "manage_power_action_page", Text: "管理功能权限页面", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerActionInsert = addPower(&PowerAction{Action: "manage_power_action_insert", Text: "管理功能权限新增", Parent: PowerManagePowerActionPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerActionUpdate = addPower(&PowerAction{Action: "manage_power_action_update", Text: "管理功能权限修改", Parent: PowerManagePowerActionPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerActionDelete = addPower(&PowerAction{Action: "manage_power_action_delete", Text: "管理功能权限删除", Parent: PowerManagePowerActionPage, ShouldLogin: true, StandAlone: false})

	// 管理数据权限 权限
	PowerManagePowerDataPage   = addPower(&PowerAction{Action: "manage_power_data_page", Text: "管理数据权限页面", Parent: PowerManagePowerRolePage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerDataInsert = addPower(&PowerAction{Action: "manage_power_data_insert", Text: "管理数据权限新增", Parent: PowerManagePowerDataPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerDataUpdate = addPower(&PowerAction{Action: "manage_power_data_update", Text: "管理数据权限修改", Parent: PowerManagePowerDataPage, ShouldLogin: true, StandAlone: false})
	PowerManagePowerDataDelete = addPower(&PowerAction{Action: "manage_power_data_delete", Text: "管理数据权限删除", Parent: PowerManagePowerDataPage, ShouldLogin: true, StandAlone: false})

	// 管理企业 权限
	PowerManageEnterprisePage   = addPower(&PowerAction{Action: "manage_enterprise_page", Text: "管理企业页面", ShouldLogin: true, StandAlone: false})
	PowerManageEnterpriseInsert = addPower(&PowerAction{Action: "manage_enterprise_insert", Text: "管理企业新增", Parent: PowerManageEnterprisePage, ShouldLogin: true, StandAlone: false})
	PowerManageEnterpriseUpdate = addPower(&PowerAction{Action: "manage_enterprise_update", Text: "管理企业修改", Parent: PowerManageEnterprisePage, ShouldLogin: true, StandAlone: false})
	PowerManageEnterpriseDelete = addPower(&PowerAction{Action: "manage_enterprise_delete", Text: "管理企业删除", Parent: PowerManageEnterprisePage, ShouldLogin: true, StandAlone: false})

	// 管理组织结构 权限
	PowerManageOrganizationPage   = addPower(&PowerAction{Action: "manage_organization_page", Text: "管理组织结构页面", ShouldLogin: true, StandAlone: false})
	PowerManageOrganizationInsert = addPower(&PowerAction{Action: "manage_organization_insert", Text: "管理组织结构新增", Parent: PowerManageOrganizationPage, ShouldLogin: true, StandAlone: false})
	PowerManageOrganizationUpdate = addPower(&PowerAction{Action: "manage_organization_update", Text: "管理组织结构修改", Parent: PowerManageOrganizationPage, ShouldLogin: true, StandAlone: false})
	PowerManageOrganizationDelete = addPower(&PowerAction{Action: "manage_organization_delete", Text: "管理组织结构删除", Parent: PowerManageOrganizationPage, ShouldLogin: true, StandAlone: false})

	// 管理群组 权限
	PowerManageGroupPage   = addPower(&PowerAction{Action: "manage_group_page", Text: "管理群组页面", ShouldLogin: true, StandAlone: false})
	PowerManageGroupInsert = addPower(&PowerAction{Action: "manage_group_insert", Text: "管理群组新增", Parent: PowerManageGroupPage, ShouldLogin: true, StandAlone: false})
	PowerManageGroupUpdate = addPower(&PowerAction{Action: "manage_group_update", Text: "管理群组修改", Parent: PowerManageGroupPage, ShouldLogin: true, StandAlone: false})
	PowerManageGroupDelete = addPower(&PowerAction{Action: "manage_group_delete", Text: "管理群组删除", Parent: PowerManageGroupPage, ShouldLogin: true, StandAlone: false})

	// 管理任务 权限
	PowerManageJobPage   = addPower(&PowerAction{Action: "manage_job_page", Text: "管理任务页面", ShouldLogin: true, StandAlone: false})
	PowerManageJobInsert = addPower(&PowerAction{Action: "manage_job_insert", Text: "管理任务新增", Parent: PowerManageJobPage, ShouldLogin: true, StandAlone: false})
	PowerManageJobUpdate = addPower(&PowerAction{Action: "manage_job_update", Text: "管理任务修改", Parent: PowerManageJobPage, ShouldLogin: true, StandAlone: false})
	PowerManageJobDelete = addPower(&PowerAction{Action: "manage_job_delete", Text: "管理任务删除", Parent: PowerManageJobPage, ShouldLogin: true, StandAlone: false})

	// 管理消息 权限
	PowerManageMessagePage   = addPower(&PowerAction{Action: "manage_message_page", Text: "管理消息页面", ShouldLogin: true, StandAlone: false})
	PowerManageMessageInsert = addPower(&PowerAction{Action: "manage_message_insert", Text: "管理消息新增", Parent: PowerManageMessagePage, ShouldLogin: true, StandAlone: false})
	PowerManageMessageUpdate = addPower(&PowerAction{Action: "manage_message_update", Text: "管理消息修改", Parent: PowerManageMessagePage, ShouldLogin: true, StandAlone: false})
	PowerManageMessageDelete = addPower(&PowerAction{Action: "manage_message_delete", Text: "管理消息删除", Parent: PowerManageMessagePage, ShouldLogin: true, StandAlone: false})

	// 管理空间 权限
	PowerManageSpacePage   = addPower(&PowerAction{Action: "manage_space_page", Text: "管理空间页面", ShouldLogin: true, StandAlone: false})
	PowerManageSpaceInsert = addPower(&PowerAction{Action: "manage_space_insert", Text: "管理空间新增", Parent: PowerManageSpacePage, ShouldLogin: true, StandAlone: false})
	PowerManageSpaceUpdate = addPower(&PowerAction{Action: "manage_space_update", Text: "管理空间修改", Parent: PowerManageSpacePage, ShouldLogin: true, StandAlone: false})
	PowerManageSpaceDelete = addPower(&PowerAction{Action: "manage_space_delete", Text: "管理空间删除", Parent: PowerManageSpacePage, ShouldLogin: true, StandAlone: false})

	// 管理日志 权限
	PowerManageLogPage   = addPower(&PowerAction{Action: "manage_log_page", Text: "管理日志页面", ShouldLogin: true, StandAlone: false})
	PowerManageLogInsert = addPower(&PowerAction{Action: "manage_log_insert", Text: "管理日志新增", Parent: PowerManageLogPage, ShouldLogin: true, StandAlone: false})
	PowerManageLogUpdate = addPower(&PowerAction{Action: "manage_log_update", Text: "管理日志修改", Parent: PowerManageLogPage, ShouldLogin: true, StandAlone: false})
	PowerManageLogDelete = addPower(&PowerAction{Action: "manage_log_delete", Text: "管理日志删除", Parent: PowerManageLogPage, ShouldLogin: true, StandAlone: false})

	// 管理登录 权限
	PowerManageLoginPage   = addPower(&PowerAction{Action: "manage_login_page", Text: "管理登录页面", ShouldLogin: true, StandAlone: false})
	PowerManageLoginInsert = addPower(&PowerAction{Action: "manage_login_insert", Text: "管理登录新增", Parent: PowerManageLoginPage, ShouldLogin: true, StandAlone: false})
	PowerManageLoginUpdate = addPower(&PowerAction{Action: "manage_login_update", Text: "管理登录修改", Parent: PowerManageLoginPage, ShouldLogin: true, StandAlone: false})
	PowerManageLoginDelete = addPower(&PowerAction{Action: "manage_login_delete", Text: "管理登录删除", Parent: PowerManageLoginPage, ShouldLogin: true, StandAlone: false})

	// 管理系统设置 权限
	PowerManageSystemSettingPage   = addPower(&PowerAction{Action: "manage_system_setting_page", Text: "管理系统设置页面", ShouldLogin: true, StandAlone: false})
	PowerManageSystemSettingUpdate = addPower(&PowerAction{Action: "manage_system_setting_update", Text: "管理系统设置修改", Parent: PowerManageSystemSettingPage, ShouldLogin: true, StandAlone: false})

	//管理系统日志 权限
	PowerManageSystemLogPage   = addPower(&PowerAction{Action: "manage_system_log_page", Text: "管理日志页面", ShouldLogin: true, StandAlone: false})
	PowerManageSystemLogInsert = addPower(&PowerAction{Action: "manage_system_log_insert", Text: "管理日志新增", Parent: PowerManageSystemLogPage, ShouldLogin: true, StandAlone: false})
	PowerManageSystemLogUpdate = addPower(&PowerAction{Action: "manage_system_log_update", Text: "管理日志修改", Parent: PowerManageSystemLogPage, ShouldLogin: true, StandAlone: false})
	PowerManageSystemLogDelete = addPower(&PowerAction{Action: "manage_system_log_delete", Text: "管理日志删除", Parent: PowerManageSystemLogPage, ShouldLogin: true, StandAlone: false})
)
