package installService

import (
	"teamide/internal/server/factory"
)

func (this_ *Service) Install() {

	for _, one := range factory.Installs {
		install(one.GetInstall())
	}

}
