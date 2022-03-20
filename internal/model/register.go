package model

import "time"

const (
	// ModuleRegister 注册模块
	ModuleRegister = "register"
	// TableRegister 注册信息表
	TableRegister = "TM_REGISTER"
)

// RegisterModel 注册模型，和注册表对应
type RegisterModel struct {
	RegisterId int64     `json:"registerId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Account    string    `json:"account,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	Ip         string    `json:"ip,omitempty"`
	SourceType int       `json:"sourceType,omitempty"`
	Source     string    `json:"source,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	Deleted    int8      `json:"deleted,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	DeleteTime time.Time `json:"deleteTime,omitempty"`
}
