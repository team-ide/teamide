package idService

import (
	"teamide/server/base"
	"teamide/server/component"
)

var (
	TABLE_ID = "TM_ID"
)

type IdService struct {
}

func (this_ *IdService) GetID(idType component.IDType) (id int64, err error) {

	var ids []int64
	ids, err = this_.GetIDs(idType, 1)

	if err != nil {
		return
	}
	id = ids[0]

	return
}

func (this_ *IdService) GetIDs(idType component.IDType, size int64) (ids []int64, err error) {

	var key = component.GetIDRedisKey(idType)
	var exists bool
	exists, err = component.Redis.Exists(key)
	if err != nil {
		return
	}
	if !exists {
		var unlock func() (err error)
		unlock, err = component.Redis.Lock(key+":lock", 10, 1000)
		if err != nil {
			return
		}
		defer unlock()

		exists, err = component.Redis.Exists(key)
		if err != nil {
			return
		}
		if !exists {
			var idInfo *base.IDEntity
			idInfo, err = get(idType)
			if err != nil {
				return
			}
			var id_ int64 = component.GetBaseID()
			if idInfo == nil {
				idInfo = &base.IDEntity{
					Type: int8(idType),
					Id:   0,
				}
			} else {
				id_ = idInfo.Id + 100
			}
			err = component.Redis.SetInt64(key, id_)
			if err != nil {
				return
			}
			idInfo.Id = id_
			err = insertOrUpdate(*idInfo)
			if err != nil {
				return
			}
		}
		unlock()
	}
	var maxId int64
	maxId, err = component.Redis.IncrBy(key, size)
	if err != nil {
		return
	}
	var minId = maxId - size + 1
	ids = []int64{}
	var id int64
	for id = minId; id <= maxId; id++ {
		if id%50 == 0 {
			idInfo := &base.IDEntity{
				Type: int8(idType),
				Id:   id + 50,
			}
			err = insertOrUpdate(*idInfo)
			if err != nil {
				return
			}
		}
		ids = append(ids, id)
	}

	return
}

func insertOrUpdate(id base.IDEntity) (err error) {
	sql := "INSERT INTO " + TABLE_ID + " (serverId, type, id) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=?"
	params := []interface{}{component.GetServerId(), id.Type, id.Id, id.Id}

	sqlParam := base.NewSqlParam(sql, params)

	_, err = component.DB.Exec(sqlParam)

	if err != nil {
		return
	}
	return
}

//查询单个ID
func get(idType component.IDType) (id *base.IDEntity, err error) {
	sql := "SELECT * FROM " + TABLE_ID + " WHERE serverId=? AND type=? "
	params := []interface{}{component.GetServerId(), idType}

	sqlParam := base.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component.DB.Query(sqlParam, base.NewIDEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		id = res[0].(*base.IDEntity)
	}
	return
}
