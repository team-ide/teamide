package model

type SqlUpdate struct {
	Database string             `json:"database,omitempty" yaml:"database,omitempty"` // 库名
	Table    string             `json:"table,omitempty" yaml:"table,omitempty"`       // 表名
	Columns  []*SqlUpdateColumn `json:"columns,omitempty" yaml:"columns,omitempty"`   // 更新字段
	Wheres   []*SqlWhere        `json:"wheres,omitempty" yaml:"wheres,omitempty"`     // 条件
}

type SqlUpdateColumn struct {
	IfScript    string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`       // 条件  满足该条件 添加
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`               // 字段名称
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`             // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`       // 必填
	IgnoreEmpty bool   `json:"ignoreEmpty,omitempty" yaml:"ignoreEmpty,omitempty"` // 忽略空值，如果忽略，则值是null、空字符串、0不设值
}

type ActionStepSqlUpdate struct {
	Base *ActionStepBase

	SqlUpdate        *SqlUpdate `json:"sqlUpdate,omitempty" yaml:"sqlUpdate,omitempty"`               // 执行 SQL UPDATE 操作
	VariableName     string     `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string     `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepSqlUpdate) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepSqlUpdate) SetBase(v *ActionStepBase) {
	this_.Base = v
}
