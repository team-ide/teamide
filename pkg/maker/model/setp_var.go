package model

type StepVarModel struct {
	*StepModel

	Var       string `json:"var,omitempty"`       // 定义变量
	Value     string `json:"value,omitempty"`     // 值
	ValueType string `json:"valueType,omitempty"` // 值类型
}
