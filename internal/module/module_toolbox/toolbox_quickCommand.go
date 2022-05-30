package module_toolbox

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/module/module_id"
	"time"
)

// GetQuickCommand 查询单个
func (this_ *ToolboxService) GetQuickCommand(quickCommandId int64) (res *ToolboxQuickCommandModel, err error) {
	res = &ToolboxQuickCommandModel{}

	sql := `SELECT * FROM ` + TableToolboxQuickCommand + ` WHERE quickCommandId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{quickCommandId}, res)
	if err != nil {
		this_.Logger.Error("GetQuickCommand Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	}
	return
}

// QueryQuickCommand 查询
func (this_ *ToolboxService) QueryQuickCommand(toolboxQuickCommand *ToolboxQuickCommandModel) (res []*ToolboxQuickCommandModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableToolboxQuickCommand + ` WHERE 1=1 `

	if toolboxQuickCommand.QuickCommandType != 0 {
		sql += " AND quickCommandType = ?"
		values = append(values, toolboxQuickCommand.QuickCommandType)
	}
	if toolboxQuickCommand.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, toolboxQuickCommand.UserId)
	}
	if toolboxQuickCommand.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", toolboxQuickCommand.Name, "%"))
	}

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("QueryQuickCommand Error", zap.Error(err))
		return
	}

	return
}

// InsertQuickCommand 新增
func (this_ *ToolboxService) InsertQuickCommand(toolboxQuickCommand *ToolboxQuickCommandModel) (rowsAffected int64, err error) {

	if toolboxQuickCommand.QuickCommandId == 0 {
		toolboxQuickCommand.QuickCommandId, err = this_.idService.GetNextID(module_id.IDTypeToolboxQuickCommand)
		if err != nil {
			return
		}
	}
	if toolboxQuickCommand.CreateTime.IsZero() {
		toolboxQuickCommand.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableToolboxQuickCommand + `(quickCommandId, quickCommandType, name, comment, option, userId, createTime) VALUES (?, ?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxQuickCommand.QuickCommandId, toolboxQuickCommand.QuickCommandType, toolboxQuickCommand.Name, toolboxQuickCommand.Comment, toolboxQuickCommand.Option, toolboxQuickCommand.UserId, toolboxQuickCommand.CreateTime})
	if err != nil {
		this_.Logger.Error("InsertQuickCommand Error", zap.Error(err))
		return
	}

	return
}

// UpdateQuickCommand 更新
func (this_ *ToolboxService) UpdateQuickCommand(toolboxQuickCommand *ToolboxQuickCommandModel) (rowsAffected int64, err error) {
	var values []interface{}

	sql := `UPDATE ` + TableToolboxQuickCommand + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if toolboxQuickCommand.QuickCommandType != 0 {
		sql += "quickCommandType=?,"
		values = append(values, toolboxQuickCommand.QuickCommandType)
	}
	if toolboxQuickCommand.Name != "" {
		sql += "name=?,"
		values = append(values, toolboxQuickCommand.Name)
	}
	if toolboxQuickCommand.Comment != "" {
		sql += "comment=?,"
		values = append(values, toolboxQuickCommand.Comment)
	}
	if toolboxQuickCommand.Option != "" {
		sql += "option=?,"
		values = append(values, toolboxQuickCommand.Option)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE quickCommandId=? "
	values = append(values, toolboxQuickCommand.QuickCommandId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateQuickCommand Error", zap.Error(err))
		return
	}

	return
}

// DeleteQuickCommand 更新
func (this_ *ToolboxService) DeleteQuickCommand(quickCommandId int64) (rowsAffected int64, err error) {

	sql := `DELETE FROM ` + TableToolboxQuickCommand + ` WHERE quickCommandId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{quickCommandId})
	if err != nil {
		this_.Logger.Error("DeleteQuickCommand Error", zap.Error(err))
		return
	}

	return
}
