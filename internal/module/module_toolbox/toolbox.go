package module_toolbox

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"time"
)

// NewToolboxService 根据库配置创建ToolboxService
func NewToolboxService(ServerContext *context.ServerContext) (res *ToolboxService) {

	idService := module_id.NewIDService(ServerContext)

	res = &ToolboxService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// ToolboxService 工具箱服务
type ToolboxService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Get 查询单个
func (this_ *ToolboxService) Get(toolboxId int64) (res *ToolboxModel, err error) {
	res = &ToolboxModel{}

	sql := `SELECT * FROM ` + TableToolbox + ` WHERE toolboxId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{toolboxId}, res)
	if err != nil {
		this_.Logger.Error("Get Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	}
	return
}

var visibilityOpen = 1

// QueryVisibility 查询 可见工具
func (this_ *ToolboxService) QueryVisibility(toolbox *ToolboxModel) (res []*ToolboxVisibilityModel, err error) {

	var values []interface{}
	sql := `SELECT toolboxId,toolboxType,groupId,name,comment,userId FROM ` + TableToolbox + ` WHERE deleted=2 `
	sql += " AND ("
	sql += " visibility = ?"
	values = append(values, visibilityOpen)
	if toolbox.UserId != 0 {
		sql += " OR userId = ?"
		values = append(values, toolbox.UserId)
	}
	sql += " ) "
	if toolbox.ToolboxType != "" {
		sql += " AND toolboxType = ?"
		values = append(values, toolbox.ToolboxType)
	}
	if toolbox.GroupId != 0 {
		sql += " AND groupId = ?"
		values = append(values, toolbox.GroupId)
	}
	if toolbox.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", toolbox.Name, "%"))
	}
	sql += " ORDER BY sequence ASC, name ASC "

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("QueryVisibility Error", zap.Error(err))
		return
	}

	return
}

// CountVisibility 查询
func (this_ *ToolboxService) CountVisibility(toolbox *ToolboxModel) (res int64, err error) {

	var values []interface{}
	sql := `SELECT COUNT(1) FROM ` + TableToolbox + ` WHERE deleted=2 `
	sql += " AND ("
	sql += " visibility = ?"
	values = append(values, visibilityOpen)
	if toolbox.UserId != 0 {
		sql += " OR userId = ?"
		values = append(values, toolbox.UserId)
	}
	sql += " ) "
	if toolbox.ToolboxType != "" {
		sql += " AND toolboxType = ?"
		values = append(values, toolbox.ToolboxType)
	}
	if toolbox.GroupId != 0 {
		sql += " AND groupId = ?"
		values = append(values, toolbox.GroupId)
	}
	if toolbox.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", toolbox.Name, "%"))
	}

	res, err = this_.DatabaseWorker.Count(sql, values)
	if err != nil {
		this_.Logger.Error("CountVisibility Error", zap.Error(err))
		return
	}

	return
}

// Query 查询 工具
func (this_ *ToolboxService) Query(toolbox *ToolboxModel) (res []*ToolboxModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableToolbox + ` WHERE deleted=2 `

	if toolbox.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, toolbox.UserId)
	}
	if toolbox.ToolboxType != "" {
		sql += " AND toolboxType = ?"
		values = append(values, toolbox.ToolboxType)
	}
	if toolbox.GroupId != 0 {
		sql += " AND groupId = ?"
		values = append(values, toolbox.GroupId)
	}
	if toolbox.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", toolbox.Name, "%"))
	}
	sql += " ORDER BY sequence ASC, name ASC "

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("Query Error", zap.Error(err))
		return
	}

	return
}

// UpdateSequence 更新顺序
func (this_ *ToolboxService) UpdateSequence(sequences map[int64]int) (err error) {
	if sequences == nil || len(sequences) == 0 {
		return
	}
	for id, sequence := range sequences {
		var values []interface{}

		sql := `UPDATE ` + TableToolbox + ` SET `
		sql += "sequence=? "
		values = append(values, sequence)

		sql += " WHERE toolboxId=? "
		values = append(values, id)

		_, err = this_.DatabaseWorker.Exec(sql, values)
		if err != nil {
			this_.Logger.Error("UpdateSequence Error", zap.Error(err))
			return
		}
	}
	return
}

// CheckUserToolboxExist 查询
func (this_ *ToolboxService) CheckUserToolboxExist(toolboxType string, name string, userId int64) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableToolbox + ` WHERE deleted=2 AND (userId = ? AND toolboxType = ? AND name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{userId, toolboxType, name})
	if err != nil {
		this_.Logger.Error("CheckUserToolboxExist Error", zap.Error(err))
		return
	}

	res = count > 0

	return
}

// GetUserToolboxByName 根据类型和名称
func (this_ *ToolboxService) GetUserToolboxByName(toolboxType string, name string, userId int64) (res *ToolboxModel, err error) {
	var list []*ToolboxModel

	sql := `SELECT * FROM ` + TableToolbox + ` WHERE deleted=2 AND (userId = ? AND toolboxType = ? AND name = ?)`
	err = this_.DatabaseWorker.Query(sql, []interface{}{userId, toolboxType, name}, &list)
	if err != nil {
		this_.Logger.Error("GetUserToolboxByName Error", zap.Error(err))
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}
	return
}

// Insert 新增
func (this_ *ToolboxService) Insert(toolbox *ToolboxModel) (rowsAffected int64, err error) {

	checked, err := this_.CheckUserToolboxExist(toolbox.ToolboxType, toolbox.Name, toolbox.UserId)
	if checked {
		err = errors.New(fmt.Sprint("工具[", toolbox.Name, "]已存在"))
		return
	}
	if toolbox.ToolboxId == 0 {
		toolbox.ToolboxId, err = this_.idService.GetNextID(module_id.IDTypeToolbox)
		if err != nil {
			return
		}
	}
	if toolbox.CreateTime.IsZero() {
		toolbox.CreateTime = time.Now()
	}

	var columns = "toolboxId, toolboxType, name, comment, visibility, option, userId, createTime"
	var values = "?, ?, ?, ?, ?, ?, ?, ?"

	if toolbox.GroupId > 0 {
		columns += ", groupId"
		values += ", " + fmt.Sprint(toolbox.GroupId)
	}

	err = this_.FormatOption(toolbox, false)
	if err != nil {
		return
	}

	sql := `INSERT INTO ` + TableToolbox + `(` + columns + `) VALUES (` + values + `) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolbox.ToolboxId, toolbox.ToolboxType, toolbox.Name, toolbox.Comment, toolbox.Visibility, toolbox.Option, toolbox.UserId, toolbox.CreateTime})
	if err != nil {
		this_.Logger.Error("Insert Error", zap.Error(err))
		return
	}

	return
}

// Open 新增
func (this_ *ToolboxService) Open(toolboxOpen *ToolboxOpenModel) (rowsAffected int64, err error) {

	if toolboxOpen.CreateTime.IsZero() {
		toolboxOpen.CreateTime = time.Now()
	}
	if toolboxOpen.Sequence == 0 {
		toolboxOpen.Sequence = time.Now().UnixMilli()
	}
	if toolboxOpen.OpenTime.IsZero() {
		toolboxOpen.OpenTime = time.Now()
	}
	if toolboxOpen.OpenId != 0 {
		if toolboxOpen.UpdateTime.IsZero() {
			toolboxOpen.UpdateTime = time.Now()
		}

		sql := `UPDATE ` + TableToolboxOpen + ` SET openTime=?,updateTime=? WHERE openId=? `
		rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpen.OpenTime, toolboxOpen.UpdateTime, toolboxOpen.OpenId})
		if err != nil {
			this_.Logger.Error("Open Error", zap.Error(err))
			return
		}
	} else {
		toolboxOpen.OpenId, err = this_.idService.GetNextID(module_id.IDTypeToolboxOpen)
		if err != nil {
			return
		}

		sql := `INSERT INTO ` + TableToolboxOpen + `(openId, userId, toolboxId, extend, sequence, createTime, openTime) VALUES (?, ?, ?, ?, ?, ?, ?) `

		rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpen.OpenId, toolboxOpen.UserId, toolboxOpen.ToolboxId, toolboxOpen.Extend, toolboxOpen.Sequence, toolboxOpen.CreateTime, toolboxOpen.OpenTime})
		if err != nil {
			this_.Logger.Error("Open Error", zap.Error(err))
			return
		}
	}

	return
}

// QueryOpens 查询
func (this_ *ToolboxService) QueryOpens(userId int64) (res []*ToolboxOpenModel, err error) {

	sql := `SELECT T_O.*,T_T.groupId toolboxGroupId,T_T.toolboxType toolboxType,T_T.name toolboxName ,T_T.comment toolboxComment
FROM ` + TableToolboxOpen + ` T_O 
LEFT JOIN ` + TableToolbox + ` T_T ON T_O.toolboxId = T_T.toolboxId 
WHERE T_O.userId=? ORDER BY T_O.sequence ASC `
	err = this_.DatabaseWorker.Query(sql, []interface{}{userId}, &res)
	if err != nil {
		this_.Logger.Error("QueryOpens Error", zap.Error(err))
		return
	}
	for _, one := range res {
		otherToolbox := this_.GetOtherToolbox(one.ToolboxId)
		if otherToolbox != nil {
			one.ToolboxType = otherToolbox.ToolboxType
			one.ToolboxName = otherToolbox.Name
		}
	}
	return
}

// GetOpen 查询单个
func (this_ *ToolboxService) GetOpen(openId int64) (res *ToolboxOpenModel, err error) {
	res = &ToolboxOpenModel{}

	sql := `SELECT T_O.*,T_T.groupId toolboxGroupId,T_T.toolboxType toolboxType,T_T.name toolboxName ,T_T.comment toolboxComment
FROM ` + TableToolboxOpen + ` T_O 
LEFT JOIN ` + TableToolbox + ` T_T ON T_O.toolboxId = T_T.toolboxId 
WHERE T_O.openId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{openId}, res)
	if err != nil {
		this_.Logger.Error("GetOpen Error", zap.Error(err))
		return
	}

	if !find {
		res = nil
	} else {
		otherToolbox := this_.GetOtherToolbox(res.ToolboxId)
		if otherToolbox != nil {
			res.ToolboxType = otherToolbox.ToolboxType
			res.ToolboxName = otherToolbox.Name
		}
	}
	return
}

// Close 更新
func (this_ *ToolboxService) Close(openId int64) (rowsAffected int64, err error) {

	sql := `DELETE FROM ` + TableToolboxOpen + ` WHERE openId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{openId})
	if err != nil {
		this_.Logger.Error("Close Error", zap.Error(err))
		return
	}
	_, err = this_.CloseOpenTabs(openId)
	if err != nil {
		this_.Logger.Error("CloseOpenTabs Error", zap.Error(err))
		return
	}
	return
}

// UpdateOpenExtend 新增
func (this_ *ToolboxService) UpdateOpenExtend(toolboxOpen *ToolboxOpenModel) (rowsAffected int64, err error) {

	if toolboxOpen.UpdateTime.IsZero() {
		toolboxOpen.UpdateTime = time.Now()
	}

	sql := `UPDATE ` + TableToolboxOpen + ` SET extend=?,updateTime=? WHERE openId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpen.Extend, toolboxOpen.UpdateTime, toolboxOpen.OpenId})
	if err != nil {
		this_.Logger.Error("UpdateOpenExtend Error", zap.Error(err))
		return
	}

	return
}

// UpdateOpenSequence 新增
func (this_ *ToolboxService) UpdateOpenSequence(toolboxOpen *ToolboxOpenModel) (rowsAffected int64, err error) {

	if toolboxOpen.UpdateTime.IsZero() {
		toolboxOpen.UpdateTime = time.Now()
	}

	sql := `UPDATE ` + TableToolboxOpen + ` SET sequence=?,updateTime=? WHERE openId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolboxOpen.Sequence, toolboxOpen.UpdateTime, toolboxOpen.OpenId})
	if err != nil {
		this_.Logger.Error("UpdateOpenSequence Error", zap.Error(err))
		return
	}

	return
}

// Update 更新
func (this_ *ToolboxService) Update(toolbox *ToolboxModel) (rowsAffected int64, err error) {
	if toolbox.Name != "" {
		var old *ToolboxModel
		old, err = this_.Get(toolbox.ToolboxId)
		if err != nil {
			return
		}
		if old == nil {
			err = errors.New("工具不存在")
			return
		}
		sql := `SELECT COUNT(1) FROM ` + TableToolbox + ` WHERE deleted=2 AND (toolboxId != ? AND userId = ? AND toolboxType = ? AND name = ?)`

		var count int64
		count, err = this_.DatabaseWorker.Count(sql, []interface{}{toolbox.ToolboxId, old.UserId, old.ToolboxType, toolbox.Name})
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.New(fmt.Sprint("工具[", toolbox.Name, "]已存在"))
			return
		}

	}

	err = this_.FormatOption(toolbox, false)
	if err != nil {
		return
	}

	var values []interface{}

	sql := `UPDATE ` + TableToolbox + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if toolbox.Name != "" {
		sql += "name=?,"
		values = append(values, toolbox.Name)
	}
	if toolbox.Comment != "" {
		sql += "comment=?,"
		values = append(values, toolbox.Comment)
	}

	if toolbox.Option != "" {
		sql += "option=?,"
		values = append(values, toolbox.Option)
	}

	if toolbox.Visibility != 0 {
		sql += "visibility=?,"
		values = append(values, toolbox.Visibility)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE toolboxId=? "
	values = append(values, toolbox.ToolboxId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("Update Error", zap.Error(err))
		return
	}

	return
}

// MoveGroup 更新
func (this_ *ToolboxService) MoveGroup(toolbox *ToolboxModel) (rowsAffected int64, err error) {

	var values []interface{}

	sql := `UPDATE ` + TableToolbox + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if toolbox.GroupId <= 0 {
		sql += "groupId=NULL,"
	} else {
		sql += "groupId=?,"
		values = append(values, toolbox.GroupId)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE toolboxId=? "
	values = append(values, toolbox.ToolboxId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		this_.Logger.Error("Update Error", zap.Error(err))
		return
	}

	return
}

// Rename 更新
func (this_ *ToolboxService) Rename(toolboxId int64, name string) (rowsAffected int64, err error) {
	_, err = this_.Update(&ToolboxModel{
		ToolboxId: toolboxId,
		Name:      name,
	})
	if err != nil {
		this_.Logger.Error("Rename Error", zap.Error(err))
		return
	}

	return
}

// Delete 更新
func (this_ *ToolboxService) Delete(toolboxId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableToolbox + ` SET deleted=?,deleteTime=? WHERE toolboxId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), toolboxId})
	if err != nil {
		this_.Logger.Error("Delete Error", zap.Error(err))
		return
	}

	return
}
