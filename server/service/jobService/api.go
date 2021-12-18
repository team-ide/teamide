package jobService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/job/page", base.PowerManageJobPage, manageJobPageRequestBean, manageJobNewBean, manageJobPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/job/insert", base.PowerManageJobInsert, manageJobNewBean, manageJobInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/job/update", base.PowerManageJobUpdate, manageJobNewBean, manageJobUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/job/delete", base.PowerManageJobDelete, manageJobNewBean, manageJobDeleteGetSqlParam))

	return
}
func manageJobNewBean() (res interface{}) {
	return
}
func manageJobPageRequestBean() (res interface{}) {
	return
}
func manageJobPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageJobInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageJobUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageJobDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
