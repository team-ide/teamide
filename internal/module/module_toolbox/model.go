package module_toolbox

import (
	"time"
)

const (
	// ModuleToolbox 工具箱模块
	ModuleToolbox = "toolbox"
	// TableToolbox 工具箱信息表
	TableToolbox        = "TM_TOOLBOX"
	TableToolboxComment = "工具箱"
	// TableToolboxOpen 工具箱打开信息表
	TableToolboxOpen        = "TM_TOOLBOX_OPEN"
	TableToolboxOpenComment = "工具箱打开记录"
	// TableToolboxOpenTab 工具箱打开标签页信息表
	TableToolboxOpenTab        = "TM_TOOLBOX_OPEN_TAB"
	TableToolboxOpenTabComment = "工具箱打开记录标签页"
	// TableToolboxGroup 工具箱分组
	TableToolboxGroup        = "TM_TOOLBOX_GROUP"
	TableToolboxGroupComment = "工具箱分组"
	// TableToolboxQuickCommand 工具箱快速命令
	TableToolboxQuickCommand        = "TM_TOOLBOX_QUICK_COMMAND"
	TableToolboxQuickCommandComment = "工具箱快速命令"
)

// ToolboxModel 工具箱模型，和工具箱表对应
type ToolboxModel struct {
	ToolboxId    int64     `json:"toolboxId,omitempty"`
	ToolboxType  string    `json:"toolboxType,omitempty"`
	GroupId      int64     `json:"groupId,omitempty"`
	Name         string    `json:"name,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	Option       string    `json:"option,omitempty"`
	UserId       int64     `json:"userId,omitempty"`
	Visibility   int       `json:"visibility,omitempty"`
	DeleteUserId int64     `json:"deleteUserId,omitempty"`
	Deleted      int8      `json:"deleted,omitempty"`
	CreateTime   time.Time `json:"createTime,omitempty"`
	UpdateTime   time.Time `json:"updateTime,omitempty"`
	DeleteTime   time.Time `json:"deleteTime,omitempty"`
}

// ToolboxOpenModel 工具箱打开模型，和工具箱打开表对应
type ToolboxOpenModel struct {
	OpenId         int64     `json:"openId,omitempty"`
	UserId         int64     `json:"userId,omitempty"`
	ToolboxId      int64     `json:"toolboxId,omitempty"`
	ToolboxType    string    `json:"toolboxType,omitempty"`
	ToolboxName    string    `json:"toolboxName,omitempty"`
	ToolboxComment string    `json:"toolboxComment,omitempty"`
	ToolboxGroupId int64     `json:"toolboxGroupId,omitempty"`
	Extend         string    `json:"extend,omitempty"`
	Sequence       int64     `json:"sequence,omitempty"`
	CreateTime     time.Time `json:"createTime,omitempty"`
	UpdateTime     time.Time `json:"updateTime,omitempty"`
	OpenTime       time.Time `json:"openTime,omitempty"`
}

// ToolboxOpenTabModel 工具箱打开模型，和工具箱打开表对应
type ToolboxOpenTabModel struct {
	TabId      int64     `json:"tabId,omitempty"`
	OpenId     int64     `json:"openId,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	ToolboxId  int64     `json:"toolboxId,omitempty"`
	Extend     string    `json:"extend,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
	OpenTime   time.Time `json:"openTime,omitempty"`
}

// ToolboxGroupModel 工具箱分组模型，和工具箱分组表对应
type ToolboxGroupModel struct {
	GroupId    int64     `json:"groupId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	Option     string    `json:"option,omitempty"`
	UserId     int64     `json:"userId,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}

// ToolboxQuickCommandModel 工具箱快速命令
type ToolboxQuickCommandModel struct {
	QuickCommandId   int64     `json:"quickCommandId,omitempty"`
	QuickCommandType int       `json:"quickCommandType,omitempty"`
	Name             string    `json:"name,omitempty"`
	Comment          string    `json:"comment,omitempty"`
	Option           string    `json:"option,omitempty"`
	UserId           int64     `json:"userId,omitempty"`
	CreateTime       time.Time `json:"createTime,omitempty"`
	UpdateTime       time.Time `json:"updateTime,omitempty"`
}
