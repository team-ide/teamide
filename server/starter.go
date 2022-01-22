package server

import (
	"teamide/server/web"

	"teamide/server/factory"
	applicationService "teamide/server/service/application"
	certificateService "teamide/server/service/certificate"
	enterpriseService "teamide/server/service/enterprise"
	groupService "teamide/server/service/group"
	idService "teamide/server/service/id"
	installService "teamide/server/service/install"
	jobService "teamide/server/service/job"
	logService "teamide/server/service/log"
	loginService "teamide/server/service/login"
	messageService "teamide/server/service/message"
	organizationService "teamide/server/service/organization"
	powerService "teamide/server/service/power"
	spaceService "teamide/server/service/space"
	systemService "teamide/server/service/system"
	toolboxService "teamide/server/service/toolbox"
	userService "teamide/server/service/user"
	wbsService "teamide/server/service/wbs"
	workspaceService "teamide/server/service/workspace"
)

func init() {
	factory.BindInstallService(&installService.Service{})
	factory.BindIdService(&idService.Service{})
	factory.BindSystemService(&systemService.Service{})
	factory.BindLogService(&logService.Service{})
	factory.BindUserService(&userService.Service{})
	factory.BindLoginService(&loginService.Service{})
	factory.BindCertificateService(&certificateService.Service{})
	factory.BindEnterpriseService(&enterpriseService.Service{})
	factory.BindGroupService(&groupService.Service{})
	factory.BindJobService(&jobService.Service{})
	factory.BindMessageService(&messageService.Service{})
	factory.BindOrganizationService(&organizationService.Service{})
	factory.BindPowerService(&powerService.Service{})
	factory.BindSpaceService(&spaceService.Service{})
	factory.BindWbsService(&wbsService.Service{})
	factory.BindWbsService(&applicationService.Service{})
	factory.BindWbsService(&workspaceService.Service{})
	factory.BindWbsService(&toolboxService.Service{})

	Init()
}
func Init() {
	factory.InstallService.Install()
}

func Start() (serverUrl string, err error) {
	return web.StartServer()
}
