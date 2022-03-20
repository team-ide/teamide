package model

import "time"

const (
	// ModuleID ID模块
	ModuleID = "id"
	// TableID ID信息表
	TableID = "TM_ID"
)

// IDModel ID模型，和ID表对应
type IDModel struct {
	IdType     int       `json:"idType,omitempty"`
	Value      int64     `json:"value,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}

type IDType int

const (
	// IDTypeUser 用户ID类型
	IDTypeUser = 1
	// IDTypeUserAuth 户授权ID类型
	IDTypeUserAuth = 2
	// IDTypeUserPassword 户授权ID类型
	IDTypeUserPassword = 3
	// IDTypeRegister 注册ID类型
	IDTypeRegister = 4
	// IDTypeLogin 登录ID类型
	IDTypeLogin = 5
)
