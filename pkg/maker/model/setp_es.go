package model

type StepEsModel struct {
	*StepModel `json:",inline"`

	Es string `json:"es,omitempty"` // ES操作

}

var (
	docTemplateStepEsName = "step_es"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepEsName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "es",
				Comment: "ES操作",
			},
		},
	})
}
