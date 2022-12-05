package module_power

import "time"

const (
	// ModulePower 权限模块
	ModulePower = "power"
	// TablePowerRole 权限角色表
	TablePowerRole        = "TM_POWER_ROLE"
	TablePowerRoleComment = "权限角色"
	// TablePowerRoute 权限路由表
	TablePowerRoute        = "TM_POWER_ROUTE"
	TablePowerRouteComment = "权限路由"
	// TablePowerUser 权限用户表
	TablePowerUser        = "TM_POWER_USER"
	TablePowerUserComment = "权限用户"
)

// PowerRoleModel 权限角色模型，和权限角色表对应
type PowerRoleModel struct {
	PowerRoleId    int64     `json:"powerRoleId,omitempty"`
	Name           string    `json:"name,omitempty"`
	RoleType       int       `json:"roleType,omitempty"`
	ExpirationTime time.Time `json:"expirationTime,omitempty"`
	CreateTime     time.Time `json:"createTime,omitempty"`
	UpdateTime     time.Time `json:"updateTime,omitempty"`
	IsSuper        bool      `json:"isSuper"`
}

// PowerRouteModel 权限路由模型，和权限路由表对应
type PowerRouteModel struct {
	PowerRouteId   int64     `json:"powerRouteId,omitempty"`
	PowerRoleId    int64     `json:"powerRoleId,omitempty"`
	Name           string    `json:"name,omitempty"`
	Route          string    `json:"route,omitempty"`
	ExpirationTime time.Time `json:"expirationTime,omitempty"`
	CreateTime     time.Time `json:"createTime,omitempty"`
	UpdateTime     time.Time `json:"updateTime,omitempty"`
}

// PowerUserModel 权限用户模型，和权限用户表对应
type PowerUserModel struct {
	PowerUserId    int64     `json:"powerUserId,omitempty"`
	PowerRoleId    int64     `json:"powerRoleId,omitempty"`
	UserId         int64     `json:"userId,omitempty"`
	ExpirationTime time.Time `json:"expirationTime,omitempty"`
	CreateTime     time.Time `json:"createTime,omitempty"`
	UpdateTime     time.Time `json:"updateTime,omitempty"`
}
