package module_power

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewPowerRoleService 根据库配置创建PowerRoleService
func NewPowerRoleService(ServerContext *context.ServerContext) (res *PowerRoleService) {

	idService := module_id.NewIDService(ServerContext)

	res = &PowerRoleService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// PowerRoleService 权限角色服务
type PowerRoleService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Insert 新增
func (this_ *PowerRoleService) Insert(powerRole *PowerRoleModel) (rowsAffected int64, err error) {

	if powerRole.PowerRoleId == 0 {
		powerRole.PowerRoleId, err = this_.idService.GetNextID(module_id.IDTypePowerRole)
		if err != nil {
			return
		}
	}
	if powerRole.CreateTime.IsZero() {
		powerRole.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TablePowerRole + `(powerRoleId, name, roleType, createTime) VALUES (?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{powerRole.PowerRoleId, powerRole.Name, powerRole.RoleType, time.Now()})
	if err != nil {
		return
	}

	return
}

// QueryByRoleType 根据 角色类型 查询 权限角色
func (this_ *PowerRoleService) QueryByRoleType(roleType int) (res []*PowerRoleModel, err error) {
	var values []interface{}
	sql := `SELECT * FROM ` + TablePowerRole + ` WHERE roleType=? `
	values = append(values, roleType)

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}
	return
}

// QueryPowerUsersByRoleType 根据 角色类型 查询 权限用户
func (this_ *PowerRoleService) QueryPowerUsersByRoleType(roleType int) (res []*PowerUserModel, err error) {
	var values []interface{}
	sql := `SELECT * FROM ` + TablePowerUser + ` WHERE powerRoleId IN `
	sql += `(SELECT powerRoleId FROM ` + TablePowerRole + ` WHERE roleType=?)`
	values = append(values, roleType)

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}
	return
}
