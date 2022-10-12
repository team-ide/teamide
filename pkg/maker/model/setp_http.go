package model

type StepHttpModel struct {
	*StepModel `json:",inline"`

	Http string `json:"http,omitempty"` // HTTP操作
	Url  string `json:"url,omitempty"`  // HTTP地址
}

var (
	docTemplateStepHttpName = "step_http"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepHttpName,
		Fields: []*docTemplateField{
			{
				Name:    "http",
				Comment: "HTTP操作",
			},
			{
				Name:    "url",
				Comment: "地址",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepHttpModel{}
		},
	})
}
