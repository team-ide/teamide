package model

import (
	"teamide/pkg/application/base"
)

type SqlWhere struct {
	IfScript   string      `json:"ifScript,omitempty" yaml:"ifScript"`           // 条件  满足该条件 添加
	Custom     bool        `json:"custom,omitempty" yaml:"custom,omitempty"`     // 是否自定义
	CustomSql  string      `json:"customSql,omitempty" yaml:"customSql"`         // 是否自定义
	TableAlias string      `json:"tableAlias,omitempty" yaml:"tableAlias"`       // 表别名
	Name       string      `json:"name,omitempty" yaml:"name,omitempty"`         // 字段名称
	Value      string      `json:"value,omitempty" yaml:"value,omitempty"`       // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required   bool        `json:"required,omitempty" yaml:"required,omitempty"` // 必填
	AndOr      string      `json:"andOr,omitempty" yaml:"andOr,omitempty"`       // AND 和 OR 运算符
	Operator   string      `json:"operator,omitempty" yaml:"operator,omitempty"` // 运算符 = < > LIKE BETWEEN
	Piece      bool        `json:"piece,omitempty" yaml:"piece,omitempty"`       // 是条件快
	Wheres     []*SqlWhere `json:"wheres,omitempty" yaml:"wheres,omitempty"`     // 条件快的条件
}

func (this_ *SqlWhere) GetAndOr() string {
	if base.IsNotEmpty(this_.AndOr) {
		return this_.AndOr
	}
	return "AND"
}

func (this_ *SqlWhere) GetOperator() Operator {
	if base.IsNotEmpty(this_.Operator) {
		return EQ
	}
	for _, operator := range operators {
		if operator.Value == this_.Operator {
			return operator
		}
	}
	return EQ
}

type Operator struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

// 运算符
var (
	EQ              = appendOperator("=", "等于")
	LIKE            = appendOperator("LIKE", "包含")
	NOT_LIKE        = appendOperator("NOT LIKE", "不包含")
	NEQ             = appendOperator("<>", "不等于")
	GT              = appendOperator(">", "大于")
	LT              = appendOperator("<", "小于")
	GTE             = appendOperator(">=", "大于或等于")
	LTE             = appendOperator("<=", "小于或等于")
	IS_NULL         = appendOperator("IS NULL", "是null")
	IS_NOT_NULL     = appendOperator("IS NOT NULL", "不是null")
	IS_EMPTY        = appendOperator("IS EMPTY", "是空字符串")
	IS_NOT_EMPTY    = appendOperator("IS NOT EMPTY", "不是空字符串")
	LIKE_BEFORE     = appendOperator("LIKE%", "开始以")
	NOT_LIKE_BEFORE = appendOperator("NOT LIKE%", "开始不是以")
	LIKE_AFTER      = appendOperator("%LIKE", "结束以")
	NOT_LIKE_AFTER  = appendOperator("NOT %LIKE", "结束不是以")
	IN              = appendOperator("IN", "IN")
	NOT_IN          = appendOperator("NOT IN", "NOT IN")
	IN_LIKE         = appendOperator("IN LIKE", "IN LIKE %...%")
	NOT_IN_LIKE     = appendOperator("NOT IN LIKE", "NOT IN LIKE %...%")
)
var (
	operators = []Operator{}
)

func appendOperator(value string, text string) Operator {
	res := Operator{value, text}
	operators = append(operators, res)
	return res
}
