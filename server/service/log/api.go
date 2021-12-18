package logService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserLogApi(appendApi)

	bindManageLogApi(appendApi)
}
