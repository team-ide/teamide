package module_power

import (
	"teamide/internal/module/module_id"
	"teamide/pkg/db"
	"time"
)

// NewPowerRouteService 根据库配置创建PowerRouteService
func NewPowerRouteService(dbWorker db.DatabaseWorker) (res *PowerRouteService) {

	idService := module_id.NewIDService(dbWorker)

	res = &PowerRouteService{
		dbWorker:  dbWorker,
		idService: idService,
	}
	return
}

// PowerRouteService 权限角色服务
type PowerRouteService struct {
	dbWorker  db.DatabaseWorker
	idService *module_id.IDService
}

// Insert 新增
func (this_ *PowerRouteService) Insert(powerRoute *PowerRouteModel) (rowsAffected int64, err error) {

	if powerRoute.PowerRouteId == 0 {
		powerRoute.PowerRouteId, err = this_.idService.GetNextID(module_id.IDTypePowerRoute)
		if err != nil {
			return
		}
	}
	if powerRoute.CreateTime.IsZero() {
		powerRoute.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TablePowerRoute + `(powerRouteId, powerRoleId, name, createTime) VALUES (?, ?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{powerRoute.PowerRouteId, powerRoute.PowerRoleId, powerRoute.Name, time.Now()})
	if err != nil {
		return
	}

	return
}
