package jobService

import (
	"server/base"
)

func (this_ *JobService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserJobApi(appendApi)

	bindManageJobApi(appendApi)
}
