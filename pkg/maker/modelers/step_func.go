package modelers

type StepFuncModel struct {
	*StepModel `json:",inline"`

	Func       string      `json:"func,omitempty"` // 调用函数
	Args       []*ArgModel `json:"args,omitempty"` // 调用参数
	SetVar     string      `json:"setVar,omitempty"`
	SetVarType string      `json:"setVarType,omitempty"`
}

var (
	docTemplateStepFuncName = "step_func"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepFuncName,
		Fields: []*docTemplateField{
			{Name: "func", Comment: "函数操作"},
			{Name: "args", Comment: "调用参数", IsList: true, StructName: docTemplateArgName},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepFuncModel{}
		},
	})
}
