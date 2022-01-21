package installService

import "teamide/server/factory"

func (this_ *InstallService) Install() {

	for _, one := range factory.Installs {
		install(one.GetInstall())
	}

}
