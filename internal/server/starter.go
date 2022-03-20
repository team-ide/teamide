package server

import (
	"teamide/internal/server/factory"
	"teamide/internal/server/service/application"
	"teamide/internal/server/service/certificate"
	"teamide/internal/server/service/enterprise"
	"teamide/internal/server/service/group"
	"teamide/internal/server/service/id"
	"teamide/internal/server/service/install"
	"teamide/internal/server/service/job"
	"teamide/internal/server/service/log"
	"teamide/internal/server/service/login"
	"teamide/internal/server/service/message"
	"teamide/internal/server/service/organization"
	"teamide/internal/server/service/power"
	"teamide/internal/server/service/space"
	"teamide/internal/server/service/system"
	"teamide/internal/server/service/toolbox"
	"teamide/internal/server/service/user"
	"teamide/internal/server/service/wbs"
	"teamide/internal/server/service/workspace"
	"teamide/internal/server/web"
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
