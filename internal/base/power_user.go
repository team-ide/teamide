package base

var (
	// 用户 权限

	// 用户资料 权限
	PowerUserProfilePage   = addPower(&PowerAction{Action: "user_profile_page", Text: "用户资料页面", ShouldLogin: true, StandAlone: false})
	PowerUserProfileUpdate = addPower(&PowerAction{Action: "user_profile_update", Text: "用户资料修改", Parent: PowerUserProfilePage, ShouldLogin: true, StandAlone: false})

	// 用户密码 权限
	PowerUserPasswordPage   = addPower(&PowerAction{Action: "user_password_page", Text: "用户密码页面", ShouldLogin: true, StandAlone: false})
	PowerUserPasswordUpdate = addPower(&PowerAction{Action: "user_password_update", Text: "用户密码修改", Parent: PowerUserPasswordPage, ShouldLogin: true, StandAlone: false})

	// 用户授权 权限
	PowerUserAuthPage    = addPower(&PowerAction{Action: "user_auth_page", Text: "用户授权页面", ShouldLogin: true, StandAlone: false})
	PowerUserAuthInsert  = addPower(&PowerAction{Action: "user_auth_insert", Text: "用户授权新增", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthUpdate  = addPower(&PowerAction{Action: "user_auth_update", Text: "用户授权修改", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthDelete  = addPower(&PowerAction{Action: "user_auth_delete", Text: "用户授权删除", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthActive  = addPower(&PowerAction{Action: "user_auth_active", Text: "用户授权激活", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthLock    = addPower(&PowerAction{Action: "user_auth_lock", Text: "用户授权锁定", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthUnlock  = addPower(&PowerAction{Action: "user_auth_unlock", Text: "用户授权解锁", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthEnable  = addPower(&PowerAction{Action: "user_auth_enable", Text: "用户授权启用", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})
	PowerUserAuthDisable = addPower(&PowerAction{Action: "user_auth_disable", Text: "用户授权禁用", Parent: PowerUserAuthPage, ShouldLogin: true, StandAlone: false})

	// 用户凭证 权限
	PowerUserCertificatePage   = addPower(&PowerAction{Action: "user_certificate_page", Text: "用户凭证页面", ShouldLogin: true, StandAlone: false})
	PowerUserCertificateInsert = addPower(&PowerAction{Action: "user_certificate_insert", Text: "用户凭证新增", Parent: PowerUserCertificatePage, ShouldLogin: true, StandAlone: false})
	PowerUserCertificateUpdate = addPower(&PowerAction{Action: "user_certificate_update", Text: "用户凭证修改", Parent: PowerUserCertificatePage, ShouldLogin: true, StandAlone: false})
	PowerUserCertificateDelete = addPower(&PowerAction{Action: "user_certificate_delete", Text: "用户凭证删除", Parent: PowerUserCertificatePage, ShouldLogin: true, StandAlone: false})

	// 用户企业 权限
	PowerUserEnterprisePage   = addPower(&PowerAction{Action: "user_enterprise_page", Text: "用户企业页面", ShouldLogin: true, StandAlone: false})
	PowerUserEnterpriseInsert = addPower(&PowerAction{Action: "user_enterprise_insert", Text: "用户企业新增", Parent: PowerUserEnterprisePage, ShouldLogin: true, StandAlone: false})
	PowerUserEnterpriseUpdate = addPower(&PowerAction{Action: "user_enterprise_update", Text: "用户企业修改", Parent: PowerUserEnterprisePage, ShouldLogin: true, StandAlone: false})
	PowerUserEnterpriseDelete = addPower(&PowerAction{Action: "user_enterprise_delete", Text: "用户企业删除", Parent: PowerUserEnterprisePage, ShouldLogin: true, StandAlone: false})

	// 用户组织结构 权限
	PowerUserOrganizationPage   = addPower(&PowerAction{Action: "user_organization_page", Text: "用户组织结构页面", ShouldLogin: true, StandAlone: false})
	PowerUserOrganizationInsert = addPower(&PowerAction{Action: "user_organization_insert", Text: "用户组织结构新增", Parent: PowerUserOrganizationPage, ShouldLogin: true, StandAlone: false})
	PowerUserOrganizationUpdate = addPower(&PowerAction{Action: "user_organization_update", Text: "用户组织结构修改", Parent: PowerUserOrganizationPage, ShouldLogin: true, StandAlone: false})
	PowerUserOrganizationDelete = addPower(&PowerAction{Action: "user_organization_delete", Text: "用户组织结构删除", Parent: PowerUserOrganizationPage, ShouldLogin: true, StandAlone: false})

	// 用户群组 权限
	PowerUserGroupPage   = addPower(&PowerAction{Action: "user_group_page", Text: "用户群组页面", ShouldLogin: true, StandAlone: false})
	PowerUserGroupInsert = addPower(&PowerAction{Action: "user_group_insert", Text: "用户群组新增", Parent: PowerUserGroupPage, ShouldLogin: true, StandAlone: false})
	PowerUserGroupUpdate = addPower(&PowerAction{Action: "user_group_update", Text: "用户群组修改", Parent: PowerUserGroupPage, ShouldLogin: true, StandAlone: false})
	PowerUserGroupDelete = addPower(&PowerAction{Action: "user_group_delete", Text: "用户群组删除", Parent: PowerUserGroupPage, ShouldLogin: true, StandAlone: false})

	// 用户任务 权限
	PowerUserJobPage   = addPower(&PowerAction{Action: "user_job_page", Text: "用户任务页面", ShouldLogin: true, StandAlone: false})
	PowerUserJobInsert = addPower(&PowerAction{Action: "user_job_insert", Text: "用户任务新增", Parent: PowerUserJobPage, ShouldLogin: true, StandAlone: false})
	PowerUserJobUpdate = addPower(&PowerAction{Action: "user_job_update", Text: "用户任务修改", Parent: PowerUserJobPage, ShouldLogin: true, StandAlone: false})
	PowerUserJobDelete = addPower(&PowerAction{Action: "user_job_delete", Text: "用户任务删除", Parent: PowerUserJobPage, ShouldLogin: true, StandAlone: false})

	// 用户空间 权限
	PowerUserSpacePage   = addPower(&PowerAction{Action: "user_space_page", Text: "用户空间页面", ShouldLogin: true, StandAlone: false})
	PowerUserSpaceInsert = addPower(&PowerAction{Action: "user_space_insert", Text: "用户空间新增", Parent: PowerUserSpacePage, ShouldLogin: true, StandAlone: false})
	PowerUserSpaceUpdate = addPower(&PowerAction{Action: "user_space_update", Text: "用户空间修改", Parent: PowerUserSpacePage, ShouldLogin: true, StandAlone: false})
	PowerUserSpaceDelete = addPower(&PowerAction{Action: "user_space_delete", Text: "用户空间删除", Parent: PowerUserSpacePage, ShouldLogin: true, StandAlone: false})

	// 用户日志 权限
	PowerUserLogPage   = addPower(&PowerAction{Action: "user_log_page", Text: "用户日志页面", ShouldLogin: true, StandAlone: false})
	PowerUserLogInsert = addPower(&PowerAction{Action: "user_log_insert", Text: "用户日志新增", Parent: PowerUserLogPage, ShouldLogin: true, StandAlone: false})
	PowerUserLogUpdate = addPower(&PowerAction{Action: "user_log_update", Text: "用户日志修改", Parent: PowerUserLogPage, ShouldLogin: true, StandAlone: false})
	PowerUserLogDelete = addPower(&PowerAction{Action: "user_log_delete", Text: "用户日志删除", Parent: PowerUserLogPage, ShouldLogin: true, StandAlone: false})

	// 用户登录 权限
	PowerUserLoginPage   = addPower(&PowerAction{Action: "user_login_page", Text: "用户登录页面", ShouldLogin: true, StandAlone: false})
	PowerUserLoginInsert = addPower(&PowerAction{Action: "user_login_insert", Text: "用户登录新增", Parent: PowerUserLoginPage, ShouldLogin: true, StandAlone: false})
	PowerUserLoginUpdate = addPower(&PowerAction{Action: "user_login_update", Text: "用户登录修改", Parent: PowerUserLoginPage, ShouldLogin: true, StandAlone: false})
	PowerUserLoginDelete = addPower(&PowerAction{Action: "user_login_delete", Text: "用户登录删除", Parent: PowerUserLoginPage, ShouldLogin: true, StandAlone: false})

	// 用户消息 权限
	PowerUserMessagePage   = addPower(&PowerAction{Action: "user_message_page", Text: "用户消息页面", ShouldLogin: true, StandAlone: false})
	PowerUserMessageInsert = addPower(&PowerAction{Action: "user_message_insert", Text: "用户消息新增", Parent: PowerUserMessagePage, ShouldLogin: true, StandAlone: false})
	PowerUserMessageUpdate = addPower(&PowerAction{Action: "user_message_update", Text: "用户消息修改", Parent: PowerUserMessagePage, ShouldLogin: true, StandAlone: false})
	PowerUserMessageDelete = addPower(&PowerAction{Action: "user_message_delete", Text: "用户消息删除", Parent: PowerUserMessagePage, ShouldLogin: true, StandAlone: false})

	// 用户设置 权限
	PowerUserSettingPage   = addPower(&PowerAction{Action: "user_setting_page", Text: "用户设置页面", ShouldLogin: true, StandAlone: false})
	PowerUserSettingUpdate = addPower(&PowerAction{Action: "user_setting_update", Text: "用户设置修改", Parent: PowerUserProfilePage, ShouldLogin: true, StandAlone: false})
)
