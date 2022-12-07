package module_log

import (
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/pkg/util"
	"time"
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

// Start 新增
func (this_ *LogService) Start(log *LogModel) (err error) {

	if log.LogId == 0 {
		log.LogId, err = this_.idService.GetNextID(module_id.IDTypeLog)
		if err != nil {
			return
		}
	}
	if log.StartTime.IsZero() {
		log.StartTime = time.Now()
	}
	if log.CreateTime.IsZero() {
		log.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableLog + `(logId, userId, ip, action, method, param, data, userAgent, startTime, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{log.LogId, log.UserId, log.Ip, log.Action, log.Method, log.Param, log.Data, log.UserAgent, log.StartTime, log.CreateTime})
	if err != nil {
		return
	}
	return
}

// End 日志结束
func (this_ *LogService) End(logId int64, startTime time.Time, errLog error) (err error) {
	errStr := ""
	status := 1
	if errLog != nil {
		status = 2
	}
	var useTime int64
	var endTime = util.Now()
	useTime = util.GetTimeTime(endTime) - util.GetTimeTime(startTime)
	sql := `UPDATE ` + TableLog + ` SET error=?,status=?,endTime=?,useTime=? WHERE logId=? `
	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{errStr, status, endTime, useTime, logId})
	if err != nil {
		return
	}

	return
}
