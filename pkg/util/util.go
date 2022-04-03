package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Chain-Zhang/pinyin"
)

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

func FormatPath(path string) string {

	var abs string
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	res := filepath.ToSlash(abs)
	return res
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

//GetNowTime 获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//Now 获取当前时间
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

func MatchEmail(email string) bool {
	pattern := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func MatchPhone(phone string) bool {
	pattern := `^\d+$` //匹配手机
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

func MatchNumber(num string) bool {
	pattern := `^\d+$` //匹配数字
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(num)
}

func MatchAccount(account string) bool {
	pattern := `^[a-zA-Z0-9_]+$` //匹配账号
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(account)
}

// 获取随机数
func RandInt(min int, max int) int {
	if max < min {
		panic(fmt.Sprint("RandInt error,min:", min, ",max:", max))
	}
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func GetPingYing(str string) (string, error) {
	// 	InitialsInCapitals: 首字母大写, 不带音调
	// WithoutTone: 全小写,不带音调
	// Tone: 全小写带音调
	str, err := pinyin.New(str).Split("").Mode(pinyin.WithoutTone).Convert()
	if err != nil {
		// 错误处理
		return "", err
	}
	return str, nil
}

func IsZero(value interface{}) (isZero bool) {
	zero := reflect.Zero(reflect.TypeOf(value)).Interface()
	isZero = reflect.DeepEqual(value, zero)
	return
}

var (
	lockMapLock = sync.Mutex{}
	lockMap     = make(map[string]*sync.Mutex)
)

func GetLock(key string) (lock *sync.Mutex) {

	lockMapLock.Lock()

	defer lockMapLock.Unlock()

	var ok bool
	lock, ok = lockMap[key]
	if ok {
		return
	}
	lock = &sync.Mutex{}
	lockMap[key] = lock
	return lock
}

func GetColumnNameByType(fieldType reflect.StructField) string {
	name := fieldType.Tag.Get("column") // field tag
	return name
}

func GetJsonNameByType(fieldType reflect.StructField) string {
	name := fieldType.Tag.Get("json") // field tag
	if strings.Index(name, ",") > 0 {
		name = name[0:strings.Index(name, ",")]
	}
	return name
}

type ColumnType struct {
	Column    string
	FieldType *reflect.StructField
	Value     *reflect.Value
}

func GetColumnTypes(bean interface{}) map[string]ColumnType {
	refType := GetRefType(bean)
	refValue := GetRefValue(bean)

	fieldCount := refType.NumField() // field count
	columnTypes := map[string]ColumnType{}
	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		fieldValue := refValue.Field(i)
		columnName := GetColumnNameByType(fieldType) // field tag
		columnTypes[columnName] = ColumnType{
			Column:    columnName,
			FieldType: &fieldType,
			Value:     &fieldValue,
		}
	}
	return columnTypes
}

func GetFieldTypeValue(t reflect.Type, v reflect.Value) interface{} {
	var value interface{}
	switch t.Name() {
	case "string":
		value = v.String()
	case "bool":
		value = v.Bool()
	case "float32", "float64":
		value = v.Float()
	case "Time", "time.Time":
		t := v.Interface().(time.Time)
		if t.UnixNano() <= 0 {
			value = nil
		} else {
			value = t
		}
	case "int", "int8", "int16", "int32", "int64":
		value = v.Int()
	default:
		if !v.IsNil() {
			value = v.Interface()
		}
	}
	return value
}

func GetStringValue(bean interface{}) string {
	var value string
	if bean == nil {
		return value
	}
	refValue := GetRefValue(bean)
	if refValue.IsZero() {
		return value
	}
	switch refValue.Type().Name() {
	case "string":
		value = refValue.String()
	case "bool":
		value = fmt.Sprint(refValue.Bool())
	case "float32", "float64":
		v := refValue.Float()
		if v != 0. {
			value = fmt.Sprint(v)
		}
	case "Time", "time.Time":
		t := refValue.Interface().(time.Time)
		if t.UnixNano() > 0 {
			value = Format(t)
		}
	case "int", "int8", "int16", "int32", "int64":
		v := refValue.Int()
		if v != 0 {
			value = fmt.Sprint(v)
		}
	default:
		if !refValue.IsNil() {
			v := refValue.Interface()
			value = fmt.Sprint(v)
		}
	}
	return value
}

func GetRefValue(bean interface{}) reflect.Value {
	if IsPtr(bean) {
		return reflect.ValueOf(bean).Elem()
	}
	return reflect.ValueOf(bean)
}

func BeanToMap(bean interface{}) (res map[string]interface{}) {
	if bean == nil {
		return nil
	}

	res = map[string]interface{}{}

	refType := GetRefType(bean)
	refValue := GetRefValue(bean)

	fieldCount := refType.NumField() // field count
	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		jsonName := GetJsonNameByType(fieldType)
		if jsonName == "" {
			continue
		}
		fieldValue := refValue.Field(i) // field vlaue
		value := GetFieldTypeValue(fieldType.Type, fieldValue)
		switch fieldType.Type.Name() {
		case "string", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool", "Time", "time.Time":
			res[jsonName] = value
		default:
			if fieldType.Type.Kind() == reflect.Array || fieldType.Type.Kind() == reflect.Slice {
				len := fieldValue.Len()
				values := []map[string]interface{}{}
				if len > 0 {
					var i int
					for i = 0; i < len; i++ {
						oneValue := fieldValue.Index(i).Interface()
						oneMap := BeanToMap(oneValue)
						values = append(values, oneMap)
					}
				}
				res[jsonName] = values
			} else if fieldType.Type.Kind() == reflect.Map {
				res[jsonName] = value
			} else {
				oneMap := BeanToMap(value)
				res[jsonName] = oneMap
			}

		}
	}
	return
}
