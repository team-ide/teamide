package module_login

import "time"

const (
	// ModuleLogin 登录模块
	ModuleLogin = "login"
	// TableLogin 登录记录表
	TableLogin        = "TM_LOGIN"
	TableLoginComment = "登录记录"
)

// LoginModel 登录模型，和登录表对应
type LoginModel struct {
	LoginId    int64     `json:"loginId,omitempty"`
	Account    string    `json:"account,omitempty"`
	Password   string    `json:"password,omitempty"`
	Ip         string    `json:"ip,omitempty"`
	SourceType int       `json:"sourceType,omitempty"`
	Source     string    `json:"source,omitempty"`
	UserAgent  string    `json:"userAgent,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	Deleted    int8      `json:"deleted,omitempty"`
	LoginTime  time.Time `json:"loginTime,omitempty"`
	LogoutTime time.Time `json:"logoutTime,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}

var SourceTypeWeb = 1
