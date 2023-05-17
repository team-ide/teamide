package module_log

import (
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-tool/util"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
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
	_, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		return
	}
	return
}
