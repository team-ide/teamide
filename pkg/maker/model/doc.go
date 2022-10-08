package model

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"sync"
	"teamide/pkg/util"
)

type docTemplate struct {
	Name         string              `json:"name"`
	Abbreviation string              `json:"abbreviation"`
	Comment      string              `json:"comment"`
	Fields       []*docTemplateField `json:"fields"`
}

type docTemplateField struct {
	Name       string      `json:"name"`
	Comment    string      `json:"comment"`
	IsList     bool        `json:"isList"`
	StructName string      `json:"structName"`
	Default    interface{} `json:"default"` // 默认值
}

type docOptions struct {
	omitEmpty  bool
	outComment bool
}

var (
	docTemplateCache     = map[string]*docTemplate{}
	docTemplateCacheLock sync.Mutex
)

func addDocTemplate(template *docTemplate) {
	docTemplateCacheLock.Lock()
	defer docTemplateCacheLock.Unlock()

	if docTemplateCache[template.Name] != nil {
		print("doc template [" + template.Name + "] already exist")
		return
	}
	docTemplateCache[template.Name] = template
}

func getDocTemplate(name string) (template *docTemplate) {
	docTemplateCacheLock.Lock()
	defer docTemplateCacheLock.Unlock()

	template = docTemplateCache[name]
	return
}

// GetDocTemplates 获取所有Doc模板
func GetDocTemplates() (templates []*docTemplate) {
	docTemplateCacheLock.Lock()
	defer docTemplateCacheLock.Unlock()

	for _, one := range docTemplateCache {
		templates = append(templates, one)
	}
	return
}

func toText(model interface{}, docTemplateName string, options *docOptions) (text string, err error) {

	if options == nil {
		options = &docOptions{}
	}

	bytes, err := json.Marshal(model)
	if err != nil {
		util.Logger.Error("model to bytes error", zap.Any("model", model), zap.Error(err))
		return
	}
	source := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, source)
	if err != nil {
		util.Logger.Error("bytes to data error", zap.Any("bytes", bytes), zap.Error(err))
		return
	}

	node := &yaml.Node{
		Kind: yaml.DocumentNode,
	}
	docStruct := getDocTemplate(docTemplateName)
	err = appendNode(node, source, docStruct, options)
	if err != nil {
		return
	}

	bytes, err = yaml.Marshal(node)
	if err != nil {
		util.Logger.Error("node to json error", zap.Any("node", node), zap.Error(err))
		return
	}

	text = string(bytes)

	return
}

func appendNode(node *yaml.Node, source map[string]interface{}, docStruct *docTemplate, options *docOptions) (err error) {

	var mapNode = &yaml.Node{
		Kind: 4,
	}

	if options.outComment {
		mapNode.HeadComment = docStruct.Comment
	}
	node.Content = append(node.Content, mapNode)

	for _, docFieldStruct := range docStruct.Fields {
		err = appendNodeField(mapNode, source[docFieldStruct.Name], docFieldStruct, options)
		if err != nil {
			return
		}
	}
	return
}

func appendNodeField(node *yaml.Node, value interface{}, docFieldStruct *docTemplateField, options *docOptions) (err error) {

	if !options.omitEmpty {
		if canNotOut(value) {
			return
		}
	}
	fieldNode := &yaml.Node{
		Value: docFieldStruct.Name,
		Kind:  8,
	}
	if options.outComment {
		if docFieldStruct.StructName != "" || docFieldStruct.IsList {
			fieldNode.HeadComment = docFieldStruct.Comment
		} else {
			fieldNode.LineComment = docFieldStruct.Comment
		}
	}

	node.Content = append(node.Content, fieldNode)

	if docFieldStruct.IsList {
		list, listOk := value.([]interface{})
		if !listOk {
			list = []interface{}{value}
		}
		if len(list) > 0 {
			err = appendNodeFieldValues(node, list, docFieldStruct, options)
			if err != nil {
				return
			}
		}
	} else {
		err = appendNodeFieldValue(node, value, docFieldStruct, options)
		if err != nil {
			return
		}
	}

	return
}

func appendNodeFieldValue(node *yaml.Node, value interface{}, docFieldStruct *docTemplateField, options *docOptions) (err error) {

	if docFieldStruct.StructName != "" {
		mapV, mapVOk := value.(map[string]interface{})
		if mapVOk {
			err = appendNodeValue(node, mapV, docFieldStruct.StructName, options)
			if err != nil {
				return
			}
			return
		}
	}

	str, _ := util.GetStringValue(value)
	node.Content = append(node.Content, &yaml.Node{
		Kind:  8,
		Value: str,
	})

	return
}

func appendNodeFieldValues(node *yaml.Node, values []interface{}, docFieldStruct *docTemplateField, options *docOptions) (err error) {
	if len(values) == 0 {
		return
	}

	var listNode = &yaml.Node{
		Kind: 2,
	}
	node.Content = append(node.Content, listNode)

	if docFieldStruct.StructName != "" {
		for _, value := range values {
			mapV, mapVOk := value.(map[string]interface{})
			if mapVOk {
				err = appendNodeValue(listNode, mapV, docFieldStruct.StructName, options)
				if err != nil {
					return
				}
			} else {
				str, _ := util.GetStringValue(value)
				listNode.Content = append(listNode.Content, &yaml.Node{
					Kind:  8,
					Value: str,
				})
			}
		}
	} else {
		for _, value := range values {
			str, _ := util.GetStringValue(value)
			listNode.Content = append(listNode.Content, &yaml.Node{
				Kind:  8,
				Value: str,
			})
		}
	}

	return
}

func appendNodeValue(node *yaml.Node, value map[string]interface{}, docTemplateName string, options *docOptions) (err error) {
	if value == nil || len(value) == 0 {
		node.Content = append(node.Content, &yaml.Node{
			Kind: 8,
		})
		return
	}
	docStruct := getDocTemplate(docTemplateName)
	if docStruct.Abbreviation != "" {
		var canNotOutCount = 0
		for _, docFieldStruct := range docStruct.Fields {
			if docFieldStruct.Name != docStruct.Abbreviation && canNotOut(value[docFieldStruct.Name]) {
				canNotOutCount++
			}
		}
		if canNotOutCount == len(docStruct.Fields)-1 {
			str, _ := util.GetStringValue(value[docStruct.Abbreviation])
			node.Content = append(node.Content, &yaml.Node{
				Kind:  8,
				Value: str,
			})
			return
		}

	}

	var mapNode = &yaml.Node{
		Kind: 4,
	}
	if options.outComment {
		mapNode.HeadComment = docStruct.Comment
	}
	node.Content = append(node.Content, mapNode)

	for _, docFieldStruct := range docStruct.Fields {
		err = appendNodeField(mapNode, value[docFieldStruct.Name], docFieldStruct, options)
		if err != nil {
			return
		}
	}
	return
}

func canNotOut(value interface{}) bool {
	if value == nil || value == "" || value == 0 || value == false || util.IsZero(value) {
		return true
	}
	return false
}

func toModel(text string, docTemplateName string, model interface{}) (err error) {
	var bs []byte
	source := map[string]interface{}{}

	err = yaml.Unmarshal([]byte(text), source)
	if err != nil {
		util.Logger.Error("text to source error", zap.Any("text", text), zap.Error(err))
		return
	}

	data, err := toData(source, docTemplateName)
	if err != nil {
		util.Logger.Error("source to data error", zap.Any("source", source), zap.Error(err))
		return
	}

	bs, err = json.Marshal(data)
	if err != nil {
		util.Logger.Error("data to bytes error", zap.Any("data", data), zap.Error(err))
		return
	}
	err = yaml.Unmarshal(bs, model)
	if err != nil {
		util.Logger.Error("data to model error", zap.Any("data", data), zap.Error(err))
		return
	}

	return
}

func toData(source map[string]interface{}, docTemplateName string) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	if source == nil {
		source = map[string]interface{}{}
	}
	docStruct := getDocTemplate(docTemplateName)
	err = appendData(data, source, docStruct)
	if err != nil {
		return
	}

	return
}

func appendData(data map[string]interface{}, source map[string]interface{}, docStruct *docTemplate) (err error) {

	for _, docFieldStruct := range docStruct.Fields {
		data[docFieldStruct.Name], err = getFieldData(source[docFieldStruct.Name], docFieldStruct)
		if err != nil {
			return
		}
	}
	return
}

func getFieldData(sourceValue interface{}, docFieldStruct *docTemplateField) (value interface{}, err error) {

	if docFieldStruct.IsList {
		list, listOk := sourceValue.([]interface{})
		if !listOk {
			list = []interface{}{sourceValue}
		}
		if len(list) > 0 {
			value, err = getFieldValues(list, docFieldStruct)
			if err != nil {
				return
			}
		}
	} else {
		value, err = getFieldValue(sourceValue, docFieldStruct)
		if err != nil {
			return
		}
	}

	return
}

func getFieldValue(sourceValue interface{}, docFieldStruct *docTemplateField) (value interface{}, err error) {

	if docFieldStruct.StructName != "" {
		value, err = getDocValue(sourceValue, docFieldStruct.StructName)
		if err != nil {
			return
		}
		return
	}
	value = sourceValue
	return
}

func getFieldValues(sourceValues []interface{}, docFieldStruct *docTemplateField) (values []interface{}, err error) {
	if len(sourceValues) == 0 {
		return
	}

	if docFieldStruct.StructName != "" {
		for _, sourceValue := range sourceValues {
			var value interface{}
			value, err = getDocValue(sourceValue, docFieldStruct.StructName)
			if err != nil {
				return
			}
			values = append(values, value)
		}
	} else {
		for _, sourceValue := range sourceValues {
			values = append(values, sourceValue)
		}
	}

	return
}

func getDocValue(sourceValue interface{}, docTemplateName string) (value map[string]interface{}, err error) {
	if sourceValue == nil {
		return
	}
	mapV, mapVOk := sourceValue.(map[string]interface{})

	docStruct := getDocTemplate(docTemplateName)
	if !mapVOk {
		if docStruct.Abbreviation == "" {
			err = errors.New("source value to struct error")
			util.Logger.Error("get struct error", zap.Any("sourceValue", sourceValue), zap.Any("struct", docStruct), zap.Error(err))
			return
		}
		mapV = map[string]interface{}{}
		mapV[docStruct.Abbreviation] = sourceValue
	}

	value = map[string]interface{}{}
	for _, docFieldStruct := range docStruct.Fields {
		var v interface{}
		v, err = getFieldData(mapV[docFieldStruct.Name], docFieldStruct)
		if err != nil {
			return
		}
		value[docFieldStruct.Name] = v
	}
	return
}
