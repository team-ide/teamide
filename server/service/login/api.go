package loginService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserLoginApi(appendApi)

	bindManageLoginApi(appendApi)
}
