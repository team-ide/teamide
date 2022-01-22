package model

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

func TextToModel(namePath string, text string, model interface{}) (name string, err error) {
	name = namePath
	if strings.HasSuffix(namePath, ".json") {
		name = strings.TrimSuffix(name, ".json")
		err = json.Unmarshal([]byte(text), model)
	} else {
		name = strings.TrimSuffix(name, ".yaml")
		name = strings.TrimSuffix(name, ".yml")
		err = yaml.Unmarshal([]byte(text), model)
	}
	return
}

func TextToModelMap(namePath string, text string) (name string, res map[string]interface{}, err error) {
	name = namePath
	res = make(map[string]interface{})
	if strings.HasSuffix(namePath, ".json") {
		name = strings.TrimSuffix(name, ".json")
		err = json.Unmarshal([]byte(text), &res)
	} else {
		name = strings.TrimSuffix(name, ".yaml")
		name = strings.TrimSuffix(name, ".yml")
		v := yaml.MapSlice{} // 用于接收解析的 yaml 数据
		err = yaml.Unmarshal([]byte(text), &v)
		if err != nil {
			return
		}
		appendMapSlice(v, res)
	}
	return
}

func appendMapSlice(mapSlice yaml.MapSlice, data map[string]interface{}) {
	if mapSlice == nil {
		return
	}
	for _, item := range mapSlice {
		key := item.Key
		switch value := item.Value.(type) { // value 表示 item.Value 转换成对应 type 的值
		case yaml.MapSlice: // item.Value 是yaml.MapSlice类型
			mapValue := make(map[string]interface{})
			appendMapSlice(value, mapValue)
			data[key.(string)] = mapValue
		case []interface{}: // []interface{} 类型
			listValue := appendMapSliceList(value)
			data[key.(string)] = listValue
		default: // 未知类型

			data[key.(string)] = value
		}
	}
}

func appendMapSliceList(list []interface{}) []interface{} {
	if list == nil {
		return nil
	}
	listValue := []interface{}{}
	for _, subint := range list {
		switch subV := subint.(type) { // value 表示 item.Value 转换成对应 type 的值
		case yaml.MapSlice: // item.Value 是yaml.MapSlice类型
			mapValue := make(map[string]interface{})
			appendMapSlice(subV, mapValue)
			listValue = append(listValue, mapValue)
		case []interface{}:
			subL := appendMapSliceList(subV)
			listValue = append(listValue, subL)
		default: // 未知类型
			listValue = append(listValue, subV)
		}
	}
	return listValue
}
func ModelToText(model map[string]interface{}) (text string, err error) {
	model = make(map[string]interface{})
	var bs []byte
	bs, err = yaml.Marshal(&model)
	if err != nil {
		return
	}
	text = string(bs)
	return
}
