package organizationService

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

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/organization/page", base.PowerManageOrganizationPage, manageOrganizationPageRequestBean, manageOrganizationNewBean, manageOrganizationPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/organization/insert", base.PowerManageOrganizationInsert, manageOrganizationNewBean, manageOrganizationInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/organization/update", base.PowerManageOrganizationUpdate, manageOrganizationNewBean, manageOrganizationUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/organization/delete", base.PowerManageOrganizationDelete, manageOrganizationNewBean, manageOrganizationDeleteGetSqlParam))

	return
}
func manageOrganizationNewBean() (res interface{}) {
	return
}
func manageOrganizationPageRequestBean() (res interface{}) {
	return
}
func manageOrganizationPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageOrganizationInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageOrganizationUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageOrganizationDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
