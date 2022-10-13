package model

type StepDaoModel struct {
	*StepModel `json:",inline"`

	Dao string `json:"dao,omitempty"` // 调用数据层
}

var (
	docTemplateStepDaoName = "step_dao"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDaoName,
		Fields: []*docTemplateField{
			{
				Name:    "dao",
				Comment: "数据层操作",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepDaoModel{}
		},
	})
}
