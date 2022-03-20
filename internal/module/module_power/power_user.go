package module_power

import (
	"teamide/internal/module/module_id"
	"teamide/pkg/db"
	"time"
)

// NewPowerUserService 根据库配置创建PowerUserService
func NewPowerUserService(dbWorker db.DatabaseWorker) (res *PowerUserService) {

	idService := module_id.NewIDService(dbWorker)

	res = &PowerUserService{
		dbWorker:  dbWorker,
		idService: idService,
	}
	return
}

// PowerUserService 权限角色服务
type PowerUserService struct {
	dbWorker  db.DatabaseWorker
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

	sql := `INSERT INTO ` + TablePowerUser + `(powerUserId, userId, createTime) VALUES (?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{powerUser.PowerUserId, powerUser.UserId, time.Now()})
	if err != nil {
		return
	}

	return
}
