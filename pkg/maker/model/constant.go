package model

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"teamide/pkg/util"
)

type ConstantModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Type    string `json:"type,omitempty"`    // 类型
	Value   string `json:"value,omitempty"`   // 值
}

func ConstantsToText(list []*ConstantModel) (text string, err error) {
	bytes, err := yaml.Marshal(list)
	if err != nil {
		util.Logger.Error("constants to yaml error", zap.Any("errors", list), zap.Error(err))
		return
	}
	text = string(bytes)
	return
}

func TextToConstants(text string) (list []*ConstantModel, err error) {
	err = yaml.Unmarshal([]byte(text), &list)
	if err != nil {
		util.Logger.Error("yaml to constants error", zap.Any("yaml", text), zap.Error(err))
		return
	}
	return
}
