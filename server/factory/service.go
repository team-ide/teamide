package factory

import (
	"teamide/server/base"
	"teamide/server/component"
)

type apiInterface interface {
	BindApi(appendApi func(apis ...*base.ApiWorker))
}

var (
	Apis = []apiInterface{}
)

type installInterface interface {
	GetInstall() (info *base.InstallInfo)
}

var (
	Installs = []installInterface{}
)

func BindCommonService(service interface{}) {
	Apis = append(Apis, service.(apiInterface))
	Installs = append(Installs, service.(installInterface))
}

type idService interface {
	GetID(idType component.IDType) (id int64, err error)
	GetIDs(idType component.IDType, size int64) (ids []int64, err error)
}

var (
	IdService idService
)

func BindIdService(service interface{}) {
	IdService = service.(idService)
	BindCommonService(service)
}

type userService interface {
	Get(userId int64) (user *base.UserEntity, err error)
	TotalInsert(userTotal *base.UserTotalBean) (err error)
	LoginByAccount(account string, password string) (res *base.UserEntity, err error)
}

var (
	UserService userService
)

func BindUserService(service interface{}) {
	UserService = service.(userService)
	BindCommonService(service)
}

type installService interface {
	Install()
}

var (
	InstallService installService
)

func BindInstallService(service interface{}) {
	InstallService = service.(installService)
	BindCommonService(service)
}

type certificateService interface {
}

var (
	CertificateService certificateService
)

func BindCertificateService(service interface{}) {
	CertificateService = service.(certificateService)
	BindCommonService(service)
}

type enterpriseService interface {
}

var (
	EnterpriseService enterpriseService
)

func BindEnterpriseService(service interface{}) {
	EnterpriseService = service.(enterpriseService)
	BindCommonService(service)
}

type groupService interface {
}

var (
	GroupService groupService
)

func BindGroupService(service interface{}) {
	GroupService = service.(groupService)
	BindCommonService(service)
}

type jobService interface {
}

var (
	JobService jobService
)

func BindJobService(service interface{}) {
	JobService = service.(jobService)
	BindCommonService(service)
}

type logService interface {
}

var (
	LogService logService
)

func BindLogService(service interface{}) {
	LogService = service.(logService)
	BindCommonService(service)
}

type loginService interface {
}

var (
	LoginService loginService
)

func BindLoginService(service interface{}) {
	LoginService = service.(loginService)
	BindCommonService(service)
}

type messageService interface {
}

var (
	MessageService messageService
)

func BindMessageService(service interface{}) {
	MessageService = service.(messageService)
	BindCommonService(service)
}

type organizationService interface {
}

var (
	OrganizationService organizationService
)

func BindOrganizationService(service interface{}) {
	OrganizationService = service.(organizationService)
	BindCommonService(service)
}

type powerService interface {
}

var (
	PowerService powerService
)

func BindPowerService(service interface{}) {
	PowerService = service.(powerService)
	BindCommonService(service)
}

type spaceService interface {
}

var (
	SpaceService spaceService
)

func BindSpaceService(service interface{}) {
	SpaceService = service.(spaceService)
	BindCommonService(service)
}

type systemService interface {
}

var (
	SystemService systemService
)

func BindSystemService(service interface{}) {
	SystemService = service.(systemService)
	BindCommonService(service)
}

type wbsService interface {
}

var (
	WbsService wbsService
)

func BindWbsService(service interface{}) {
	WbsService = service.(wbsService)
	BindCommonService(service)
}

type applicationService interface {
}

var (
	ApplicationService applicationService
)

func BindApplicationService(service interface{}) {
	ApplicationService = service.(applicationService)
	BindCommonService(service)
}

type workspaceService interface {
}

var (
	WorkspaceService workspaceService
)

func BindWorkspaceService(service interface{}) {
	WorkspaceService = service.(workspaceService)
	BindCommonService(service)
}

type toolboxService interface {
}

var (
	ToolboxService workspaceService
)

func BindToolboxService(service interface{}) {
	ToolboxService = service.(toolboxService)
	BindCommonService(service)
}
