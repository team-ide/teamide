package model

type StepZkModel struct {
	*StepModel `json:",inline"`

	Zk string `json:"zk,omitempty"` // ZK操作
}

var (
	docTemplateStepZkName = "step_zk"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepZkName,
		Fields: []*docTemplateField{
			{
				Name:    "zk",
				Comment: "ZK操作",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepZkModel{}
		},
	})
}
