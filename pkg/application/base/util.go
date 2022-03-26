package base

import (
	"encoding/json"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

//生成UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	slice := strings.Split(uuid, "-")
	var uuidNew string
	for _, str := range slice {
		uuidNew += str
	}
	return uuidNew
}

//获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//获取当前时间
func Now() time.Time {
	return time.Now()
}

func NowStr() string {
	now := time.Now()
	return Format(now)
}

func Format(date time.Time) string {
	return date.Format("2006-01-02 15:04:05.000")
}

func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func IsEmptyObj(obj interface{}) bool {
	if obj == nil || obj == "" || obj == 0 {
		return true
	}
	return false
}

func IsNotEmptyObj(obj interface{}) bool {
	return !IsEmptyObj(obj)
}

func ToJSON(data interface{}) string {
	if data != nil {
		bs, _ := json.Marshal(data)
		if bs != nil {
			return string(bs)
		}
	}
	return ""
}

func ToBean(bytes []byte, req interface{}) (err error) {
	err = json.Unmarshal(bytes, req)
	return
}

func YamlToBean(bytes []byte, req interface{}) (err error) {
	err = yaml.Unmarshal(bytes, req)
	return
}

func YamlToMap(bytes []byte) (res yaml.MapSlice, err error) {
	res = yaml.MapSlice{}
	err = yaml.Unmarshal(bytes, &res)
	return
}

func ToInt64(value interface{}) (res int64, ok bool) {
	if value == nil {
		return
	}
	refType := GetRefType(value)
	refValue := GetRefValue(value)
	switch refType.Name() {
	case "int", "int8", "int16", "int32", "int64":
		ok = true
		res = refValue.Int()
	}
	return
}

func ToFloat64(value interface{}) (res float64, ok bool) {
	if value == nil {
		return
	}
	refType := GetRefType(value)
	refValue := GetRefValue(value)
	switch refType.Name() {
	case "float32", "float64":
		ok = true
		res = refValue.Float()
	}
	return
}
func IsPtr(obj interface{}) bool {

	return reflect.ValueOf(obj).Kind() == reflect.Ptr
}

func AppendLine(content *string, line string, tab int) {
	for i := 0; i < tab; i++ {
		*content += "    "
	}
	*content += line
	*content += "\n"
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetRefType(obj interface{}) reflect.Type {
	if IsPtr(obj) {
		return reflect.TypeOf(obj).Elem()
	}
	return reflect.TypeOf(obj)
}

func IsInt(obj interface{}) bool {
	refType := GetRefType((obj))
	switch refType.Name() {
	case "int", "int8", "int16", "int32", "int64":
		return true
	}
	return false
}

func IsFloat(obj interface{}) bool {
	refType := GetRefType((obj))
	switch refType.Name() {
	case "float32", "float64":
		return true
	}
	return false
}

func GetRefValue(obj interface{}) reflect.Value {
	if IsPtr(obj) {
		return reflect.ValueOf(obj).Elem()
	}
	return reflect.ValueOf(obj)
}

func GetFieldTypeValue(t reflect.Type, v reflect.Value) interface{} {
	switch t.Name() {
	case "string":
		return v.String()
	case "bool":
		return v.Bool()
	case "float32", "float64":
		return v.Float()
	case "Time", "time.Time":
		t := v.Interface().(time.Time)
		if t.UnixNano() <= 0 {
			return nil
		} else {
			return t
		}
	case "int", "int8", "int16", "int32", "int64":
		return v.Int()
	default:
		if !v.IsNil() {
			return v.Interface()
		}
	}
	return nil
}
