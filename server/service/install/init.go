package installService

import (
	"server/service/enterpriseService"
	"server/service/groupService"
	"server/service/idService"
	"server/service/jobService"
	"server/service/logService"
	"server/service/loginService"
	"server/service/messageService"
	"server/service/organizationService"
	"server/service/powerService"
	"server/service/spaceService"
	"server/service/systemService"
	"server/service/userService"
	"server/service/wbsService"
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
