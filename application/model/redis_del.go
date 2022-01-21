package model

type RedisDel struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"` // 建
}

type ServiceStepRedisDel struct {
	Base *ServiceStepBase

	RedisDel         *RedisDel `json:"redisDel,omitempty" yaml:"redisDel,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ServiceStepRedisDel) GetBase() *ServiceStepBase {
	return this_.Base
}
func (this_ *ServiceStepRedisDel) SetBase(v *ServiceStepBase) {
	this_.Base = v
}
