package groupService

import (
	"server/base"
	"server/baseService"
)

func BindApi() (workers []*base.ApiWorker) {
	workers = append(workers, bindUserApi()...)
	workers = append(workers, bindManageApi()...)
	return
}

func bindUserApi() (workers []*base.ApiWorker) {

	return
}

func bindManageApi() (workers []*base.ApiWorker) {

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/group/page", base.PowerManageGroupPage, manageGroupPageRequestBean, manageGroupNewBean, manageGroupPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/group/insert", base.PowerManageGroupInsert, manageGroupNewBean, manageGroupInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/group/update", base.PowerManageGroupUpdate, manageGroupNewBean, manageGroupUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/group/delete", base.PowerManageGroupDelete, manageGroupNewBean, manageGroupDeleteGetSqlParam))

	return
}
func manageGroupNewBean() (res interface{}) {
	return
}
func manageGroupPageRequestBean() (res interface{}) {
	return
}
func manageGroupPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageGroupInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageGroupUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageGroupDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
