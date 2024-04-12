package modelers

import (
	"github.com/team-ide/go-tool/util"
	"strings"
)

type StepDbWhere struct {
	If         string         `json:"if,omitempty"`         // 条件  满足该条件 添加
	Custom     bool           `json:"custom,omitempty"`     // 是否自定义
	Sql        string         `json:"sql,omitempty"`        // 自定义SQL
	TableAlias string         `json:"tableAlias,omitempty"` // 表别名
	Name       string         `json:"name,omitempty"`       // 字段名称
	Value      string         `json:"value,omitempty"`      // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required   bool           `json:"required,omitempty"`   // 必填
	And        bool           `json:"and,omitempty"`        // AND  默认
	Or         bool           `json:"or,omitempty"`         // OR
	Operator   string         `json:"operator,omitempty"`   // 运算符 = < > LIKE BETWEEN
	Piece      bool           `json:"piece,omitempty"`      // 是条件块
	Wheres     []*StepDbWhere `json:"wheres,omitempty"`     // 条件快的条件
}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateStepDbWhereName,
		Abbreviation: "name",
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "tableAlias", Comment: "表别名"},
			{Name: "name", Comment: "字段名称"},
			{Name: "value", Comment: "字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值"},
			{Name: "required", Comment: "必填"},
			{Name: "and", Comment: "AND 默认"},
			{Name: "or", Comment: "OR"},
			{Name: "operator", Comment: "运算符 = < > LIKE BETWEEN"},
			{Name: "piece", Comment: "是条件块"},
			{Name: "wheres", Comment: "条件块的条件", IsList: true, StructName: docTemplateStepDbWhereName},
		},
		newModel: func() interface{} {
			return &StepDbWhere{}
		},
		newModels: func() interface{} {
			var vs []*StepDbWhere
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbWhere)
			vs = append(vs, value.(*StepDbWhere))
			return vs
		},
	})
}

func (this_ *StepDbWhere) GetAndOr() string {
	if this_.Or {
		return "OR"
	}
	return "AND"
}

func (this_ *StepDbWhere) GetOperator() *Operator {
	if util.IsNotEmpty(this_.Operator) {
		return OperatorEqual
	}
	for _, operator := range operators {
		if strings.EqualFold(
			strings.ReplaceAll(operator.Value, " ", ""),
			strings.ReplaceAll(this_.Operator, " ", ""),
		) {
			return operator
		}
	}
	return OperatorEqual
}

type Operator struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

// 运算符
var (
	OperatorEqual         = appendOperator("=", "等于")
	OperatorLike          = appendOperator("LIKE", "包含")
	OperatorNotLike       = appendOperator("NOT LIKE", "不包含")
	OperatorNotEqual      = appendOperator("<>", "不等于")
	OperatorGt            = appendOperator(">", "大于")
	OperatorLt            = appendOperator("<", "小于")
	OperatorGte           = appendOperator(">=", "大于或等于")
	OperatorLte           = appendOperator("<=", "小于或等于")
	OperatorIsNull        = appendOperator("IS NULL", "是null")
	OperatorIsNotNull     = appendOperator("IS NOT NULL", "不是null")
	OperatorIsEmpty       = appendOperator("IS EMPTY", "是空字符串")
	OperatorIsNotEmpty    = appendOperator("IS NOT EMPTY", "不是空字符串")
	OperatorLikeBefore    = appendOperator("LIKE%", "开始以")
	OperatorNotLikeBefore = appendOperator("NOT LIKE%", "开始不是以")
	OperatorLikeAfter     = appendOperator("%LIKE", "结束以")
	OperatorNotLikeAfter  = appendOperator("NOT %LIKE", "结束不是以")
	OperatorIn            = appendOperator("IN", "IN")
	OperatorNotIn         = appendOperator("NOT IN", "NOT IN")
	OperatorInLike        = appendOperator("IN LIKE", "IN LIKE %...%")
	OperatorNotInLike     = appendOperator("NOT IN LIKE", "NOT IN LIKE %...%")
)
var (
	operators []*Operator
)

func appendOperator(value string, text string) *Operator {
	res := &Operator{value, text}
	operators = append(operators, res)
	return res
}
