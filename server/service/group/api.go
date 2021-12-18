package groupService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserGroupApi(appendApi)

	bindManageGroupApi(appendApi)
}
