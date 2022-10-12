package model

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

type ServiceModel struct {
	Name    string        `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string        `json:"comment,omitempty"` // 说明
	Note    string        `json:"note,omitempty"`    // 注释
	Args    []*ServiceArg `json:"args"`              //入参
	Steps   []interface{} `json:"steps,omitempty"`   // 阶段
	Return  string        `json:"return,omitempty"`  // 返回
}

type ServiceArg struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Type    string `json:"type,omitempty"`    // 类型
}

var (
	docTemplateServiceName = "service"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    docTemplateServiceName,
		Comment: "服务文件，该文件用于动作处理，如数据库、redis、文件等地方",
		Fields: []*docTemplateField{
			{
				Name:    "name",
				Comment: "结构体名称",
			},
			{
				Name:    "comment",
				Comment: "结构体说明",
			},
			{
				Name:    "note",
				Comment: "结构体源码注释",
			},
			{
				Name:       "steps",
				Comment:    "阶段",
				IsList:     true,
				StructName: docTemplateStepName,
			},
			{
				Name:    "return",
				Comment: "返回值",
			},
		},
	})
}

func ServiceToText(model interface{}) (text string, err error) {
	text, err = toText(model, docTemplateServiceName, &docOptions{
		outComment: true,
		omitEmpty:  false,
	})
	if err != nil {
		util.Logger.Error("service model to text error", zap.Any("model", model), zap.Error(err))
		return
	}
	return
}

func TextToService(text string) (model *ServiceModel, err error) {
	model = &ServiceModel{}
	err = toModel(text, docTemplateServiceName, model)
	if err != nil {
		util.Logger.Error("text to service model error", zap.Any("text", text), zap.Error(err))
		return
	}
	return
}
