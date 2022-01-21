package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

type ServiceStep interface {
	GetBase() *ServiceStepBase
	SetBase(*ServiceStepBase)
}

type ServiceStepBase struct {
	Name               string           `json:"name,omitempty" yaml:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment            string           `json:"comment,omitempty" yaml:"comment,omitempty"`                       // 注释
	If                 string           `json:"if,omitempty" yaml:"if,omitempty"`                                 // 满足 if 条件
	Validates          []*ValidateModel `json:"validates,omitempty" yaml:"validates,omitempty"`                   // 验证
	Variables          []*VariableModel `json:"variables,omitempty" yaml:"variables,omitempty"`                   // 变量
	Steps              []ServiceStep    `json:"steps,omitempty" yaml:"steps,omitempty"`                           // 子阶段
	Return             bool             `json:"return,omitempty" yaml:"return,omitempty"`                         // 是否退出
	ReturnVariableName string           `json:"returnVariableName,omitempty" yaml:"returnVariableName,omitempty"` // 是否退出
	TryError           *ErrorModel      `json:"tryError,omitempty" yaml:"tryError"`
}

func (this_ *ServiceStepBase) GetBase() *ServiceStepBase {
	return this_
}

func (this_ *ServiceStepBase) SetBase(v *ServiceStepBase) {

}

type ServiceStepLock struct {
	Base *ServiceStepBase

	Lock *LockModel `json:"lock,omitempty" yaml:"lock,omitempty"` // 锁
}

type LockModel struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"` // 锁名称 只用作定义 如需主动解锁 必须定义名称 不自动解锁 最 所有阶段执行结束 自动解锁
	Type string `json:"type,omitempty" yaml:"type,omitempty"` // 锁类型 不填写：直接使用内存锁，redis：使用redis key设置分布式锁
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`   // 锁Key
}

func (this_ *ServiceStepLock) GetBase() *ServiceStepBase {
	return this_.Base
}
func (this_ *ServiceStepLock) SetBase(v *ServiceStepBase) {
	this_.Base = v
}

type ServiceStepUnlock struct {
	Base *ServiceStepBase

	Unlock *UnlockModel `json:"unlock,omitempty" yaml:"unlock,omitempty"` // 解锁
}

type UnlockModel struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"` // 锁名称 只用作定义 如需主动解锁 必须定义名称 不自动解锁 最 所有阶段执行结束 自动解锁
}

func (this_ *ServiceStepUnlock) GetBase() *ServiceStepBase {
	return this_.Base
}
func (this_ *ServiceStepUnlock) SetBase(v *ServiceStepBase) {
	this_.Base = v
}

type ServiceStepError struct {
	Base *ServiceStepBase

	Error *ErrorModel `json:"error,omitempty" yaml:"error,omitempty"` // 异常
}

func (this_ *ServiceStepError) GetBase() *ServiceStepBase {
	return this_.Base
}
func (this_ *ServiceStepError) SetBase(v *ServiceStepBase) {
	this_.Base = v
}

type ServiceStepService struct {
	Base *ServiceStepBase

	Service          *StepServiceModel `json:"service,omitempty" yaml:"service,omitempty"`                   // 解锁
	VariableName     string            `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string            `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

type StepServiceModel struct {
	Name          string           `json:"name,omitempty" yaml:"name,omitempty"`                   // 服务名称
	CallVariables []*VariableModel `json:"callVariables,omitempty" yaml:"callVariables,omitempty"` // 变量
}

func (this_ *ServiceStepService) GetBase() *ServiceStepBase {
	return this_.Base
}
func (this_ *ServiceStepService) SetBase(v *ServiceStepBase) {
	this_.Base = v
}

func getServiceStepLockByMap(value interface{}) (step ServiceStep, err error) {
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
		err = errors.New(fmt.Sprint("[", v, "] to service step lock error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepLock{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return step, err
}

func getServiceStepUnlockByMap(value interface{}) (step ServiceStep, err error) {
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
		err = errors.New(fmt.Sprint("[", v, "] to service step unlock error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepUnlock{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getServiceStepErrorByMap(value interface{}) (step ServiceStep, err error) {
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
		err = errors.New(fmt.Sprint("[", v, "] to service step error error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepError{}
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
func getServiceStepServiceByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	var callVariables []*VariableModel
	switch v := value.(type) {
	case map[string]interface{}:
		if v["service"] == nil {
			return
		}
		var serviceMap map[string]interface{}
		switch data := v["service"].(type) {
		case map[string]interface{}:
			serviceMap = data
		default:
			serviceMap = map[string]interface{}{
				"name": data,
			}
		}
		if serviceMap["callVariables"] != nil {
			callVariables, err = getVariablesByValue(serviceMap["callVariables"])
			if err != nil {
				return
			}
			delete(serviceMap, "callVariables")
		}
		v["service"] = serviceMap
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step service error"))
	}
	if err != nil {
		return
	}
	stepService := &ServiceStepService{}
	err = base.ToBean([]byte(base.ToJSON(value)), stepService)
	if err != nil {
		return
	}
	stepService.Service.CallVariables = callVariables
	step = stepService
	return
}
