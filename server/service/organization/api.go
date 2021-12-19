package organizationService

import (
	"server/base"
)

func (this_ *OrganizationService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserOrganizationApi(appendApi)

	bindManageOrganizationApi(appendApi)
}
