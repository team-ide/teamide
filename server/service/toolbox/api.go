package toolboxService

import (
	"teamide/server/base"
)

func (this_ *Service) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	// appendApi(&base.ApiWorker{Apis: []string{"toolbox/list"}, Power: base.PowerApplicationList, Do: apiList})
}
