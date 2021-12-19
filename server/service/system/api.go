package systemService

import (
	"server/base"
)

func (this_ *SystemService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindManageSystemSettingApi(appendApi)
	bindManageSystemLogApi(appendApi)
}
