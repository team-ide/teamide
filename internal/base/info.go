package base

import "os/user"

var (
	// StandAloneUserId 单机版本 用户 ID
	StandAloneUserId   int64  = 1
	SystemUserUid      string // 用户的 ID
	SystemUserGid      string // 用户所属组的 ID，如果属于多个组，那么此 ID 为主组的 ID
	SystemUserUsername string // 用户名
	SystemUserName     string // 属组名称，如果属于多个组，那么此名称为主组的名称
	SystemUserHomeDir  string // 用户的宿主目录

	// SuperRoleType 超管 角色 类型
	SuperRoleType = 1
	// SuperRoleName 超管 角色 名称
	SuperRoleName = "超管"
	// SuperUserAccount 超管用户账号
	SuperUserAccount = "admin"
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	SystemUserUsername = u.Username
	SystemUserName = u.Name
	SystemUserHomeDir = u.HomeDir
	SystemUserGid = u.Gid
	SystemUserUid = u.Uid
}
