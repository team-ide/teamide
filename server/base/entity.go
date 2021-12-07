package base

import (
	"time"
)

type IDEntity struct {
	ServerId int64 `json:"serverId" column:"serverId"`
	Type     int8  `json:"type" column:"type"`
	Id       int64 `json:"id" column:"id"`
}

func NewIDEntity() (res *IDEntity) {
	res = &IDEntity{}
	return res
}
func NewIDEntityInterface() interface{} {
	return NewIDEntity()
}

type UserEntity struct {
	ServerId     int64     `json:"serverId" column:"serverId"`
	UserId       int64     `json:"userId" column:"userId"`
	Name         string    `json:"name" column:"name"`
	Avatar       string    `json:"avatar" column:"avatar"`
	Account      string    `json:"account" column:"account"`
	Email        string    `json:"email" column:"email"`
	ActivedState int8      `json:"activedState" column:"activedState"`
	LockedState  int8      `json:"lockedState" column:"lockedState"`
	EnabledState int8      `json:"enabledState" column:"enabledState"`
	CreateTime   time.Time `json:"createTime" column:"createTime"`
	UpdateTime   time.Time `json:"updateTime" column:"updateTime"`
}

func NewUserEntity() (res *UserEntity) {
	res = &UserEntity{}
	return res
}
func NewUserEntityInterface() interface{} {
	return NewUserEntity()
}

type UserMetadataEntity struct {
	ServerId       int64     `json:"serverId" column:"serverId"`
	MetadataId     int64     `json:"metadataId" column:"metadataId"`
	UserId         int64     `json:"userId" column:"userId"`
	MetadataStruct int       `json:"metadataStruct" column:"metadataStruct"`
	MetadataField  int       `json:"metadataField" column:"metadataField"`
	MetadataValue  string    `json:"metadataValue" column:"metadataValue"`
	ParentId       int64     `json:"parentId" column:"parentId"`
	CreateTime     time.Time `json:"createTime" column:"createTime"`
	UpdateTime     time.Time `json:"updateTime" column:"updateTime"`
}

func NewUserMetadataEntity() (res *UserMetadataEntity) {
	res = &UserMetadataEntity{}
	return res
}
func NewUserMetadataEntityInterface() interface{} {
	return NewUserMetadataEntity()
}

type UserAuthEntity struct {
	ServerId     int64     `json:"serverId" column:"serverId"`
	AuthId       int64     `json:"authId" column:"authId"`
	UserId       int64     `json:"userId" column:"userId"`
	AuthType     int8      `json:"authType" column:"authType"`
	OpenId       string    `json:"openId" column:"openId"`
	Name         string    `json:"name" column:"name"`
	Avatar       string    `json:"avatar" column:"avatar"`
	ActivedState int8      `json:"activedState" column:"activedState"`
	LockedState  int8      `json:"lockedState" column:"lockedState"`
	EnabledState int8      `json:"enabledState" column:"enabledState"`
	CreateTime   time.Time `json:"createTime" column:"createTime"`
	UpdateTime   time.Time `json:"updateTime" column:"updateTime"`
}

func NewUserAuthEntity() (res *UserAuthEntity) {
	res = &UserAuthEntity{}
	return res
}
func NewUserAuthEntityInterface() interface{} {
	return NewUserAuthEntity()
}

type UserPasswordEntity struct {
	ServerId   int64     `json:"serverId" column:"serverId"`
	UserId     int64     `json:"userId" column:"userId"`
	Salt       string    `json:"salt" column:"salt"`
	Password   string    `json:"password" column:"password"`
	CreateTime time.Time `json:"createTime" column:"createTime"`
	UpdateTime time.Time `json:"updateTime" column:"updateTime"`
}

func NewUserPasswordEntity() (res *UserPasswordEntity) {
	res = &UserPasswordEntity{}
	return res
}
func NewUserPasswordEntityInterface() interface{} {
	return NewUserPasswordEntity()
}
