package installService

import "server/factory"

func (this_ *InstallService) Install() {

	install(factory.IdService.GetInstall())
	install(factory.UserService.GetInstall())
	install(factory.WbsService.GetInstall())
	install(factory.LogService.GetInstall())
	install(factory.EnterpriseService.GetInstall())
	install(factory.OrganizationService.GetInstall())
	install(factory.JobService.GetInstall())
	install(factory.PowerService.GetInstall())
	install(factory.SpaceService.GetInstall())
	install(factory.LoginService.GetInstall())
	install(factory.SystemService.GetInstall())
	install(factory.MessageService.GetInstall())
	install(factory.GroupService.GetInstall())
	install(factory.CertificateService.GetInstall())

}
