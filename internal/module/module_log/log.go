package module_log

import (
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewLogService 根据库配置创建LogService
func NewLogService(ServerContext *context.ServerContext) (res *LogService, err error) {

	idService := module_id.NewIDService(ServerContext)

	res = &LogService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	err = res.init()
	return
}

// LogService 日志服务
type LogService struct {
	*context.ServerContext
	idService *module_id.IDService
}

func (this_ *LogService) init() (err error) {

	this_.cleanTask()
	// 每天 2 点执行
	_, _ = this_.CronHandler.AddFunc("0 0 2 * * ?", this_.cleanTask)
	return
}

func (this_ *LogService) cleanTask() {
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

	sql += "DELETE FROM " + TableLog + " WHERE createTime<? "
	values = append(values, deleteBeforeTime)
	deleteCount, _ = this_.DatabaseWorker.Exec(sql, values)

	_, _ = this_.DatabaseWorker.GetDb().Exec("VACUUM")
	return
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
	log.UseTime = int(util.GetMilliByTime(log.EndTime) - util.GetMilliByTime(log.StartTime))

	sql := `INSERT INTO ` + TableLog + `(logId, loginId, userId, userName, userAccount, ip, action, method, param, data, userAgent, status, error, useTime, startTime, endTime, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{log.LogId, log.LoginId, log.UserId, log.UserName, log.UserAccount, log.Ip, log.Action, log.Method, log.Param, log.Data, log.UserAgent, log.Status, log.Error, log.UseTime, log.StartTime, log.EndTime, log.CreateTime})
	if err != nil {
		return
	}
	return
}

type LogPage struct {
	*worker.Page
	DataList []*LogModel `json:"dataList"`
}

// QueryPage 新增
func (this_ *LogService) QueryPage(log *LogModel, page *LogPage) (err error) {
	var sql string
	var values []interface{}

	sql += "SELECT * FROM " + TableLog + " WHERE 1=1"
	if log.UserId != 0 {
		sql += " AND userId=?"
		values = append(values, log.UserId)
	}
	if log.Action != "" {
		sql += " AND action=?"
		values = append(values, log.Action)
	}
	if !log.StartTime.IsZero() {
		sql += " AND (startTime>=? OR endTime>=?)"
		values = append(values, log.StartTime, log.StartTime)
	}
	if !log.EndTime.IsZero() {
		sql += " AND (startTime<=? OR endTime<=?)"
		values = append(values, log.EndTime, log.EndTime)
	}
	sql += " ORDER BY createTime DESC"
	page.DataList = []*LogModel{}
	err = this_.DatabaseWorker.QueryPage(sql, values, &page.DataList, page.Page)
	if err != nil {
		return
	}
	return
}

func (this_ *LogService) clean(log *LogModel) (err error) {
	var sql string
	var values []interface{}

	sql += "DELETE FROM " + TableLog + " WHERE 1=1"
	if log.UserId != 0 {
		sql += " AND userId=?"
		values = append(values, log.UserId)
	}
	_, err = this_.DatabaseWorker.GetDb().Exec(sql, values)
	if err != nil {
		return
	}
	_, _ = this_.DatabaseWorker.GetDb().Exec("VACUUM")
	return
}
