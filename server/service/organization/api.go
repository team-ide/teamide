package organizationService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserOrganizationApi(appendApi)

	bindManageOrganizationApi(appendApi)
}
