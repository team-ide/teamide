package groupService

import (
	"server/base"
)

func (this_ *GroupService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserGroupApi(appendApi)

	bindManageGroupApi(appendApi)
}
