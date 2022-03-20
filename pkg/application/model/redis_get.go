package model

type RedisGet struct {
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`             // 建
	ValueType string `json:"valueType,omitempty" yaml:"valueType,omitempty"` // 值类型
}

type ActionStepRedisGet struct {
	Base *ActionStepBase

	RedisGet         *RedisGet `json:"redisGet,omitempty" yaml:"redisGet,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepRedisGet) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepRedisGet) SetBase(v *ActionStepBase) {
	this_.Base = v
}
