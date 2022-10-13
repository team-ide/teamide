package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func IsEmpty(obj interface{}) (res bool) { // 是否为空  null或空字符串
	if obj == nil || obj == "" || obj == 0 {
		res = true
	}
	return
}

func IsNotEmpty(obj interface{}) (res bool) { // 是否不为空  取 是否为空 反
	return !IsEmpty(obj)
}

func IsTrue(obj interface{}) (res bool) { // 是否为真 入：true、"true"、1、"1"为真
	if obj == true || obj == "true" || obj == 1 || obj == "1" {
		res = true
	}
	return
}

func IsFalse(obj interface{}) (res bool) { // 是否为假  取 是否为真 反
	return !IsTrue(obj)
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func AppendLine(content *string, line string, tab int) {
	for i := 0; i < tab; i++ {
		*content += "    "
	}
	*content += line
	*content += "\n"
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
