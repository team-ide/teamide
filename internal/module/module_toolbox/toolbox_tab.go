package module_toolbox

import (
	"go.uber.org/zap"
	"teamide/internal/module/module_id"
	"time"
)

// OpenTab 新增
func (this_ *ToolboxService) OpenTab(toolboxOpenTab *ToolboxOpenTabModel) (rowsAffected int64, err error) {

	if toolboxOpenTab.CreateTime.IsZero() {
		toolboxOpenTab.CreateTime = time.Now()
	}
	if toolboxOpenTab.OpenTime.IsZero() {
		toolboxOpenTab.OpenTime = time.Now()
	}
	if toolboxOpenTab.TabId != 0 {
		if toolboxOpenTab.UpdateTime.IsZero() {
			toolboxOpenTab.UpdateTime = time.Now()
		}

		sql := `UPDATE ` + TableToolboxOpenTab + ` SET openTime=?,updateTime=? WHERE tabId=? `
		rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpenTab.OpenTime, toolboxOpenTab.UpdateTime, toolboxOpenTab.TabId})
		if err != nil {
			this_.Logger.Error("OpenTab Error", zap.Error(err))
			return
		}
	} else {
		toolboxOpenTab.TabId, err = this_.idService.GetNextID(module_id.IDTypeToolboxOpenTab)
		if err != nil {
			return
		}

		sql := `INSERT INTO ` + TableToolboxOpenTab + `(tabId, openId, userId, toolboxId, extend, createTime, openTime) VALUES (?, ?, ?, ?, ?, ?, ?) `

		rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpenTab.TabId, toolboxOpenTab.OpenId, toolboxOpenTab.UserId, toolboxOpenTab.ToolboxId, toolboxOpenTab.Extend, toolboxOpenTab.CreateTime, toolboxOpenTab.OpenTime})
		if err != nil {
			this_.Logger.Error("OpenTab Error", zap.Error(err))
			return
		}
	}

	return
}

// QueryOpenTabs 查询
func (this_ *ToolboxService) QueryOpenTabs(openId int64) (res []*ToolboxOpenTabModel, err error) {

	sql := `SELECT * FROM ` + TableToolboxOpenTab + ` WHERE openId=? ORDER BY createTime ASC `
	err = this_.DatabaseWorker.Query(sql, []interface{}{openId}, &res)
	if err != nil {
		this_.Logger.Error("QueryOpenTabs Error", zap.Error(err))
		return
	}

	return
}

// CloseTab 更新
func (this_ *ToolboxService) CloseTab(tabId int64) (rowsAffected int64, err error) {

	sql := `DELETE FROM ` + TableToolboxOpenTab + ` WHERE tabId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{tabId})
	if err != nil {
		this_.Logger.Error("CloseTab Error", zap.Error(err))
		return
	}

	return
}

// CloseOpenTabs 更新
func (this_ *ToolboxService) CloseOpenTabs(openId int64) (rowsAffected int64, err error) {

	sql := `DELETE FROM ` + TableToolboxOpenTab + ` WHERE openId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{openId})
	if err != nil {
		this_.Logger.Error("CloseOpenTabs Error", zap.Error(err))
		return
	}

	return
}

// UpdateOpenTabExtend 新增
func (this_ *ToolboxService) UpdateOpenTabExtend(toolboxOpenTab *ToolboxOpenTabModel) (rowsAffected int64, err error) {
	sql := `UPDATE ` + TableToolboxOpenTab + ` SET extend=?,updateTime=? WHERE tabId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpenTab.Extend, toolboxOpenTab.UpdateTime, toolboxOpenTab.TabId})
	if err != nil {
		this_.Logger.Error("UpdateOpenTabExtend Error", zap.Error(err))
		return
	}

	return
}
