package base

import (
	"github.com/team-ide/go-tool/util"
	"io/ioutil"
)

func AppendLine(content *string, line string, tab int) {
	for i := 0; i < tab; i++ {
		*content += "    "
	}
	*content += line
	*content += "\n"
}

// GetMd5String 获取MD5字符串
func GetMd5String(str string) string {
	return util.GetMD5(str)
}

// EncodePassword 加密密码
func EncodePassword(salt string, password string) (res string) {
	res = util.GetMD5(salt + password)
	return
}

var (
	TempDir = ""
)

// GetTempDir 获取临时目录
func GetTempDir() (dir string, err error) {
	if TempDir != "" {
		dir = TempDir
		return
	}
	dir, err = ioutil.TempDir("teamide/temp", "temp")
	return
}
