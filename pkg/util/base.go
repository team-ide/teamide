package util

import (
	"crypto/md5"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"strings"
)

//GenerateUUID 生成UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	slice := strings.Split(uuid, "-")
	var uuidNew string
	for _, str := range slice {
		uuidNew += str
	}
	return uuidNew
}

//EncodePassword 加密密码
func EncodePassword(salt string, password string) (res string) {
	res = GetMd5String(salt + password)
	return
}

//GetMd5String 获取MD5字符串
func GetMd5String(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		log.Fatal(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

var (
	TempDir = ""
)
