package modelers

type StepServiceModel struct {
	*StepModel `json:",inline"`

	Service    string      `json:"service,omitempty"` // 调用服务
	Args       []*ArgModel `json:"args,omitempty"`    // 调用参数
	SetVar     string      `json:"setVar,omitempty"`
	SetVarType string      `json:"setVarType,omitempty"`
}

var (
	docTemplateStepServiceName = "step_service"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepServiceName,
		Fields: []*docTemplateField{
			{Name: "service", Comment: "service操作"},
			{Name: "args", Comment: "调用参数", IsList: true, StructName: docTemplateArgName},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepServiceModel{}
		},
	})
}
