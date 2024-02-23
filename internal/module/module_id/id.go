package module_id

import (
	"fmt"
	"go.uber.org/zap"
	"teamide/internal/context"
	"teamide/internal/module/module_lock"
	"time"
)

// NewIDService 根据库配置创建IDService
func NewIDService(ServerContext *context.ServerContext) (res *IDService) {
	res = &IDService{
		ServerContext: ServerContext,
	}
	return
}

// IDService ID服务
type IDService struct {
	*context.ServerContext
}

// GetNextID 根据类型获取一个ID
func (this_ *IDService) GetNextID(idType IDType) (id int64, err error) {
	ids, err := this_.GetNextIDs(idType, 1)
	if err != nil {
		this_.Logger.Error("GetNextID Error", zap.Error(err))
		return
	}
	id = ids[0]
	return
}

// GetNextIDs 根据类型获取一组ID
func (this_ *IDService) GetNextIDs(idType IDType, size int64) (ids []int64, err error) {
	locker := module_lock.GetLock(fmt.Sprintf("ID:GetNextIDs:%d", idType))
	locker.Lock()
	defer locker.Unlock()

	var id int64
	id, err = this_.getID(idType)
	if err != nil {
		return
	}
	var index int64
	for {
		if index == size {
			break
		}
		index++
		ids = append(ids, id+index)
	}
	err = this_.updateID(idType, id+size)
	if err != nil {
		return
	}

	return
}

// GetID 查询ID
func (this_ *IDService) getID(idType IDType) (id int64, err error) {

	sql := `SELECT value FROM ` + TableID + ` WHERE idType=? `

	find, err := this_.DatabaseWorker.QueryOne(sql, []interface{}{int(idType)}, &id)
	if err != nil {
		return
	}

	if !find {
		id = 0

		sql = `INSERT INTO ` + TableID + ` (idType, value, createTime) VALUES(?, ?, ?)`
		_, err = this_.DatabaseWorker.Exec(sql, []interface{}{int(idType), id, time.Now()})
		if err != nil {
			return
		}

	}
	return
}

// updateID 修改ID
func (this_ *IDService) updateID(idType IDType, id int64) (err error) {

	sql := `UPDATE ` + TableID + ` SET value=?,updateTime=? WHERE idType=?`
	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{id, time.Now(), int(idType)})
	if err != nil {
		return
	}

	return
}
