package model

type StepCommandModel struct {
	*StepModel `json:",inline"`

	Command string `json:"command,omitempty"` // 执行命令
}

var (
	docTemplateStepCommandName = "step_command"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepCommandName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "command",
				Comment: "命令操作",
			},
		},
	})
}
