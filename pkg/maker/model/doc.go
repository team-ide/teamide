package model

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	"sync"
	"teamide/pkg/util"
)

type docTemplate struct {
	Name           string `json:"name"`
	Inline         string `json:"inline"`
	inlineNewModel func() interface{}
	Abbreviation   string              `json:"abbreviation"`
	Comment        string              `json:"comment"`
	Fields         []*docTemplateField `json:"fields"`
	newModel       func() interface{}
	newModels      func() interface{}
	appendModel    func(values interface{}, value interface{})
}

type docTemplateField struct {
	Name       string            `json:"name"`
	Comment    string            `json:"comment"`
	IsList     bool              `json:"isList"`
	StructName string            `json:"structName"`
	Default    interface{}       `json:"default"` // 默认值
	Sons       []*docTemplateSon `json:"sons"`
}

type docTemplateSon struct {
	MatchKey   string      `json:"matchKey"`
	MatchValue interface{} `json:"matchValue"`
	StructName string      `json:"structName"`
	newModel   func() interface{}
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
		if len(docFieldStruct.Sons) > 0 {
			err = appendNodeSonsFieldValue(mapNode, source, docFieldStruct.Sons, options)
			if err != nil {
				return
			}
			continue
		}
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
		err = appendNodeFieldValues(node, list, docFieldStruct, options)
		if err != nil {
			return
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
	mapV, mapVOk := value.(map[string]interface{})
	if docFieldStruct.StructName != "" {
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

func getFieldSon(value map[string]interface{}, sons []*docTemplateSon) (son *docTemplateSon, err error) {
	for _, one := range sons {
		if one.StructName == "" {
			continue
		}
		findValue, find := value[one.MatchKey]
		if !find {
			continue
		}
		if one.MatchValue != nil && one.MatchValue != findValue {
			continue
		}
		docStruct := getDocTemplate(one.StructName)
		if docStruct == nil {
			err = errors.New("doc template [" + one.StructName + "] is not exist")
			return
		}
		son = one
		break
	}
	return
}

func appendNodeSonsFieldValue(node *yaml.Node, value map[string]interface{}, sons []*docTemplateSon, options *docOptions) (err error) {
	//util.Logger.Info("append son field", zap.Any("value", value))
	son, err := getFieldSon(value, sons)
	if err != nil {
		return
	}
	if son == nil {
		return
	}
	docStruct := getDocTemplate(son.StructName)
	for _, docFieldStruct := range docStruct.Fields {
		//util.Logger.Info("append son field", zap.Any("key", docFieldStruct.Name))
		if len(docFieldStruct.Sons) > 0 {
			err = appendNodeSonsFieldValue(node, value, docFieldStruct.Sons, options)
			if err != nil {
				return
			}
			continue
		}
		err = appendNodeField(node, value[docFieldStruct.Name], docFieldStruct, options)
		if err != nil {
			return
		}
	}
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
		if len(docFieldStruct.Sons) > 0 {
			err = appendNodeSonsFieldValue(mapNode, value, docFieldStruct.Sons, options)
			if err != nil {
				return
			}
			continue
		}
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
	source := map[string]interface{}{}

	err = yaml.Unmarshal([]byte(text), source)
	if err != nil {
		util.Logger.Error("text to source error", zap.Any("text", text), zap.Error(err))
		return
	}

	err = appendData(source, model, docTemplateName)
	if err != nil {
		util.Logger.Error("append data error", zap.Any("source", source), zap.Error(err))
		return
	}

	return
}

func appendData(source map[string]interface{}, data interface{}, docTemplateName string) (err error) {
	if source == nil {
		source = map[string]interface{}{}
	}
	docStruct := getDocTemplate(docTemplateName)
	err = appendDataByStruct(source, data, docStruct)
	if err != nil {
		return
	}

	return
}

func appendDataByStruct(source map[string]interface{}, data interface{}, docStruct *docTemplate) (err error) {

	for _, docFieldStruct := range docStruct.Fields {
		var v interface{}
		v, err = getFieldData(source[docFieldStruct.Name], docFieldStruct)
		if err != nil {
			return
		}
		setFieldValue(data, docFieldStruct.Name, v)
	}
	return
}

func setFieldValue(data interface{}, name string, value interface{}) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("set field error", zap.Any("name", name), zap.Any("data", data), zap.Any("value", value), zap.Any("error", e))
		}
	}()
	var fV reflect.Value
	dV := reflect.ValueOf(data) // 取得struct变量的指针
	dT := reflect.TypeOf(data)
	num := dT.Elem().NumField()
	for i := 0; i < num; i++ {
		field := dT.Elem().Field(i)
		jsonKey := field.Tag.Get("json")
		if jsonKey != "" {
			if strings.Contains(jsonKey, ",") {
				jsonKey = strings.Split(jsonKey, ",")[0]
			}
			if strings.EqualFold(jsonKey, name) {
				fV = dV.Elem().Field(i)
				break
			}
		}
		if strings.EqualFold(field.Name, name) {
			fV = dV.Elem().Field(i)
			break
		}
	}
	if !canNotOut(value) {
		fV.Set(reflect.ValueOf(value))
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

		docStruct := getDocTemplate(docFieldStruct.StructName)
		var mList interface{}
		if docStruct.newModels != nil {
			mList = docStruct.newModels()
		}

		for _, sourceValue := range sourceValues {
			var value interface{}
			value, err = getDocValue(sourceValue, docFieldStruct.StructName)
			if err != nil {
				return
			}
			if value != nil {
				if docStruct.newModels != nil {
					docStruct.appendModel(mList, value)
				} else {
					values = append(values, value)
				}
			}
		}
	} else {
		for _, sourceValue := range sourceValues {
			values = append(values, sourceValue)
		}
	}

	return
}

func getDocValue(sourceValue interface{}, docTemplateName string) (value interface{}, err error) {
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

	var sonNewModel func() interface{}
	var inline string
	var inlineValue interface{}
	for _, docFieldStruct := range docStruct.Fields {
		if len(docFieldStruct.Sons) > 0 {
			var sonDocStruct *docTemplate
			sonDocStruct, err = getSonInfo(mapV, docFieldStruct)
			if err != nil {
				return
			}
			if sonDocStruct == nil {
				continue
			}
			if sonDocStruct.Inline != "" {
				if sonDocStruct.inlineNewModel != nil {
					inlineValue = sonDocStruct.inlineNewModel()
				} else {
					inlineValue = docStruct.inlineNewModel()
				}
			}
			inline = sonDocStruct.Inline
			sonNewModel = sonDocStruct.newModel
		}
	}
	if sonNewModel != nil {
		value = sonNewModel()
	} else {
		value = docStruct.newModel()
	}
	for _, docFieldStruct := range docStruct.Fields {
		var v interface{}

		var sonDocStruct *docTemplate
		if len(docFieldStruct.Sons) > 0 {
			sonDocStruct, err = getSonInfo(mapV, docFieldStruct)
			if err != nil {
				return
			}
			if sonDocStruct == nil {
				continue
			}
		}

		if sonDocStruct != nil {
			for _, sonDocFieldStruct := range sonDocStruct.Fields {
				v, err = getFieldData(mapV[sonDocFieldStruct.Name], sonDocFieldStruct)
				if err != nil {
					return
				}
				if !canNotOut(v) {
					setFieldValue(value, sonDocFieldStruct.Name, v)
				}
			}
		} else {
			v, err = getFieldData(mapV[docFieldStruct.Name], docFieldStruct)
			if err != nil {
				return
			}
			if !canNotOut(v) {
				if inlineValue != nil {
					setFieldValue(inlineValue, docFieldStruct.Name, v)
				} else {
					setFieldValue(value, docFieldStruct.Name, v)
				}
			}
		}
	}

	if inline != "" && inlineValue != nil {
		setFieldValue(value, inline, inlineValue)
	}

	return
}

func getSonInfo(sourceValue map[string]interface{}, docFieldStruct *docTemplateField) (sonDocStruct *docTemplate, err error) {
	if len(docFieldStruct.Sons) > 0 {
		var son *docTemplateSon
		son, err = getFieldSon(sourceValue, docFieldStruct.Sons)
		if err != nil {
			return
		}
		if son != nil {
			sonDocStruct = getDocTemplate(son.StructName)
		}
	}
	return
}
