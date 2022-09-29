package model

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"teamide/pkg/util"
)

type modelNodeStruct struct {
	Name         string                  `json:"name"`
	Comment      string                  `json:"comment"`
	Abbreviation string                  `json:"abbreviation"`
	Fields       []*modelNodeFieldStruct `json:"fields"`
}

type modelNodeFieldStruct struct {
	Name    string           `json:"name"`
	Comment string           `json:"comment"`
	IsList  bool             `json:"isList"`
	Struct  *modelNodeStruct `json:"struct"`
}

type docOptions struct {
	omitEmpty  bool
	outComment bool
}

func toText(data map[string]interface{}, docStruct *modelNodeStruct, options *docOptions) (text string, err error) {
	if data == nil {
		data = map[string]interface{}{}
	}
	if options == nil {
		options = &docOptions{}
	}

	node := &yaml.Node{
		Kind: yaml.DocumentNode,
	}
	err = appendNode(node, data, docStruct, options)
	if err != nil {
		return
	}

	bytes, err := yaml.Marshal(node)
	if err != nil {
		util.Logger.Error("node to json error", zap.Any("node", node), zap.Error(err))
		return
	}

	text = string(bytes)

	return
}

func appendNode(node *yaml.Node, data map[string]interface{}, docStruct *modelNodeStruct, options *docOptions) (err error) {

	var mapNode = &yaml.Node{
		Kind: 4,
	}

	if options.outComment {
		mapNode.HeadComment = docStruct.Comment
	}
	node.Content = append(node.Content, mapNode)

	for _, docFieldStruct := range docStruct.Fields {
		err = appendNodeField(mapNode, data[docFieldStruct.Name], docFieldStruct, options)
		if err != nil {
			return
		}
	}
	return
}

func appendNodeField(node *yaml.Node, value interface{}, docFieldStruct *modelNodeFieldStruct, options *docOptions) (err error) {

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
		if docFieldStruct.Struct != nil || docFieldStruct.IsList {
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

func appendNodeFieldValue(node *yaml.Node, value interface{}, docFieldStruct *modelNodeFieldStruct, options *docOptions) (err error) {

	if docFieldStruct.Struct != nil {
		mapV, mapVOk := value.(map[string]interface{})
		if mapVOk {
			err = appendNodeValue(node, mapV, docFieldStruct.Struct, options)
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

func appendNodeFieldValues(node *yaml.Node, values []interface{}, docFieldStruct *modelNodeFieldStruct, options *docOptions) (err error) {
	if len(values) == 0 {
		return
	}

	var listNode = &yaml.Node{
		Kind: 2,
	}
	node.Content = append(node.Content, listNode)

	if docFieldStruct.Struct != nil {
		for _, value := range values {
			mapV, mapVOk := value.(map[string]interface{})
			if mapVOk {
				err = appendNodeValue(listNode, mapV, docFieldStruct.Struct, options)
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

func appendNodeValue(node *yaml.Node, value map[string]interface{}, docStruct *modelNodeStruct, options *docOptions) (err error) {
	if value == nil || len(value) == 0 {
		node.Content = append(node.Content, &yaml.Node{
			Kind: 8,
		})
		return
	}
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
