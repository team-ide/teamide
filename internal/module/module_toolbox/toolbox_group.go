package module_toolbox

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/module/module_id"
	"time"
)

// GetGroup 查询单个
func (this_ *ToolboxService) GetGroup(groupId int64) (res *ToolboxGroupModel, err error) {
	res = &ToolboxGroupModel{}

	sql := `SELECT * FROM ` + TableToolboxGroup + ` WHERE groupId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{groupId}, res)
	if err != nil {
		this_.Logger.Error("GetGroup Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	}
	return
}

// UpdateGroupSequence 更新顺序
func (this_ *ToolboxService) UpdateGroupSequence(sequences map[int64]int) (err error) {
	if sequences == nil || len(sequences) == 0 {
		return
	}
	for id, sequence := range sequences {
		var values []interface{}

		sql := `UPDATE ` + TableToolboxGroup + ` SET `

		sql += "sequence=? "
		values = append(values, sequence)

		sql += " WHERE groupId=? "
		values = append(values, id)

		_, err = this_.DatabaseWorker.Exec(sql, values)
		if err != nil {
			this_.Logger.Error("UpdateGroupSequence Error", zap.Error(err))
			return
		}
	}
	return
}

// QueryGroup 查询
func (this_ *ToolboxService) QueryGroup(toolboxGroup *ToolboxGroupModel) (res []*ToolboxGroupModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableToolboxGroup + ` WHERE 1=1 `

	if toolboxGroup.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, toolboxGroup.UserId)
	}
	if toolboxGroup.Name != "" {
		sql += " AND name LIKE ?"
		values = append(values, fmt.Sprint("%", toolboxGroup.Name, "%"))
	}

	sql += " ORDER BY sequence ASC "

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("QueryGroup Error", zap.Error(err))
		return
	}

	return
}

// CheckUserGroupExist 查询
func (this_ *ToolboxService) CheckUserGroupExist(name string, userId int64) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableToolboxGroup + ` WHERE userId = ?  AND name = ?`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{userId, name})
	if err != nil {
		this_.Logger.Error("CheckUserToolboxGroupExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// InsertGroup 新增
func (this_ *ToolboxService) InsertGroup(toolboxGroup *ToolboxGroupModel) (rowsAffected int64, err error) {

	checked, err := this_.CheckUserGroupExist(toolboxGroup.Name, toolboxGroup.UserId)
	if checked {
		err = errors.New(fmt.Sprint("工具分组[", toolboxGroup.Name, "]已存在"))
		return
	}
	if toolboxGroup.GroupId == 0 {
		toolboxGroup.GroupId, err = this_.idService.GetNextID(module_id.IDTypeToolboxGroup)
		if err != nil {
			return
		}
	}
	if toolboxGroup.CreateTime.IsZero() {
		toolboxGroup.CreateTime = time.Now()
	}

	sql := `INSERT INTO ` + TableToolboxGroup + `(groupId, name, comment, option, userId, createTime) VALUES (?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxGroup.GroupId, toolboxGroup.Name, toolboxGroup.Comment, toolboxGroup.Option, toolboxGroup.UserId, toolboxGroup.CreateTime})
	if err != nil {
		this_.Logger.Error("InsertGroup Error", zap.Error(err))
		return
	}

	return
}

// UpdateGroup 更新
func (this_ *ToolboxService) UpdateGroup(toolboxGroup *ToolboxGroupModel) (rowsAffected int64, err error) {
	if toolboxGroup.Name != "" {
		var old *ToolboxGroupModel
		old, err = this_.GetGroup(toolboxGroup.GroupId)
		if err != nil {
			return
		}
		if old == nil {
			err = errors.New("工具分组不存在")
			return
		}
		sql := `SELECT COUNT(1) FROM ` + TableToolboxGroup + ` WHERE groupId != ? AND userId = ? AND name = ?`

		var count int64
		count, err = this_.DatabaseWorker.Count(sql, []interface{}{toolboxGroup.GroupId, old.UserId, toolboxGroup.Name})
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.New(fmt.Sprint("工具分组[", toolboxGroup.Name, "]已存在"))
			return
		}

	}
	var values []interface{}

	sql := `UPDATE ` + TableToolboxGroup + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if toolboxGroup.Name != "" {
		sql += "name=?,"
		values = append(values, toolboxGroup.Name)
	}
	if toolboxGroup.Comment != "" {
		sql += "comment=?,"
		values = append(values, toolboxGroup.Comment)
	}
	if toolboxGroup.Option != "" {
		sql += "option=?,"
		values = append(values, toolboxGroup.Option)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE groupId=? "
	values = append(values, toolboxGroup.GroupId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("UpdateGroup Error", zap.Error(err))
		return
	}

	return
}

// DeleteGroup 更新
func (this_ *ToolboxService) DeleteGroup(groupId int64) (rowsAffected int64, err error) {
	sql := `UPDATE ` + TableToolbox + ` SET groupId=NULL,updateTime=? WHERE groupId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{time.Now(), groupId})
	if err != nil {
		this_.Logger.Error("DeleteGroup Trim Toolbox GroupId Error", zap.Error(err))
		return
	}

	sql = `DELETE FROM ` + TableToolboxGroup + ` WHERE groupId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{groupId})
	if err != nil {
		this_.Logger.Error("DeleteGroup Error", zap.Error(err))
		return
	}

	return
}
