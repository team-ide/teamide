package module_user

import (
	"teamide/internal/module/module_id"
	"teamide/pkg/db"
	"time"
)

// NewUserAuthService 根据库配置创建UserAuthService
func NewUserAuthService(dbWorker db.DatabaseWorker) (res *UserAuthService) {

	idService := module_id.NewIDService(dbWorker)

	res = &UserAuthService{
		dbWorker:  dbWorker,
		idService: idService,
	}
	return
}

// UserAuthService 用户授权服务
type UserAuthService struct {
	dbWorker  db.DatabaseWorker
	idService *module_id.IDService
}

// Insert 新增
func (this_ *UserAuthService) Insert(userAuth *UserAuthModel) (rowsAffected int64, err error) {

	if userAuth.AuthId == 0 {
		userAuth.AuthId, err = this_.idService.GetNextID(module_id.IDTypeUserAuth)
		if err != nil {
			return
		}
	}
	if userAuth.CreateTime.IsZero() {
		userAuth.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableUserAuth + `(authId, userId, createTime) VALUES (?, ?, ?) `

	rowsAffected, err = this_.dbWorker.Exec(sql, []interface{}{userAuth.AuthId, userAuth.UserId, time.Now()})
	if err != nil {
		return
	}

	return
}
