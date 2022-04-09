package module_toolbox

import (
	"time"
)

const (
	// ModuleToolbox 工具箱模块
	ModuleToolbox = "toolbox"
	// TableToolbox 工具箱信息表
	TableToolbox = "TM_TOOLBOX"
	// TableToolboxOpen 工具箱打开信息表
	TableToolboxOpen = "TM_TOOLBOX_OPEN"
)

// ToolboxModel 工具箱模型，和工具箱表对应
type ToolboxModel struct {
	ToolboxId    int64     `json:"toolboxId,omitempty"`
	ToolboxType  string    `json:"toolboxType,omitempty"`
	Name         string    `json:"name,omitempty"`
	Option       string    `json:"option,omitempty"`
	UserId       int64     `json:"userId,omitempty"`
	DeleteUserId int64     `json:"deleteUserId,omitempty"`
	Deleted      int8      `json:"deleted,omitempty"`
	CreateTime   time.Time `json:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`
	DeleteTime   time.Time `json:"deleteTime,omitempty"`
}

// ToolboxOpenModel 工具箱打开模型，和工具箱打开表对应
type ToolboxOpenModel struct {
	OpenId     int64     `json:"openId,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	ToolboxId  int64     `json:"toolboxId,omitempty"`
	Extend     string    `json:"extend,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	OpenTime   time.Time `json:"openTime,omitempty"`
}
