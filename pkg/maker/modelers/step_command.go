package modelers

type StepCommandModel struct {
	*StepModel `json:",inline"`

	Command string `json:"command,omitempty"` // 执行命令
}

var (
	docTemplateStepCommandName = "step_command"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepCommandName,
		Fields: []*docTemplateField{
			{Name: "command", Comment: "命令操作"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepCommandModel{}
		},
	})
}
