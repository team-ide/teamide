package model

type FileGet struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty"` // 目录
}

type ActionStepFileGet struct {
	Base *ActionStepBase

	FileGet          *FileGet `json:"fileGet,omitempty" yaml:"fileGet,omitempty"`                   // 执行 SQL DELETE 操作
	VariableName     string   `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string   `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepFileGet) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepFileGet) SetBase(v *ActionStepBase) {
	this_.Base = v
}
