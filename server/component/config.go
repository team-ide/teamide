package component

import (
	"fmt"
	"os"
	"teamide/server/base"
	"teamide/server/config"
	"teamide/util"
)

func GetUserApps(jwt *base.JWTBean) (appsDir string) {
	appsDir = fmt.Sprint(config.Config.Server.Data, "/apps/", jwt.ServerId, "-", jwt.UserId)

	appsDir = util.GetAbsolutePath(appsDir)

	var exist bool
	exist, _ = util.PathExists(appsDir)
	if !exist {
		os.MkdirAll(appsDir, 0777)
	}
	return
}
