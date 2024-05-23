package modelers

type StepVarModel struct {
	*StepModel `json:",inline"`

	Var   string `json:"var,omitempty"`   // 定义变量
	Value string `json:"value,omitempty"` // 值
	Type  string `json:"type,omitempty"`  // 值类型
}

var (
	docTemplateStepVarName = "step_var"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepVarName,
		Fields: []*docTemplateField{
			{Name: "var", Comment: "定义变量操作"},
			{Name: "value", Comment: "变量值"},
			{Name: "type", Comment: "变量值类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepVarModel{}
		},
	})
}
