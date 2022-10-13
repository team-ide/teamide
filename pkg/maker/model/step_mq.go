package model

type StepMqModel struct {
	*StepModel `json:",inline"`

	Mq    string `json:"mq,omitempty"`    // MQ操作
	Topic string `json:"topic,omitempty"` // MQ topic

}

var (
	docTemplateStepMqName = "step_mq"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepMqName,
		Fields: []*docTemplateField{
			{
				Name:    "mq",
				Comment: "MQ操作",
			},
			{
				Name:    "topic",
				Comment: "MQ主题",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepMqModel{}
		},
	})
}
