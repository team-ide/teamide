package model

type SqlDelete struct {
	Database string      `json:"database,omitempty" yaml:"database,omitempty"` // 库名
	Table    string      `json:"table,omitempty" yaml:"table,omitempty"`       // 表名
	Wheres   []*SqlWhere `json:"wheres,omitempty" yaml:"wheres,omitempty"`     // 条件
}

type ActionStepSqlDelete struct {
	Base *ActionStepBase

	SqlDelete        *SqlDelete `json:"sqlDelete,omitempty" yaml:"sqlDelete,omitempty"`               // 执行 SQL DELETE 操作
	VariableName     string     `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string     `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepSqlDelete) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepSqlDelete) SetBase(v *ActionStepBase) {
	this_.Base = v
}
