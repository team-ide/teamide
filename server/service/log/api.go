package logService

import (
	"server/base"
)

func (this_ *LogService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserLogApi(appendApi)

	bindManageLogApi(appendApi)
}
