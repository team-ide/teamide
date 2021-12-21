package modelcoder

import (
	"encoding/json"
	"strings"
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
