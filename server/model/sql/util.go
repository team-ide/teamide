package sql

import "strings"

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

func GetColumnValue(data map[string]interface{}, name string, valueScript string) interface{} {
	var res interface{}
	if IsEmpty(valueScript) {
		res = data[name]
	} else {
		res = GetScriptValue(data, valueScript)
	}
	return res
}

func IfScriptValue(data map[string]interface{}, ifScript string) bool {
	if IsEmpty(ifScript) {
		return true
	}
	value := GetScriptValue(data, ifScript)
	if value == nil {
		return false
	}
	if value == true || value == "1" || value == "true" {
		return true
	}
	return false
}

func GetScriptValue(data map[string]interface{}, script string) interface{} {
	var res interface{}
	if IsEmpty(script) {
		res = nil
	} else {
		res = script
	}
	return res
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
