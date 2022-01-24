package model

type SqlInsert struct {
	Database string             `json:"database,omitempty" yaml:"database,omitempty"` // 库名
	Table    string             `json:"table,omitempty" yaml:"table,omitempty"`       // 表名
	Columns  []*SqlInsertColumn `json:"columns,omitempty" yaml:"columns,omitempty"`   // 新增字段
}

type SqlInsertColumn struct {
	IfScript      string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`           // 条件  满足该条件 添加
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`                   // 字段名称
	Value         string `json:"value,omitempty" yaml:"value,omitempty"`                 // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required      bool   `json:"required,omitempty" yaml:"required,omitempty"`           // 必填
	AutoIncrement bool   `json:"autoIncrement,omitempty" yaml:"autoIncrement,omitempty"` // 自增列
	IgnoreEmpty   bool   `json:"ignoreEmpty,omitempty" yaml:"ignoreEmpty,omitempty"`     // 忽略空值，如果忽略，则值是null、空字符串、0不设值
}

type ActionStepSqlInsert struct {
	Base *ActionStepBase

	SqlInsert        *SqlInsert `json:"sqlInsert,omitempty" yaml:"sqlInsert,omitempty"`               // 执行 SQL INSERT 操作
	VariableName     string     `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string     `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ActionStepSqlInsert) GetBase() *ActionStepBase {
	return this_.Base
}

func (this_ *ActionStepSqlInsert) SetBase(v *ActionStepBase) {
	this_.Base = v
}
