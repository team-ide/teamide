package model

type StepErrorModel struct {
	*StepModel `json:",inline"`

	Error string `json:"error,omitempty"` // 异常
}

var (
	docTemplateStepErrorName = "step_error"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepErrorName,
		Fields: []*docTemplateField{
			{
				Name:    "error",
				Comment: "异常操作",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepErrorModel{}
		},
	})
}
