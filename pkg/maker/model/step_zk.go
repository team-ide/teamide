package model

type StepZkModel struct {
	*StepModel `json:",inline"`

	Zk         string `json:"zk,omitempty"` // ZK操作
	Path       string `json:"path,omitempty"`
	Value      string `json:"value,omitempty"`
	Ephemeral  bool   `json:"ephemeral,omitempty"`
	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

var (
	docTemplateStepZkName = "step_zk"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepZkName,
		Fields: []*docTemplateField{
			{Name: "zk", Comment: "ZK操作"},
			{Name: "path", Comment: "路径"},
			{Name: "value", Comment: "操作的Value"},
			{Name: "ephemeral", Comment: "是否是临时"},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
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
