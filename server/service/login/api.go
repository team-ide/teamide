package loginService

import (
	"server/base"
)

func (this_ *LoginService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserLoginApi(appendApi)

	bindManageLoginApi(appendApi)
}
