package powerService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindManagePowerRoleApi(appendApi)
	bindManagePowerActionApi(appendApi)
	bindManagePowerDataApi(appendApi)
}
