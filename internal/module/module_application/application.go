package module_application

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
	"teamide/pkg/util"
)

// NewApplicationService 根据库配置创建ApplicationService
func NewApplicationService(ServerContext *context.ServerContext) (res *ApplicationService) {

	idService := module_id.NewIDService(ServerContext)

	res = &ApplicationService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// ApplicationService 工具箱服务
type ApplicationService struct {
	*context.ServerContext
	idService *module_id.IDService
}

func (this_ *ApplicationService) GetUserApps(JWT *base.JWTBean) string {
	appsDir := fmt.Sprintf("%sapps/%d", this_.ServerConfig.Server.Data, JWT.UserId)

	exist, err := util.PathExists(appsDir)
	if err != nil {
		this_.Logger.Error("目录检测异常", zap.Error(err))
		return appsDir
	}
	if !exist {
		err = os.MkdirAll(appsDir, os.ModeDir)
		if err != nil {
			this_.Logger.Error("创建目录异常", zap.Error(err))
		}
	}
	return appsDir
}

func (this_ *ApplicationService) GetAppsDir(requestBean *base.RequestBean) string {
	return this_.GetUserApps(requestBean.JWT)
}

func (this_ *ApplicationService) GetAppPath(requestBean *base.RequestBean, name string) string {
	appPath := this_.GetAppsDir(requestBean) + "/" + name
	appPath = util.GetAbsolutePath(appPath)
	return appPath
}
