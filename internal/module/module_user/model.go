package module_user

import "time"

const (
	// ModuleUser 用户模块
	ModuleUser = "user"
	// TableUser 用户信息表
	TableUser = "TM_USER"
	// TableUserAuth 用户授权表
	TableUserAuth = "TM_USER_AUTH"
	// TableUserPassword 用户密码表
	TableUserPassword = "TM_USER_PASSWORD"
)

// UserModel 用户模型，和用户表对应
type UserModel struct {
	UserId     int64     `json:"userId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	Account    string    `json:"account,omitempty"`
	Email      string    `json:"email,omitempty"`
	Activated  int8      `json:"activated,omitempty"` // 激活 用户在注册时候使用邮箱激活，或管理员激活，未激活状态可以登录但不可以使用系统功能
	Locked     int8      `json:"locked,omitempty"`    // 锁定 账号异常系统自动锁定，如登录异常，系统可以自动解锁或管理员解锁
	Enabled    int8      `json:"enabled,omitempty"`   // 启用/禁用 管理员可以禁用用户，用户无法登录和使用系统，需要管理员启用
	Deleted    int8      `json:"deleted,omitempty"`   // 删除 已删除用户不可再使用
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}

// UserAuthModel 用户授权模型，和用户表对应
type UserAuthModel struct {
	AuthId     int64     `json:"authId,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	AuthType   int8      `json:"authType,omitempty"`
	OpenId     string    `json:"openId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	Homepage   string    `json:"homepage,omitempty"`
	Deleted    int8      `json:"deleted,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}

// UserPasswordModel 用户密码模型，和用户表对应
type UserPasswordModel struct {
	UserId     int64     `json:"userId,omitempty"`
	Salt       string    `json:"salt,omitempty"`
	Password   string    `json:"password,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}
