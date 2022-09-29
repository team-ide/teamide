package model

import (
	"encoding/json"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"teamide/pkg/util"
)

type StructModel struct {
	Name    string         `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string         `json:"comment,omitempty"` // 说明
	Note    string         `json:"note,omitempty"`    // 注释
	Parent  string         `json:"parent,omitempty"`  // 父结构体
	Fields  []*StructField `json:"fields,omitempty"`  // 结构体字段
}

type StructField struct {
	Name          string `json:"name,omitempty"`          // 字段名称，同一个结构体中唯一
	Comment       string `json:"comment,omitempty"`       // 说明
	Note          string `json:"note,omitempty"`          // 注释
	JsonName      string `json:"jsonName,omitempty"`      // 映射 JSON 字段 默认和字段名称一致
	JsonOmitempty bool   `json:"jsonOmitempty,omitempty"` // 映射 JSON 字段 省略空值
	IsList        bool   `json:"isList,omitempty"`        // 是否是列表
	DataType      string `json:"dataType,omitempty"`      // 数据类型
	Default       string `json:"default,omitempty"`       // 默认值
}

func (this_ *StructField) otherEmpty() (res bool) {
	res = true
	if this_.Comment != "" ||
		this_.Note != "" ||
		(this_.JsonName != "" && this_.JsonName != this_.Name) ||
		this_.JsonOmitempty ||
		this_.IsList ||
		this_.DataType != "" ||
		this_.Default != "" {
		res = false
	}
	return
}

var (
	DocStructStruct = &modelDocStruct{
		Name:    "struct",
		Comment: "结构体文件，该文件用于入参、出参、函数调用、数据存储等地方",
		Fields: []*modelDocFieldStruct{
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
				Name:    "parent",
				Comment: "父级结构体，源码将继承该结构体",
			},
			{
				Name:    "fields",
				Comment: "这是结构体字段",
				IsList:  true,
				Struct: &modelDocStruct{
					Comment:      "结构体字段",
					Abbreviation: "name",
					Fields: []*modelDocFieldStruct{
						{
							Name:    "name",
							Comment: "字段名称",
						},
						{
							Name:    "name",
							Comment: "字段名称",
						},
					},
				},
			},
		},
	}
)

func StructToText(model *StructModel) (text string, err error) {
	if model == nil {
		model = &StructModel{}
		return
	}

	//doc := newModelDoc()
	//doc.comment = "结构体文件，该文件用于入参、出参、函数调用、数据存储等地方"
	//doc.addField(newModelDocFieldByValue("name", model.Name, "结构体名称"))
	//doc.addField(newModelDocFieldByValue("comment", model.Comment, "结构体说明"))
	//doc.addField(newModelDocFieldByValue("note", model.Note, "结构体源码注释"))
	//doc.addField(newModelDocFieldByValue("parent", model.Parent, "父级结构体，源码将继承该结构体"))
	//
	//field := newModelDocFieldByValues("fields", "这是结构体字段")
	//doc.addField(field)
	//
	//for _, one := range model.Fields {
	//	if one.otherEmpty() {
	//		field.values = append(field.values, &modelDocFieldValue{
	//			value:   one.Name,
	//			comment: "结构体字段，字段名称",
	//		})
	//	} else {
	//		oneDoc := newModelDoc()
	//		oneDoc.comment = "结构体字段"
	//		oneDoc.addField(newModelDocFieldByValue("name", one.Name, "字段名称"))
	//		oneDoc.addField(newModelDocFieldByValue("comment", one.Comment, "字段说明"))
	//		oneDoc.addField(newModelDocFieldByValue("note", one.Note, "字段源码注释"))
	//		oneDoc.addField(newModelDocFieldByValue("jsonName", one.JsonName, "序列化JSON名称"))
	//		oneDoc.addField(newModelDocFieldByValue("jsonOmitempty", one.JsonOmitempty, "序列化JSON，省略空值"))
	//		oneDoc.addField(newModelDocFieldByValue("isList", one.IsList, "是集合"))
	//		oneDoc.addField(newModelDocFieldByValue("dataType", one.IsList, "数据类型，基础类型或结构体名称"))
	//		oneDoc.addField(newModelDocFieldByValue("default", one.IsList, "创建对象该字段默认的值"))
	//		field.values = append(field.values, oneDoc)
	//	}
	//}

	bytes, err := json.Marshal(model)
	if err != nil {
		util.Logger.Error("model to bytes error", zap.Any("model", model), zap.Error(err))
		return
	}
	data := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, data)
	if err != nil {
		util.Logger.Error("bytes to data error", zap.Any("bytes", bytes), zap.Error(err))
		return
	}
	text, err = toText(data, DocStructStruct)
	if err != nil {
		util.Logger.Error("struct to text error", zap.Any("model", model), zap.Error(err))
	}
	return
}

func TextToStruct(text string) (model *StructModel, err error) {
	var bs []byte
	data := map[string]interface{}{}

	model = &StructModel{}

	err = yaml.Unmarshal([]byte(text), data)
	if err != nil {
		util.Logger.Error("text to data error", zap.Any("text", text), zap.Error(err))
		return
	}

	fieldsData := data["fields"]
	delete(data, "fields")

	bs, err = json.Marshal(data)
	if err != nil {
		util.Logger.Error("data to bytes error", zap.Any("data", data), zap.Error(err))
		return
	}
	err = yaml.Unmarshal(bs, model)
	if err != nil {
		util.Logger.Error("data to struct error", zap.Any("data", data), zap.Error(err))
		return
	}
	if fieldsData != nil {
		toList(
			fieldsData,
			func(str string) interface{} {
				return &StructField{Name: str}
			},
			func(one interface{}) {
				model.Fields = append(model.Fields, one.(*StructField))
			})
	}

	//util.Logger.Info("text to model success", zap.Any("data", data), zap.Any("model", model))

	return
}
