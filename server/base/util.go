package base

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Chain-Zhang/pinyin"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
)

var (
	//全局的JSON转换对象
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

func ToJSON(data interface{}) string {
	if data != nil {
		bs, _ := JSON.Marshal(data)
		if bs != nil {
			return string(bs)
		}
	}
	return ""
}

func ToBean(bytes []byte, req interface{}) (err error) {
	err = JSON.Unmarshal(bytes, req)
	return
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

//string是否在某个string数组中
func StringInArray(target string, str_array []string) bool {
	if str_array == nil {
		return false
	}
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

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

//复制属性
func CopyProperties(dst interface{}, src interface{}) (err error) {

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
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

func EncodePassword(salt string, password string) (res string) {
	res = GetMd5String(salt + password)
	return
}
func GetMd5String(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		log.Fatal(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
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

func IsPtr(bean interface{}) bool {

	return reflect.ValueOf(bean).Kind() == reflect.Ptr
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
		value = v.Interface()
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

func GetRefType(bean interface{}) reflect.Type {
	if IsPtr(bean) {
		return reflect.TypeOf(bean).Elem()
	}
	return reflect.TypeOf(bean)
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
		fieldType := refType.Field(i)   // field type
		fieldValue := refValue.Field(i) // field vlaue
		jsonName := GetJsonNameByType(fieldType)
		if jsonName == "" {
			continue
		}
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
