package server

import (
	"os"
	"os/signal"
	"syscall"
	"teamide/server/base"
	"teamide/server/component"
	"teamide/server/config"
	"teamide/server/web"

	"teamide/server/factory"
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
	userService "teamide/server/service/user"
	wbsService "teamide/server/service/wbs"
)

func init() {
	factory.BindInstallService(&installService.InstallService{})
	factory.BindIdService(&idService.IdService{})
	factory.BindSystemService(&systemService.SystemService{})
	factory.BindLogService(&logService.LogService{})
	factory.BindUserService(&userService.UserService{})
	factory.BindLoginService(&loginService.LoginService{})
	factory.BindCertificateService(&certificateService.CertificateService{})
	factory.BindEnterpriseService(&enterpriseService.EnterpriseService{})
	factory.BindGroupService(&groupService.GroupService{})
	factory.BindJobService(&jobService.JobService{})
	factory.BindMessageService(&messageService.MessageService{})
	factory.BindOrganizationService(&organizationService.OrganizationService{})
	factory.BindPowerService(&powerService.PowerService{})
	factory.BindSpaceService(&spaceService.SpaceService{})
	factory.BindWbsService(&wbsService.WbsService{})

	Init()
}
func Init() {
	base.Init()
	config.Init()
	component.Init()
	factory.InstallService.Install()
	// service.CheckModel()
	web.Init()
}

var (
	done = make(chan os.Signal, 1)
)

func StartServer() {
	web.StartServer()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
