package model

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
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
		err = yaml.Unmarshal([]byte(text), &res)
	}
	if err != nil {
		return
	}
	return
}

func ToText(model interface{}) (text string, err error) {
	var bs []byte
	bs, err = yaml.Marshal(&model)
	if err != nil {
		return
	}
	text = string(bs)
	return
}

func MapToModel(modelType *ModelType, data map[string]interface{}) (model interface{}, err error) {

	var userJSON = true
	switch modelType {
	case MODEL_TYPE_CONSTANT:
		model = &ConstantModel{}
	case MODEL_TYPE_ERROR:
		model = &ErrorModel{}
	case MODEL_TYPE_STRUCT:
		model = &StructModel{}
	case MODEL_TYPE_SERVER_WEB:
		model = &ServerWebModel{}
	case MODEL_TYPE_DICTIONARY:
		model = &DictionaryModel{}
	case MODEL_TYPE_DATASOURCE_DATABASE:
		model = &DatasourceDatabase{}
	case MODEL_TYPE_DATASOURCE_REDIS:
		model = &DatasourceRedis{}
	case MODEL_TYPE_DATASOURCE_KAFKA:
		model = &DatasourceKafka{}
	case MODEL_TYPE_DATASOURCE_ZOOKEEPER:
		model = &DatasourceZookeeper{}
	case MODEL_TYPE_ACTION:
		var m *ActionModel
		m, err = MapToActionModel(data)
		if err != nil {
			return
		}
		model = m
		userJSON = false
	case MODEL_TYPE_TEST:
		var m *TestModel
		m, err = MapToTestModel(data)
		if err != nil {
			return
		}
		model = m
		userJSON = false
	}
	if userJSON {
		var bs []byte
		bs, err = json.Marshal(data)
		if err != nil {
			return
		}
		err = json.Unmarshal(bs, model)
		if err != nil {
			return
		}
	}
	return
}
