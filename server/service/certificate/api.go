package certificateService

import (
	"server/base"
)

func BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserCertificateApi(appendApi)

	bindManageCertificateApi(appendApi)
}
