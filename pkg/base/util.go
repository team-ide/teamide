package base

import (
	"encoding/json"
	"reflect"
	"strings"
)

/*
转换为大驼峰命名法则
首字母大写，“_” 忽略后大写
*/
func Marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}

	return s
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

func IsZero(value interface{}) (isZero bool) {
	zero := reflect.Zero(reflect.TypeOf(value)).Interface()
	isZero = reflect.DeepEqual(value, zero)
	return
}

type ColumnType struct {
	Column    string
	FieldType *reflect.StructField
	Value     *reflect.Value
}

func GetRefValue(bean interface{}) reflect.Value {
	if IsPtr(bean) {
		return reflect.ValueOf(bean).Elem()
	}
	return reflect.ValueOf(bean)
}

func GetRefType(bean interface{}) reflect.Type {
	if IsPtr(bean) {
		return reflect.TypeOf(bean).Elem()
	}
	return reflect.TypeOf(bean)
}

func IsPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
