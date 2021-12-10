package installService

import (
	"server/enterpriseService"
	"server/groupService"
	"server/idService"
	"server/jobService"
	"server/logService"
	"server/messageService"
	"server/powerService"
	"server/settingService"
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
	Install(jobService.GetInstall())
	Install(powerService.GetInstall())
	Install(settingService.GetInstall())
	Install(spaceService.GetInstall())
	Install(systemService.GetInstall())
	Install(messageService.GetInstall())
	Install(groupService.GetInstall())
}
