package spaceService

import (
	"server/base"
)

func (this_ *SpaceService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserSpaceApi(appendApi)

	bindManageSpaceApi(appendApi)
}
