package model

import (
	"go.uber.org/zap"
	"teamide/pkg/util"
)

type DaoModel struct {
	Name    string        `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string        `json:"comment,omitempty"` // 说明
	Note    string        `json:"note,omitempty"`    // 注释
	Args    []*ArgModel   `json:"args,omitempty"`    //入参
	Steps   []interface{} `json:"steps,omitempty"`   // 阶段
	Return  string        `json:"return,omitempty"`  // 返回
}

var (
	docTemplateDaoName = "dao"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    docTemplateDaoName,
		Comment: "数据层，处理数据库等落地数据",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "数据层名称"},
			{Name: "comment", Comment: "数据层说明"},
			{Name: "note", Comment: "数据层源码注释"},
			{Name: "args", Comment: "参数", IsList: true, StructName: docTemplateArgName},
			{Name: "steps", Comment: "阶段", IsList: true, StructName: docTemplateStepName},
			{Name: "return", Comment: "返回值"},
		},
	})
}

func DaoToText(model interface{}) (text string, err error) {
	text, err = toText(model, docTemplateDaoName, &docOptions{
		outComment: true,
		omitEmpty:  false,
	})
	if err != nil {
		util.Logger.Error("dao model to text error", zap.Any("model", model), zap.Error(err))
		return
	}
	return
}

func TextToDao(text string) (model *DaoModel, err error) {
	model = &DaoModel{}
	err = toModel(text, docTemplateDaoName, model)
	if err != nil {
		util.Logger.Error("text to dao model error", zap.Any("text", text), zap.Error(err))
		return
	}
	return
}
