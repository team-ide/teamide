package loginService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/login/page", base.PowerManageLoginPage, manageLoginPageRequestBean, manageLoginNewBean, manageLoginPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/login/insert", base.PowerManageLoginInsert, manageLoginNewBean, manageLoginInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/login/update", base.PowerManageLoginUpdate, manageLoginNewBean, manageLoginUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/login/delete", base.PowerManageLoginDelete, manageLoginNewBean, manageLoginDeleteGetSqlParam))

	return
}
func manageLoginNewBean() (res interface{}) {
	return
}
func manageLoginPageRequestBean() (res interface{}) {
	return
}
func manageLoginPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageLoginInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageLoginUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageLoginDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
