package systemService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindManageSystemSettingApi(appendApi)
	bindManageSystemLogApi(appendApi)
}
