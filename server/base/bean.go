package base

import (
	"fmt"
	"strings"
)

type OBean struct {
	Text    string      `json:"text" column:"text,omitempty"`
	Value   interface{} `json:"value" column:"value,omitempty"`
	Comment string      `json:"comment" column:"comment,omitempty"`
	Color   string      `json:"color" column:"color,omitempty"`
}

func NewOBean(text string, value interface{}) (res OBean) {
	res = OBean{
		Text:  text,
		Value: value,
	}
	return
}

type UserTotalBean struct {
	User       *UserEntity         `json:"user,omitempty"`
	Password   *UserPasswordEntity `json:"password,omitempty"`
	Persona    *UserPersonaBean    `json:"persona,omitempty"`
	Enterprise *UserEnterpriseBean `json:"enterprise,omitempty"`
}

type UserPersonaBean struct {
	Name   string  `json:"name,omitempty"`
	Age    int     `json:"age,omitempty"`
	Sex    int8    `json:"sex,omitempty"`
	Photo  string  `json:"photo,omitempty"`
	Height float32 `json:"height,omitempty"`
	Weight float32 `json:"weight,omitempty"`
}

type UserEnterpriseBean struct {
	Name   string                  `json:"name,omitempty"`
	Salary float32                 `json:"salary,omitempty"`
	Orgs   []UserEnterpriseOrgBean `json:"orgs,omitempty"`
}

type UserEnterpriseOrgBean struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Position string `json:"position,omitempty"`
}

type InstallInfo struct {
	Module string              `json:"module,omitempty"`
	Stages []*InstallStageInfo `json:"stages,omitempty"`
}

type InstallStageInfo struct {
	Stage    string   `json:"stage,omitempty"`
	SqlParam SqlParam `json:"sqlParam,omitempty"`
}

type SqlParam struct {
	PageIndex   int64         `json:"pageIndex,omitempty"`
	PageSize    int64         `json:"pageSize,omitempty"`
	Sql         string        `json:"sql,omitempty"`
	Params      []interface{} `json:"params,omitempty"`
	CountSql    string        `json:"countSql,omitempty"`
	CountParams []interface{} `json:"countParams,omitempty"`
}

func (this_ SqlParam) ToExecSql() string {
	strs := strings.Split(this_.Sql, `?`)
	res := ""
	for index, str := range strs {
		res += str
		if index >= len(strs)-1 {
			continue
		}
		param := this_.Params[index]
		res += fmt.Sprint("'", param, "'")
	}
	return res
}

func NewSqlParam(sql_ string, params []interface{}) (sqlParam SqlParam) {
	if params == nil {
		params = []interface{}{}
	}
	sqlParam = SqlParam{
		Sql:    sql_,
		Params: params,
	}
	return
}
