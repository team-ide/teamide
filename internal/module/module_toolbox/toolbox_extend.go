package module_toolbox

import (
	"encoding/json"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"teamide/internal/module/module_id"
	"time"
)

func (this_ *ToolboxService) FormatExtend(extendType string, value string) (res map[string]interface{}, err error) {
	res, err = util.JsonToMap(value)
	return
}

// GetExtend 查询单个
func (this_ *ToolboxService) GetExtend(extendId int64) (res *ToolboxExtendModel, err error) {
	res = &ToolboxExtendModel{}

	sql := `SELECT * FROM ` + TableToolboxExtend + ` WHERE extendId=? `
	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{extendId}, res)
	if err != nil {
		this_.Logger.Error("GetExtend Error", zap.Error(err))
		return
	}
	if !find {
		res = nil
	} else {
		res.Extend, err = this_.FormatExtend(res.ExtendType, res.Value)
	}
	return
}

// SaveExtend 新增
func (this_ *ToolboxService) SaveExtend(data *ToolboxExtendModel) (err error) {

	if data.CreateTime.IsZero() {
		data.CreateTime = time.Now()
	}
	var bs []byte
	if data.Extend != nil {
		filePath := util.GetStringValue(data.Extend["filePath"])
		if filePath != "" && data.Extend["fileSize"] == nil {
			if fs, _ := os.Stat(this_.GetFilesFile(filePath)); fs != nil {
				data.Extend["fileSize"] = fs.Size()
			}
		}
		bs, err = json.Marshal(data.Extend)
		if err != nil {
			return
		}
	}
	data.Value = string(bs)

	if data.ExtendId != 0 {
		if data.UpdateTime.IsZero() {
			data.UpdateTime = time.Now()
		}

		sql := `UPDATE ` + TableToolboxExtend + ` SET name=?,value=?,updateTime=? WHERE extendId=? `
		_, err = this_.DatabaseWorker.Exec(sql, []interface{}{data.Name, data.Value, data.UpdateTime, data.ExtendId})
		if err != nil {
			this_.Logger.Error("SaveExtend Error", zap.Error(err))
			return
		}
	} else {
		data.ExtendId, err = this_.idService.GetNextID(module_id.IDTypeToolboxExtend)
		if err != nil {
			return
		}

		sql := `INSERT INTO ` + TableToolboxExtend + `(extendId, toolboxId, extendType, name, value, userId, createTime) VALUES (?, ?, ?, ?, ?, ?, ?) `

		_, err = this_.DatabaseWorker.Exec(sql, []interface{}{data.ExtendId, data.ToolboxId, data.ExtendType, data.Name, data.Value, data.UserId, data.CreateTime})
		if err != nil {
			this_.Logger.Error("SaveExtend Error", zap.Error(err))
			return
		}
	}

	return
}

// QueryExtends 查询
func (this_ *ToolboxService) QueryExtends(q *ToolboxExtendModel) (res []*ToolboxExtendModel, err error) {

	sql := `SELECT * FROM ` + TableToolboxExtend + ` WHERE 1=1  `
	var values []interface{}
	if q.ToolboxId > 0 {
		sql += " AND toolboxId=?"
		values = append(values, q.ToolboxId)
	}
	if q.ExtendType != "" {
		sql += " AND extendType=?"
		values = append(values, q.ExtendType)
	}
	if q.UserId > 0 {
		sql += " AND userId=?"
		values = append(values, q.UserId)
	}

	sql += " ORDER BY createTime DESC"
	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		this_.Logger.Error("QueryExtends Error", zap.Error(err))
		return
	}
	for _, one := range res {
		one.Extend, _ = this_.FormatExtend(one.ExtendType, one.Value)
	}

	return
}

// DeleteExtend 删除
func (this_ *ToolboxService) DeleteExtend(extendId int64) (rowsAffected int64, err error) {
	find, _ := this_.GetExtend(extendId)
	if find != nil {
		filePath := util.GetStringValue(find.Extend["filePath"])
		if filePath != "" {
			filePath = this_.GetFilesFile(filePath)
			if e, _ := util.PathExists(filePath); e {
				_ = os.Remove(filePath)
			}
		}
	}
	sql := `DELETE FROM ` + TableToolboxExtend + ` WHERE extendId=? `
	rowsAffected, err = this_.DatabaseWorker.Exec(sql, []interface{}{extendId})
	if err != nil {
		this_.Logger.Error("DeleteExtend Error", zap.Error(err))
		return
	}

	return
}
