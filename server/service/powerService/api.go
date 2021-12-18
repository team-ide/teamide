package powerService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/power/role/page", base.PowerManagePowerRolePage, managePowerRolePageRequestBean, managePowerRoleNewBean, managePowerRolePageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/power/role/insert", base.PowerManagePowerRoleInsert, managePowerRoleNewBean, managePowerRoleInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/power/role/update", base.PowerManagePowerRoleUpdate, managePowerRoleNewBean, managePowerRoleUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/power/role/delete", base.PowerManagePowerRoleDelete, managePowerRoleNewBean, managePowerRoleDeleteGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/power/action/page", base.PowerManagePowerActionPage, managePowerActionPageRequestBean, managePowerActionNewBean, managePowerActionPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/power/action/insert", base.PowerManagePowerActionInsert, managePowerActionNewBean, managePowerActionInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/power/action/update", base.PowerManagePowerActionUpdate, managePowerActionNewBean, managePowerActionUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/power/action/delete", base.PowerManagePowerActionDelete, managePowerActionNewBean, managePowerActionDeleteGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/power/data/page", base.PowerManagePowerDataPage, managePowerDataPageRequestBean, managePowerDataNewBean, managePowerDataPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/power/data/insert", base.PowerManagePowerDataInsert, managePowerDataNewBean, managePowerDataInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/power/data/update", base.PowerManagePowerDataUpdate, managePowerDataNewBean, managePowerDataUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/power/data/delete", base.PowerManagePowerDataDelete, managePowerDataNewBean, managePowerDataDeleteGetSqlParam))

	return
}
func managePowerRoleNewBean() (res interface{}) {
	return
}
func managePowerRolePageRequestBean() (res interface{}) {
	return
}
func managePowerRolePageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func managePowerRoleInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func managePowerRoleUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func managePowerRoleDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func managePowerActionNewBean() (res interface{}) {
	return
}
func managePowerActionPageRequestBean() (res interface{}) {
	return
}
func managePowerActionPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func managePowerActionInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func managePowerActionUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func managePowerActionDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func managePowerDataNewBean() (res interface{}) {
	return
}
func managePowerDataPageRequestBean() (res interface{}) {
	return
}
func managePowerDataPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func managePowerDataInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func managePowerDataUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func managePowerDataDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
