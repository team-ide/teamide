package userService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserProfileApi(appendApi)
	bindUserAuthApi(appendApi)
	bindUserPasswordApi(appendApi)
	bindUserMessageApi(appendApi)
	bindUserCertificateApi(appendApi)
	bindUserSettingApi(appendApi)

	bindManageUserApi(appendApi)
	bindManageUserPasswordApi(appendApi)
	bindManageUserAuthApi(appendApi)
	bindManageUserLockApi(appendApi)

}
