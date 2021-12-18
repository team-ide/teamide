package sqlModel

import (
	"fmt"
	"strings"
)

type Where struct {
	IfScript    string   `json:"ifScript"`    // 条件  满足该条件 添加
	Custom      bool     `json:"custom"`      // 是否自定义
	CustomSql   string   `json:"customSql"`   // 是否自定义
	TableAlias  string   `json:"tableAlias"`  // 表别名
	Name        string   `json:"name"`        // 字段名称
	ValueScript string   `json:"valueScript"` // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required    bool     `json:"required"`    // 必填
	AndOr       string   `json:"andOr"`       // AND 和 OR 运算符
	Operator    string   `json:"operator"`    // 运算符 = < > LIKE BETWEEN
	Piece       bool     `json:"piece"`       // 是条件快
	Wheres      []*Where `json:"wheres"`      // 条件快的条件
}

func (this_ *Where) GetAndOr() string {
	if IsNotEmpty(this_.AndOr) {
		return this_.AndOr
	}
	return "AND"
}

func (this_ *Where) GetOperator() Operator {
	if IsNotEmpty(this_.Operator) {
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
	Value string `json:"value"`
	Text  string `json:"text"`
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

func getWhereSqlParam(data map[string]interface{}, wheres []*Where) (whereSql string, params []interface{}, err error) {
	if len(wheres) == 0 {
		return
	}
	var whereSql_ string
	var params_ []interface{}
	err = appendWhereSql(data, wheres, &whereSql_, &params_)
	if err != nil {
		return
	}
	if IsNotEmpty(whereSql_) {
		whereSql_ = strings.TrimPrefix(whereSql_, " AND ")
		whereSql_ = strings.TrimPrefix(whereSql_, " OR ")

		whereSql = "" + whereSql_
		params = params_
	}
	return
}

func appendWhereSql(data map[string]interface{}, wheres []*Where, whereSql *string, params *[]interface{}) (err error) {
	if len(wheres) == 0 {
		return
	}

	for _, where := range wheres {

		if !IfScriptValue(data, where.IfScript) {
			continue
		}
		if where.Piece {
			var pieceWhereSql string
			var pieceParams []interface{}
			pieceWhereSql, pieceParams, err = getPieceWhereSql(data, where.Wheres)
			if err != nil {
				return
			}
			if IsNotEmpty(pieceWhereSql) {
				*whereSql += " " + where.GetAndOr() + " " + pieceWhereSql
				*params = append(*params, pieceParams...)
			}
			continue
		}
		if where.Custom {
			*whereSql += where.CustomSql
			continue
		}
		if IsEmpty(where.Name) {
			continue
		}

		value := GetColumnValue(data, where.Name, where.ValueScript)
		if IsEmptyObj(value) {
			continue
		}

		wrapColumn := WrapColumnName(where.TableAlias, where.Name)
		*whereSql += " " + where.GetAndOr() + " "
		operator := where.GetOperator()
		switch operator {
		case IS_NULL:
			*whereSql += wrapColumn + " IS NULL"
		case IS_NOT_NULL:
			*whereSql += wrapColumn + " IS NOT NULL"
		case IS_EMPTY:
			*whereSql += wrapColumn + " = ''"
		case IS_NOT_EMPTY:
			*whereSql += wrapColumn + " <> ''"
		case LIKE:
			*whereSql += wrapColumn + " LIKE ?"
			setValue := fmt.Sprint("%", value, "%")
			*params = append(*params, setValue)
		case NOT_LIKE:
			*whereSql += wrapColumn + " NOT LIKE ?"
			setValue := fmt.Sprint("%", value, "%")
			*params = append(*params, setValue)
		case LIKE_BEFORE:
			*whereSql += wrapColumn + " LIKE ?%"
			setValue := fmt.Sprint(value, "%")
			*params = append(*params, setValue)
		case NOT_LIKE_BEFORE:
			*whereSql += wrapColumn + " NOT LIKE ?%"
			setValue := fmt.Sprint(value, "%")
			*params = append(*params, setValue)
		case LIKE_AFTER:
			*whereSql += wrapColumn + " LIKE %?"
			setValue := fmt.Sprint("%", value)
			*params = append(*params, setValue)
		case NOT_LIKE_AFTER:
			*whereSql += wrapColumn + " NOT LIKE %?"
			setValue := fmt.Sprint("%", value)
			*params = append(*params, setValue)
		case IN:
			*whereSql += wrapColumn + " IN (?)"
			setValue := "0,2,3,4"
			*params = append(*params, setValue)
		case NOT_IN:
			*whereSql += wrapColumn + " NOT IN (?)"
			setValue := "0,2,3,4"
			*params = append(*params, setValue)
		case IN_LIKE:
			*whereSql += wrapColumn + " IN LIKE (?)"
			setValue := "0,2,3,4"
			*params = append(*params, setValue)
		case NOT_IN:
			*whereSql += wrapColumn + " NOT IN LIKE (?)"
			setValue := "0,2,3,4"
			*params = append(*params, setValue)
		default:
			*whereSql += wrapColumn + " " + operator.Value + " ?"
			*params = append(*params, value)
		}
	}
	return
}

func getPieceWhereSql(data map[string]interface{}, wheres []*Where) (whereSql string, params []interface{}, err error) {
	if len(wheres) == 0 {
		return
	}
	whereSql_ := ""
	params_ := []interface{}{}
	err = appendWhereSql(data, wheres, &whereSql_, &params_)
	if err != nil {
		return
	}
	if IsNotEmpty(whereSql_) {
		whereSql_ = strings.TrimPrefix(whereSql_, " AND ")
		whereSql_ = strings.TrimPrefix(whereSql_, " OR ")

		whereSql = "(" + whereSql_ + ")"
		params = params_
	}
	return
}
