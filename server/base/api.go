package base

import "github.com/gin-gonic/gin"

type RequestBean struct {
	JWT  *JWTBean
	Path string
}

type PageBean struct {
	PageIndex int64
	PageSize  int64
	Total     int64
	TotalPage int64
	Value     interface{}
}

func (page *PageBean) Init() {
	page.TotalPage = (page.Total + page.PageSize - 1) / page.PageSize
}

type JWTBean struct {
	Sign     string `json:"sign"`
	ServerId int64  `json:"serverId"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Time     int64  `json:"time"`
}

type ApiWorker struct {
	Apis    []string
	Power   *PowerAction
	Do      func(request *RequestBean, c *gin.Context) (res interface{}, err error)
	DoOther func(request *RequestBean, c *gin.Context)
}

type PowerAction struct {
	Action      string `json:"action"`
	Text        string `json:"text"`
	ShouldLogin bool   `json:"shouldLogin"`
	Parent      *PowerAction
}

var (
	powers []*PowerAction

	// 基础权限
	PowerRegister  = addPower(&PowerAction{Action: "register", Text: "注册"})
	PowerData      = addPower(&PowerAction{Action: "data", Text: "数据"})
	PowerSession   = addPower(&PowerAction{Action: "session", Text: "会话"})
	PowerLogin     = addPower(&PowerAction{Action: "login", Text: "登录"})
	PowerLogout    = addPower(&PowerAction{Action: "logout", Text: "登出"})
	PowerAutoLogin = addPower(&PowerAction{Action: "auto_login", Text: "自动登录"})

	// 用户 权限

	// 用户资料 权限
	PowerUserProfilePage   = addPower(&PowerAction{Action: "user_profile_page", Text: "用户资料页面", ShouldLogin: true})
	PowerUserProfileUpdate = addPower(&PowerAction{Action: "user_profile_update", Text: "用户资料修改", Parent: PowerUserProfilePage, ShouldLogin: true})
	// 用户密码 权限
	PowerUserPasswordPage   = addPower(&PowerAction{Action: "user_password_page", Text: "用户密码页面", ShouldLogin: true})
	PowerUserPasswordUpdate = addPower(&PowerAction{Action: "user_password_update", Text: "用户密码修改", Parent: PowerUserPasswordPage, ShouldLogin: true})

	// 用户授权 权限
	PowerUserAuthPage    = addPower(&PowerAction{Action: "user_auth_page", Text: "用户授权页面", ShouldLogin: true})
	PowerUserAuthInsert  = addPower(&PowerAction{Action: "user_auth_insert", Text: "用户授权新增", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthUpdate  = addPower(&PowerAction{Action: "user_auth_update", Text: "用户授权修改", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthDelete  = addPower(&PowerAction{Action: "user_auth_delete", Text: "用户授权删除", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthActive  = addPower(&PowerAction{Action: "user_auth_active", Text: "用户授权激活", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthLock    = addPower(&PowerAction{Action: "user_auth_lock", Text: "用户授权锁定", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthUnlock  = addPower(&PowerAction{Action: "user_auth_unlock", Text: "用户授权解锁", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthEnable  = addPower(&PowerAction{Action: "user_auth_enable", Text: "用户授权启用", Parent: PowerUserAuthPage, ShouldLogin: true})
	PowerUserAuthDisable = addPower(&PowerAction{Action: "user_auth_disable", Text: "用户授权禁用", Parent: PowerUserAuthPage, ShouldLogin: true})

	// 用户凭证 权限
	PowerUserCertificatePage   = addPower(&PowerAction{Action: "user_certificate_page", Text: "用户凭证页面", ShouldLogin: true})
	PowerUserCertificateInsert = addPower(&PowerAction{Action: "user_certificate_insert", Text: "用户凭证新增", Parent: PowerUserCertificatePage, ShouldLogin: true})
	PowerUserCertificateUpdate = addPower(&PowerAction{Action: "user_certificate_update", Text: "用户凭证修改", Parent: PowerUserCertificatePage, ShouldLogin: true})
	PowerUserCertificateDelete = addPower(&PowerAction{Action: "user_certificate_delete", Text: "用户凭证删除", Parent: PowerUserCertificatePage, ShouldLogin: true})

	// 用户消息 权限
	PowerUserMessagePage   = addPower(&PowerAction{Action: "user_message_page", Text: "用户消息页面", ShouldLogin: true})
	PowerUserMessageInsert = addPower(&PowerAction{Action: "user_message_insert", Text: "用户消息新增", Parent: PowerUserMessagePage, ShouldLogin: true})
	PowerUserMessageUpdate = addPower(&PowerAction{Action: "user_message_update", Text: "用户消息修改", Parent: PowerUserMessagePage, ShouldLogin: true})
	PowerUserMessageDelete = addPower(&PowerAction{Action: "user_message_delete", Text: "用户消息删除", Parent: PowerUserMessagePage, ShouldLogin: true})

	// 用户设置 权限
	PowerUserSettingPage   = addPower(&PowerAction{Action: "user_setting_page", Text: "用户设置页面", ShouldLogin: true})
	PowerUserSettingUpdate = addPower(&PowerAction{Action: "user_setting_update", Text: "用户设置修改", Parent: PowerUserProfilePage, ShouldLogin: true})

	// 系统管理 权限

	// 管理用户 权限
	PowerManageUserPage    = addPower(&PowerAction{Action: "manage_user_page", Text: "管理用户页面", ShouldLogin: true})
	PowerManageUserInsert  = addPower(&PowerAction{Action: "manage_user_insert", Text: "管理用户新增", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserUpdate  = addPower(&PowerAction{Action: "manage_user_update", Text: "管理用户修改", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserDelete  = addPower(&PowerAction{Action: "manage_user_delete", Text: "管理用户删除", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserActive  = addPower(&PowerAction{Action: "manage_user_active", Text: "管理用户激活", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserLock    = addPower(&PowerAction{Action: "manage_user_lock", Text: "管理用户锁定", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserUnlock  = addPower(&PowerAction{Action: "manage_user_unlock", Text: "管理用户解锁", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserEnable  = addPower(&PowerAction{Action: "manage_user_enable", Text: "管理用户启用", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserDisable = addPower(&PowerAction{Action: "manage_user_disable", Text: "管理用户禁用", Parent: PowerManageUserPage, ShouldLogin: true})
	// 管理用户密码 权限
	PowerManageUserPasswordReset  = addPower(&PowerAction{Action: "manage_user_password_reset", Text: "管理用户密码重置", Parent: PowerManageUserPage, ShouldLogin: true})
	PowerManageUserPasswordUpdate = addPower(&PowerAction{Action: "manage_user_password_update", Text: "管理用户密码修改", Parent: PowerManageUserPage, ShouldLogin: true})

	// 管理锁定用户 权限
	PowerManageUserLockPage   = addPower(&PowerAction{Action: "manage_user_lock_page", Text: "管理锁定用户页面", ShouldLogin: true})
	PowerManageUserLockInsert = addPower(&PowerAction{Action: "manage_user_lock_insert", Text: "管理锁定用户新增", Parent: PowerManageUserLockPage, ShouldLogin: true})
	PowerManageUserLockUpdate = addPower(&PowerAction{Action: "manage_user_lock_update", Text: "管理锁定用户修改", Parent: PowerManageUserLockPage, ShouldLogin: true})
	PowerManageUserLockDelete = addPower(&PowerAction{Action: "manage_user_lock_delete", Text: "管理锁定用户删除", Parent: PowerManageUserLockPage, ShouldLogin: true})

	// 管理用户授权 权限
	PowerManageUserAuthPage    = addPower(&PowerAction{Action: "manage_user_auth_page", Text: "用户授权页面", ShouldLogin: true})
	PowerManageUserAuthInsert  = addPower(&PowerAction{Action: "manage_user_auth_insert", Text: "用户授权新增", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthUpdate  = addPower(&PowerAction{Action: "manage_user_auth_update", Text: "用户授权修改", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthDelete  = addPower(&PowerAction{Action: "manage_user_auth_delete", Text: "用户授权删除", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthActive  = addPower(&PowerAction{Action: "manage_user_auth_active", Text: "用户授权激活", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthLock    = addPower(&PowerAction{Action: "manage_user_auth_lock", Text: "用户授权锁定", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthUnlock  = addPower(&PowerAction{Action: "manage_user_auth_unlock", Text: "用户授权解锁", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthEnable  = addPower(&PowerAction{Action: "manage_user_auth_enable", Text: "用户授权启用", Parent: PowerManageUserAuthPage, ShouldLogin: true})
	PowerManageUserAuthDisable = addPower(&PowerAction{Action: "manage_user_auth_disable", Text: "用户授权禁用", Parent: PowerManageUserAuthPage, ShouldLogin: true})

	// 管理角色权限 权限
	PowerManagePowerRolePage   = addPower(&PowerAction{Action: "manage_power_role_page", Text: "管理权限角色页面", ShouldLogin: true})
	PowerManagePowerRoleInsert = addPower(&PowerAction{Action: "manage_power_role_insert", Text: "管理权限角色新增", Parent: PowerManagePowerRolePage, ShouldLogin: true})
	PowerManagePowerRoleUpdate = addPower(&PowerAction{Action: "manage_power_role_update", Text: "管理权限角色修改", Parent: PowerManagePowerRolePage, ShouldLogin: true})
	PowerManagePowerRoleDelete = addPower(&PowerAction{Action: "manage_power_role_delete", Text: "管理权限角色删除", Parent: PowerManagePowerRolePage, ShouldLogin: true})
	// 管理功能权限 权限
	PowerManagePowerActionPage   = addPower(&PowerAction{Action: "manage_power_action_page", Text: "管理功能权限页面", Parent: PowerManagePowerRolePage, ShouldLogin: true})
	PowerManagePowerActionInsert = addPower(&PowerAction{Action: "manage_power_action_insert", Text: "管理功能权限新增", Parent: PowerManagePowerActionPage, ShouldLogin: true})
	PowerManagePowerActionUpdate = addPower(&PowerAction{Action: "manage_power_action_update", Text: "管理功能权限修改", Parent: PowerManagePowerActionPage, ShouldLogin: true})
	PowerManagePowerActionDelete = addPower(&PowerAction{Action: "manage_power_action_delete", Text: "管理功能权限删除", Parent: PowerManagePowerActionPage, ShouldLogin: true})
	// 管理数据权限 权限
	PowerManagePowerDataPage   = addPower(&PowerAction{Action: "manage_power_data_page", Text: "管理数据权限页面", Parent: PowerManagePowerRolePage, ShouldLogin: true})
	PowerManagePowerDataInsert = addPower(&PowerAction{Action: "manage_power_data_insert", Text: "管理数据权限新增", Parent: PowerManagePowerDataPage, ShouldLogin: true})
	PowerManagePowerDataUpdate = addPower(&PowerAction{Action: "manage_power_data_update", Text: "管理数据权限修改", Parent: PowerManagePowerDataPage, ShouldLogin: true})
	PowerManagePowerDataDelete = addPower(&PowerAction{Action: "manage_power_data_delete", Text: "管理数据权限删除", Parent: PowerManagePowerDataPage, ShouldLogin: true})

	// 管理企业 权限
	PowerManageEnterprisePage   = addPower(&PowerAction{Action: "manage_enterprise_page", Text: "管理企业页面", ShouldLogin: true})
	PowerManageEnterpriseInsert = addPower(&PowerAction{Action: "manage_enterprise_insert", Text: "管理企业新增", Parent: PowerManageEnterprisePage, ShouldLogin: true})
	PowerManageEnterpriseUpdate = addPower(&PowerAction{Action: "manage_enterprise_update", Text: "管理企业修改", Parent: PowerManageEnterprisePage, ShouldLogin: true})
	PowerManageEnterpriseDelete = addPower(&PowerAction{Action: "manage_enterprise_delete", Text: "管理企业删除", Parent: PowerManageEnterprisePage, ShouldLogin: true})

	// 管理组织结构 权限
	PowerManageOrganizationPage   = addPower(&PowerAction{Action: "manage_organization_page", Text: "管理组织结构页面", ShouldLogin: true})
	PowerManageOrganizationInsert = addPower(&PowerAction{Action: "manage_organization_insert", Text: "管理组织结构新增", Parent: PowerManageOrganizationPage, ShouldLogin: true})
	PowerManageOrganizationUpdate = addPower(&PowerAction{Action: "manage_organization_update", Text: "管理组织结构修改", Parent: PowerManageOrganizationPage, ShouldLogin: true})
	PowerManageOrganizationDelete = addPower(&PowerAction{Action: "manage_organization_delete", Text: "管理组织结构删除", Parent: PowerManageOrganizationPage, ShouldLogin: true})

	// 管理群组 权限
	PowerManageGroupPage   = addPower(&PowerAction{Action: "manage_group_page", Text: "管理群组页面", ShouldLogin: true})
	PowerManageGroupInsert = addPower(&PowerAction{Action: "manage_group_insert", Text: "管理群组新增", Parent: PowerManageGroupPage, ShouldLogin: true})
	PowerManageGroupUpdate = addPower(&PowerAction{Action: "manage_group_update", Text: "管理群组修改", Parent: PowerManageGroupPage, ShouldLogin: true})
	PowerManageGroupDelete = addPower(&PowerAction{Action: "manage_group_delete", Text: "管理群组删除", Parent: PowerManageGroupPage, ShouldLogin: true})

	// 管理任务 权限
	PowerManageJobPage   = addPower(&PowerAction{Action: "manage_job_page", Text: "管理任务页面", ShouldLogin: true})
	PowerManageJobInsert = addPower(&PowerAction{Action: "manage_job_insert", Text: "管理任务新增", Parent: PowerManageJobPage, ShouldLogin: true})
	PowerManageJobUpdate = addPower(&PowerAction{Action: "manage_job_update", Text: "管理任务修改", Parent: PowerManageJobPage, ShouldLogin: true})
	PowerManageJobDelete = addPower(&PowerAction{Action: "manage_job_delete", Text: "管理任务删除", Parent: PowerManageJobPage, ShouldLogin: true})

	// 管理日志 权限
	PowerManageLogPage   = addPower(&PowerAction{Action: "manage_log_page", Text: "管理日志页面", ShouldLogin: true})
	PowerManageLogInsert = addPower(&PowerAction{Action: "manage_log_insert", Text: "管理日志新增", Parent: PowerManageLogPage, ShouldLogin: true})
	PowerManageLogUpdate = addPower(&PowerAction{Action: "manage_log_update", Text: "管理日志修改", Parent: PowerManageLogPage, ShouldLogin: true})
	PowerManageLogDelete = addPower(&PowerAction{Action: "manage_log_delete", Text: "管理日志删除", Parent: PowerManageLogPage, ShouldLogin: true})

	// 管理登录 权限
	PowerManageLoginPage   = addPower(&PowerAction{Action: "manage_login_page", Text: "管理登录页面", ShouldLogin: true})
	PowerManageLoginInsert = addPower(&PowerAction{Action: "manage_login_insert", Text: "管理登录新增", Parent: PowerManageLoginPage, ShouldLogin: true})
	PowerManageLoginUpdate = addPower(&PowerAction{Action: "manage_login_update", Text: "管理登录修改", Parent: PowerManageLoginPage, ShouldLogin: true})
	PowerManageLoginDelete = addPower(&PowerAction{Action: "manage_login_delete", Text: "管理登录删除", Parent: PowerManageLoginPage, ShouldLogin: true})

	// 管理系统设置 权限
	PowerManageSystemSettingPage   = addPower(&PowerAction{Action: "manage_system_setting_page", Text: "管理系统设置页面", ShouldLogin: true})
	PowerManageSystemSettingUpdate = addPower(&PowerAction{Action: "manage_system_setting_update", Text: "管理系统设置修改", Parent: PowerManageSystemSettingPage, ShouldLogin: true})

	//管理系统日志 权限
	PowerManageSystemLogPage   = addPower(&PowerAction{Action: "manage_system_log_page", Text: "管理日志页面", ShouldLogin: true})
	PowerManageSystemLogInsert = addPower(&PowerAction{Action: "manage_system_log_insert", Text: "管理日志新增", Parent: PowerManageSystemLogPage, ShouldLogin: true})
	PowerManageSystemLogUpdate = addPower(&PowerAction{Action: "manage_system_log_update", Text: "管理日志修改", Parent: PowerManageSystemLogPage, ShouldLogin: true})
	PowerManageSystemLogDelete = addPower(&PowerAction{Action: "manage_system_log_delete", Text: "管理日志删除", Parent: PowerManageSystemLogPage, ShouldLogin: true})
)

func addPower(power *PowerAction) *PowerAction {
	powers = append(powers, power)
	return power
}

func GetPowers() (ps []*PowerAction) {

	ps = append(ps, powers...)

	return
}
