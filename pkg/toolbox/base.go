package toolbox

import (
	"io/ioutil"
	"teamide/pkg/util"
	"time"
)

//GetNowTime 获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//GetTempDir 获取临时目录
func GetTempDir() (dir string, err error) {
	if util.TempDir != "" {
		dir = util.TempDir
		return
	}
	dir, err = ioutil.TempDir("toolbox/temp", "temp")
	return
}
