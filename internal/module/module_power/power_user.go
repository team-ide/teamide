package module_power

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewPowerUserService 根据库配置创建PowerUserService
func NewPowerUserService(ServerContext *context.ServerContext) (res *PowerUserService) {

	idService := module_id.NewIDService(ServerContext)

	res = &PowerUserService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// PowerUserService 权限角色服务
type PowerUserService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Insert 新增
func (this_ *PowerUserService) Insert(powerUser *PowerUserModel) (rowsAffected int64, err error) {

	if powerUser.PowerUserId == 0 {
		powerUser.PowerUserId, err = this_.idService.GetNextID(module_id.IDTypePowerUser)
		if err != nil {
			return
		}
	}
	if powerUser.CreateTime.IsZero() {
		powerUser.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TablePowerUser + `(powerUserId, userId, powerRoleId, createTime) VALUES (?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{powerUser.PowerUserId, powerUser.UserId, powerUser.PowerRoleId, time.Now()})
	if err != nil {
		return
	}

	return
}

// QueryPowerRolesByUserId 根据 用户ID 查询角色
func (this_ *PowerUserService) QueryPowerRolesByUserId(userId int64) (res []*PowerRoleModel, err error) {
	var values []interface{}
	sql := `SELECT * FROM ` + TablePowerRole + ` WHERE powerRoleId IN `
	sql += `(SELECT powerRoleId FROM ` + TablePowerUser + ` WHERE userId=?)`
	values = append(values, userId)

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}
	return
}
