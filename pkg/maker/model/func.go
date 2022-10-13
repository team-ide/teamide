package model

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

type FuncModel struct {
	Name       string      `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Comment    string      `json:"comment,omitempty"`    // 说明
	Note       string      `json:"note,omitempty"`       // 注释
	Args       []*ArgModel `json:"args,omitempty"`       // 入参
	Func       string      `json:"func,omitempty"`       // 函数内容
	ReturnType string      `json:"returnType,omitempty"` // 返回类型
}

var (
	docTemplateFuncName = "func"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    docTemplateFuncName,
		Comment: "函数",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "函数名称"},
			{Name: "comment", Comment: "函数说明"},
			{Name: "note", Comment: "函数注释"},
			{Name: "args", Comment: "参数", IsList: true, StructName: docTemplateArgName},
			{Name: "func", Comment: "函数内容"},
			{Name: "returnType", Comment: "返回类型"},
		},
	})
}

func FuncToText(model interface{}) (text string, err error) {
	text, err = toText(model, docTemplateFuncName, &docOptions{
		outComment: true,
		omitEmpty:  false,
	})
	if err != nil {
		util.Logger.Error("func model to text error", zap.Any("model", model), zap.Error(err))
		return
	}
	return
}

func TextToFunc(text string) (model *FuncModel, err error) {
	model = &FuncModel{}
	err = toModel(text, docTemplateFuncName, model)
	if err != nil {
		util.Logger.Error("text to func model error", zap.Any("text", text), zap.Error(err))
		return
	}
	return
}
