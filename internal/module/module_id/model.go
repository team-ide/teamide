package module_id

import "time"

const (
	// ModuleID ID模块
	ModuleID = "id"
	// TableID ID信息表
	TableID        = "TM_ID"
	TableIDComment = "ID信息"
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
	IDTypeUser = 1001
	// IDTypeUserAuth 户授权ID类型
	IDTypeUserAuth = 1002
	// IDTypeRegister 注册ID类型
	IDTypeRegister = 2001

	// IDTypeLogin 登录ID类型
	IDTypeLogin = 3001

	// IDTypePowerRole 权限角色ID类型
	IDTypePowerRole = 4001
	// IDTypePowerRoute 权限路由ID类型
	IDTypePowerRoute = 4002
	// IDTypePowerUser 权限用户ID类型
	IDTypePowerUser = 4003

	// IDTypeToolbox 工具箱ID类型
	IDTypeToolbox = 5001
	// IDTypeToolboxOpen 工具箱开启ID类型
	IDTypeToolboxOpen = 5002
	// IDTypeToolboxOpenTab 工具箱开启标签页ID类型
	IDTypeToolboxOpenTab = 5003
	// IDTypeToolboxGroup 工具箱分组ID类型
	IDTypeToolboxGroup = 5004
	// IDTypeToolboxQuickCommand 工具箱快速命令ID类型
	IDTypeToolboxQuickCommand = 5005

	// IDTypeNode 节点
	IDTypeNode = 6001
	// IDTypeNodeNetProxy 节点
	IDTypeNodeNetProxy = 6002

	// IDTypeLog 日志
	IDTypeLog = 7001

	// IDTypeTerminalLog 控制台日志
	IDTypeTerminalLog = 8001
)
