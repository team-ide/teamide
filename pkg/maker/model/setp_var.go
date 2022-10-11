package model

type StepVarModel struct {
	*StepModel `json:",inline"`

	Var       string `json:"var,omitempty"`       // 定义变量
	Value     string `json:"value,omitempty"`     // 值
	ValueType string `json:"valueType,omitempty"` // 值类型
}

var (
	docTemplateStepVarName = "step_var"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepVarName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "var",
				Comment: "定义变量操作",
			},
			{
				Name:    "value",
				Comment: "变量值",
			},
			{
				Name:    "valueType",
				Comment: "变量值类型",
			},
		},
	})
}
