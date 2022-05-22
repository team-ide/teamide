package toolbox

import (
	"io/ioutil"
	"time"
)

//GetNowTime 获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}

var (
	TempDir = ""
)

//GetTempDir 获取临时目录
func GetTempDir() (dir string, err error) {
	if TempDir != "" {
		dir = TempDir
		return
	}
	dir, err = ioutil.TempDir("toolbox/temp", "temp")
	return
}
