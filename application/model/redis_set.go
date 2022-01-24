package model

type RedisSet struct {
	Key   string `json:"key,omitempty" yaml:"key,omitempty"`     // 建
	Value string `json:"value,omitempty" yaml:"value,omitempty"` // 值
}

type ActionStepRedisSet struct {
	Base *ActionStepBase

	RedisSet         *RedisSet `json:"redisSet,omitempty" yaml:"redisSet,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepRedisSet) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepRedisSet) SetBase(v *ActionStepBase) {
	this_.Base = v
}
