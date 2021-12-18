package messageService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserMessageApi(appendApi)

	bindManageMessageApi(appendApi)
}
