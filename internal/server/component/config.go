package component

import (
	"fmt"
	"os"
	"teamide/internal/server/base"
	"teamide/internal/server/config"
	"teamide/pkg/util"
)

func GetUserApps(jwt *base.JWTBean) (appsDir string) {
	appsDir = fmt.Sprint(config.Config.Server.Data, "/apps/", jwt.UserId)

	appsDir = util.GetAbsolutePath(appsDir)

	var exist bool
	exist, _ = util.PathExists(appsDir)
	if !exist {
		os.MkdirAll(appsDir, 0777)
	}
	return
}
