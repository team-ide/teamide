package model

import (
	"errors"
	"fmt"
	"teamide/pkg/application/base"
)

type ValidateModel struct {
	Name      string               `json:"name,omitempty" yaml:"name,omitempty"`           // 名称
	Comment   string               `json:"comment,omitempty" yaml:"comment,omitempty"`     // 注释说明
	Required  bool                 `json:"required,omitempty" yaml:"required,omitempty"`   // 必填
	Min       int                  `json:"min,omitempty" yaml:"min,omitempty"`             // 最小数
	Max       int                  `json:"max,omitempty" yaml:"max,omitempty"`             // 最大数
	MinLength int                  `json:"minLength,omitempty" yaml:"minLength,omitempty"` // 最小长度
	MaxLength int                  `json:"maxLength,omitempty" yaml:"maxLength,omitempty"` // 最大长度
	Pattern   string               `json:"pattern,omitempty" yaml:"pattern,omitempty"`     // 正则匹配
	Error     string               `json:"error,omitempty" yaml:"error,omitempty"`         // 异常
	ErrorCode string               `json:"errorCode,omitempty" yaml:"errorCode,omitempty"` // 异常码码
	ErrorMsg  string               `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`   // 异常信息
	Rules     []*ValidateRuleModel `json:"rules,omitempty" yaml:"rules,omitempty"`         // 变量
	TryError  *ErrorModel          `json:"tryError,omitempty" yaml:"tryError,omitempty"`
}

type ValidateRuleModel struct {
	Required  bool   `json:"required,omitempty" yaml:"required,omitempty"`   // 必填
	Min       int    `json:"min,omitempty" yaml:"min,omitempty"`             // 最小数
	Max       int    `json:"max,omitempty" yaml:"max,omitempty"`             // 最大数
	MinLength int    `json:"minLength,omitempty" yaml:"minLength,omitempty"` // 最小长度
	MaxLength int    `json:"maxLength,omitempty" yaml:"maxLength,omitempty"` // 最大长度
	Pattern   string `json:"pattern,omitempty" yaml:"pattern,omitempty"`     // 正则匹配
	Error     string `json:"error,omitempty" yaml:"error,omitempty"`         // 异常
	ErrorCode string `json:"errorCode,omitempty" yaml:"errorCode,omitempty"` // 异常码码
	ErrorMsg  string `json:"errorMsg,omitempty" yaml:"errorMsg,omitempty"`   // 异常信息
}

func getValidatesByValue(value interface{}) (validates []*ValidateModel, err error) {
	if value == nil {
		return
	}
	values, valuesOk := value.([]interface{})
	if !valuesOk || len(values) == 0 {
		return
	}

	for _, valuesOne := range values {
		switch v := valuesOne.(type) {
		case map[string]interface{}:
			validateMap := v
			if len(v) == 1 {
				for mapKey, mapValue := range v {
					switch subV := mapValue.(type) {
					case map[string]interface{}:
						validateMap = subV
					default:
						validateMap["pattern"] = fmt.Sprint(subV)
					}
					if validateMap["name"] == nil {
						validateMap["name"] = mapKey
					}
				}
			}
			formatTryErrorByMap(validateMap)
			validate := &ValidateModel{}
			err = base.ToBean([]byte(base.ToJSON(validateMap)), validate)
			if err != nil {
				return
			}
			validates = append(validates, validate)
		default:
			err = errors.New(fmt.Sprint("[", v, "] to variable error"))
		}
	}
	return
}
