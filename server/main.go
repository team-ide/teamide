package main

import (
	"os"
	"os/signal"
	"server/base"
	"server/component"
	"server/config"
	"server/web"
	"syscall"

	"server/factory"
	certificateService "server/service/certificate"
	enterpriseService "server/service/enterprise"
	groupService "server/service/group"
	idService "server/service/id"
	installService "server/service/install"
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

func main() {
	web.StartServer()
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
