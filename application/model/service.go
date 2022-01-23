package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

type ServiceModel struct {
	Name              string           `json:"name,omitempty" yaml:"name,omitempty"`               // 名称，同一个应用中唯一
	Comment           string           `json:"comment,omitempty" yaml:"name,omitempty"`            // 注释说明
	Description       string           `json:"description,omitempty" yaml:"description,omitempty"` // 注释说明
	Api               *ServiceApi      `json:"api,omitempty" yaml:"api,omitempty"`
	InVariables       []*VariableModel `json:"inVariables,omitempty" yaml:"inVariables,omitempty"` // 输入变量
	OutVariable       *VariableModel   `json:"outVariable,omitempty" yaml:"outVariable,omitempty"` // 输出变量
	Steps             []ServiceStep    `json:"steps,omitempty" yaml:"steps,omitempty"`
	ServiceJavascript string           `json:"-" yaml:"-"` // Javascript
	WebApiJavascript  string           `json:"-" yaml:"-"` // Javascript
}

type ServiceApi struct {
	Request  *ApiRequest  `json:"request,omitempty" yaml:"request,omitempty"`   //
	Response *ApiResponse `json:"response,omitempty" yaml:"response,omitempty"` //
}

func TextToServiceModel(namePath string, text string) (model *ServiceModel, err error) {
	var modelMap map[string]interface{}
	var name string
	name, modelMap, err = TextToModelMap(namePath, text)
	if err != nil {
		return
	}
	model = &ServiceModel{
		Name: name,
	}
	for key, value := range modelMap {
		switch key {
		case "api":
			api := &ServiceApi{}
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

func getStepsByValue(value interface{}) (steps []ServiceStep, err error) {
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

		var step ServiceStep
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

func getStepByValue(valuesOneMap map[string]interface{}) (step ServiceStep, err error) {
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

	var subSteps []ServiceStep
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
		step, err = getServiceStepLockByMap(valuesOneMap)
		if err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepUnlockByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepErrorByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepSqlSelectByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepSqlInsertByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepSqlUpdateByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepSqlDeleteByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepFileSaveByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepFileGetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepRedisSetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepRedisGetByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepRedisDelByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepRedisExpireByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepRedisExpireatByMap(valuesOneMap); err != nil {
			return
		}
	}
	if step == nil {
		if step, err = getServiceStepServiceByMap(valuesOneMap); err != nil {
			return
		}
	}
	baseStep := &ServiceStepBase{
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
