package service

import (
	"base"
	"db"
	"redis"
)

func GetID(idType base.IDType) (id int64, err error) {

	var ids []int64
	ids, err = GetIDs(idType, 1)

	if err != nil {
		return
	}
	id = ids[0]

	return
}

func GetIDs(idType base.IDType, size int64) (ids []int64, err error) {

	var key = base.GetIDRedisKey(idType)
	var exists bool
	exists, err = redis.RedisService.Exists(key)
	if err != nil {
		return
	}
	if !exists {
		var unlock func() (err error)
		unlock, err = redis.RedisService.Lock(key+":lock", 10, 1000)
		if err != nil {
			return
		}
		defer unlock()

		exists, err = redis.RedisService.Exists(key)
		if err != nil {
			return
		}
		if !exists {
			var idInfo *base.IDEntity
			idInfo, err = IDGet(idType)
			if err != nil {
				return
			}
			var id_ int64 = 0
			if idInfo == nil {
				idInfo = &base.IDEntity{
					Type: int8(idType),
					Id:   0,
				}
			} else {
				id_ = idInfo.Id + 100
			}
			err = redis.RedisService.SetInt64(key, id_)
			if err != nil {
				return
			}
			idInfo.Id = id_
			err = IDInsertOrUpdate(*idInfo)
			if err != nil {
				return
			}
		}
		unlock()
	}
	var maxId int64
	maxId, err = redis.RedisService.IncrBy(key, size)
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
			err = IDInsertOrUpdate(*idInfo)
			if err != nil {
				return
			}
		}
		ids = append(ids, id)
	}

	return
}
func IDInsert(id base.IDEntity) (err error) {

	sqlParam := db.InsertSqlByBean(db.TABLE_ID, id)

	// fmt.Println("IDInsert:", base.ToJSON(sqlParam))
	_, err = db.DBService.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func IDBatchInsert(ids []interface{}) (err error) {

	sqlParam := db.InsertSqlByBean(db.TABLE_ID, ids...)

	_, err = db.DBService.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func IDUpdate(id base.IDEntity) (err error) {
	sql := "UPDATE " + db.TABLE_ID + " SET id=? WHERE type=? "
	params := []interface{}{id.Id, id.Type}

	sqlParam := db.NewSqlParam(sql, params)

	// fmt.Println("IDUpdate:", base.ToJSON(sqlParam))
	_, err = db.DBService.Exec(sqlParam)

	if err != nil {
		return
	}
	return
}

func IDInsertOrUpdate(id base.IDEntity) (err error) {
	sql := "INSERT INTO " + db.TABLE_ID + " (type, id) VALUES (?, ?) ON DUPLICATE KEY UPDATE id=?"
	params := []interface{}{id.Type, id.Id, id.Id}

	sqlParam := db.NewSqlParam(sql, params)

	// fmt.Println("IDUpdate:", base.ToJSON(sqlParam))
	_, err = db.DBService.Exec(sqlParam)

	if err != nil {
		return
	}
	return
}

//查询单个ID
func IDGet(idType base.IDType) (id *base.IDEntity, err error) {
	sql := "SELECT * FROM " + db.TABLE_ID + " WHERE type=? "
	params := []interface{}{idType}

	sqlParam := db.NewSqlParam(sql, params)

	var res []interface{}
	res, err = db.DBService.Query(sqlParam, base.NewIDEntityInterface)

	if err != nil {
		return
	}
	if len(res) > 0 {
		id = res[0].(*base.IDEntity)
	}
	return
}
