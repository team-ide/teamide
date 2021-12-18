package sql

import (
	"fmt"
	"strings"
)

type SqlParam struct {
	PageIndex int64         `json:"pageIndex,omitempty"`
	PageSize  int64         `json:"pageSize,omitempty"`
	Sql       string        `json:"sql,omitempty"`
	Params    []interface{} `json:"params,omitempty"`
}

func (this_ SqlParam) ToExecSql() string {
	strs := strings.Split(this_.Sql, `?`)
	res := ""
	for index, str := range strs {
		res += str
		println("index:", index, ",str:", str)
		if index >= len(strs)-1 {
			continue
		}
		param := this_.Params[index]
		res += fmt.Sprint("'", param, "'")
	}
	return res
}
