package jobService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindManageJobApi(appendApi)
}
