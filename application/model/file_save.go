package model

type FileSave struct {
	Dir    string `json:"dir,omitempty" yaml:"dir,omitempty"`       // 目录
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`     // 文件名
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`     // 文件类型
	Reader string `json:"reader,omitempty" yaml:"reader,omitempty"` // 值
	Bytes  string `json:"bytes,omitempty" yaml:"bytes,omitempty"`   // 值
}

type ServiceStepFileSave struct {
	Base *ServiceStepBase

	FileSave         *FileSave `json:"fileSave,omitempty" yaml:"fileSave,omitempty"`                 // 执行 SQL DELETE 操作
	VariableName     string    `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string    `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ServiceStepFileSave) GetBase() *ServiceStepBase {
	return this_.Base
}

func (this_ *ServiceStepFileSave) SetBase(v *ServiceStepBase) {
	this_.Base = v
}
