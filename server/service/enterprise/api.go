package enterpriseService

import (
	"server/base"
)

func (this_ *EnterpriseService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserEnterpriseApi(appendApi)

	bindManageEnterpriseApi(appendApi)
}
