package util

import (
	"encoding/json"
	"reflect"
	"strings"
)

func ToStruct(value interface{}, structValue interface{}) (err error) {
	if value == nil {
		return
	}
	// 序列化
	arr, err := json.Marshal(value)
	if err != nil {
		return
	}
	// 反序列化
	err = json.Unmarshal(arr, &structValue)
	if err != nil {
		return
	}
	return
}

func GetStructFieldTypes(struct_ interface{}) map[string]string {
	refType := GetRefType(struct_)

	fieldCount := refType.NumField() // field count
	res := map[string]string{}
	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i)                // field type
		columnName := GetNameByFieldTType(fieldType) // field tag
		typeName := fieldType.Type.Name()
		res[columnName] = typeName
		res[fieldType.Name] = typeName
	}
	return res
}

func GetNameByFieldTType(fieldType reflect.StructField) string {
	name := fieldType.Tag.Get("json") // field tag
	if strings.Index(name, ",") > 0 {
		name = name[0:strings.Index(name, ",")]
	}
	return name
}

func GetRefType(v interface{}) reflect.Type {
	if IsPtr(v) {
		return reflect.TypeOf(v).Elem()
	}
	return reflect.TypeOf(v)
}

func IsPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
