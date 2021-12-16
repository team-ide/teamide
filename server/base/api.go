package base

import "github.com/gin-gonic/gin"

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

type RequestBean struct {
	JWT  *JWTBean
	Path string
}

type JWTBean struct {
	ServerId int64  `json:"serverId"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Time     int64  `json:"time"`
}

var (
	powers []*PowerAction

	// 基础权限
	PowerRegister  = &PowerAction{Action: "register", Text: "注册"}
	PowerData      = &PowerAction{Action: "data", Text: "数据"}
	PowerSession   = &PowerAction{Action: "session", Text: "会话"}
	PowerLogin     = &PowerAction{Action: "login", Text: "登录"}
	PowerLogout    = &PowerAction{Action: "logout", Text: "登出"}
	PowerAutoLogin = &PowerAction{Action: "auto_login", Text: "自动登录"}

	// 管理用户权限
	PowerManageUserPage    = &PowerAction{Action: "manage_user_page", Text: "管理用户页面", ShouldLogin: true}
	PowerManageUserInsert  = &PowerAction{Action: "manage_user_insert", Text: "管理用户新增", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserUpdate  = &PowerAction{Action: "manage_user_update", Text: "管理用户修改", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserDelete  = &PowerAction{Action: "manage_user_delete", Text: "管理用户删除", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserActive  = &PowerAction{Action: "manage_user_active", Text: "管理用户激活", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserLock    = &PowerAction{Action: "manage_user_lock", Text: "管理用户锁定", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserUnlock  = &PowerAction{Action: "manage_user_unlock", Text: "管理用户解锁", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserEnable  = &PowerAction{Action: "manage_user_enable", Text: "管理用户启用", Parent: PowerManageUserPage, ShouldLogin: true}
	PowerManageUserDisable = &PowerAction{Action: "manage_user_disable", Text: "管理用户禁用", Parent: PowerManageUserPage, ShouldLogin: true}

	// 管理设置权限
	PowerManageSettingPage  = &PowerAction{Action: "manage_setting_page", Text: "管理设置页面", ShouldLogin: true}
	PowerManageSettingUser  = &PowerAction{Action: "manage_setting_user", Text: "管理设置用户", Parent: PowerManageSettingPage, ShouldLogin: true}
	PowerManageSettingLogin = &PowerAction{Action: "manage_setting_login", Text: "管理设置登录", Parent: PowerManageSettingPage, ShouldLogin: true}
)

func init() {

	// 添加 基础权限
	powers = append(powers, PowerRegister)
	powers = append(powers, PowerData)
	powers = append(powers, PowerSession)
	powers = append(powers, PowerLogin)
	powers = append(powers, PowerLogout)
	powers = append(powers, PowerAutoLogin)

	// 添加 管理用户权限
	powers = append(powers, PowerManageUserPage)
	powers = append(powers, PowerManageUserInsert)
	powers = append(powers, PowerManageUserUpdate)
	powers = append(powers, PowerManageUserDelete)
	powers = append(powers, PowerManageUserActive)
	powers = append(powers, PowerManageUserLock)
	powers = append(powers, PowerManageUserUnlock)
	powers = append(powers, PowerManageUserEnable)
	powers = append(powers, PowerManageUserDisable)

	// 添加 管理设置权限
	powers = append(powers, PowerManageSettingPage)
	powers = append(powers, PowerManageSettingUser)
	powers = append(powers, PowerManageSettingLogin)
}

func GetPowers() (ps []*PowerAction) {

	ps = append(ps, powers...)

	return
}
