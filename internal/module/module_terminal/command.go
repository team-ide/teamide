package module_terminal

import (
	"errors"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewTerminalCommandService 根据库配置创建TerminalCommandService
func NewTerminalCommandService(ServerContext *context.ServerContext) (res *TerminalCommandService) {

	idService := module_id.NewIDService(ServerContext)

	res = &TerminalCommandService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	res.init()
	return
}

// TerminalCommandService 日志服务
type TerminalCommandService struct {
	*context.ServerContext
	idService *module_id.IDService
}

func (this_ *TerminalCommandService) init() {

	this_.cleanDeprecatedLog()
	return
}

func (this_ *TerminalCommandService) cleanDeprecatedLog() {
	var sql string
	var values []interface{}

	sql += "DELETE FROM " + TableTerminalLog + "  "
	_, _ = this_.DatabaseWorker.Exec(sql, values)

	_, _ = this_.DatabaseWorker.GetDb().Exec("VACUUM")
	return
}

// Insert 新增
func (this_ *TerminalCommandService) Insert(command *TerminalCommandModel) (err error) {

	if command.TerminalCommandId == 0 {
		command.TerminalCommandId, err = this_.idService.GetNextID(module_id.IDTypeTerminalCommand)
		if err != nil {
			return
		}
	}
	if command.CreateTime.IsZero() {
		command.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableTerminalCommand + `(terminalCommandId, loginId, workerId, userId, userName, userAccount, ip, place, placeId, command, userAgent, createTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{
		command.TerminalCommandId,
		command.LoginId,
		command.WorkerId,
		command.UserId,
		command.UserName,
		command.UserAccount,
		command.Ip,
		command.Place,
		command.PlaceId,
		command.Command,
		command.UserAgent,
		command.CreateTime,
	})
	if err != nil {
		return
	}
	return
}

// Query 查询
func (this_ *TerminalCommandService) Query(command *TerminalCommandModel) (list []*TerminalCommand, err error) {

	var sqlInfo = "SELECT command FROM " + TableTerminalCommand + " WHERE 1=1 "
	var values []interface{}

	if command.Place != "" {
		sqlInfo += " AND place=? "
		values = append(values, command.Place)
	}
	if command.PlaceId != "" {
		sqlInfo += " AND placeId=? "
		values = append(values, command.PlaceId)
	}
	if command.UserId != 0 {
		sqlInfo += " AND userId=? "
		values = append(values, command.UserId)
	}
	if command.WorkerId != "" {
		sqlInfo += " AND workerId=? "
		values = append(values, command.WorkerId)
	}

	err = this_.DatabaseWorker.Query(sqlInfo, values, &list)
	if err != nil {
		return
	}
	return
}

// Count 查询
func (this_ *TerminalCommandService) Count(command *TerminalCommandModel) (res int64, err error) {

	var sqlInfo = "SELECT COUNT(*) FROM " + TableTerminalCommand + " WHERE 1=1 "
	var values []interface{}

	if command.Place != "" {
		sqlInfo += " AND place=? "
		values = append(values, command.Place)
	}
	if command.PlaceId != "" {
		sqlInfo += " AND placeId=? "
		values = append(values, command.PlaceId)
	}
	if command.UserId != 0 {
		sqlInfo += " AND userId=? "
		values = append(values, command.UserId)
	}
	if command.WorkerId != "" {
		sqlInfo += " AND workerId=? "
		values = append(values, command.WorkerId)
	}

	res, err = this_.DatabaseWorker.Count(sqlInfo, values)
	if err != nil {
		return
	}
	return
}

// Clean 清理
func (this_ *TerminalCommandService) Clean(command *TerminalCommandModel) (list []*TerminalCommandModel, err error) {

	var sqlInfo = "DELETE FROM " + TableTerminalCommand + " WHERE 1=1 "
	var values []interface{}

	if command.Place != "" {
		sqlInfo += " AND place=? "
		values = append(values, command.Place)
	}
	if command.PlaceId != "" {
		sqlInfo += " AND placeId=? "
		values = append(values, command.PlaceId)
	}
	if command.UserId != 0 {
		sqlInfo += " AND userId=? "
		values = append(values, command.UserId)
	}
	if command.WorkerId != "" {
		sqlInfo += " AND workerId=? "
		values = append(values, command.WorkerId)
	}
	if len(values) == 0 {
		err = errors.New("清理历史命令需要最少一个参数")
		return
	}

	_, err = this_.DatabaseWorker.Exec(sqlInfo, values)
	if err != nil {
		return
	}
	return
}
