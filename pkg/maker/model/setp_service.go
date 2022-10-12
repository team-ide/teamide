package model

type StepServiceModel struct {
	*StepModel `json:",inline"`

	Service string `json:"service,omitempty"` // 调用服务
}

var (
	docTemplateStepServiceName = "step_service"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepServiceName,
		Fields: []*docTemplateField{
			{
				Name:    "service",
				Comment: "service操作",
			},
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
