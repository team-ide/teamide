package base

import (
	"reflect"
)

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

func IsPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}
