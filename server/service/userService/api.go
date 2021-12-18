package userService

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

	workers = append(workers, baseService.NewQueryOneApiWorker("/user/profile/page", base.PowerUserProfilePage, nil, userNewBean, userProfileQueryOne))
	workers = append(workers, baseService.NewExecApiWorker("/user/profile/update", base.PowerUserProfileUpdate, userNewBean, userProfileUpdateGetSqlParam))

	workers = append(workers, baseService.NewQueryOneApiWorker("/user/security/page", base.PowerUserSecurityPage, nil, userNewBean, userSecurityQueryOne))
	workers = append(workers, baseService.NewExecApiWorker("/user/security/update", base.PowerUserSecurityUpdate, userNewBean, userSecurityUpdateGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/user/auth/page", base.PowerUserAuthPage, userAuthPageRequestBean, userAuthNewBean, userAuthPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/user/auth/insert", base.PowerUserAuthInsert, userAuthNewBean, userAuthInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/user/auth/update", base.PowerUserAuthUpdate, userAuthNewBean, userAuthUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/delete", base.PowerUserAuthDelete, userAuthNewBean, userAuthDeleteGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/active", base.PowerUserAuthActive, userAuthNewBean, userAuthActiveGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/lock", base.PowerUserAuthLock, userAuthNewBean, userAuthLockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/unlock", base.PowerUserAuthUnlock, userAuthNewBean, userAuthUnlockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/enable", base.PowerUserAuthEnable, userAuthNewBean, userAuthEnableGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/user/auth/disable", base.PowerUserAuthDisable, userAuthNewBean, userAuthDisableGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/user/certificate/page", base.PowerUserCertificatePage, userAuthPageRequestBean, userAuthNewBean, userAuthPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/user/certificate/insert", base.PowerUserCertificateInsert, userAuthNewBean, userAuthInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/user/certificate/update", base.PowerUserCertificateUpdate, userAuthNewBean, userAuthUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/user/certificate/delete", base.PowerUserCertificateDelete, userAuthNewBean, userAuthDeleteGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/user/message/page", base.PowerUserMessagePage, userAuthPageRequestBean, userAuthNewBean, userAuthPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/user/message/insert", base.PowerUserMessageInsert, userAuthNewBean, userAuthInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/user/message/update", base.PowerUserMessageUpdate, userAuthNewBean, userAuthUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/user/message/delete", base.PowerUserMessageDelete, userAuthNewBean, userAuthDeleteGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/user/setting/page", base.PowerUserSettingPage, userAuthPageRequestBean, userAuthNewBean, userAuthPageGetSqlParam))
	workers = append(workers, baseService.NewUpdateApiWorker("/user/setting/update", base.PowerUserSettingUpdate, userAuthNewBean, userAuthUpdateGetTableBean))
	return
}

func bindManageApi() (workers []*base.ApiWorker) {

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/user/page", base.PowerManageUserPage, manageUserPageRequestBean, manageUserNewBean, manageUserPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/user/insert", base.PowerManageUserInsert, manageUserNewBean, manageUserInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/user/update", base.PowerManageUserUpdate, manageUserNewBean, manageUserUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/delete", base.PowerManageUserDelete, manageUserNewBean, manageUserDeleteGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/active", base.PowerManageUserActive, manageUserNewBean, manageUserActiveGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/lock", base.PowerManageUserLock, manageUserNewBean, manageUserLockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/unlock", base.PowerManageUserUnlock, manageUserNewBean, manageUserUnlockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/enable", base.PowerManageUserEnable, manageUserNewBean, manageUserEnableGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/disable", base.PowerManageUserDisable, manageUserNewBean, manageUserDisableGetSqlParam))

	workers = append(workers, baseService.NewExecApiWorker("/manage/user/password/reset", base.PowerManageUserPasswordReset, manageUserNewBean, manageUserPasswordResetGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/password/update", base.PowerManageUserPasswordUpdate, manageUserNewBean, manageUserPasswordUpdateGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/user/auth/page", base.PowerManageUserAuthPage, manageUserAuthPageRequestBean, manageUserAuthNewBean, manageUserAuthPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/user/auth/insert", base.PowerManageUserAuthInsert, manageUserAuthNewBean, manageUserAuthInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/user/auth/update", base.PowerManageUserAuthUpdate, manageUserAuthNewBean, manageUserAuthUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/delete", base.PowerManageUserAuthDelete, manageUserAuthNewBean, manageUserAuthDeleteGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/active", base.PowerManageUserAuthActive, manageUserAuthNewBean, manageUserAuthActiveGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/lock", base.PowerManageUserAuthLock, manageUserAuthNewBean, manageUserAuthLockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/unlock", base.PowerManageUserAuthUnlock, manageUserAuthNewBean, manageUserAuthUnlockGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/enable", base.PowerManageUserAuthEnable, manageUserAuthNewBean, manageUserAuthEnableGetSqlParam))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/auth/disable", base.PowerManageUserAuthDisable, manageUserAuthNewBean, manageUserAuthDisableGetSqlParam))

	workers = append(workers, baseService.NewQueryPageApiWorker("/manage/user/lock/page", base.PowerManageUserLockPage, manageUserLockPageRequestBean, manageUserLockNewBean, manageUserLockPageGetSqlParam))
	workers = append(workers, baseService.NewInsertApiWorker("/manage/user/lock/insert", base.PowerManageUserLockInsert, manageUserLockNewBean, manageUserLockInsertGetTableBean))
	workers = append(workers, baseService.NewUpdateApiWorker("/manage/user/lock/update", base.PowerManageUserLockUpdate, manageUserLockNewBean, manageUserLockUpdateGetTableBean))
	workers = append(workers, baseService.NewExecApiWorker("/manage/user/lock/delete", base.PowerManageUserLockDelete, manageUserLockNewBean, manageUserLockDeleteGetSqlParam))

	return
}
func userNewBean() (res interface{}) {

	return
}

func userProfileQueryOne(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
func userProfileUpdateGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userSecurityQueryOne(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
func userSecurityUpdateGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthNewBean() (res interface{}) {
	return
}
func userAuthPageRequestBean() (res interface{}) {
	return
}
func userAuthPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func userAuthInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func userAuthUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func userAuthDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthActiveGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthLockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthUnlockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthDisableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func userAuthEnableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserNewBean() (res interface{}) {
	return
}
func manageUserPageRequestBean() (res interface{}) {
	return
}
func manageUserPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageUserInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageUserUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageUserDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserActiveGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserLockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserUnlockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserDisableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserEnableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserPasswordResetGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserPasswordUpdateGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthNewBean() (res interface{}) {
	return
}
func manageUserAuthPageRequestBean() (res interface{}) {
	return
}
func manageUserAuthPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageUserAuthUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageUserAuthDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthActiveGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthLockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthUnlockGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthDisableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserAuthEnableGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}

func manageUserLockNewBean() (res interface{}) {
	return
}
func manageUserLockPageRequestBean() (res interface{}) {
	return
}
func manageUserLockPageGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, countSqlParam base.SqlParam, err error) {
	return
}

func manageUserLockInsertGetTableBean(requestBean interface{}) (table string, bean interface{}, err error) {
	return
}

func manageUserLockUpdateGetTableBean(requestBean interface{}) (table string, keys []string, bean interface{}, err error) {
	return
}

func manageUserLockDeleteGetSqlParam(requestBean interface{}) (sqlParam base.SqlParam, err error) {
	return
}
