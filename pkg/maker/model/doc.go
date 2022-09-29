package model

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"teamide/pkg/util"
)

type modelDocStruct struct {
	Name         string                 `json:"name"`
	Comment      string                 `json:"comment"`
	Abbreviation string                 `json:"abbreviation"`
	Fields       []*modelDocFieldStruct `json:"fields"`
}

type modelDocFieldStruct struct {
	Name    string          `json:"name"`
	Comment string          `json:"comment"`
	IsList  bool            `json:"isList"`
	Struct  *modelDocStruct `json:"struct"`
}

func toText(data map[string]interface{}, docStruct *modelDocStruct) (text string, err error) {
	if data == nil {
		data = map[string]interface{}{}
	}

	doc, err := toDoc(data, docStruct)
	if err != nil {
		return
	}

	text, err = docToText(doc, false)
	if err != nil {
		return
	}

	return
}

func toDoc(data map[string]interface{}, docStruct *modelDocStruct) (doc *modelDoc, err error) {
	doc = newModelDoc()
	doc.comment = docStruct.Comment

	for _, docFieldStruct := range docStruct.Fields {
		err = appendDocField(doc, data[docFieldStruct.Name], docFieldStruct)
		if err != nil {
			return
		}
	}
	return
}

func appendDocField(doc *modelDoc, value interface{}, docFieldStruct *modelDocFieldStruct) (err error) {
	field := &modelDocField{}
	doc.addField(field)

	field.name = docFieldStruct.Name
	field.comment = docFieldStruct.Comment
	field.isValues = docFieldStruct.IsList

	list, listOk := value.([]interface{})
	if docFieldStruct.Struct != nil {
		if field.isValues {
			if listOk {
				field.values = list
			} else {
				field.values = []interface{}{value}
			}
		} else {
			field.value = value
		}
	} else {
		if field.isValues {
			if listOk {
				field.values = list
			} else {
				field.values = []interface{}{value}
			}
		} else {
			field.value = value
		}
	}

	return
}

func getDocFieldValue(doc *modelDoc, value interface{}, docFieldStruct *modelDocFieldStruct) (err error) {
	field := &modelDocField{}
	doc.addField(field)

	field.name = docFieldStruct.Name
	field.comment = docFieldStruct.Comment
	field.isValues = docFieldStruct.IsList

	list, listOk := value.([]interface{})
	if docFieldStruct.Struct != nil {
		if field.isValues {
			if listOk {
				field.values = list
			} else {
				field.values = []interface{}{value}
			}
		} else {
			field.value = value
		}
	} else {
		if field.isValues {
			if listOk {
				field.values = list
			} else {
				field.values = []interface{}{value}
			}
		} else {
			field.value = value
		}
	}

	return
}

type modelDoc struct {
	comment string
	fields  []*modelDocField
}

func (this_ *modelDoc) addField(field *modelDocField) (res *modelDoc) {
	if field == nil {
		return
	}
	res = this_
	this_.fields = append(this_.fields, field)

	return
}

type modelDocField struct {
	name     string
	value    interface{}
	comment  string
	values   []interface{}
	isValue  bool
	isValues bool
}

type modelDocFieldValue struct {
	value   interface{}
	comment string
}

func newModelDoc() (doc *modelDoc) {
	doc = &modelDoc{}
	return
}

func newModelDocFieldByValue(name string, value interface{}, comment string) (field *modelDocField) {
	field = &modelDocField{
		name:    name,
		value:   value,
		comment: comment,
		isValue: true,
	}
	return
}

func newModelDocFieldByValues(name string, comment string) (field *modelDocField) {
	field = &modelDocField{
		name:     name,
		comment:  comment,
		isValues: true,
	}
	return
}

func docToText(doc *modelDoc, omitEmpty bool) (text string, err error) {
	if doc == nil {
		return
	}

	documentNode := &yaml.Node{
		Kind:    yaml.DocumentNode,
		Content: []*yaml.Node{},
	}

	appendDocContent(documentNode, doc, omitEmpty)

	data, err := yaml.Marshal(documentNode)
	if err != nil {
		util.Logger.Error("docToText error", zap.Any("node", documentNode), zap.Error(err))
		return
	}

	text = string(data)

	return
}

func appendDocContent(node *yaml.Node, doc *modelDoc, omitEmpty bool) {
	if doc == nil {
		return
	}

	var mapNode = &yaml.Node{
		Content:     []*yaml.Node{},
		Kind:        4,
		HeadComment: doc.comment,
	}
	node.Content = append(node.Content, mapNode)

	for _, field := range doc.fields {
		appendFieldContent(mapNode, field, omitEmpty)
	}
}

func appendFieldContent(node *yaml.Node, field *modelDocField, omitEmpty bool) {
	if field == nil {
		return
	}
	if field.isValue {
		if omitEmpty {
			if field.value == "" || field.value == 0 || field.value == false || util.IsZero(field.value) {
				return
			}
		}
	}
	node.Content = append(node.Content, &yaml.Node{
		Value:       field.name,
		Kind:        8,
		LineComment: field.comment,
	})

	if field.isValues {
		if len(field.values) > 0 {
			var listNode = &yaml.Node{
				Content: []*yaml.Node{},
				Kind:    2,
			}
			node.Content = append(node.Content, listNode)

			for _, value := range field.values {
				appendFieldValue(listNode, value, omitEmpty)
			}
		}
	} else {
		appendFieldValue(node, field.value, omitEmpty)
	}
}

func appendFieldValue(node *yaml.Node, value interface{}, omitEmpty bool) {

	docV, docVOk := value.(*modelDoc)

	if docVOk {
		appendDocContent(node, docV, omitEmpty)
	} else {
		vValue, vValueOk := value.(*modelDocFieldValue)
		if vValueOk {
			str, _ := util.GetStringValue(vValue.value)
			node.Content = append(node.Content, &yaml.Node{
				Kind:        8,
				Value:       str,
				LineComment: vValue.comment,
			})
		} else {
			str, _ := util.GetStringValue(value)
			node.Content = append(node.Content, &yaml.Node{
				Kind:  8,
				Value: str,
			})
		}
	}
}
