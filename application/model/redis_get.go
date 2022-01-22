package model

type RedisGet struct {
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`             // 建
	ValueType string `json:"valueType,omitempty" yaml:"valueType,omitempty"` // 值类型
}

type ServiceStepRedisGet struct {
	Base *ServiceStepBase

	RedisGet         *RedisGet `json:"redisGet,omitempty" yaml:"redisGet,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ServiceStepRedisGet) GetBase() *ServiceStepBase {
	return this_.Base
}

func (this_ *ServiceStepRedisGet) SetBase(v *ServiceStepBase) {
	this_.Base = v
}
