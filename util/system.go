package util

import "os/user"

var (
	SystemUser_Uid      string // 用户的 ID
	SystemUser_Gid      string // 用户所属组的 ID，如果属于多个组，那么此 ID 为主组的 ID
	SystemUser_Username string // 用户名
	SystemUser_Name     string // 属组名称，如果属于多个组，那么此名称为主组的名称
	SystemUser_HomeDir  string // 用户的宿主目录
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	SystemUser_Username = u.Username
	SystemUser_Name = u.Name
	SystemUser_HomeDir = u.HomeDir
	SystemUser_Gid = u.Gid
	SystemUser_Uid = u.Uid
}
