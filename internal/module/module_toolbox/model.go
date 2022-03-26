package module_toolbox

import (
	"time"
)

const (
	// ModuleToolbox 工具箱模块
	ModuleToolbox = "toolbox"
	// TableToolbox 工具箱信息表
	TableToolbox = "TM_TOOLBOX"
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
