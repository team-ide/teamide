package model

type StepScriptModel struct {
	*StepModel `json:",inline"`

	Script string `json:"script,omitempty"` // 脚本
}

var (
	docTemplateStepScriptName = "step_script"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepScriptName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "script",
				Comment: "script操作",
			},
		},
	})
}
