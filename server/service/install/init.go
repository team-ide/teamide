package installService

import (
	certificateService "server/service/certificate"
	enterpriseService "server/service/enterprise"
	groupService "server/service/group"
	idService "server/service/id"
	jobService "server/service/job"
	logService "server/service/log"
	loginService "server/service/login"
	messageService "server/service/message"
	organizationService "server/service/organization"
	powerService "server/service/power"
	spaceService "server/service/space"
	systemService "server/service/system"
	userService "server/service/user"
	wbsService "server/service/wbs"
)

func Init() {

	Install(idService.GetInstall())
	Install(userService.GetInstall())
	Install(wbsService.GetInstall())
	Install(logService.GetInstall())
	Install(enterpriseService.GetInstall())
	Install(organizationService.GetInstall())
	Install(jobService.GetInstall())
	Install(powerService.GetInstall())
	Install(spaceService.GetInstall())
	Install(loginService.GetInstall())
	Install(systemService.GetInstall())
	Install(messageService.GetInstall())
	Install(groupService.GetInstall())
	Install(certificateService.GetInstall())

}
