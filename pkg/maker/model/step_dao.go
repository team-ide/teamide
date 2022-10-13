package model

type StepDaoModel struct {
	*StepModel `json:",inline"`

	Dao        string      `json:"dao,omitempty"`  // 调用数据层
	Args       []*ArgModel `json:"args,omitempty"` // 调用参数
	SetVar     string      `json:"setVar,omitempty"`
	SetVarType string      `json:"setVarType,omitempty"`
}

var (
	docTemplateStepDaoName = "step_dao"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDaoName,
		Fields: []*docTemplateField{
			{Name: "dao", Comment: "数据层操作"},
			{Name: "args", Comment: "调用参数", IsList: true, StructName: docTemplateArgName},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
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
