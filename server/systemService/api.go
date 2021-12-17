package systemService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/system/setting/page", base.PowerManageSystemSettingPage, manageSystemSettingPageRequestBean, manageSystemSettingNewBean, manageSystemSettingPageGetSqlParam))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/system/setting/update", base.PowerManageSystemSettingUpdate, manageSystemSettingNewBean, manageSystemSettingUpdateGetTableBean))

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/system/log/page", base.PowerManageSystemLogPage, manageSystemLogPageRequestBean, manageSystemLogNewBean, manageSystemLogPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/system/log/insert", base.PowerManageSystemLogInsert, manageSystemLogNewBean, manageSystemLogInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/system/log/update", base.PowerManageSystemLogUpdate, manageSystemLogNewBean, manageSystemLogUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/system/log/delete", base.PowerManageSystemLogDelete, manageSystemLogNewBean, manageSystemLogDeleteGetSqlParam))

	return
}

func manageSystemSettingNewBean() (res interface{}) {
	return
}
func manageSystemSettingPageRequestBean() (res interface{}) {
	return
}
func manageSystemSettingPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}
func manageSystemSettingUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageSystemLogNewBean() (res interface{}) {
	return
}
func manageSystemLogPageRequestBean() (res interface{}) {
	return
}
func manageSystemLogPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}
func manageSystemLogInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}
func manageSystemLogUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}
func manageSystemLogDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
