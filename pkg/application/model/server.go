package model

import (
	"teamide/pkg/application/base"
)

func TextToServerWebModel(namePath string, text string) (model *ServerWebModel, err error) {

	var name string
	var modelMap map[string]interface{}
	name, modelMap, err = TextToModelMap(namePath, text)
	if err != nil {
		return
	}
	var token *ServerWebToken
	if modelMap["token"] != nil {
		tokenMap, tokenMapOk := modelMap["token"].(map[string]interface{})
		if tokenMapOk {
			var variables []*VariableModel
			variables, err = getVariablesByValue(tokenMap["variables"])
			if err != nil {
				return
			}
			delete(tokenMap, "variables")

			var validates []*ValidateModel
			validates, err = getValidatesByValue(tokenMap["validates"])
			if err != nil {
				return
			}
			delete(tokenMap, "validates")

			token = &ServerWebToken{
				Variables: variables,
				Validates: validates,
			}

			err = base.ToBean([]byte(base.ToJSON(tokenMap)), token)
			if err != nil {
				return
			}

		}
		delete(modelMap, "token")
	}

	model = &ServerWebModel{
		Name:  name,
		Token: token,
	}

	err = base.ToBean([]byte(base.ToJSON(modelMap)), model)
	if err != nil {
		return
	}
	if name == "default" {
		name = ""
	}
	model.Name = name
	return
}

func ServerWebModelToText(model *ServerWebModel) (text string, err error) {
	text, err = ToText(model)
	if err != nil {
		return
	}
	return
}
