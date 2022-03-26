package module_toolbox

import (
	"errors"
	"fmt"
	"strings"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/pkg/util"
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

	sql := `SELECT * FROM ` + TableToolbox + ` WHERE toolboxId=? `
	list, err := this_.DatabaseWorker.Query(sql, []interface{}{toolboxId}, util.GetStructFieldTypes(ToolboxModel{}))
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = &ToolboxModel{}
		err = util.ToStruct(list[0], res)
	} else {
		res = nil
	}
	return
}

// Query 查询
func (this_ *ToolboxService) Query(toolbox *ToolboxModel) (res []*ToolboxModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableToolbox + ` WHERE deleted=2 `
	if toolbox.ToolboxType != "" {
		sql += " AND toolboxType = ?"
		values = append(values, toolbox.ToolboxType)
	}
	if toolbox.UserId != 0 {
		sql += " AND userId = ?"
		values = append(values, toolbox.UserId)
	}
	if toolbox.Name != "" {
		sql += " AND name like ?"
		values = append(values, fmt.Sprint("%", toolbox.Name, "%"))
	}

	list, err := this_.DatabaseWorker.Query(sql, values, util.GetStructFieldTypes(ToolboxModel{}))
	if err != nil {
		return
	}

	err = util.ToStruct(list, &res)
	if err != nil {
		return
	}

	return
}

// CheckUserToolboxExist 查询
func (this_ *ToolboxService) CheckUserToolboxExist(toolboxType string, name string, userId int64) (res bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + TableToolbox + ` WHERE deleted=2 AND (userId = ? AND toolboxType = ? AND name = ?)`

	count, err := this_.DatabaseWorker.Count(sql, []interface{}{userId, toolboxType, name})
	if err != nil {
		return
	}

	res = count > 0

	return
}

// GetUserToolboxByName 根据类型和名称
func (this_ *ToolboxService) GetUserToolboxByName(toolboxType string, name string, userId int64) (res *ToolboxModel, err error) {

	sql := `SELECT * FROM ` + TableToolbox + ` WHERE deleted=2 AND (userId = ? AND toolboxType = ? AND name = ?)`
	list, err := this_.DatabaseWorker.Query(sql, []interface{}{userId, toolboxType, name}, util.GetStructFieldTypes(ToolboxModel{}))
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = &ToolboxModel{}
		err = util.ToStruct(list[0], res)
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

	sql := `INSERT INTO ` + TableToolbox + `(toolboxId, toolboxType, name, option, userId, createTime) VALUES (?, ?, ?, ?, ?, ?) `

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{toolbox.ToolboxId, toolbox.ToolboxType, toolbox.Name, toolbox.Option, toolbox.UserId, toolbox.CreateTime})
	if err != nil {
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
	var values []interface{}

	sql := `UPDATE ` + TableToolbox + ` SET `

	sql += "updateTime=?,"
	values = append(values, time.Now())

	if toolbox.Name != "" {
		sql += "name=?,"
		values = append(values, toolbox.Name)
	}
	if toolbox.Option != "" {
		sql += "option=?,"
		values = append(values, toolbox.Option)
	}

	sql = strings.TrimSuffix(sql, ",")

	sql += " WHERE toolboxId=? "
	values = append(values, toolbox.ToolboxId)

	rowsAffected, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
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
		return
	}

	return
}

// Delete 更新
func (this_ *ToolboxService) Delete(toolboxId int64) (rowsAffected int64, err error) {

	sql := `UPDATE ` + TableToolbox + ` SET deleted=?,deleteTime=? WHERE toolboxId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{1, time.Now(), toolboxId})
	if err != nil {
		return
	}

	return
}
