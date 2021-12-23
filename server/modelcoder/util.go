package modelcoder

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

func WrapTableName(database string, table string) string {
	wrapTable := ""
	if IsNotEmpty(database) {
		wrapTable += "`" + strings.TrimSpace(database) + "`."
	}
	if IsNotEmpty(table) {
		wrapTable += "`" + strings.TrimSpace(table) + "`"
	}
	return wrapTable
}

func WrapColumnName(alias string, name string) string {
	wrapColumn := ""
	if IsNotEmpty(alias) {
		wrapColumn += "`" + strings.TrimSpace(alias) + "`."
	}
	if IsNotEmpty(name) {
		wrapColumn += "`" + strings.TrimSpace(name) + "`"
	}
	return wrapColumn
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

func IsPtr(obj interface{}) bool {

	return reflect.ValueOf(obj).Kind() == reflect.Ptr
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
