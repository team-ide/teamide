package module_id

import (
	"teamide/internal/module/module_lock"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

// NewIDService 根据库配置创建IDService
func NewIDService(dbWorker db.DatabaseWorker) (res *IDService) {
	res = &IDService{
		dbWorker: dbWorker,
	}
	return
}

// IDService ID服务
type IDService struct {
	dbWorker db.DatabaseWorker
}

// GetNextID 根据类型获取一个ID
func (this_ *IDService) GetNextID(idType IDType) (id int64, err error) {
	ids, err := this_.GetNextIDs(idType, 1)
	if err != nil {
		return
	}
	id = ids[0]
	return
}

// GetNextIDs 根据类型获取一组ID
func (this_ *IDService) GetNextIDs(idType IDType, size int64) (ids []int64, err error) {
	locker := module_lock.GetLock("ID:GetNextIDs")
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
	list, err := this_.dbWorker.Query(sql, []interface{}{idType}, util.GetStructFieldTypes(IDModel{}))
	if err != nil {
		return
	}

	if len(list) == 0 {
		id = 0

		sql = `INSERT INTO ` + TableID + ` (idType, value, createTime) VALUES(?, ?, ?)`
		_, err = this_.dbWorker.Exec(sql, []interface{}{idType, id, time.Now()})
		if err != nil {
			return
		}

	} else {
		id = list[0]["value"].(int64)
	}
	return
}

// updateID 修改ID
func (this_ *IDService) updateID(idType IDType, id int64) (err error) {

	sql := `UPDATE ` + TableID + ` SET value=?,updateTime=? WHERE idType=?`
	_, err = this_.dbWorker.Exec(sql, []interface{}{id, time.Now(), idType})
	if err != nil {
		return
	}

	return
}
