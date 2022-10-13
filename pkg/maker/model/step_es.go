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
		Name: docTemplateStepEsName,
		Fields: []*docTemplateField{
			{Name: "es", Comment: "ES操作"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepEsModel{}
		},
	})
}
