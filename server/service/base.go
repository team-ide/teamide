package service

import (
	"server/base"
	"server/db"
)

func Insert(table string, one interface{}) (err error) {

	sqlParam := db.InsertSqlByBean(table, one)

	_, err = db.DBService.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func BatchInsert(table string, list []interface{}) (err error) {

	sqlParam := db.InsertSqlByBean(table, list...)

	_, err = db.DBService.Insert(sqlParam)

	if err != nil {
		return
	}
	return
}

func Query(table string, one interface{}, newBean func() interface{}) (users []*base.UserEntity, err error) {
	sql := "SELECT * FROM " + table + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := db.NewSqlParam(sql, params)

	AppendWhere(one, &sqlParam)

	var res []interface{}
	_, err = db.DBService.Query(sqlParam, newBean)

	if err != nil {
		return
	}
	users = []*base.UserEntity{}
	for _, one := range res {
		user := one.(*base.UserEntity)
		users = append(users, user)
	}
	return
}

func Count(table string, one interface{}, newBean func() interface{}) (count int64, err error) {
	sql := "SELECT COUNT(*) FROM " + table + " WHERE 1=1 "
	params := []interface{}{}

	sqlParam := db.NewSqlParam(sql, params)

	AppendWhere(one, &sqlParam)

	count, err = db.DBService.Count(sqlParam)
	if err != nil {
		return
	}
	return
}
func AppendWhere(one interface{}, sqlParam *db.SqlParam) {

}
