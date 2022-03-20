package idService

import (
	base2 "teamide/internal/server/base"
	component2 "teamide/internal/server/component"
)

func (this_ *Service) GetID(idType component2.IDType) (id int64, err error) {

	var ids []int64
	ids, err = this_.GetIDs(idType, 1)

	if err != nil {
		return
	}
	id = ids[0]

	return
}

func (this_ *Service) GetIDs(idType component2.IDType, size int64) (ids []int64, err error) {

	var key = component2.GetIDRedisKey(idType)
	var exists bool
	exists, err = component2.Redis.Exists(key)
	if err != nil {
		return
	}
	if !exists {
		var unlock func() (err error)
		unlock, err = component2.Redis.Lock(key+":lock", 10, 1000)
		if err != nil {
			return
		}
		defer unlock()

		exists, err = component2.Redis.Exists(key)
		if err != nil {
			return
		}
		if !exists {
			var idInfo *base2.IDEntity
			idInfo, err = get(idType)
			if err != nil {
				return
			}
			var id_ int64 = component2.GetBaseID()
			if idInfo == nil {
				idInfo = &base2.IDEntity{
					Type: int8(idType),
					Id:   0,
				}
			} else {
				id_ = idInfo.Id + 100
			}
			err = component2.Redis.SetInt64(key, id_)
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
	maxId, err = component2.Redis.IncrBy(key, size)
	if err != nil {
		return
	}
	var minId = maxId - size + 1
	ids = []int64{}
	var id int64
	for id = minId; id <= maxId; id++ {
		if id%50 == 0 {
			idInfo := &base2.IDEntity{
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

func insertOrUpdate(id base2.IDEntity) (err error) {
	sql := "INSERT INTO " + TABLE_ID + " (type, id) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=?"
	params := []interface{}{id.Type, id.Id, id.Id}

	sqlParam := base2.NewSqlParam(sql, params)

	_, err = component2.DB.Exec(sqlParam)

	if err != nil {
		return
	}
	return
}

//查询单个ID
func get(idType component2.IDType) (id *base2.IDEntity, err error) {
	sql := "SELECT * FROM " + TABLE_ID + " WHERE type=? "
	params := []interface{}{idType}

	sqlParam := base2.NewSqlParam(sql, params)

	var res []interface{}
	res, err = component2.DB.Query(sqlParam, base2.NewIDEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		id = res[0].(*base2.IDEntity)
	}
	return
}
