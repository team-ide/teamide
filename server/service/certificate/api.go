package certificateService

import (
	"server/base"
)

func (this_ *CertificateService) BindApi(appendApi func(apis ...*base.ApiWorker)) {
	bindUserCertificateApi(appendApi)

	bindManageCertificateApi(appendApi)
}
