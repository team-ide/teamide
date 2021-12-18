package logService

import (
	"server/base"
	"server/service/baseService"
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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/log/page", base.PowerManageLogPage, manageLogPageRequestBean, manageLogNewBean, manageLogPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/log/insert", base.PowerManageLogInsert, manageLogNewBean, manageLogInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/log/update", base.PowerManageLogUpdate, manageLogNewBean, manageLogUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/log/delete", base.PowerManageLogDelete, manageLogNewBean, manageLogDeleteGetSqlParam))

	return
}

func manageLogNewBean() (res interface{}) {
	return
}
func manageLogPageRequestBean() (res interface{}) {
	return
}
func manageLogPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageLogInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageLogUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageLogDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
