package module_power

import (
	"teamide/internal/module/module_id"
	"teamide/pkg/db"
	"time"
)

// NewPowerRoleService 根据库配置创建PowerRoleService
func NewPowerRoleService(dbWorker db.DatabaseWorker) (res *PowerRoleService) {

	idService := module_id.NewIDService(dbWorker)

	res = &PowerRoleService{
		dbWorker:  dbWorker,
		idService: idService,
	}
	return
}

// PowerRoleService 权限角色服务
type PowerRoleService struct {
	dbWorker  db.DatabaseWorker
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

	sql := `INSERT INTO ` + TablePowerRole + `(powerRoleId, name, createTime) VALUES (?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{powerRole.PowerRoleId, powerRole.Name, time.Now()})
	if err != nil {
		return
	}

	return
}
