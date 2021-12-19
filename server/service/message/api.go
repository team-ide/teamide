package messageService

import (
	"server/base"
)

func (this_ *MessageService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserMessageApi(appendApi)

	bindManageMessageApi(appendApi)
}
