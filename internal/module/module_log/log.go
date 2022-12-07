package module_log

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/pkg/util"
)

// NewLogService 根据库配置创建LogService
func NewLogService(ServerContext *context.ServerContext) (res *LogService) {

	idService := module_id.NewIDService(ServerContext)

	res = &LogService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// LogService 日志服务
type LogService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Insert 新增
func (this_ *LogService) Insert(log *LogModel, errLog error) (err error) {

	if log.LogId == 0 {
		log.LogId, err = this_.idService.GetNextID(module_id.IDTypeLog)
		if err != nil {
			return
		}
	}

	log.Status = 1
	if errLog != nil {
		log.Status = 2
		log.Error = errLog.Error()
	}
	log.UseTime = int(util.GetTimeTime(log.EndTime) - util.GetTimeTime(log.StartTime))

	sql := `INSERT INTO ` + TableLog + `(logId, loginId, userId, userName, userAccount, ip, action, method, param, data, userAgent, status, error, useTime, startTime, endTime, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{log.LogId, log.LoginId, log.UserId, log.UserName, log.UserAccount, log.Ip, log.Action, log.Method, log.Param, log.Data, log.UserAgent, log.Status, log.Error, log.UseTime, log.StartTime, log.EndTime, log.CreateTime})
	if err != nil {
		return
	}
	return
}
