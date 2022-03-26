package application

import (
	"encoding/json"
	"teamide/pkg/application/model"
)

func GetContextByText(text string) (context *model.ModelContext, err error) {
	var contextMap map[string]interface{}
	err = json.Unmarshal([]byte(text), &contextMap)
	if err != nil {
		return
	}
	context = &model.ModelContext{}
	if contextMap["actions"] != nil {
		list, listOk := contextMap["actions"].([]interface{})
		if listOk {
			for _, one := range list {
				oneMap, oneMapOk := one.(map[string]interface{})
				if !oneMapOk {
					continue
				}
				var oneModel *model.ActionModel
				oneModel, err = model.MapToActionModel(oneMap)
				if err != nil {
					return
				}
				context.Actions = append(context.Actions, oneModel)
			}
		}

		delete(contextMap, "actions")
	}
	if contextMap["tests"] != nil {
		list, listOk := contextMap["tests"].([]interface{})
		if listOk {
			for _, one := range list {
				oneMap, oneMapOk := one.(map[string]interface{})
				if !oneMapOk {
					continue
				}
				var oneModel *model.TestModel
				oneModel, err = model.MapToTestModel(oneMap)
				if err != nil {
					return
				}
				context.Tests = append(context.Tests, oneModel)
			}
		}

		delete(contextMap, "tests")
	}
	var bs []byte
	bs, err = json.Marshal(contextMap)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, context)
	if err != nil {
		return
	}
	return
}

func GetTextByContext(context *model.ModelContext) (text string, err error) {
	var bs []byte
	bs, err = json.Marshal(context)
	if err != nil {
		return
	}
	text = string(bs)
	return
}
