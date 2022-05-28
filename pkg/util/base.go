package util

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
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

func GetStringValue(value interface{}) (valueString string, err error) {
	if value == nil {
		return "", nil
	}

	switch v := value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case uint:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case uint8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case uint16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case uint32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatInt(int64(v), 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		if v {
			return "1", nil
		}
		return "0", nil
	case time.Time:
		if v.IsZero() {
			return "", nil
		}
		valueString = v.Format("2006-01-02 15:04:05")
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		var bs []byte
		bs, err = json.Marshal(value)
		if err != nil {
			return
		}
		valueString = string(bs)
		break
	}
	return
}
