package enterpriseService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindManageEnterpriseApi(appendApi)
}
