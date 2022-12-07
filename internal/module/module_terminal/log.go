package module_terminal

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewTerminalLogService 根据库配置创建TerminalLogService
func NewTerminalLogService(ServerContext *context.ServerContext) (res *TerminalLogService) {

	idService := module_id.NewIDService(ServerContext)

	res = &TerminalLogService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// TerminalLogService 日志服务
type TerminalLogService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Insert 新增
func (this_ *TerminalLogService) Insert(log *TerminalLogModel) (err error) {

	if log.TerminalLogId == 0 {
		log.TerminalLogId, err = this_.idService.GetNextID(module_id.IDTypeTerminalLog)
		if err != nil {
			return
		}
	}
	if log.CreateTime.IsZero() {
		log.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableTerminalLog + `(terminalLogId, loginId, userId, userName, userAccount, ip, place, placeId, command, userAgent, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{log.TerminalLogId, log.LoginId, log.UserId, log.UserName, log.UserAccount, log.Ip, log.Place, log.PlaceId, log.Command, log.UserAgent, log.CreateTime})
	if err != nil {
		return
	}
	return
}
