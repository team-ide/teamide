package installService

import (
	"server/enterpriseService"
	"server/groupService"
	"server/idService"
	"server/jobService"
	"server/logService"
	"server/loginService"
	"server/messageService"
	"server/organizationService"
	"server/powerService"
	"server/spaceService"
	"server/systemService"
	"server/userService"
	"server/wbsService"
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
}
