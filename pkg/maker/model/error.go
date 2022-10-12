package model

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"teamide/pkg/util"
)

type ErrorModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Code    string `json:"code,omitempty"`    // 错误码
	Msg     string `json:"msg,omitempty"`     // 错误信息
}

func ErrorsToText(list []*ErrorModel) (text string, err error) {
	bytes, err := yaml.Marshal(list)
	if err != nil {
		util.Logger.Error("errors to yaml error", zap.Any("errors", list), zap.Error(err))
		return
	}
	text = string(bytes)
	return
}

func TextToErrors(text string) (list []*ErrorModel, err error) {
	err = yaml.Unmarshal([]byte(text), &list)
	if err != nil {
		util.Logger.Error("yaml to errors error", zap.Any("yaml", text), zap.Error(err))
		return
	}
	return
}
