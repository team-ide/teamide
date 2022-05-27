package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
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

//GetTempDir 获取临时目录
func GetTempDir() (dir string, err error) {
	if TempDir != "" {
		dir = TempDir
		return
	}
	dir, err = ioutil.TempDir("toolbox/temp", "temp")
	return
}

func CopyBytes(dst io.Writer, src io.Reader, call func(readSize int64, writeSize int64)) (err error) {
	var buf = make([]byte, 32*1024)
	var errInvalidWrite = errors.New("invalid write result")
	var ErrShortWrite = errors.New("short write")
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			call(int64(nr), 0)
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			call(0, int64(nw))
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return
}
