package base

import (
	"github.com/team-ide/go-tool/util"
)

// GetMd5String 获取MD5字符串
func GetMd5String(str string) string {
	return util.GetMD5(str)
}

// EncodePassword 加密密码
func EncodePassword(salt string, password string) (res string) {
	res = util.GetMD5(salt + password)
	return
}
