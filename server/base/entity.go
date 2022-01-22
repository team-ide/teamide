package base

import (
	"time"
)

type IDEntity struct {
	ServerId int64 `json:"serverId" column:"serverId,omitempty"`
	Type     int8  `json:"type" column:"type,omitempty"`
	Id       int64 `json:"id" column:"id,omitempty"`
}

func NewIDEntity() (res *IDEntity) {
	res = &IDEntity{}
	return res
}
func NewIDEntityInterface() interface{} {
	return NewIDEntity()
}

type UserEntity struct {
	ServerId     int64     `json:"serverId" column:"serverId,omitempty"`
	UserId       int64     `json:"userId" column:"userId,omitempty"`
	Name         string    `json:"name" column:"name,omitempty"`
	Avatar       string    `json:"avatar" column:"avatar,omitempty"`
	Account      string    `json:"account" column:"account,omitempty"`
	Email        string    `json:"email" column:"email,omitempty"`
	ActivedState int8      `json:"activedState" column:"activedState,omitempty"`
	LockedState  int8      `json:"lockedState" column:"lockedState,omitempty"`
	EnabledState int8      `json:"enabledState" column:"enabledState,omitempty"`
	CreateTime   time.Time `json:"createTime" column:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime" column:"updateTime,omitempty"`
}

func NewUserEntity() (res *UserEntity) {
	res = &UserEntity{}
	return res
}
func NewUserEntityInterface() interface{} {
	return NewUserEntity()
}

type UserMetadataEntity struct {
	ServerId       int64     `json:"serverId" column:"serverId,omitempty"`
	MetadataId     int64     `json:"metadataId" column:"metadataId,omitempty"`
	UserId         int64     `json:"userId" column:"userId,omitempty"`
	MetadataStruct int       `json:"metadataStruct" column:"metadataStruct,omitempty"`
	MetadataField  int       `json:"metadataField" column:"metadataField,omitempty"`
	MetadataValue  string    `json:"metadataValue" column:"metadataValue,omitempty"`
	ParentId       int64     `json:"parentId" column:"parentId,omitempty"`
	CreateTime     time.Time `json:"createTime" column:"createTime,omitempty"`
	UpdateTime     time.Time `json:"updateTime" column:"updateTime,omitempty"`
}

func NewUserMetadataEntity() (res *UserMetadataEntity) {
	res = &UserMetadataEntity{}
	return res
}
func NewUserMetadataEntityInterface() interface{} {
	return NewUserMetadataEntity()
}

type UserAuthEntity struct {
	ServerId     int64     `json:"serverId" column:"serverId,omitempty"`
	AuthId       int64     `json:"authId" column:"authId,omitempty"`
	UserId       int64     `json:"userId" column:"userId,omitempty"`
	AuthType     int8      `json:"authType" column:"authType,omitempty"`
	OpenId       string    `json:"openId" column:"openId,omitempty"`
	Name         string    `json:"name" column:"name,omitempty"`
	Avatar       string    `json:"avatar" column:"avatar,omitempty"`
	ActivedState int8      `json:"activedState" column:"activedState,omitempty"`
	LockedState  int8      `json:"lockedState" column:"lockedState,omitempty"`
	EnabledState int8      `json:"enabledState" column:"enabledState,omitempty"`
	CreateTime   time.Time `json:"createTime" column:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime" column:"updateTime,omitempty"`
}

func NewUserAuthEntity() (res *UserAuthEntity) {
	res = &UserAuthEntity{}
	return res
}
func NewUserAuthEntityInterface() interface{} {
	return NewUserAuthEntity()
}

type UserPasswordEntity struct {
	ServerId   int64     `json:"serverId" column:"serverId,omitempty"`
	UserId     int64     `json:"userId" column:"userId,omitempty"`
	Salt       string    `json:"salt" column:"salt,omitempty"`
	Password   string    `json:"password" column:"password,omitempty"`
	CreateTime time.Time `json:"createTime" column:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime" column:"updateTime,omitempty"`
}

func NewUserPasswordEntity() (res *UserPasswordEntity) {
	res = &UserPasswordEntity{}
	return res
}
func NewUserPasswordEntityInterface() interface{} {
	return NewUserPasswordEntity()
}
