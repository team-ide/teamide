package enterpriseService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/enterprise/page", base.PowerManageEnterprisePage, manageEnterprisePageRequestBean, manageEnterpriseNewBean, manageEnterprisePageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/enterprise/insert", base.PowerManageEnterpriseInsert, manageEnterpriseNewBean, manageEnterpriseInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/enterprise/update", base.PowerManageEnterpriseUpdate, manageEnterpriseNewBean, manageEnterpriseUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/enterprise/delete", base.PowerManageEnterpriseDelete, manageEnterpriseNewBean, manageEnterpriseDeleteGetSqlParam))

	return
}
func manageEnterpriseNewBean() (res interface{}) {
	return
}
func manageEnterprisePageRequestBean() (res interface{}) {
	return
}
func manageEnterprisePageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageEnterpriseInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageEnterpriseUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageEnterpriseDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
