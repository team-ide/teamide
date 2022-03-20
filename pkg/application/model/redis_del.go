package model

type RedisDel struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"` // 建
}

type ActionStepRedisDel struct {
	Base *ActionStepBase

	RedisDel         *RedisDel `json:"redisDel,omitempty" yaml:"redisDel,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepRedisDel) GetBase() *ActionStepBase {
	return this_.Base
}
func (this_ *ActionStepRedisDel) SetBase(v *ActionStepBase) {
	this_.Base = v
}
