package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

type ActionStep interface {
	GetBase() *ActionStepBase
	SetBase(*ActionStepBase)
}

type ActionStepBase struct {
	Name               string           `json:"name,omitempty" yaml:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment            string           `json:"comment,omitempty" yaml:"comment,omitempty"`                       // 注释
	If                 string           `json:"if,omitempty" yaml:"if,omitempty"`                                 // 满足 if 条件
	Validates          []*ValidateModel `json:"validates,omitempty" yaml:"validates,omitempty"`                   // 验证
	Variables          []*VariableModel `json:"variables,omitempty" yaml:"variables,omitempty"`                   // 变量
	Steps              []ActionStep     `json:"steps,omitempty" yaml:"steps,omitempty"`                           // 子阶段
	Return             bool             `json:"return,omitempty" yaml:"return,omitempty"`                         // 是否退出
	ReturnVariableName string           `json:"returnVariableName,omitempty" yaml:"returnVariableName,omitempty"` // 是否退出
	TryError           *ErrorModel      `json:"tryError,omitempty" yaml:"tryError"`
}

func (this_ *ActionStepBase) GetBase() *ActionStepBase {
	return this_
}

func (this_ *ActionStepBase) SetBase(v *ActionStepBase) {

}

type ActionStepLock struct {
	Base *ActionStepBase

	Lock *LockModel `json:"lock,omitempty" yaml:"lock,omitempty"` // 锁
}

type LockModel struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"` // 锁名称 只用作定义 如需主动解锁 必须定义名称 不自动解锁 最 所有阶段执行结束 自动解锁
	Type string `json:"type,omitempty" yaml:"type,omitempty"` // 锁类型 不填写：直接使用内存锁，redis：使用redis key设置分布式锁
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`   // 锁Key
}

func (this_ *ActionStepLock) GetBase() *ActionStepBase {
	return this_.Base
}
func (this_ *ActionStepLock) SetBase(v *ActionStepBase) {
	this_.Base = v
}

type ActionStepUnlock struct {
	Base *ActionStepBase

	Unlock *UnlockModel `json:"unlock,omitempty" yaml:"unlock,omitempty"` // 解锁
}

type UnlockModel struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"` // 锁名称 只用作定义 如需主动解锁 必须定义名称 不自动解锁 最 所有阶段执行结束 自动解锁
}

func (this_ *ActionStepUnlock) GetBase() *ActionStepBase {
	return this_.Base
}
func (this_ *ActionStepUnlock) SetBase(v *ActionStepBase) {
	this_.Base = v
}

type ActionStepError struct {
	Base *ActionStepBase

	Error *ErrorModel `json:"error,omitempty" yaml:"error,omitempty"` // 异常
}

func (this_ *ActionStepError) GetBase() *ActionStepBase {
	return this_.Base
}
func (this_ *ActionStepError) SetBase(v *ActionStepBase) {
	this_.Base = v
}

type ActionStepAction struct {
	Base *ActionStepBase

	Action           *StepActionModel `json:"action,omitempty" yaml:"action,omitempty"`                     // 解锁
	VariableName     string           `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string           `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

type StepActionModel struct {
	Name          string           `json:"name,omitempty" yaml:"name,omitempty"`                   // 服务名称
	CallVariables []*VariableModel `json:"callVariables,omitempty" yaml:"callVariables,omitempty"` // 变量
}

func (this_ *ActionStepAction) GetBase() *ActionStepBase {
	return this_.Base
}
func (this_ *ActionStepAction) SetBase(v *ActionStepBase) {
	this_.Base = v
}

func getActionStepLockByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["lock"] == nil {
			return
		}
		switch data := v["lock"].(type) {
		case map[string]interface{}:
			v["lock"] = data
		default:
			v["lock"] = map[string]interface{}{
				"name": data,
			}
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step lock error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepLock{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return step, err
}

func getActionStepUnlockByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["unlock"] == nil {
			return
		}
		switch data := v["unlock"].(type) {
		case map[string]interface{}:
			v["unlock"] = data
		default:
			v["unlock"] = map[string]interface{}{
				"name": data,
			}
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step unlock error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepUnlock{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getActionStepErrorByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["error"] == nil {
			return
		}
		switch data := v["error"].(type) {
		case map[string]interface{}:
			v["error"] = data
		default:
			v["error"] = map[string]interface{}{
				"name": data,
			}
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step error error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepError{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func formatTryErrorByMap(value interface{}) (err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["tryError"] == nil {
			return
		}
		switch data := v["tryError"].(type) {
		case map[string]interface{}:
			v["tryError"] = data
		default:
			v["tryError"] = map[string]interface{}{
				"name": data,
			}
		}
	}
	return
}
func getActionStepActionByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	var callVariables []*VariableModel
	switch v := value.(type) {
	case map[string]interface{}:
		if v["action"] == nil {
			return
		}
		var actionMap map[string]interface{}
		switch data := v["action"].(type) {
		case map[string]interface{}:
			actionMap = data
		default:
			actionMap = map[string]interface{}{
				"name": data,
			}
		}
		if actionMap["callVariables"] != nil {
			callVariables, err = getVariablesByValue(actionMap["callVariables"])
			if err != nil {
				return
			}
			delete(actionMap, "callVariables")
		}
		v["action"] = actionMap
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step action error"))
	}
	if err != nil {
		return
	}
	stepAction := &ActionStepAction{}
	err = base.ToBean([]byte(base.ToJSON(value)), stepAction)
	if err != nil {
		return
	}
	stepAction.Action.CallVariables = callVariables
	step = stepAction
	return
}
