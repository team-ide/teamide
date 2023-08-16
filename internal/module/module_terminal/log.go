package module_terminal

import (
	"go.uber.org/zap"
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
	res.init()
	return
}

// TerminalLogService 日志服务
type TerminalLogService struct {
	*context.ServerContext
	idService *module_id.IDService
}

func (this_ *TerminalLogService) init() {

	this_.cleanTask()
	// 每天 2 点执行
	_, _ = this_.CronHandler.AddFunc("0 0 2 * * ?", this_.cleanTask)
	return
}

func (this_ *TerminalLogService) cleanTask() {
	saveDays := this_.ServerConfig.LogDataSaveDays
	var deleteCount int64
	this_.Logger.Info("log data clean task start", zap.Any("saveDays", saveDays))
	defer func() {
		this_.Logger.Info("log data clean task end", zap.Any("saveDays", saveDays), zap.Any("deleteCount", deleteCount))
	}()
	if saveDays <= 0 {
		return
	}
	deleteBeforeTime := time.Now().AddDate(0, 0, -saveDays)
	this_.Logger.Info("log data clean task info", zap.Any("saveDays", saveDays), zap.Any("deleteBeforeTime", deleteBeforeTime))
	var sql string
	var values []interface{}

	sql += "DELETE FROM " + TableTerminalLog + " WHERE createTime<? "
	values = append(values, deleteBeforeTime)
	deleteCount, _ = this_.DatabaseWorker.Exec(sql, values)

	_, _ = this_.DatabaseWorker.GetDb().Exec("VACUUM")
	return
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

	sql := `INSERT INTO ` + TableTerminalLog + `(terminalLogId, loginId, workerId, userId, userName, userAccount, ip, place, placeId, command, userAgent, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{log.TerminalLogId, log.LoginId, log.WorkerId, log.UserId, log.UserName, log.UserAccount, log.Ip, log.Place, log.PlaceId, log.Command, log.UserAgent, log.CreateTime})
	if err != nil {
		return
	}
	return
}
