package factory

import (
	"server/base"
	"server/component"
)

type idService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
	GetID(idType component.IDType) (id int64, err error)
	GetIDs(idType component.IDType, size int64) (ids []int64, err error)
}

var (
	IdService idService
)

func BindIdService(service interface{}) {
	IdService = service.(idService)
}

type userService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
	Get(userId int64) (user *base.UserEntity, err error)
	TotalInsert(userTotal *base.UserTotalBean) (err error)
	LoginByAccount(account string, password string) (res *base.UserEntity, err error)
}

var (
	UserService userService
)

func BindUserService(service interface{}) {
	UserService = service.(userService)
}

type installService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
	Install()
}

var (
	InstallService installService
)

func BindInstallService(service interface{}) {
	InstallService = service.(installService)
}

type certificateService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	CertificateService certificateService
)

func BindCertificateService(service interface{}) {
	CertificateService = service.(certificateService)
}

type enterpriseService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	EnterpriseService enterpriseService
)

func BindEnterpriseService(service interface{}) {
	EnterpriseService = service.(enterpriseService)
}

type groupService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	GroupService groupService
)

func BindGroupService(service interface{}) {
	GroupService = service.(groupService)
}

type jobService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	JobService jobService
)

func BindJobService(service interface{}) {
	JobService = service.(jobService)
}

type logService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	LogService logService
)

func BindLogService(service interface{}) {
	LogService = service.(logService)
}

type loginService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	LoginService loginService
)

func BindLoginService(service interface{}) {
	LoginService = service.(loginService)
}

type messageService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	MessageService messageService
)

func BindMessageService(service interface{}) {
	MessageService = service.(messageService)
}

type organizationService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	OrganizationService organizationService
)

func BindOrganizationService(service interface{}) {
	OrganizationService = service.(organizationService)
}

type powerService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	PowerService powerService
)

func BindPowerService(service interface{}) {
	PowerService = service.(powerService)
}

type spaceService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	SpaceService spaceService
)

func BindSpaceService(service interface{}) {
	SpaceService = service.(spaceService)
}

type systemService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	SystemService systemService
)

func BindSystemService(service interface{}) {
	SystemService = service.(systemService)
}

type wbsService interface {
	GetInstall() (info *base.InstallInfo)
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	WbsService wbsService
)

func BindWbsService(service interface{}) {
	WbsService = service.(wbsService)
}
