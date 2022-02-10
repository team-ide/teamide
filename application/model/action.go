package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

type ActionModel struct {
	Name             string           `json:"name,omitempty" yaml:"name,omitempty"`               // 名称，同一个应用中唯一
	Comment          string           `json:"comment,omitempty" yaml:"comment,omitempty"`         // 注释说明
	Description      string           `json:"description,omitempty" yaml:"description,omitempty"` // 注释说明
	Api              *ActionApi       `json:"api,omitempty" yaml:"api,omitempty"`
	InVariables      []*VariableModel `json:"inVariables,omitempty" yaml:"inVariables,omitempty"` // 输入变量
	OutVariable      *VariableModel   `json:"outVariable,omitempty" yaml:"outVariable,omitempty"` // 输出变量
	Steps            []ActionStep     `json:"steps,omitempty" yaml:"steps,omitempty"`
	ActionJavascript string           `json:"-" yaml:"-"` // Javascript
	WebApiJavascript string           `json:"-" yaml:"-"` // Javascript
}

type ActionApi struct {
	Request  *ApiRequest  `json:"request,omitempty" yaml:"request,omitempty"`   //
	Response *ApiResponse `json:"response,omitempty" yaml:"response,omitempty"` //
}

func TextToActionModel(namePath string, text string) (model *ActionModel, err error) {
	var modelMap map[string]interface{}
	var name string
	name, modelMap, err = TextToModelMap(namePath, text)
	if err != nil {
		return
	}
	model, err = MapToActionModel(modelMap)
	if err != nil {
		return
	}
	model.Name = name
	return
}

func MapToActionModel(modelMap map[string]interface{}) (model *ActionModel, err error) {

	model = &ActionModel{}
	for key, value := range modelMap {
		switch key {
		case "api":
			api := &ActionApi{}
			err = base.ToBean([]byte(base.ToJSON(value)), api)
			if err != nil {
				return
			}
			model.Api = api
			delete(modelMap, "api")
		case "inVariables":
			model.InVariables, err = getVariablesByValue(value)
			if err != nil {
				return
			}
			delete(modelMap, "inVariables")
		case "outVariable":
			outVariable := &VariableModel{}
			err = base.ToBean([]byte(base.ToJSON(value)), outVariable)
			if err != nil {
				return
			}
			model.OutVariable = outVariable
			delete(modelMap, "outVariable")
		case "steps":
			model.Steps, err = getStepsByValue(value)
			if err != nil {
				return
			}
			delete(modelMap, "steps")
		}
	}

	err = base.ToBean([]byte(base.ToJSON(modelMap)), model)
	if err != nil {
		return
	}

	return
}

func getStepsByValue(value interface{}) (steps []ActionStep, err error) {
	if value == nil {
		return
	}
	values, valuesOk := value.([]interface{})
	if !valuesOk {
		return
	}
	for _, valuesOne := range values {
		valuesOneMap, valuesOneMapOk := valuesOne.(map[string]interface{})
		if !valuesOneMapOk || len(valuesOneMap) == 0 {
			err = errors.New(fmt.Sprint("[", valuesOne, "] to step error"))
			return
		}

		var step ActionStep
		step, err = getStepByValue(valuesOneMap)

		if err != nil {
			return
		}
		if step != nil {
			steps = append(steps, step)
		}
	}
	return
}

func getStepByValue(valuesOneMap map[string]interface{}) (step ActionStep, err error) {
	if valuesOneMap == nil {
		return
	}

	formatTryErrorByMap(valuesOneMap)

	var isReturn bool
	var returnVariableName string
	if valuesOneMap["return"] != nil {
		obj := valuesOneMap["return"]
		switch v := obj.(type) {
		case bool:
			isReturn = true
		case string:
			isReturn = true
			returnVariableName = v
		default:
			isReturn = true
			returnVariableName = fmt.Sprint(v)
		}
	}
	delete(valuesOneMap, "return")

	var subSteps []ActionStep
	subSteps, err = getStepsByValue(valuesOneMap["steps"])
	if err != nil {
		return
	}
	delete(valuesOneMap, "steps")

	var variables []*VariableModel
	variables, err = getVariablesByValue(valuesOneMap["variables"])
	if err != nil {
		return
	}
	delete(valuesOneMap, "variables")

	var validates []*ValidateModel
	validates, err = getValidatesByValue(valuesOneMap["validates"])
	if err != nil {
		return
	}
	delete(valuesOneMap, "validates")

	step = nil
	if step == nil {
		step, err = getActionStepLockByMap(valuesOneMap)
		if err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepUnlockByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepErrorByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepSqlSelectByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepSqlInsertByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepSqlUpdateByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepSqlDeleteByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepFileSaveByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepFileGetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepRedisSetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepRedisGetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepRedisDelByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepRedisExpireByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepRedisExpireatByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getActionStepActionByMap(valuesOneMap); err != nil {
			return
		}
	}
	baseStep := &ActionStepBase{
		Validates:          validates,
		Variables:          variables,
		Steps:              subSteps,
		Return:             isReturn,
		ReturnVariableName: returnVariableName,
	}
	err = base.ToBean([]byte(base.ToJSON(valuesOneMap)), baseStep)
	if err != nil {
		return
	}
	if step == nil {
		step = baseStep
	}
	if step != nil {
		step.SetBase(baseStep)
	}
	return
}

func ActionModelToText(model *ActionModel) (text string, err error) {
	text, err = ModelToText(model)
	if err != nil {
		return
	}
	return
}
