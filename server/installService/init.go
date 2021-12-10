package installService

import (
	"server/idService"
	"server/userService"
	"server/wbsService"
)

func Init() {
	Install(idService.GetInstall())
	Install(userService.GetInstall())
	Install(wbsService.GetInstall())
}
