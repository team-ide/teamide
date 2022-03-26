package model

// 设置一个key在"timestamp"(时间戳(秒))之后过期
type RedisExpireat struct {
	Key       string `json:"key,omitempty" yaml:"key,omitempty"`             // 建
	Timestamp string `json:"timestamp,omitempty" yaml:"timestamp,omitempty"` // 时间戳
}

type ActionStepRedisExpireat struct {
	Base *ActionStepBase

	RedisExpireat    *RedisExpireat `json:"redisExpireat,omitempty" yaml:"redisExpireat,omitempty"`       // 执行 SQL DELETE 操作
	VariableName     string         `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string         `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepRedisExpireat) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepRedisExpireat) SetBase(v *ActionStepBase) {
	this_.Base = v
}
